package helper

import (
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
	resp(0, msg, data...)
}

// resp 返回
func resp(code int, msg string, data ...interface{}) {
	resp := Resp{
		Code: code,
		Msg:  msg,
		Data: data,
	}

	if len(data) == 1 {
		resp.Data = data[0]
	}

	if len(data) == 0 {
		resp.Data = &EmptyData{}
	}
	// 设置返回格式是json
	res := NewErrorWithData(200, code, data, msg)

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
		error = NewError(200, 1, msg)
		panic(error)
	}
}
