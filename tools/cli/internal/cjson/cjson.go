package cjson

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
)

// MarshalIndent marshals v to deterministic JSON with stable map key ordering.
// It supports the JSON-compatible subset produced by our YAML normalization.
func MarshalIndent(v any, prefix, indent string) ([]byte, error) {
	var buf bytes.Buffer
	w := &writer{w: &buf, prefix: prefix, indent: indent}
	if err := w.writeValue(v, 0); err != nil {
		return nil, err
	}
	buf.WriteByte('\n')
	return buf.Bytes(), nil
}

type writer struct {
	w      *bytes.Buffer
	prefix string
	indent string
}

func (w *writer) nl(depth int) {
	w.w.WriteByte('\n')
	w.w.WriteString(w.prefix)
	for i := 0; i < depth; i++ {
		w.w.WriteString(w.indent)
	}
}

func (w *writer) writeValue(v any, depth int) error {
	switch x := v.(type) {
	case nil:
		w.w.WriteString("null")
		return nil
	case bool:
		if x {
			w.w.WriteString("true")
		} else {
			w.w.WriteString("false")
		}
		return nil
	case string:
		b, _ := json.Marshal(x)
		w.w.Write(b)
		return nil
	case float64, float32, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		b, err := json.Marshal(x)
		if err != nil {
			return err
		}
		w.w.Write(b)
		return nil
	case []any:
		w.w.WriteByte('[')
		if len(x) == 0 {
			w.w.WriteByte(']')
			return nil
		}
		for i, it := range x {
			w.nl(depth + 1)
			if err := w.writeValue(it, depth+1); err != nil {
				return err
			}
			if i != len(x)-1 {
				w.w.WriteByte(',')
			}
		}
		w.nl(depth)
		w.w.WriteByte(']')
		return nil
	case map[string]any:
		keys := make([]string, 0, len(x))
		for k := range x {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		w.w.WriteByte('{')
		if len(keys) == 0 {
			w.w.WriteByte('}')
			return nil
		}
		for i, k := range keys {
			w.nl(depth + 1)
			kb, _ := json.Marshal(k)
			w.w.Write(kb)
			w.w.WriteString(": ")
			if err := w.writeValue(x[k], depth+1); err != nil {
				return err
			}
			if i != len(keys)-1 {
				w.w.WriteByte(',')
			}
		}
		w.nl(depth)
		w.w.WriteByte('}')
		return nil
	default:
		return fmt.Errorf("unsupported type for canonical json: %T", v)
	}
}


