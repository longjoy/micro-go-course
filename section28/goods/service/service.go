package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/longjoy/micro-go-course/section28/goods/common"
	"github.com/longjoy/micro-go-course/section28/goods/pkg/discovery"
	"github.com/longjoy/micro-go-course/section28/goods/pkg/loadbalancer"
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
	HealthCheck() string
}

func NewGoodsServiceImpl(discoveryClient *discovery.DiscoveryClient, loadbalancer loadbalancer.LoadBalancer) Service {
	return &GoodsDetailServiceImpl{
		discoveryClient: discoveryClient,
		loadbalancer:    loadbalancer,
	}
}

type GoodsDetailServiceImpl struct {
	discoveryClient    *discovery.DiscoveryClient
	loadbalancer       loadbalancer.LoadBalancer
	callCommentService int
}

func (service *GoodsDetailServiceImpl) GetGoodsDetail(ctx context.Context, id string) (GoodsDetailVO, error) {
	detail := GoodsDetailVO{Id: id, Name: "Name"}
	commentResult, _ := service.GetGoodsComments(ctx, id)
	detail.Comments = commentResult.Detail
	return detail, nil
}

var ErrNotServiceInstances = errors.New("instances are not existed")
var ErrLoadBalancer = errors.New("loadbalancer select instance error")

func (service *GoodsDetailServiceImpl) DiscoveryService(ctx context.Context, serviceName string) ([]*discovery.InstanceInfo, error) {

	instances, err := service.discoveryClient.DiscoverServices(ctx, serviceName)

	if err != nil {
		log.Printf("get service info err: %s", err)
	}
	if instances == nil || len(instances) == 0 {
		return nil, ErrNotServiceInstances
	}
	return instances, nil
}

func (service *GoodsDetailServiceImpl) GetGoodsComments(ctx context.Context, id string) (common.CommentResult, error) {
	var result common.CommentResult
	serviceName := "comment"

	instances, err := service.discoveryClient.DiscoverServices(ctx, serviceName)

	if err != nil {
		log.Printf("get service info err: %s", err)
	}
	if instances == nil || len(instances) == 0 {
		log.Printf("no instance")
		return result, ErrNotServiceInstances
	}

	selectedInstance, err2 := service.loadbalancer.SelectService(instances)

	log.Print("select instance info :" + selectedInstance.Address + ":" + strconv.Itoa(selectedInstance.Port))

	if err2 != nil {
		log.Printf("loadbalancer get selected instance  err: %s", err2)
		return result, ErrLoadBalancer
	}

	call_err := hystrix.Do(serviceName, func() error {
		requestUrl := url.URL{
			Scheme:   "http",
			Host:     selectedInstance.Address + ":" + strconv.Itoa(selectedInstance.Port),
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

	if call_err == nil {
		return result, nil
	} else {
		return result, call_err
	}
}

func (service *GoodsDetailServiceImpl) InitConfig(ctx context.Context) {
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

func (service *GoodsDetailServiceImpl) HealthCheck() string {
	return "OK"
}
