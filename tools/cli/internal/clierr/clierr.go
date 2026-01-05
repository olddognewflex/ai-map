package clierr

// ExitError is returned by commands to signal a specific process exit code.
// The caller should not treat this as "unexpected".
type ExitError struct {
	Code int
}

func (e ExitError) Error() string {
	// Keep empty to avoid accidental duplicate printing.
	return ""
}


