package controller

type Result struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func ResultOk() Result {
	return Result{
		Code:    200,
		Message: "success",
	}
}

func ResultMsg(message string) Result {
	return Result{
		Code:    200,
		Message: message,
	}
}

func ResultData(data any) Result {
	return Result{
		Code:    200,
		Message: "success",
		Data:    data,
	}
}

func ResultMsgData(message string, data any) Result {
	return Result{
		Code:    200,
		Message: message,
		Data:    data,
	}
}

func ResultErr(message string) Result {
	return Result{
		Code:    500,
		Message: message,
	}
}

func ResultErrCode(code int, message string) Result {
	return Result{
		Code:    code,
		Message: message,
	}
}
