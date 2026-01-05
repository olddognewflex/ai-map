package main

import (
	"io"
	"os"
	"path/filepath"

	"github.com/olddognewflex/ai-map/tools/cli/internal/cli"
	"github.com/olddognewflex/ai-map/tools/cli/internal/typesgen"
	"github.com/spf13/cobra"
)

func newTypesCmd(stdout, stderr io.Writer) *cobra.Command {
	var lang string
	var outPath string
	var pkg string

	cmd := &cobra.Command{
		Use:   "types [--lang go] [--pkg NAME] [--out FILE]",
		Short: "Generate types (MVP)",
		Long:  "Generates types for AI-Map v1.0 (MVP; not schema-derived yet).",
		RunE: func(cmd *cobra.Command, args []string) error {
			if lang != "go" {
				return cli.ExitError{Code: cli.ExitUsageOrConfig, Msg: "error: unsupported --lang (only \"go\" is supported)"}
			}

			out := typesgen.GenerateGo(typesgen.Options{Package: pkg})

			if outPath == "" {
				_, _ = stdout.Write(out)
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
			if err := os.WriteFile(absOut, out, 0o644); err != nil {
				return cli.ExitError{Code: cli.ExitInternalError, Msg: "error: cannot write output: " + err.Error()}
			}
			return nil
		},
	}

	cmd.SetOut(stdout)
	cmd.SetErr(stderr)
	cmd.Flags().StringVar(&lang, "lang", "go", "Target language (go only for now)")
	cmd.Flags().StringVar(&pkg, "pkg", "aimap", "Go package name (go only)")
	cmd.Flags().StringVar(&outPath, "out", "", "Output file (defaults to stdout)")
	return cmd
}


