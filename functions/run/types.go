package run

// LocalProc represents the result of executing a local process
// from a Starlark script.
type LocalProc struct {
	Pid      int64
	Error    string
	Result   string
	ExitCode int64
}
