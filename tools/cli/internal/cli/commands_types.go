package cli

import (
	"fmt"
	"io"
)

func runTypes(stdout, stderr io.Writer, args []string) int {
	fs := newFlagSet("types", stderr)
	var sel fileSelection
	addFileSelectionFlags(fs, &sel)
	var lang string
	fs.StringVar(&lang, "lang", "go", "Target language (go only for now)")

	help, rest, err := parseCommon(fs, args)
	if help {
		fmt.Fprint(stdout, typesHelpText())
		return ExitOK
	}
	if err != nil {
		fmt.Fprint(stderr, typesHelpText())
		return ExitUsageOrConfig
	}

	inputs, err := selectInputFiles(sel, rest)
	if err != nil {
		fmt.Fprintf(stderr, "error: %s\n", err)
		return ExitUsageOrConfig
	}
	if err := ensureInputsOrDirSelected(sel, inputs); err != nil {
		fmt.Fprintf(stderr, "error: %s\n\n", err)
		fmt.Fprint(stderr, typesHelpText())
		return ExitUsageOrConfig
	}

	_ = lang
	_ = inputs
	fmt.Fprintln(stderr, "types: not implemented yet (MVP type generation will be added later)")
	return ExitInternalError
}

func typesHelpText() string {
	return "" +
		"Usage:\n" +
		"  ai-map types [--lang go] [--dir DIR] [--recursive] [files...]\n\n" +
		"Generates types from AI-Map YAML (MVP).\n\n" +
		"Flags:\n" +
		"  --lang string       Target language (go only for now) (default \"go\")\n" +
		"  --dir string        Directory to scan for *.yml|*.yaml (non-recursive by default)\n" +
		"  --recursive         Scan directories recursively (off by default)\n" +
		"  -h, --help          Show help\n"
}


