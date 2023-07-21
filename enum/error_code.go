package enum

type ErrorCode = int

const (
	Success ErrorCode = iota
	ParseRequestFailed
	ParamsError
	SqlError
	TokenCheckFailed
	TokenTimeout

	SystemError ErrorCode = 999
)
