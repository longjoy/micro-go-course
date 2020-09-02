package user

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

type LoginForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResult struct {
	Ret bool  `json:"ret"`
	Err error `json:"err"`
}

type Endpoints struct {
	UserEndpoint endpoint.Endpoint
}

func MakeUserEndpoint(svc UserService) endpoint.Endpoint {

	return func(ctx context.Context, form interface{}) (result interface{}, err error) {
		req := form.(LoginForm)
		ret, err := svc.CheckPassword(ctx, req.Username, req.Password)
		return LoginResult{Ret: ret, Err: err}, nil
	}
}
