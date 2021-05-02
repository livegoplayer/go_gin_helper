package helper

import (
	"net/http"
	"strings"
)

// Resp 返回
type Resp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type EmptyData struct {
}

// ErrorResp 错误返回值
func ErrorResp(code int, msg string, data ...interface{}) {
	//如果isabort 组织输出
	resp(code, msg, data...)
}

// SuccessResp 正确返回值
func SuccessResp(msg string, data ...interface{}) {
	//如果isabort 组织输出
	resp(COMMON_STATUS, msg, data...)
}

// resp 返回
func resp(code int, msg string, data ...interface{}) {

	var res *Error
	if len(data) == 1 {
		res = NewErrorWithData(http.StatusOK, code, data[0], msg)
	}

	if len(data) == 0 {
		res = NewErrorWithData(http.StatusOK, code, EmptyData{}, msg)
	}

	if len(data) > 1 {
		res = NewErrorWithData(http.StatusOK, code, data, msg)
	}

	panic(res)
}

//负责调用panic触发外部的panic处理函数
func CheckError(error error, message ...string) {
	var msg string

	if len(message) == 0 {
		msg = ""
	}

	if error != nil {
		msg = strings.Join(message, " ")
		msg = msg + " error:" + error.Error()
		error = NewError(http.StatusOK, COMMON_ERROR, msg)
		panic(error)
	}
}

func AuthResp(msg string, url string) {

	authData := make(map[string]interface{})
	authData["redirect_url"] = url
	resp(AUTH_ERROR, msg, authData)
}
