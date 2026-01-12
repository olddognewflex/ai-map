package input

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const (
	// MaxYAMLBytes is a safety cap for reading YAML inputs.
	MaxYAMLBytes int64 = 2 << 20 // 2 MiB
)

type Selection struct {
	Dir       string
	Recursive bool
}

func SelectFiles(sel Selection, args []string) ([]string, error) {
	if sel.Dir != "" && len(args) > 0 {
		return nil, fmt.Errorf("provide either --dir or file paths, not both")
	}

	var files []string
	if sel.Dir != "" {
		root, err := filepath.Abs(sel.Dir)
		if err != nil {
			return nil, fmt.Errorf("invalid --dir: %w", err)
		}
		entries, err := scanYAMLFiles(root, sel.Recursive)
		if err != nil {
			return nil, err
		}
		files = append(files, entries...)
	} else {
		for _, p := range args {
			if strings.TrimSpace(p) == "" {
				continue
			}
			abs, err := filepath.Abs(p)
			if err != nil {
				return nil, fmt.Errorf("invalid path %q: %w", p, err)
			}
			files = append(files, abs)
		}
	}

	sort.Strings(files)

	for _, f := range files {
		info, err := os.Stat(f)
		if err != nil {
			return nil, fmt.Errorf("cannot stat %q: %w", f, err)
		}
		if info.IsDir() {
			return nil, fmt.Errorf("%q is a directory; use --dir to scan directories", f)
		}
	}

	return files, nil
}

func EnsureSelected(sel Selection, inputs []string) error {
	if sel.Dir == "" && len(inputs) == 0 {
		return errors.New("no input files provided")
	}
	return nil
}

func ReadFileWithLimit(path string, maxBytes int64) ([]byte, error) {
	if maxBytes <= 0 {
		return nil, errors.New("maxBytes must be > 0")
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if info.Size() > maxBytes {
		return nil, fmt.Errorf("file too large (%d bytes > %d): %s", info.Size(), maxBytes, path)
	}
	return os.ReadFile(path)
}

func scanYAMLFiles(root string, recursive bool) ([]string, error) {
	info, err := os.Stat(root)
	if err != nil {
		return nil, fmt.Errorf("cannot stat --dir %q: %w", root, err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("--dir %q is not a directory", root)
	}

	var out []string
	if !recursive {
		entries, err := os.ReadDir(root)
		if err != nil {
			return nil, fmt.Errorf("cannot read dir %q: %w", root, err)
		}
		for _, e := range entries {
			if e.IsDir() {
				continue
			}
			if isYAMLName(e.Name()) {
				out = append(out, filepath.Join(root, e.Name()))
			}
		}
		return out, nil
	}

	err = filepath.WalkDir(root, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			return nil
		}
		if isYAMLName(d.Name()) {
			out = append(out, path)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("cannot walk dir %q: %w", root, err)
	}
	return out, nil
}

func isYAMLName(name string) bool {
	n := strings.ToLower(name)
	return strings.HasSuffix(n, ".yml") || strings.HasSuffix(n, ".yaml")
}


