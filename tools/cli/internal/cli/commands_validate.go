package cli

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/olddognewflex/ai-map/tools/cli/internal/validate"
)

func runValidate(stdout, stderr io.Writer, args []string) int {
	fs := newFlagSet("validate", stderr)
	var sel fileSelection
	addFileSelectionFlags(fs, &sel)
	var schemaPath string
	fs.StringVar(&schemaPath, "schema", defaultSchemaPath(), "Path to JSON Schema")

	help, rest, err := parseCommon(fs, args)
	if help {
		fmt.Fprint(stdout, validateHelpText())
		return ExitOK
	}
	if err != nil {
		fmt.Fprint(stderr, validateHelpText())
		return ExitUsageOrConfig
	}

	inputs, err := selectInputFiles(sel, rest)
	if err != nil {
		fmt.Fprintf(stderr, "error: %s\n", err)
		return ExitUsageOrConfig
	}
	if err := ensureInputsOrDirSelected(sel, inputs); err != nil {
		fmt.Fprintf(stderr, "error: %s\n\n", err)
		fmt.Fprint(stderr, validateHelpText())
		return ExitUsageOrConfig
	}

	v, err := validate.New(validate.Options{
		MaxBytes:   MaxYAMLBytes,
		SchemaPath: schemaPath,
	})
	if err != nil {
		fmt.Fprintf(stderr, "error: %s\n", err)
		fmt.Fprintln(stderr, "hint: if the repo does not yet include spec/ai-map.schema.json, provide --schema PATH (or approve adding a placeholder under spec/).")
		return ExitUsageOrConfig
	}

	var failed bool
	for _, p := range inputs {
		res, err := v.ValidateFile(p)
		if err != nil {
			fmt.Fprintf(stderr, "%s: error: %s\n", p, err)
			return ExitInternalError
		}
		if res.OK {
			continue
		}
		failed = true
		fmt.Fprintf(stderr, "%s: invalid\n", p)
		for _, e := range res.Errors {
			fmt.Fprintf(stderr, "  - %s\n", trimTrailingNewline(e))
		}
	}
	if failed {
		return ExitCheckFailed
	}
	return ExitOK
}

func validateHelpText() string {
	return "" +
		"Usage:\n" +
		"  ai-map validate [--schema FILE] [--dir DIR] [--recursive] [files...]\n\n" +
		"Validates one or more YAML files against the AI-Map JSON Schema.\n\n" +
		"Flags:\n" +
		"  --schema string     Path to JSON Schema (default \"spec/ai-map.schema.json\")\n" +
		"  --dir string        Directory to scan for *.yml|*.yaml (non-recursive by default)\n" +
		"  --recursive         Scan directories recursively (off by default)\n" +
		"  -h, --help          Show help\n"
}

func defaultSchemaPath() string {
	// Intentionally relative: keep behavior simple and predictable.
	// Users can override with --schema and absolute paths.
	return filepath.Join("spec", "ai-map.schema.json")
}


