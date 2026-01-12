package lint

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

type Severity string

const (
	SeverityError Severity = "error"
	SeverityWarn  Severity = "warn"
)

type Issue struct {
	Severity Severity
	Message  string
	Path     string
}

type Result struct {
	Issues []Issue
}

func (r Result) HasErrors() bool {
	for _, it := range r.Issues {
		if it.Severity == SeverityError {
			return true
		}
	}
	return false
}

type Options struct {
	MaxBytes int64
}

func LintYAMLBytes(b []byte) Result {
	var doc any
	if err := yaml.Unmarshal(b, &doc); err != nil {
		return Result{Issues: []Issue{{
			Severity: SeverityError,
			Message:  fmt.Sprintf("YAML parse error: %s", err),
		}}}
	}
	m, ok := doc.(map[string]any)
	if !ok {
		// yaml.v3 often uses map[any]any; tolerate and convert at this boundary.
		m2, ok2 := doc.(map[any]any)
		if !ok2 {
			return Result{Issues: []Issue{{Severity: SeverityError, Message: "top-level document must be a mapping/object"}}}
		}
		m = make(map[string]any, len(m2))
		for k, v := range m2 {
			ks, ok := k.(string)
			if !ok {
				return Result{Issues: []Issue{{Severity: SeverityError, Message: fmt.Sprintf("non-string top-level key %T", k)}}}
			}
			m[ks] = v
		}
	}

	var issues []Issue

	// Required: version (spec says required)
	if _, ok := m["version"]; !ok {
		issues = append(issues, Issue{Severity: SeverityError, Path: "version", Message: "missing required field"})
	} else {
		switch v := m["version"].(type) {
		case int, int64, uint64, uint, float64:
			// ok-ish; we don't enforce exact '1' here (spec version might evolve).
		default:
			issues = append(issues, Issue{Severity: SeverityError, Path: "version", Message: fmt.Sprintf("must be a number, got %T", v)})
		}
	}

	// Required: system (spec says required)
	sys, ok := m["system"]
	if !ok {
		issues = append(issues, Issue{Severity: SeverityError, Path: "system", Message: "missing required field"})
	} else {
		sm, ok := asStringMap(sys)
		if !ok {
			issues = append(issues, Issue{Severity: SeverityError, Path: "system", Message: "must be an object"})
		} else {
			// system.name: required, non-empty, slug-ish
			if name, ok := sm["name"].(string); !ok || strings.TrimSpace(name) == "" {
				issues = append(issues, Issue{Severity: SeverityError, Path: "system.name", Message: "missing or empty"})
			} else if strings.ContainsAny(name, " \t\r\n") {
				issues = append(issues, Issue{Severity: SeverityWarn, Path: "system.name", Message: "should not contain whitespace"})
			}

			// system.type: if present, must be one of known values
			if tRaw, ok := sm["type"]; ok {
				t, ok := tRaw.(string)
				if !ok {
					issues = append(issues, Issue{Severity: SeverityError, Path: "system.type", Message: "must be a string"})
				} else if !isAllowedType(t) {
					issues = append(issues, Issue{Severity: SeverityWarn, Path: "system.type", Message: "unknown value (expected one of service|webapp|library|infra|monorepo)"})
				}
			}
		}
	}

	return Result{Issues: issues}
}

func asStringMap(v any) (map[string]any, bool) {
	switch m := v.(type) {
	case map[string]any:
		return m, true
	case map[any]any:
		out := make(map[string]any, len(m))
		for k, v := range m {
			ks, ok := k.(string)
			if !ok {
				return nil, false
			}
			out[ks] = v
		}
		return out, true
	default:
		return nil, false
	}
}

func isAllowedType(s string) bool {
	switch s {
	case "service", "webapp", "library", "infra", "monorepo":
		return true
	default:
		return false
	}
}


