package endpoint

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/longjoy/micro-go-course/section28/comment/service"
	"log"
)

type CommentsEndpoints struct {
	CommentsListEndpoint endpoint.Endpoint
	HealthCheckEndpoint  endpoint.Endpoint
}

// 服务发现请求结构体
type CommentsListRequest struct {
	Id string
}

// 服务发现响应结构体
type CommentsListResponse struct {
	Detail service.CommentListVO `json:"detail"`
	Error  string                `json:"error"`
}

func MakeCommentsListEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(CommentsListRequest)
		detail, err := svc.GetCommentsList(ctx, req.Id)
		var errString = ""
		if err != nil {
			errString = err.Error()
		}
		return &CommentsListResponse{
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
