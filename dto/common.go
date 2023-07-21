package dto

import "github.com/canteen_management/enum"

type Response struct {
	Code enum.ErrorCode `json:"code"` // status code
	Msg  string         `json:"msg"`
	Data interface{}    `json:"data"`

	NotEscapeHtml bool `json:"-"`
}

func GetInitResponse() *Response {
	return &Response{Code: 0, NotEscapeHtml: false}
}
