package dto

import (
	"github.com/canteen_management/config"
	"github.com/canteen_management/enum"
	"github.com/canteen_management/model"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

type Request struct {
}

type Response struct {
	Code    enum.ErrorCode `json:"code"` // status code
	Success bool           `json:"success"`
	Msg     string         `json:"message"`
	Data    interface{}    `json:"data"`

	NotEscapeHtml bool `json:"-"`
}

type PaginationQ interface {
	FixPagination()
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

type PaginationS interface {
	Format()
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
	ID             string       `json:"id"`
	Price          float64      `json:"price,omitempty"`
	Name           string       `json:"name"`
	DishID         uint32       `json:"dish_id,omitempty"`
	Picture        string       `json:"picture,omitempty"`
	SelectedNumber int32        `json:"selected_number"`
	Children       []*OrderNode `json:"children,omitempty"`
}

type GoodsNode struct {
	ID             string       `json:"id"`
	Price          float64      `json:"price,omitempty"`
	Name           string       `json:"name"`
	Picture        string       `json:"picture"`
	GoodsID        uint32       `json:"goods_id,omitempty"`
	SelectedNumber int32        `json:"selected_number"`
	Left           float64      `json:"left"`
	BatchSize      float64      `json:"batch_size"`
	BatchUnit      string       `json:"batch_unit"`
	Children       []*GoodsNode `json:"children,omitempty"`
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

type PurchaseGoodsBase struct {
	ID           uint32  `json:"id"`
	GoodsID      uint32  `json:"goods_id"`
	Name         string  `json:"name"`
	Picture      string  `json:"picture"`
	GoodsTypeID  uint32  `json:"goods_type_id"`
	ExpectNumber float64 `json:"expect_number"`
}
