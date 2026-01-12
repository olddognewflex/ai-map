package main

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/olddognewflex/ai-map/tools/cli/internal/cli"
	"github.com/olddognewflex/ai-map/tools/cli/internal/input"
	"github.com/olddognewflex/ai-map/tools/cli/internal/validate"
	"github.com/spf13/cobra"
)

func newValidateCmd(stdout, stderr io.Writer) *cobra.Command {
	var sel input.Selection
	var schemaPath string

	cmd := &cobra.Command{
		Use:   "validate [--dir DIR] [--recursive] [files...]",
		Short: "Validate YAML files against the JSON Schema",
		RunE: func(cmd *cobra.Command, args []string) error {
			inputs, err := input.SelectFiles(sel, args)
			if err != nil {
				return cli.ExitError{Code: cli.ExitUsageOrConfig, Msg: "error: " + err.Error()}
			}
			if err := input.EnsureSelected(sel, inputs); err != nil {
				_ = cmd.Help()
				return cli.ExitError{Code: cli.ExitUsageOrConfig, Msg: "error: " + err.Error()}
			}

			v, err := validate.New(validate.Options{
				MaxBytes:    input.MaxYAMLBytes,
				SchemaPath:  schemaPath,
			})
			if err != nil {
				return cli.ExitError{
					Code: cli.ExitUsageOrConfig,
					Msg:  fmt.Sprintf("error: %s\nhint: if the repo does not yet include spec/ai-map.schema.json, provide --schema PATH (or approve adding a placeholder under spec/).", err),
				}
			}

			var failed bool
			for _, p := range inputs {
				res, err := v.ValidateFile(p)
				if err != nil {
					fmt.Fprintf(stderr, "%s: error: %s\n", p, err)
					return cli.ExitError{Code: cli.ExitInternalError}
				}
				if res.OK {
					continue
				}
				failed = true
				fmt.Fprintf(stderr, "%s: invalid\n", p)
				for _, e := range res.Errors {
					fmt.Fprintf(stderr, "  - %s\n", cli.TrimTrailingNewline(e))
				}
			}
			if failed {
				return cli.ExitError{Code: cli.ExitCheckFailed}
			}
			return nil
		},
	}

	cmd.SetOut(stdout)
	cmd.SetErr(stderr)

	cmd.Flags().StringVar(&schemaPath, "schema", filepath.Join("spec", "ai-map.schema.json"), "Path to JSON Schema")
	cmd.Flags().StringVar(&sel.Dir, "dir", "", "Directory to scan for *.yml|*.yaml (non-recursive by default)")
	cmd.Flags().BoolVar(&sel.Recursive, "recursive", false, "Scan directories recursively (off by default)")
	return cmd
}


