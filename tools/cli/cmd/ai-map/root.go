package main

import (
	"errors"
	"fmt"
	"io"

	"github.com/olddognewflex/ai-map/tools/cli/internal/cli"
	"github.com/spf13/cobra"
)

func newRootCmd(stdout, stderr io.Writer) *cobra.Command {
	root := &cobra.Command{
		Use:   "ai-map",
		Short: "Tooling for the AI-Map spec",
		Long:  "ai-map: tooling for the AI-Map spec",
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = cmd.Help()
			return cli.ExitError{Code: cli.ExitUsageOrConfig}
		},
		SilenceErrors: true,
		SilenceUsage:  true,
	}
	// Keep the command surface area minimal and stable.
	root.CompletionOptions.DisableDefaultCmd = true
	root.SetOut(stdout)
	root.SetErr(stderr)

	// Make flag errors follow our exit-code contract.
	root.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		// Keep message first, then help (deterministic, user-friendly).
		if err != nil {
			fmt.Fprintf(stderr, "error: %s\n\n", err)
		}
		_ = cmd.Help()
		return cli.ExitError{Code: cli.ExitUsageOrConfig}
	})

	root.AddCommand(newValidateCmd(stdout, stderr))
	root.AddCommand(newLintCmd(stdout, stderr))
	root.AddCommand(newRenderCmd(stdout, stderr))
	root.AddCommand(newTypesCmd(stdout, stderr))
	root.AddCommand(newConformanceCmd(stdout, stderr))
	root.AddCommand(newScaffoldCmd(stdout, stderr))
	root.AddCommand(newVersionCmd(stdout, stderr))

	return root
}

func exitCodeFromError(err error) (code int, msg string, ok bool) {
	if err == nil {
		return 0, "", true
	}
	var ee cli.ExitError
	if errors.As(err, &ee) {
		return ee.Code, ee.Msg, true
	}
	return 0, "", false
}


