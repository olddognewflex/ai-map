package main

import (
	"io"
	"os"
	"path/filepath"

	"github.com/olddognewflex/ai-map/tools/cli/internal/cli"
	"github.com/olddognewflex/ai-map/tools/cli/internal/input"
	"github.com/olddognewflex/ai-map/tools/cli/internal/render"
	"github.com/spf13/cobra"
)

func newRenderCmd(stdout, stderr io.Writer) *cobra.Command {
	var sel input.Selection
	var outPath string
	var title string

	cmd := &cobra.Command{
		Use:   "render [--out FILE] [--title TITLE] [--dir DIR] [--recursive] [files...]",
		Short: "Render AI-Map docs (Markdown)",
		RunE: func(cmd *cobra.Command, args []string) error {
			inputs, err := input.SelectFiles(sel, args)
			if err != nil {
				return cli.ExitError{Code: cli.ExitUsageOrConfig, Msg: "error: " + err.Error()}
			}
			if err := input.EnsureSelected(sel, inputs); err != nil {
				_ = cmd.Help()
				return cli.ExitError{Code: cli.ExitUsageOrConfig, Msg: "error: " + err.Error()}
			}

			if outPath != "" && len(inputs) != 1 {
				return cli.ExitError{Code: cli.ExitUsageOrConfig, Msg: "error: --out requires exactly one input file"}
			}

			outBytes, err := func() ([]byte, error) {
				if len(inputs) == 1 {
					b, err := input.ReadFileWithLimit(inputs[0], input.MaxYAMLBytes)
					if err != nil {
						return nil, err
					}
					return render.MarkdownFromYAML(b, render.Options{Title: title})
				}
				var all []byte
				for i, p := range inputs {
					b, err := input.ReadFileWithLimit(p, input.MaxYAMLBytes)
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
				return cli.ExitError{Code: cli.ExitCheckFailed, Msg: "error: " + err.Error()}
			}

			if outPath == "" {
				if _, err := stdout.Write(outBytes); err != nil {
					return cli.ExitError{Code: cli.ExitInternalError, Msg: "error: " + err.Error()}
				}
				return nil
			}

			absOut, err := filepath.Abs(outPath)
			if err != nil {
				return cli.ExitError{Code: cli.ExitUsageOrConfig, Msg: "error: invalid --out: " + err.Error()}
			}
			if _, err := os.Stat(absOut); err == nil {
				return cli.ExitError{Code: cli.ExitUsageOrConfig, Msg: "error: refusing to overwrite existing file: " + absOut}
			}
			if err := os.MkdirAll(filepath.Dir(absOut), 0o755); err != nil {
				return cli.ExitError{Code: cli.ExitInternalError, Msg: "error: cannot create output dir: " + err.Error()}
			}
			if err := os.WriteFile(absOut, outBytes, 0o644); err != nil {
				return cli.ExitError{Code: cli.ExitInternalError, Msg: "error: cannot write output: " + err.Error()}
			}
			return nil
		},
	}

	cmd.SetOut(stdout)
	cmd.SetErr(stderr)
	cmd.Flags().StringVar(&outPath, "out", "", "Output file (defaults to stdout)")
	cmd.Flags().StringVar(&title, "title", "", "Document title (optional)")
	cmd.Flags().StringVar(&sel.Dir, "dir", "", "Directory to scan for *.yml|*.yaml (non-recursive by default)")
	cmd.Flags().BoolVar(&sel.Recursive, "recursive", false, "Scan directories recursively (off by default)")
	return cmd
}


