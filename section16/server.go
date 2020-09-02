package main

import (
	service "github.com/longjoy/micro-go-course/section16/service"
	"log"
	"net"
	"net/http"
	"net/rpc"
)

func main() {
	stringService := new(service.StringService)
	rpc.Register(stringService)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", "127.0.0.1:1234")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	http.Serve(l, nil)
}
