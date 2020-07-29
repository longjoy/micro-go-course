package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/longjoy/micro-go-course/section10/user/dao"
	"github.com/longjoy/micro-go-course/section10/user/endpoint"
	"github.com/longjoy/micro-go-course/section10/user/redis"
	"github.com/longjoy/micro-go-course/section10/user/service"
	"github.com/longjoy/micro-go-course/section10/user/transport"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main()  {


	var (
		// 服务地址和服务名
		servicePort = flag.Int("service.port", 10086, "service port")

		//waitTime = flag.Int("wait.time", 10, "wait time")



	)

	flag.Parse()

	time.Sleep(10 * time.Second) // 延时启动，等待 MySQL 和 Redis 准备好

	ctx := context.Background()
	errChan := make(chan error)

	err := dao.InitMysql("127.0.0.1", "3306", "root", "123456", "user")
	if err != nil{
		log.Fatal(err)
	}

	err = redis.InitRedis("127.0.0.1","6379", "" )
	if err != nil{
		log.Fatal(err)
	}

	userService := service.MakeUserServiceImpl(&dao.UserDAOImpl{})

	userEndpoints := &endpoint.UserEndpoints{
		endpoint.MakeRegisterEndpoint(userService),
		endpoint.MakeLoginEndpoint(userService),
	}

	r := transport.MakeHttpHandler(ctx, userEndpoints)



	go func() {
		errChan <- http.ListenAndServe(":"  + strconv.Itoa(*servicePort), r)
	}()

	go func() {
		// 监控系统信号，等待 ctrl + c 系统信号通知服务关闭
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	error := <-errChan
	log.Println(error)

}
