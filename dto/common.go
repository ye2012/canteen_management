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

type TableColumnInfo struct {
	Name      string             `json:"name"`
	DataIndex string             `json:"data_index"`
	Hide      bool               `json:"hide"`
	MergeRow  bool               `json:"merge_row"`
	Children  []*TableColumnInfo `json:"children,omitempty"`
}

type TableRowColumnInfo struct {
	ID    uint32 `json:"id"`
	Value string `json:"value,omitempty"`
}

type TableRowInfo = map[string]*TableRowColumnInfo

type ModifyTableRowInfo = map[string]interface{}

type OrderNode struct {
	ID       string       `json:"id"`
	Name     string       `json:"name"`
	Children []*OrderNode `json:"children,omitempty"`
	OrderNodeMap
}

type OrderNodeMap = map[string]interface{}
