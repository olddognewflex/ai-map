package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/olddognewflex/ai-map/tools/cli/internal/cli"
	"github.com/olddognewflex/ai-map/tools/cli/internal/input"
	"github.com/olddognewflex/ai-map/tools/cli/internal/validate"
	"github.com/spf13/cobra"
)

func newConformanceCmd(stdout, stderr io.Writer) *cobra.Command {
	var updateGolden bool
	var repoRoot string
	var schemaPath string
	cmd := &cobra.Command{
		Use:   "conformance [--repo-root DIR] [--schema FILE] [--update-golden]",
		Short: "Run fixtures and golden tests",
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = updateGolden // golden support will land later; flag is reserved and safe.

			root := strings.TrimSpace(repoRoot)
			if root == "" {
				root = "."
			}
			absRoot, err := filepath.Abs(root)
			if err != nil {
				return cli.ExitError{Code: cli.ExitUsageOrConfig, Msg: "error: invalid --repo-root: " + err.Error()}
			}
			if st, err := os.Stat(absRoot); err != nil || !st.IsDir() {
				return cli.ExitError{Code: cli.ExitUsageOrConfig, Msg: "error: --repo-root is not a directory: " + absRoot}
			}

			// Locate fixtures under repoRoot/spec/examples/{valid,invalid}
			validDir := filepath.Join(absRoot, "spec", "examples", "valid")
			invalidDir := filepath.Join(absRoot, "spec", "examples", "invalid")

			validFiles, _ := listYAMLFilesIfDir(validDir)
			invalidFiles, _ := listYAMLFilesIfDir(invalidDir)

			if len(validFiles) == 0 && len(invalidFiles) == 0 {
				fmt.Fprintf(stderr, "conformance: no fixtures found under %s; skipping\n", filepath.Join(absRoot, "spec", "examples"))
				return nil
			}

			sp := strings.TrimSpace(schemaPath)
			if sp == "" {
				sp = filepath.Join(absRoot, "spec", "ai-map.schema.json")
			}

			v, err := validate.New(validate.Options{
				MaxBytes:   input.MaxYAMLBytes,
				SchemaPath: sp,
			})
			if err != nil {
				return cli.ExitError{Code: cli.ExitUsageOrConfig, Msg: "error: cannot load schema: " + err.Error()}
			}

			var failed bool
			// Deterministic ordering.
			sort.Strings(validFiles)
			sort.Strings(invalidFiles)

			for _, p := range validFiles {
				res, err := v.ValidateFile(p)
				if err != nil {
					fmt.Fprintf(stderr, "%s: error: %s\n", p, err)
					return cli.ExitError{Code: cli.ExitInternalError}
				}
				if !res.OK {
					failed = true
					fmt.Fprintf(stderr, "%s: expected valid but was invalid\n", p)
					for _, e := range res.Errors {
						fmt.Fprintf(stderr, "  - %s\n", cli.TrimTrailingNewline(e))
					}
				}
			}
			for _, p := range invalidFiles {
				res, err := v.ValidateFile(p)
				if err != nil {
					fmt.Fprintf(stderr, "%s: error: %s\n", p, err)
					return cli.ExitError{Code: cli.ExitInternalError}
				}
				if res.OK {
					failed = true
					fmt.Fprintf(stderr, "%s: expected invalid but was valid\n", p)
				}
			}
			if failed {
				return cli.ExitError{Code: cli.ExitCheckFailed}
			}
			return nil
		},
	}
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)
	cmd.Flags().StringVar(&repoRoot, "repo-root", ".", "Repository root (used to locate spec/examples)")
	cmd.Flags().StringVar(&schemaPath, "schema", "", "Path to JSON Schema (defaults to <repo-root>/spec/ai-map.schema.json)")
	cmd.Flags().BoolVar(&updateGolden, "update-golden", false, "Update golden files (off by default)")
	return cmd
}

func listYAMLFilesIfDir(dir string) ([]string, error) {
	st, err := os.Stat(dir)
	if err != nil || !st.IsDir() {
		return nil, err
	}
	ents, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var out []string
	for _, e := range ents {
		if e.IsDir() {
			continue
		}
		n := strings.ToLower(e.Name())
		if strings.HasSuffix(n, ".yml") || strings.HasSuffix(n, ".yaml") {
			out = append(out, filepath.Join(dir, e.Name()))
		}
	}
	return out, nil
}
