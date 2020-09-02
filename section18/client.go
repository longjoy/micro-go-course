package main

import (
	"context"
	"fmt"
	"github.com/longjoy/micro-go-course/section18/pb"
	"google.golang.org/grpc"
)

func main() {
	serviceAddress := "127.0.0.1:8080"
	conn, err := grpc.Dial(serviceAddress, grpc.WithInsecure())
	if err != nil {
		panic("connect error")
	}
	defer conn.Close()
	userClient := pb.NewUserServiceClient(conn)
	stringReq := &pb.LoginRequest{Username: "admin", Password: "admin"}
	reply, _ := userClient.CheckPassword(context.Background(), stringReq)
	fmt.Printf("CheckPassword ret is %s\n", reply.Ret)
}
