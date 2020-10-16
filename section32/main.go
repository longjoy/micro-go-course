package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/longjoy/micro-go-course/section32/config"
	"github.com/longjoy/micro-go-course/section32/endpoint"
	"github.com/longjoy/micro-go-course/section32/service"
	"github.com/longjoy/micro-go-course/section32/transport"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func main() {

	var (
		servicePort = flag.Int("service.port", 10099, "service port")
	)

	flag.Parse()

	ctx := context.Background()
	errChan := make(chan error)


	var tokenService service.ResourceServerTokenService
	var tokenEnhancer service.TokenEnhancer
	var tokenStore service.TokenStore
	var srv service.Service


	tokenEnhancer = service.NewJwtTokenEnhancer("secret")
	tokenStore = service.NewJwtTokenStore(tokenEnhancer.(*service.JwtTokenEnhancer))
	tokenService = service.NewTokenService(tokenStore, tokenEnhancer)




	srv = service.NewCommonService()


	indexEndpoint := endpoint.MakeIndexEndpoint(srv)
	sampleEndpoint := endpoint.MakeSampleEndpoint(srv)
	sampleEndpoint = endpoint.MakeOAuth2AuthorizationMiddleware(config.KitLogger)(sampleEndpoint)
	adminEndpoint := endpoint.MakeAdminEndpoint(srv)
	adminEndpoint = endpoint.MakeOAuth2AuthorizationMiddleware(config.KitLogger)(adminEndpoint)
	adminEndpoint = endpoint.MakeAuthorityAuthorizationMiddleware("Admin", config.KitLogger)(adminEndpoint)

	//创建健康检查的Endpoint
	healthEndpoint := endpoint.MakeHealthCheckEndpoint(srv)

	endpts := endpoint.OAuth2Endpoints{
		HealthCheckEndpoint: healthEndpoint,
		IndexEndpoint: indexEndpoint,
		SampleEndpoint:sampleEndpoint,
		AdminEndpoint:adminEndpoint,
	}

	//创建http.Handler
	r := transport.MakeHttpHandler(ctx, endpts, tokenService, config.KitLogger)

	go func() {
		config.Logger.Println("Http Server start at port:" + strconv.Itoa(*servicePort))
		handler := r
		errChan <- http.ListenAndServe(":"  + strconv.Itoa(*servicePort), handler)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	error := <-errChan
	config.Logger.Println(error)
}
