package main

import (
	"fmt"
	"io"

	"github.com/olddognewflex/ai-map/tools/cli/internal/version"
	"github.com/spf13/cobra"
)

func newVersionCmd(stdout, stderr io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintf(stdout, "ai-map %s\n", version.Version)
			if version.Commit != "" {
				fmt.Fprintf(stdout, "commit %s\n", version.Commit)
			}
			if version.Date != "" {
				fmt.Fprintf(stdout, "date   %s\n", version.Date)
			}
			return nil
		},
	}
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)
	return cmd
}


