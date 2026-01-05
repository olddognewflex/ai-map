package cli

import (
	"fmt"
	"io"
)

func runConformance(stdout, stderr io.Writer, args []string) int {
	fs := newFlagSet("conformance", stderr)
	var updateGolden bool
	fs.BoolVar(&updateGolden, "update-golden", false, "Update golden files (off by default)")

	help, _, err := parseCommon(fs, args)
	if help {
		fmt.Fprint(stdout, conformanceHelpText())
		return ExitOK
	}
	if err != nil {
		fmt.Fprint(stderr, conformanceHelpText())
		return ExitUsageOrConfig
	}

	_ = updateGolden
	fmt.Fprintln(stderr, "conformance: not implemented yet (fixture runner + golden tests will be added after validate/lint)")
	return ExitInternalError
}

func conformanceHelpText() string {
	return "" +
		"Usage:\n" +
		"  ai-map conformance [--update-golden]\n\n" +
		"Runs conformance fixtures and compares output to golden files.\n\n" +
		"Flags:\n" +
		"  --update-golden     Update golden files (off by default)\n" +
		"  -h, --help          Show help\n"
}


