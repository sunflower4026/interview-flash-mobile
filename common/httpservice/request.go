package httpservice

import (
	"github.com/gin-gonic/gin"
)

type PaginationRequest struct {
	Page  int `json:"current_page,omitempty" form:"page"`
	Limit int `json:"limit,omitempty" form:"limit"`
}

func (p *PaginationRequest) Validate() PaginationRequest {
	// if no limit is provided, set default limit to 10 and page to 1
	if p.Limit == 0 {
		p.Limit = 10
	}

	if p.Page == 0 {
		p.Page = 1
	}

	return *p
}

func (p *PaginationRequest) GetOffset() int {
	return (p.Page - 1) * p.Limit
}

type OrderRequest struct {
	OrderBy string `json:"order_by,omitempty" form:"order_by"`
	Sort    string `json:"sort,omitempty" form:"sort"`
}

func (o *OrderRequest) Validate() OrderRequest {
	if o.OrderBy == "" {
		o.OrderBy = "id"
	}

	if o.Sort == "" {
		o.Sort = "asc"
	}

	return *o
}

func GetPaginationRequest(ctx *gin.Context) (PaginationRequest, error) {
	var request PaginationRequest
	err := ctx.ShouldBindQuery(&request)
	return request.Validate(), err
}

func GetOrderRequest(ctx *gin.Context) (OrderRequest, error) {
	var request OrderRequest
	err := ctx.ShouldBindQuery(&request)
	return request.Validate(), err
}

func GetFilterRequest(ctx *gin.Context) (map[string]interface{}, error) {
	filter := make(map[string]interface{})
	for key, values := range ctx.Request.URL.Query() {
		if len(values) > 0 && values[0] != "" && key != "page" && key != "limit" && key != "order_by" && key != "sort" {
			if len(values) > 1 {
				filter[key] = values
			} else {
				filter[key] = values[0]
			}
		}
	}
	return filter, nil
}
