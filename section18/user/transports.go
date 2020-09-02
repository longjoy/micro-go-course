package user

import (
	"context"
	"errors"
	"github.com/go-kit/kit/transport/grpc"
	"github.com/longjoy/micro-go-course/section18/pb"
)

var (
	ErrorBadRequest = errors.New("invalid request parameter")
)

type grpcServer struct {
	checkPassword grpc.Handler
}

func (s *grpcServer) CheckPassword(ctx context.Context, r *pb.LoginRequest) (*pb.LoginResponse, error) {
	_, resp, err := s.checkPassword.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.LoginResponse), nil
}

func NewUserServer(ctx context.Context, endpoints Endpoints) pb.UserServiceServer {
	return &grpcServer{
		checkPassword: grpc.NewServer(
			endpoints.UserEndpoint,
			DecodeLoginRequest,
			EncodeLoginResponse,
		),
	}
}

func DecodeLoginRequest(ctx context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.LoginRequest)
	return LoginForm{
		Username: req.Username,
		Password: req.Password,
	}, nil
}

func EncodeLoginResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(LoginResult)
	retStr := "fail"
	if resp.Ret {
		retStr = "success"
	}
	errStr := ""
	if resp.Err != nil {
		errStr = resp.Err.Error()
	}
	return &pb.LoginResponse{
		Ret: retStr,
		Err: errStr,
	}, nil
}
