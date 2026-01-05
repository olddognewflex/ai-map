package main

import (
	"io"

	"github.com/olddognewflex/ai-map/tools/cli/internal/cli"
	"github.com/spf13/cobra"
)

func newConformanceCmd(stdout, stderr io.Writer) *cobra.Command {
	var updateGolden bool
	cmd := &cobra.Command{
		Use:   "conformance [--update-golden]",
		Short: "Run fixtures and golden tests",
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = updateGolden
			// MVP placeholder: fixtures and golden tests are not present in this repo yet.
			_, _ = stderr.Write([]byte("conformance: not implemented yet (fixtures/golden tests will land later)\n"))
			return cli.ExitError{Code: cli.ExitUsageOrConfig}
		},
	}
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)
	cmd.Flags().BoolVar(&updateGolden, "update-golden", false, "Update golden files (off by default)")
	return cmd
}


