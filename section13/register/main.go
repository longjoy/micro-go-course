package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/google/uuid"
	"github.com/longjoy/micro-go-course/section13/register/discovery"
	"github.com/longjoy/micro-go-course/section13/register/endpoint"
	service "github.com/longjoy/micro-go-course/section13/register/service"
	"github.com/longjoy/micro-go-course/section13/register/transport"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

/**
 * @Author : dixuanhuang
 * @File : main.go
 * @Date : 2020/8/10 9:19 下午
 * @Description:
**/



func main()  {



	consulAddr := flag.String("consul.addr", "localhost", "consul address")
	consulPort := flag.Int("consul.port", 8500, "consul port")
	serviceName := flag.String("service.name", "register", "service name")
	serviceAddr := flag.String("service.addr", "localhost", "service addr")
	servicePort := flag.Int("service.port", 12312, "service port")

	flag.Parse()

	client := discovery.NewDiscoveryClient(*consulAddr, *consulPort)

	errChan := make(chan error)


	srv := service.NewRegisterServiceImpl(client)


	endpoints := endpoint.RegisterEndpoints{
		DiscoveryEndpoint: endpoint.MakeDiscoveryEndpoint(srv),
		HealthCheckEndpoint: endpoint.MakeHealthCheckEndpoint(srv),
	}

	handler := transport.MakeHttpHandler(context.Background(), &endpoints)

	go func() {
		errChan <- http.ListenAndServe(":" + strconv.Itoa(*servicePort), handler)
	}()

	go func() {
		// 监控系统信号，等待 ctrl + c 系统信号通知服务关闭
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	instanceId := *serviceName + "-" + uuid.New().String()
	err := client.Register(context.Background(), *serviceName, instanceId, "/health", *serviceAddr, *servicePort, nil, nil)


	if err != nil{
		log.Printf("register service err : %s", err)
		os.Exit(-1)
	}

	err = <-errChan
	log.Printf("listen err : %s", err)
	client.Deregister(context.Background(), instanceId)

}

func init(){
	file := "./" +"register.log"
	logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		panic(err)
	}
	log.SetOutput(logFile) // 将文件设置为log输出的文件
	log.SetFlags(log.Ldate|log.Lshortfile)
}