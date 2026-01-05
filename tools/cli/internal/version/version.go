package version

// These are intended to be set via -ldflags at build time.
// Example:
//   go build -ldflags "-X github.com/olddognewflex/ai-map/tools/cli/internal/version.Version=v0.1.0 -X github.com/olddognewflex/ai-map/tools/cli/internal/version.Commit=abc123 -X github.com/olddognewflex/ai-map/tools/cli/internal/version.Date=2026-01-05"
var (
	Version = "dev"
	Commit  = ""
	Date    = ""
)


