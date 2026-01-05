package cli

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/olddognewflex/ai-map/tools/cli/internal/version"
)

func Run(argv []string) int {
	if len(argv) == 0 {
		return ExitInternalError
	}

	// Global help.
	if len(argv) == 1 {
		printRootHelp(os.Stdout)
		return ExitUsageOrConfig
	}
	if isHelpFlag(argv[1]) {
		printRootHelp(os.Stdout)
		return ExitOK
	}

	cmd := argv[1]
	args := argv[2:]

	switch cmd {
	case "validate":
		return runValidate(os.Stdout, os.Stderr, args)
	case "lint":
		return runLint(os.Stdout, os.Stderr, args)
	case "render":
		return runRender(os.Stdout, os.Stderr, args)
	case "types":
		return runTypes(os.Stdout, os.Stderr, args)
	case "conformance":
		return runConformance(os.Stdout, os.Stderr, args)
	case "scaffold":
		return runScaffold(os.Stdout, os.Stderr, args)
	case "version":
		fmt.Fprintf(os.Stdout, "ai-map %s\n", version.Version)
		if version.Commit != "" {
			fmt.Fprintf(os.Stdout, "commit %s\n", version.Commit)
		}
		if version.Date != "" {
			fmt.Fprintf(os.Stdout, "date   %s\n", version.Date)
		}
		return ExitOK
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", cmd)
		printRootHelp(os.Stderr)
		return ExitUsageOrConfig
	}
}

func isHelpFlag(s string) bool {
	return s == "-h" || s == "--help" || s == "help"
}

func newFlagSet(name string, out io.Writer) *flag.FlagSet {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.SetOutput(out) // errors from parsing go here; we keep help separate.
	return fs
}

func parseCommon(fs *flag.FlagSet, args []string) (help bool, rest []string, err error) {
	var h bool
	fs.BoolVar(&h, "help", false, "Show help")
	fs.BoolVar(&h, "h", false, "Show help (shorthand)")

	if err := fs.Parse(args); err != nil {
		// flag already wrote a parse error; we just return.
		return false, nil, err
	}
	if h {
		return true, fs.Args(), nil
	}
	// Also accept --help/-h without defining both.
	for _, a := range args {
		if isHelpFlag(a) {
			return true, fs.Args(), nil
		}
	}
	return false, fs.Args(), nil
}

func ensureInputsOrDirSelected(sel fileSelection, inputs []string) error {
	if sel.Dir == "" && len(inputs) == 0 {
		return fmt.Errorf("no input files provided")
	}
	return nil
}

func trimTrailingNewline(s string) string {
	return strings.TrimRight(s, "\r\n")
}


