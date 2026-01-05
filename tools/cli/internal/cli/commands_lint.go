package cli

import (
	"fmt"
	"io"

	"github.com/olddognewflex/ai-map/tools/cli/internal/lint"
)

func runLint(stdout, stderr io.Writer, args []string) int {
	fs := newFlagSet("lint", stderr)
	var sel fileSelection
	addFileSelectionFlags(fs, &sel)

	help, rest, err := parseCommon(fs, args)
	if help {
		fmt.Fprint(stdout, lintHelpText())
		return ExitOK
	}
	if err != nil {
		fmt.Fprint(stderr, lintHelpText())
		return ExitUsageOrConfig
	}

	inputs, err := selectInputFiles(sel, rest)
	if err != nil {
		fmt.Fprintf(stderr, "error: %s\n", err)
		return ExitUsageOrConfig
	}
	if err := ensureInputsOrDirSelected(sel, inputs); err != nil {
		fmt.Fprintf(stderr, "error: %s\n\n", err)
		fmt.Fprint(stderr, lintHelpText())
		return ExitUsageOrConfig
	}

	var hadErrors bool
	for _, p := range inputs {
		b, err := readFileWithLimit(p, MaxYAMLBytes)
		if err != nil {
			fmt.Fprintf(stderr, "%s: error: %s\n", p, err)
			return ExitUsageOrConfig
		}

		res := lint.LintYAMLBytes(b)
		for _, is := range res.Issues {
			if is.Path != "" {
				fmt.Fprintf(stderr, "%s: %s: %s (%s)\n", p, is.Severity, is.Message, is.Path)
			} else {
				fmt.Fprintf(stderr, "%s: %s: %s\n", p, is.Severity, is.Message)
			}
			if is.Severity == lint.SeverityError {
				hadErrors = true
			}
		}
	}
	if hadErrors {
		return ExitCheckFailed
	}
	return ExitOK
}

func lintHelpText() string {
	return "" +
		"Usage:\n" +
		"  ai-map lint [--dir DIR] [--recursive] [files...]\n\n" +
		"Runs opinionated checks on AI-Map YAML files.\n\n" +
		"Flags:\n" +
		"  --dir string        Directory to scan for *.yml|*.yaml (non-recursive by default)\n" +
		"  --recursive         Scan directories recursively (off by default)\n" +
		"  -h, --help          Show help\n"
}
