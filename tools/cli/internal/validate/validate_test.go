package validate

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidator_OKAndFail(t *testing.T) {
	td := t.TempDir()

	schemaPath := filepath.Join(td, "schema.json")
	if err := os.WriteFile(schemaPath, []byte(`{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "properties": {
    "version": { "type": "integer" }
  },
  "required": ["version"]
}`), 0o644); err != nil {
		t.Fatalf("write schema: %v", err)
	}

	okYAML := filepath.Join(td, "ok.yaml")
	if err := os.WriteFile(okYAML, []byte("version: 1\n"), 0o644); err != nil {
		t.Fatalf("write yaml: %v", err)
	}

	badYAML := filepath.Join(td, "bad.yaml")
	if err := os.WriteFile(badYAML, []byte("system: {}\n"), 0o644); err != nil {
		t.Fatalf("write yaml: %v", err)
	}

	v, err := New(Options{MaxBytes: 1 << 20, SchemaPath: schemaPath})
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	res1, err := v.ValidateFile(okYAML)
	if err != nil {
		t.Fatalf("ValidateFile(ok): %v", err)
	}
	if !res1.OK {
		t.Fatalf("expected ok, got errors: %#v", res1.Errors)
	}

	res2, err := v.ValidateFile(badYAML)
	if err != nil {
		t.Fatalf("ValidateFile(bad): %v", err)
	}
	if res2.OK {
		t.Fatalf("expected failure")
	}
	if len(res2.Errors) == 0 {
		t.Fatalf("expected errors")
	}
}


