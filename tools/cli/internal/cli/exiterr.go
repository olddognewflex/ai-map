package cli

// ExitError is returned by commands to signal a specific process exit code.
// Keep Msg empty if the command already wrote a user-facing message to stderr.
type ExitError struct {
	Code int
	Msg  string
}

func (e ExitError) Error() string {
	if e.Msg == "" {
		return "exit"
	}
	return e.Msg
}


