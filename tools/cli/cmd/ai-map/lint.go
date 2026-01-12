package main

import (
	"fmt"
	"io"

	"github.com/olddognewflex/ai-map/tools/cli/internal/cli"
	"github.com/olddognewflex/ai-map/tools/cli/internal/input"
	"github.com/olddognewflex/ai-map/tools/cli/internal/lint"
	"github.com/spf13/cobra"
)

func newLintCmd(stdout, stderr io.Writer) *cobra.Command {
	var sel input.Selection

	cmd := &cobra.Command{
		Use:   "lint [--dir DIR] [--recursive] [files...]",
		Short: "Run opinionated checks",
		RunE: func(cmd *cobra.Command, args []string) error {
			inputs, err := input.SelectFiles(sel, args)
			if err != nil {
				return cli.ExitError{Code: cli.ExitUsageOrConfig, Msg: "error: " + err.Error()}
			}
			if err := input.EnsureSelected(sel, inputs); err != nil {
				_ = cmd.Help()
				return cli.ExitError{Code: cli.ExitUsageOrConfig, Msg: "error: " + err.Error()}
			}

			var hadErrors bool
			for _, p := range inputs {
				b, err := input.ReadFileWithLimit(p, input.MaxYAMLBytes)
				if err != nil {
					return cli.ExitError{Code: cli.ExitUsageOrConfig, Msg: fmt.Sprintf("%s: error: %s", p, err)}
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
				return cli.ExitError{Code: cli.ExitCheckFailed}
			}
			return nil
		},
	}

	cmd.SetOut(stdout)
	cmd.SetErr(stderr)
	cmd.Flags().StringVar(&sel.Dir, "dir", "", "Directory to scan for *.yml|*.yaml (non-recursive by default)")
	cmd.Flags().BoolVar(&sel.Recursive, "recursive", false, "Scan directories recursively (off by default)")
	return cmd
}


