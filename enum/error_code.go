package enum

type ErrorCode = int

const (
	Success ErrorCode = iota
	ParseRequestFailed
	ParamsError
	SqlError
	TokenCheckFailed
	TokenTimeout

	OrderTimeLimit = 100

	SystemError ErrorCode = 999
)

var (
	errorMessageMap = map[ErrorCode]string{
		Success:            "成功",
		ParseRequestFailed: "参数解析失败",
		ParamsError:        "参数错误",
		SqlError:           "数据库错误",
		TokenCheckFailed:   "token检查失败",
		TokenTimeout:       "token已过期",
	}
)

func GetMessage(code ErrorCode) string {
	msg, ok := errorMessageMap[code]
	if ok {
		return msg
	}
	return ""
}
