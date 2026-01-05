package main

import (
	"io"

	"github.com/olddognewflex/ai-map/tools/cli/internal/cli"
	"github.com/spf13/cobra"
)

func newScaffoldCmd(stdout, stderr io.Writer) *cobra.Command {
	var outDir string
	var name string
	cmd := &cobra.Command{
		Use:   "scaffold --out DIR [--name NAME]",
		Short: "Create a new agent map folder skeleton",
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = outDir
			_ = name
			_, _ = stderr.Write([]byte("scaffold: not implemented yet (safe skeleton writer will land later)\n"))
			return cli.ExitError{Code: cli.ExitUsageOrConfig}
		},
	}
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)
	cmd.Flags().StringVar(&outDir, "out", "", "Output directory to create (required)")
	cmd.Flags().StringVar(&name, "name", "", "System name (optional; used in template)")
	_ = cmd.MarkFlagRequired("out")
	return cmd
}


