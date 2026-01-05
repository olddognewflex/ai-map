package main

import (
	"fmt"
	"os"

	"github.com/olddognewflex/ai-map/tools/cli/internal/cli"
)

func main() {
	cmd := newRootCmd(os.Stdout, os.Stderr)
	cmd.SetArgs(os.Args[1:])
	if err := cmd.Execute(); err != nil {
		if code, msg, ok := exitCodeFromError(err); ok {
			if msg != "" {
				fmt.Fprintln(os.Stderr, msg)
			}
			os.Exit(code)
		}
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(cli.ExitInternalError)
	}
	os.Exit(cli.ExitOK)
}
