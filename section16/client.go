package main

import (
	"fmt"
	service "github.com/longjoy/micro-go-course/section16/service"
	"log"
	"net/rpc"
)

func main() {

	client, err := rpc.DialHTTP("tcp", "127.0.0.1:1234")

	if err != nil {
		log.Fatal("dialing:", err)
	}

	stringReq := &service.StringRequest{"A", "B"}
	var reply string
	err = client.Call("StringService.Concat", stringReq, &reply)
	fmt.Printf("StringService Concat : %s concat %s = %s\n", stringReq.A, stringReq.B, reply)
	if err != nil {
		log.Fatal("Concat error:", err)
	}
	// 异步的调用方式
	call := client.Go("StringService.Concat", stringReq, &reply, nil)
	_ = <-call.Done
	fmt.Printf("StringService Concat : %s concat %s = %s\n", stringReq.A, stringReq.B, reply)
}
