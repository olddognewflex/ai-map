package main

import (
	"os"

	"github.com/olddognewflex/ai-map/tools/cli/internal/cli"
)

func main() {
	os.Exit(cli.Run(os.Args))
}


