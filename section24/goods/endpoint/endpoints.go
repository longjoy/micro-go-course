package endpoint

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/longjoy/micro-go-course/section24/goods/service"
)

type GoodsEndpoints struct {
	GoodsDetailEndpoint endpoint.Endpoint
}

// 服务发现请求结构体
type GoodsDetailRequest struct {
	Id string
}

// 服务发现响应结构体
type GoodsDetailResponse struct {
	Detail service.GoodsDetailVO `json:"detail"`
	Error  string                `json:"error"`
}

// 创建服务发现的 Endpoint
func MakeGoodsDetailEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {

		req := request.(GoodsDetailRequest)
		detail, err := svc.GetGoodsDetail(ctx, req.Id)
		var errString = ""

		if err != nil {
			errString = err.Error()
			return &GoodsDetailResponse{
				Detail: detail,
				Error:  errString,
			}, nil
		}
		return &GoodsDetailResponse{
			Detail: detail,
			Error:  errString,
		}, nil
	}
}
