package exitcode

// Exit code contract:
//   0: success
//   1: validation/lint/test failures
//   2: usage/config errors
//  >=3: unexpected/internal errors
const (
	OK            = 0
	CheckFailed   = 1
	UsageOrConfig = 2
	InternalError = 3
)


