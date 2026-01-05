package validate

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/santhosh-tekuri/jsonschema/v5"
	"gopkg.in/yaml.v3"
)

type Options struct {
	// MaxBytes caps YAML input reads.
	MaxBytes int64
	// SchemaPath points to a JSON Schema file on disk.
	SchemaPath string
}

type Result struct {
	OK     bool
	Errors []string
}

type Validator struct {
	schema *jsonschema.Schema
	opt    Options
}

func New(opt Options) (*Validator, error) {
	if strings.TrimSpace(opt.SchemaPath) == "" {
		return nil, errors.New("schema path is required")
	}
	if opt.MaxBytes <= 0 {
		return nil, errors.New("max bytes must be > 0")
	}

	s, err := loadSchemaFromFile(opt.SchemaPath)
	if err != nil {
		return nil, err
	}
	return &Validator{schema: s, opt: opt}, nil
}

func (v *Validator) ValidateFile(path string) (Result, error) {
	b, err := readFileWithLimit(path, v.opt.MaxBytes)
	if err != nil {
		return Result{}, err
	}

	var doc any
	if err := yaml.Unmarshal(b, &doc); err != nil {
		return Result{OK: false, Errors: []string{fmt.Sprintf("YAML parse error: %s", err)}}, nil
	}
	jsonReady, err := yamlToJSONReady(doc)
	if err != nil {
		return Result{OK: false, Errors: []string{fmt.Sprintf("YAML normalization error: %s", err)}}, nil
	}

	// Validate expects JSON-compatible types.
	if err := v.schema.Validate(jsonReady); err != nil {
		return Result{OK: false, Errors: flattenSchemaError(err)}, nil
	}
	return Result{OK: true}, nil
}

func loadSchemaFromFile(schemaPath string) (*jsonschema.Schema, error) {
	abs, err := filepath.Abs(schemaPath)
	if err != nil {
		return nil, fmt.Errorf("invalid schema path: %w", err)
	}
	info, err := os.Stat(abs)
	if err != nil {
		return nil, fmt.Errorf("cannot stat schema %q: %w", schemaPath, err)
	}
	if info.IsDir() {
		return nil, fmt.Errorf("schema path is a directory: %q", schemaPath)
	}

	u := (&url.URL{Scheme: "file", Path: filepath.ToSlash(abs)}).String()
	compiler := jsonschema.NewCompiler()

	// Security/determinism: disallow network loads.
	compiler.LoadURL = func(raw string) (io.ReadCloser, error) {
		pu, err := url.Parse(raw)
		if err != nil {
			return nil, err
		}
		if pu.Scheme != "" && pu.Scheme != "file" {
			return nil, fmt.Errorf("disallowed schema ref scheme %q in %q", pu.Scheme, raw)
		}

		// Resolve relative file paths relative to the referring schema file directory.
		p := pu.Path
		if p == "" {
			p = raw
		}
		if !filepath.IsAbs(p) {
			p = filepath.Join(filepath.Dir(abs), p)
		}
		p = filepath.Clean(p)
		return os.Open(p)
	}

	f, err := os.Open(abs)
	if err != nil {
		return nil, fmt.Errorf("cannot open schema: %w", err)
	}
	defer f.Close()

	schemaBytes, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("cannot read schema: %w", err)
	}
	if err := compiler.AddResource(u, bytes.NewReader(schemaBytes)); err != nil {
		return nil, fmt.Errorf("cannot add schema: %w", err)
	}
	s, err := compiler.Compile(u)
	if err != nil {
		return nil, fmt.Errorf("cannot compile schema: %w", err)
	}
	return s, nil
}

// readFileWithLimit is intentionally duplicated in this package to avoid exporting internal/cli helpers.
func readFileWithLimit(path string, maxBytes int64) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if info.Size() > maxBytes {
		return nil, fmt.Errorf("file too large (%d bytes > %d): %s", info.Size(), maxBytes, path)
	}
	return os.ReadFile(path)
}

func yamlToJSONReady(v any) (any, error) {
	switch x := v.(type) {
	case nil, bool, string:
		return x, nil
	case int, int8, int16, int32, int64:
		return x, nil
	case uint, uint8, uint16, uint32, uint64:
		return x, nil
	case float32, float64:
		return x, nil
	case []any:
		out := make([]any, 0, len(x))
		for _, it := range x {
			cv, err := yamlToJSONReady(it)
			if err != nil {
				return nil, err
			}
			out = append(out, cv)
		}
		return out, nil
	case map[string]any:
		out := make(map[string]any, len(x))
		for k, v := range x {
			cv, err := yamlToJSONReady(v)
			if err != nil {
				return nil, err
			}
			out[k] = cv
		}
		return out, nil
	case map[any]any:
		out := make(map[string]any, len(x))
		for k, v := range x {
			ks, ok := k.(string)
			if !ok {
				return nil, fmt.Errorf("non-string map key %T", k)
			}
			cv, err := yamlToJSONReady(v)
			if err != nil {
				return nil, err
			}
			out[ks] = cv
		}
		return out, nil
	default:
		// Fallback: attempt JSON roundtrip for scalar-ish types (e.g. time.Time) but keep it deterministic.
		b, err := json.Marshal(x)
		if err != nil {
			return nil, fmt.Errorf("unsupported YAML value %T", x)
		}
		var out any
		if err := json.Unmarshal(b, &out); err != nil {
			return nil, fmt.Errorf("unsupported YAML value %T", x)
		}
		return out, nil
	}
}

func flattenSchemaError(err error) []string {
	// jsonschema/v5 returns a rich error type; stringifying is stable enough for an MVP.
	// We normalize line endings later at the caller boundary by writing \n only.
	if err == nil {
		return nil
	}
	// Some errors can be multi-error (ValidationError with Causes); try to unwrap.
	var ve *jsonschema.ValidationError
	if errors.As(err, &ve) {
		var out []string
		out = append(out, fmt.Sprintf("%s: %s", ve.InstanceLocation, ve.Message))
		for _, c := range ve.Causes {
			out = append(out, fmt.Sprintf("%s: %s", c.InstanceLocation, c.Message))
		}
		return out
	}
	return []string{err.Error()}
}


