package main

import (
	"context"
	"flag"
	"github.com/go-kit/kit/log"
	"github.com/longjoy/micro-go-course/section18/pb"
	"github.com/longjoy/micro-go-course/section18/user"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"net"
	"os"
	"time"
)

func main() {
	flag.Parse()

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	ctx := context.Background()
	// 建立 service
	var svc user.UserService
	svc = user.UserServiceImpl{}

	// 建立 endpoint
	endpoint := user.MakeUserEndpoint(svc)
	// 构造限流中间件
	ratebucket := rate.NewLimiter(rate.Every(time.Second*1), 100)
	endpoint = user.NewTokenBucketLimitterWithBuildIn(ratebucket)(endpoint)

	endpts := user.Endpoints{
		UserEndpoint: endpoint,
	}
	// 使用 transport 构造 UserServiceServer
	handler := user.NewUserServer(ctx, endpts)
	// 监听端口，建立 gRPC 网络服务器，注册 RPC 服务
	ls, _ := net.Listen("tcp", "127.0.0.1:8080")
	gRPCServer := grpc.NewServer()
	pb.RegisterUserServiceServer(gRPCServer, handler)
	gRPCServer.Serve(ls)
}
