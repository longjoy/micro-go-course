package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/longjoy/micro-go-course/section24/goods/common"
	"io/ioutil"
	"net/http"
	"net/url"
)

type GoodsDetailVO struct {
	Id       string
	Name     string
	Comments common.CommentListVO
}

type Service interface {
	GetGoodsDetail(ctx context.Context, id string) (GoodsDetailVO, error)
}

func NewGoodsServiceImpl() Service {
	return &GoodsDetailServiceImpl{}
}

type GoodsDetailServiceImpl struct{}

func (service *GoodsDetailServiceImpl) GetGoodsDetail(ctx context.Context, id string) (GoodsDetailVO, error) {
	detail := GoodsDetailVO{Id: id, Name: "Name"}
	var err error
	commonResult, err := GetGoodsComments(id)
	detail.Comments = commonResult.Detail
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
