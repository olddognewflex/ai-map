package render

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/olddognewflex/ai-map/tools/cli/internal/cjson"
	"gopkg.in/yaml.v3"
)

type Options struct {
	Title string
}

func MarkdownFromYAML(yamlBytes []byte, opt Options) ([]byte, error) {
	var doc any
	if err := yaml.Unmarshal(yamlBytes, &doc); err != nil {
		return nil, fmt.Errorf("YAML parse error: %w", err)
	}
	jsonReady, err := yamlToJSONReady(doc)
	if err != nil {
		return nil, err
	}

	title := strings.TrimSpace(opt.Title)
	if title == "" {
		title = "AI-Map"
	}

	canon, err := cjson.MarshalIndent(jsonReady, "", "  ")
	if err != nil {
		return nil, err
	}

	var b bytes.Buffer
	b.WriteString("# ")
	b.WriteString(title)
	b.WriteString("\n\n")
	b.WriteString("## Raw (canonical JSON)\n\n")
	b.WriteString("```json\n")
	b.Write(canon)
	b.WriteString("```\n")
	return b.Bytes(), nil
}

func yamlToJSONReady(v any) (any, error) {
	// Keep this small and deterministic; we only need JSON-compatible output for rendering.
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
		return nil, fmt.Errorf("unsupported YAML value %T", v)
	}
}


