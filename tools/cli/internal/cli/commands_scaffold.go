package cli

import (
	"fmt"
	"io"
)

func runScaffold(stdout, stderr io.Writer, args []string) int {
	fs := newFlagSet("scaffold", stderr)
	var outDir string
	fs.StringVar(&outDir, "out", "", "Output directory to create (required)")
	var name string
	fs.StringVar(&name, "name", "", "System name (optional; used in template)")

	help, _, err := parseCommon(fs, args)
	if help {
		fmt.Fprint(stdout, scaffoldHelpText())
		return ExitOK
	}
	if err != nil {
		fmt.Fprint(stderr, scaffoldHelpText())
		return ExitUsageOrConfig
	}

	if outDir == "" {
		fmt.Fprint(stderr, "error: --out is required\n\n")
		fmt.Fprint(stderr, scaffoldHelpText())
		return ExitUsageOrConfig
	}

	_ = name
	fmt.Fprintln(stderr, "scaffold: not implemented yet (safe skeleton writer will be added later)")
	return ExitInternalError
}

func scaffoldHelpText() string {
	return "" +
		"Usage:\n" +
		"  ai-map scaffold --out DIR [--name NAME]\n\n" +
		"Creates a new AI-Map folder skeleton safely.\n\n" +
		"Flags:\n" +
		"  --out string        Output directory to create (required)\n" +
		"  --name string       System name (optional; used in template)\n" +
		"  -h, --help          Show help\n"
}


