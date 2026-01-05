package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/olddognewflex/ai-map/tools/cli/internal/render"
)

func runRender(stdout, stderr io.Writer, args []string) int {
	fs := newFlagSet("render", stderr)
	var sel fileSelection
	addFileSelectionFlags(fs, &sel)
	var outPath string
	fs.StringVar(&outPath, "out", "", "Output file (defaults to stdout)")
	var title string
	fs.StringVar(&title, "title", "", "Document title (optional)")

	help, rest, err := parseCommon(fs, args)
	if help {
		fmt.Fprint(stdout, renderHelpText())
		return ExitOK
	}
	if err != nil {
		fmt.Fprint(stderr, renderHelpText())
		return ExitUsageOrConfig
	}

	inputs, err := selectInputFiles(sel, rest)
	if err != nil {
		fmt.Fprintf(stderr, "error: %s\n", err)
		return ExitUsageOrConfig
	}
	if err := ensureInputsOrDirSelected(sel, inputs); err != nil {
		fmt.Fprintf(stderr, "error: %s\n\n", err)
		fmt.Fprint(stderr, renderHelpText())
		return ExitUsageOrConfig
	}

	if outPath != "" && len(inputs) != 1 {
		fmt.Fprintln(stderr, "error: --out requires exactly one input file")
		return ExitUsageOrConfig
	}

	// If title not provided, derive from system.name if available; otherwise from filename.
	docTitle := strings.TrimSpace(title)
	if docTitle == "" && len(inputs) == 1 {
		docTitle = "AI-Map: " + filepath.Base(inputs[0])
	}

	outBytes, err := func() ([]byte, error) {
		// For multiple inputs to stdout, concatenate with a deterministic separator.
		if len(inputs) == 1 {
			b, err := readFileWithLimit(inputs[0], MaxYAMLBytes)
			if err != nil {
				return nil, err
			}
			return render.MarkdownFromYAML(b, render.Options{Title: docTitle})
		}
		var all []byte
		for i, p := range inputs {
			b, err := readFileWithLimit(p, MaxYAMLBytes)
			if err != nil {
				return nil, err
			}
			md, err := render.MarkdownFromYAML(b, render.Options{Title: "AI-Map: " + filepath.Base(p)})
			if err != nil {
				return nil, err
			}
			if i > 0 {
				all = append(all, []byte("\n---\n\n")...)
			}
			all = append(all, md...)
		}
		return all, nil
	}()
	if err != nil {
		fmt.Fprintf(stderr, "error: %s\n", err)
		return ExitCheckFailed
	}

	if outPath == "" {
		if _, err := stdout.Write(outBytes); err != nil {
			fmt.Fprintf(stderr, "error: %s\n", err)
			return ExitInternalError
		}
		return ExitOK
	}

	absOut, err := filepath.Abs(outPath)
	if err != nil {
		fmt.Fprintf(stderr, "error: invalid --out: %s\n", err)
		return ExitUsageOrConfig
	}
	if _, err := os.Stat(absOut); err == nil {
		fmt.Fprintf(stderr, "error: refusing to overwrite existing file: %s\n", absOut)
		return ExitUsageOrConfig
	}
	if err := os.MkdirAll(filepath.Dir(absOut), 0o755); err != nil {
		fmt.Fprintf(stderr, "error: cannot create output dir: %s\n", err)
		return ExitInternalError
	}
	if err := os.WriteFile(absOut, outBytes, 0o644); err != nil {
		fmt.Fprintf(stderr, "error: cannot write output: %s\n", err)
		return ExitInternalError
	}
	return ExitOK
}

func renderHelpText() string {
	return "" +
		"Usage:\n" +
		"  ai-map render [--out FILE] [--dir DIR] [--recursive] [files...]\n\n" +
		"Renders AI-Map YAML to Markdown documentation.\n\n" +
		"Flags:\n" +
		"  --out string        Output file (defaults to stdout)\n" +
		"  --title string      Document title (optional)\n" +
		"  --dir string        Directory to scan for *.yml|*.yaml (non-recursive by default)\n" +
		"  --recursive         Scan directories recursively (off by default)\n" +
		"  -h, --help          Show help\n"
}


