package core

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

const (
	PageNoParam   = "pageNo"
	PageSizeParam = "pageSize"
)

type HttpPageResponse struct {
	Data       any `json:"data"`
	PageNo     int `json:"pageNo"`
	PageSize   int `json:"pageSize"`
	TotalCount int `json:"totalCount"`
}

type HttpPageCondition struct {
	Offset   int
	Limit    int
	PageNo   int
	PageSize int
}

// BuildOffsetLimitWithGinContext 返回分页查询的 offset 和 limit 值
func BuildOffsetLimitWithGinContext(context *gin.Context) (*HttpPageCondition, error) {
	return BuildOffsetLimit(context.Param(PageNoParam), context.Param(PageSizeParam))
}

func BuildOffsetLimit(pageNo string, pageSize string) (*HttpPageCondition, error) {
	pn, err := strconv.Atoi(pageNo)
	if err != nil {
		return nil, err
	}
	ps, err := strconv.Atoi(pageSize)
	if err != nil {
		return nil, err
	}
	return &HttpPageCondition{
		Offset:   (pn - 1) * ps,
		Limit:    ps,
		PageNo:   pn,
		PageSize: ps,
	}, nil
}
