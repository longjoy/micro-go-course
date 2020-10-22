package endpoint

import (
	"context"
	"errors"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	. "github.com/longjoy/micro-go-course/section32/model"
	"github.com/longjoy/micro-go-course/section32/service"
)

// CalculateEndpoint define endpoint
type OAuth2Endpoints struct {
	IndexEndpoint 		endpoint.Endpoint
	SampleEndpoint 		endpoint.Endpoint
	AdminEndpoint		endpoint.Endpoint
	HealthCheckEndpoint 		endpoint.Endpoint
}



func MakeOAuth2AuthorizationMiddleware(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {

		return func(ctx context.Context, request interface{}) (response interface{}, err error) {

			if err, ok := ctx.Value(OAuth2ErrorKey).(error); ok{
				return nil, err
			}
			if _, ok := ctx.Value(OAuth2DetailsKey).(*OAuth2Details); !ok{
				return  nil, ErrInvalidUserRequest
			}
			return next(ctx, request)
		}
	}
}
func MakeAuthorityAuthorizationMiddleware(authority string, logger log.Logger) endpoint.Middleware  {
	return func(next endpoint.Endpoint) endpoint.Endpoint {

		return func(ctx context.Context, request interface{}) (response interface{}, err error) {

			if err, ok := ctx.Value(OAuth2ErrorKey).(error); ok{
				return nil, err
			}
			if details, ok := ctx.Value(OAuth2DetailsKey).(*OAuth2Details); !ok{
				return  nil, ErrInvalidClientRequest
			}else {
				for _, value := range details.User.Authorities{
					if value == authority{
						return next(ctx, request)
					}
				}
				return nil, ErrNotPermit
			}
		}
	}
}

const (

	OAuth2DetailsKey       = "OAuth2Details"
	OAuth2ErrorKey         = "OAuth2Error"

)


var (
	ErrInvalidClientRequest = errors.New("invalid client message")
	ErrInvalidUserRequest = errors.New("invalid user message")
	ErrNotPermit = errors.New("not permit")
)

type IndexRequest struct {
}

type IndexResponse struct {
	Result string `json:"result"`
	Error string `json:"error"`
}

type SampleRequest struct {
}

type SampleResponse struct {
	Result string `json:"result"`
	Error string `json:"error"`
}

type AdminRequest struct {
}

type AdminResponse struct {
	Result string `json:"result"`
	Error string `json:"error"`
}


func MakeIndexEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		result := svc.Index()
		return &SampleResponse{
			Result:result,
		}, nil
	}

}

func MakeSampleEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		result := svc.Sample(ctx.Value(OAuth2DetailsKey).(*OAuth2Details).User.Username)
		return &SampleResponse{
			Result:result,
		}, nil
	}

}

func MakeAdminEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		result := svc.Admin(ctx.Value(OAuth2DetailsKey).(*OAuth2Details).User.Username)
		return &AdminResponse{
			Result:result,
		}, nil
	}
}




// HealthRequest 健康检查请求结构
type HealthRequest struct{}

// HealthResponse 健康检查响应结构
type HealthResponse struct {
	Status bool `json:"status"`
}

// MakeHealthCheckEndpoint 创建健康检查Endpoint
func MakeHealthCheckEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		status := svc.HealthCheck()
		return HealthResponse{
			Status:status,
		}, nil
	}
}
