package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/longjoy/micro-go-course/section24/goods/endpoint"
	"github.com/longjoy/micro-go-course/section24/goods/service"
	"github.com/longjoy/micro-go-course/section24/goods/transport"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func main() {

	servicePort := flag.Int("service.port", 12312, "service port")

	flag.Parse()

	errChan := make(chan error)

	srv := service.NewGoodsServiceImpl()

	endpoints := endpoint.GoodsEndpoints{
		GoodsDetailEndpoint: endpoint.MakeGoodsDetailEndpoint(srv),
	}
	handler := transport.MakeHttpHandler(context.Background(), &endpoints)

	hystrix.ConfigureCommand("Comments", hystrix.CommandConfig{
		RequestVolumeThreshold: 4,
	})

	go func() {
		errChan <- http.ListenAndServe(":"+strconv.Itoa(*servicePort), handler)
	}()

	go func() {
		// 监控系统信号，等待 ctrl + c 系统信号通知服务关闭
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	err := <-errChan
	log.Printf("listen err : %s", err)
}
