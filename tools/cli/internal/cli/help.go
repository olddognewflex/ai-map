package cli

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

func printRootHelp(w io.Writer) {
	fmt.Fprint(w, rootHelpText())
}

func rootHelpText() string {
	// Keep deterministic order.
	cmds := []string{"validate", "lint", "render", "types", "conformance", "scaffold", "version"}
	sort.Strings(cmds)

	var b strings.Builder
	b.WriteString("ai-map: tooling for the AI-Map spec\n\n")
	b.WriteString("Usage:\n")
	b.WriteString("  ai-map <command> [flags] [args]\n\n")
	b.WriteString("Commands:\n")
	for _, c := range cmds {
		b.WriteString(fmt.Sprintf("  %-12s %s\n", c, shortHelp(c)))
	}
	b.WriteString("\n")
	b.WriteString("Global flags:\n")
	b.WriteString("  -h, --help   Show help\n\n")
	b.WriteString("Run 'ai-map <command> --help' for command-specific help.\n")
	return b.String()
}

func shortHelp(cmd string) string {
	switch cmd {
	case "validate":
		return "Validate YAML files against the JSON Schema"
	case "lint":
		return "Run opinionated checks"
	case "render":
		return "Render AI-Map docs (Markdown)"
	case "types":
		return "Generate types (MVP)"
	case "conformance":
		return "Run fixtures and golden tests"
	case "scaffold":
		return "Create a new agent map folder skeleton"
	case "version":
		return "Print version information"
	default:
		return ""
	}
}


