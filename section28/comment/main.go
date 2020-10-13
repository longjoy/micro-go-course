package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/google/uuid"
	"github.com/longjoy/micro-go-course/section28/comment/endpoint"
	"github.com/longjoy/micro-go-course/section28/comment/pkg/discovery"
	"github.com/longjoy/micro-go-course/section28/comment/service"
	"github.com/longjoy/micro-go-course/section28/comment/transport"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func main() {

	consulAddr := flag.String("consul.addr", "localhost", "consul address")
	consulPort := flag.Int("consul.port", 8500, "consul port")
	servicePort := flag.Int("service.port", 13312, "service port")
	serviceName := flag.String("service.name", "comment", "service name")
	serviceAddr := flag.String("service.addr", "127.0.0.1", "service addr")
	flag.Parse()

	client := discovery.NewDiscoveryClient(*consulAddr, *consulPort)

	instanceId := *serviceName + "-" + uuid.New().String()
	err := client.Register(context.Background(), *serviceName, instanceId, "/health", *serviceAddr, *servicePort, nil, nil)

	errChan := make(chan error)

	srv := service.NewGoodsServiceImpl()

	endpoints := endpoint.CommentsEndpoints{
		CommentsListEndpoint: endpoint.MakeCommentsListEndpoint(srv),
		HealthCheckEndpoint:  endpoint.MakeHealthCheckEndpoint(srv),
	}

	handler := transport.MakeHttpHandler(context.Background(), &endpoints)

	go func() {
		errChan <- http.ListenAndServe(":"+strconv.Itoa(*servicePort), handler)
	}()

	go func() {
		// 监控系统信号，等待 ctrl + c 系统信号通知服务关闭
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	err = <-errChan
	log.Printf("listen err : %s", err)
	client.Deregister(context.Background(), instanceId)
}
