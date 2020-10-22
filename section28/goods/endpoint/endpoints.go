package endpoint

import (
	"context"
	"errors"
	"github.com/go-kit/kit/endpoint"
	"github.com/longjoy/micro-go-course/section28/goods/service"
	"golang.org/x/time/rate"
	"log"
)

type GoodsEndpoints struct {
	GoodsDetailEndpoint endpoint.Endpoint
	HealthCheckEndpoint endpoint.Endpoint
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
func MakeGoodsDetailEndpoint(svc service.Service, limiter *rate.Limiter) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {

		if !limiter.Allow() {
			// Allow返回false，表示桶内不足一个令牌，应该被限流，默认返回 ErrLimiExceed 异常
			return nil, errors.New("ErrLimitExceed")
		}
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

// HealthRequest 健康检查请求结构
type HealthRequest struct{}

// HealthResponse 健康检查响应结构
type HealthResponse struct {
	Status string `json:"status"`
}

// MakeHealthCheckEndpoint 创建健康检查Endpoint
func MakeHealthCheckEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		log.Printf("helthcheck")
		status := svc.HealthCheck()
		return HealthResponse{
			Status: status,
		}, nil
	}
}
