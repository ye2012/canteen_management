package dto

import "github.com/canteen_management/enum"

type Response struct {
	Code    enum.ErrorCode `json:"code"` // status code
	Success bool           `json:"success"`
	Msg     string         `json:"errorMessage"`
	Data    interface{}    `json:"data"`

	NotEscapeHtml bool `json:"-"`
}

func GetInitResponse() *Response {
	return &Response{Code: 0, Success: true, NotEscapeHtml: false}
}
