package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/longjoy/micro-go-course/section25/goods/common"
	"go.etcd.io/etcd/clientv3"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type GoodsDetailVO struct {
	Id       string
	Name     string
	Comments common.CommentListVO
}

type Service interface {
	GetGoodsDetail(ctx context.Context, id string) (GoodsDetailVO, error)
	InitConfig(ctx context.Context)
}

func NewGoodsServiceImpl() Service {
	return &GoodsDetailServiceImpl{}
}

type GoodsDetailServiceImpl struct {
	callCommentService int
}

func (service *GoodsDetailServiceImpl) GetGoodsDetail(ctx context.Context, id string) (GoodsDetailVO, error) {
	detail := GoodsDetailVO{Id: id, Name: "Name"}

	if service.callCommentService != 0 {
		commentResult, _ := GetGoodsComments(id)
		detail.Comments = commentResult.Detail

	}

	var err error
	if err != nil {
		return detail, err
	}
	return detail, nil
}

func GetGoodsComments(id string) (common.CommentResult, error) {
	var result common.CommentResult
	serviceName := "Comments"
	err := hystrix.Do(serviceName, func() error {
		requestUrl := url.URL{
			Scheme:   "http",
			Host:     "127.0.0.1" + ":" + "8081",
			Path:     "/comments/detail",
			RawQuery: "id=" + id,
		}
		resp, err := http.Get(requestUrl.String())
		if err != nil {
			return err
		}
		body, _ := ioutil.ReadAll(resp.Body)
		jsonErr := json.Unmarshal(body, &result)
		if jsonErr != nil {
			return jsonErr
		}
		return nil
	}, func(e error) error {
		// 断路器打开时的处理逻辑，本示例是直接返回错误提示
		return errors.New("Http errors！")
	})
	if err == nil {
		return result, nil
	} else {
		return result, err
	}
}

func (service *GoodsDetailServiceImpl) InitConfig(ctx context.Context) {
	log.Printf("InitConfig")
	cli, _ := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	})
	// get
	resp, _ := cli.Get(ctx, "call_service_d")
	for _, ev := range resp.Kvs {
		fmt.Printf("%s:%s\n", ev.Key, ev.Value)
		if string(ev.Key) == "call_service_d" {
			service.callCommentService, _ = strconv.Atoi(string(ev.Value))
		}
	}

	rch := cli.Watch(context.Background(), "call_service_d") // <-chan WatchResponse
	for wresp := range rch {
		for _, ev := range wresp.Events {
			fmt.Printf("Type: %s Key:%s Value:%s\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			if string(ev.Kv.Key) == "call_service_d" {
				service.callCommentService, _ = strconv.Atoi(string(ev.Kv.Value))
			}
		}
	}
}
