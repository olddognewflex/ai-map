package cli

import "strings"

// TrimTrailingNewline removes trailing \r and \n characters for stable, single-line messages.
func TrimTrailingNewline(s string) string {
	return strings.TrimRight(s, "\r\n")
}


