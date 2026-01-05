package cli

// Exit code contract:
//   0: success
//   1: validation/lint/test failures
//   2: usage/config errors
//  >=3: unexpected/internal errors
const (
	ExitOK             = 0
	ExitCheckFailed    = 1
	ExitUsageOrConfig  = 2
	ExitInternalError  = 3
)


