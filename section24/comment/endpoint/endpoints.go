package endpoint

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/longjoy/micro-go-course/section24/comment/service"
)

type CommentsEndpoints struct {
	CommentsListEndpoint endpoint.Endpoint
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

// 创建服务发现的 Endpoint
func MakeCommentsListEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		println("MakeCommentsListEndpoint")
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
