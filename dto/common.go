package dto

import (
	"github.com/canteen_management/config"
	"github.com/canteen_management/enum"
	"github.com/canteen_management/model"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

type Response struct {
	Code    enum.ErrorCode `json:"code"` // status code
	Success bool           `json:"success"`
	Msg     string         `json:"message"`
	Data    interface{}    `json:"data"`

	NotEscapeHtml bool `json:"-"`
}

type PaginationReq struct {
	Page     int32 `json:"page"`
	PageSize int32 `json:"page_size"`
}

func (pr *PaginationReq) FixPagination() {
	if pr.Page == 0 {
		pr.Page = 1
	}
	if pr.PageSize < 1 {
		pr.PageSize = 100
	}
}

type PaginationRes struct {
	Page        int32 `json:"page"`
	PageSize    int32 `json:"page_size"`
	TotalNumber int32 `json:"total_number"`
	TotalPage   int32 `json:"total_page"`
}

func (pr *PaginationRes) Format() {
	extraPage := int32(1)
	if pr.TotalNumber%pr.PageSize == 0 {
		extraPage = 0
	}
	pr.TotalPage = pr.TotalNumber/pr.PageSize + extraPage
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

type OrderNode struct {
	ID       string       `json:"id"`
	Price    float64      `json:"price,omitempty"`
	Name     string       `json:"name"`
	DishID   uint32       `json:"dish_id,omitempty"`
	Children []*OrderNode `json:"children,omitempty"`
}

type RequestChecker interface {
	CheckParams() error
}

type CustomContextInfo struct {
	Token    *model.TokenDAO
	ParamMap map[string]string
	Session  *sessions.Session
}

func GetCustomContextInfo(c *gin.Context) *CustomContextInfo {
	customInfo, ok := c.Get(config.CustomKey)
	if !ok {
		return &CustomContextInfo{}
	}
	return customInfo.(*CustomContextInfo)
}
