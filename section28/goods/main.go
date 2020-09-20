package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/google/uuid"
	"github.com/longjoy/micro-go-course/section28/goods/endpoint"
	"github.com/longjoy/micro-go-course/section28/goods/pkg/discovery"
	"github.com/longjoy/micro-go-course/section28/goods/pkg/loadbalancer"
	"github.com/longjoy/micro-go-course/section28/goods/service"
	"github.com/longjoy/micro-go-course/section28/goods/transport"
	"golang.org/x/time/rate"
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
	servicePort := flag.Int("service.port", 12313, "service port")
	serviceName := flag.String("service.name", "good", "service name")
	serviceAddr := flag.String("service.addr", "localhost", "service addr")
	flag.Parse()

	flag.Parse()

	errChan := make(chan error)
	client := discovery.NewDiscoveryClient(*consulAddr, *consulPort)

	instanceId := *serviceName + "-" + uuid.New().String()

	err := client.Register(context.Background(), *serviceName, instanceId, "/health", *serviceAddr, *servicePort, nil, nil)

	loadbalancer := loadbalancer.NewRandomLoadBalancer()

	srv := service.NewGoodsServiceImpl(client, loadbalancer)

	limiter := rate.NewLimiter(20, 20)

	endpoints := endpoint.GoodsEndpoints{
		GoodsDetailEndpoint: endpoint.MakeGoodsDetailEndpoint(srv, limiter),
		HealthCheckEndpoint: endpoint.MakeHealthCheckEndpoint(srv),
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
