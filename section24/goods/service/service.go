package service

import (
	"context"
	"encoding/json"
	"github.com/longjoy/micro-go-course/section24/goods/common"
	"io/ioutil"
	"net/http"
	"net/url"
)


type GoodsDetailVO struct {
	Id string
	Name string
	Comments common.CommentListVO
}


type Service interface {

	GetGoodsDetail(ctx context.Context, id string) (GoodsDetailVO, error)

}


func NewGoodsServiceImpl() Service  {
	return &GoodsDetailServiceImpl{}
}

type GoodsDetailServiceImpl struct {}


func (service *GoodsDetailServiceImpl) GetGoodsDetail(ctx context.Context, id string) (GoodsDetailVO, error)  {
	println("GetGoodsDetail")
	detail := GoodsDetailVO{Id: id, Name:"Name"}
	var err error
	detail.Comments, err = GetGoodsComments(id)
	if err != nil {
		return detail, err
	}
	return detail,nil
}

func GetGoodsComments(id string) (common.CommentListVO, error)  {
	var result common.CommentListVO
	requestUrl := url.URL{
		Scheme:"http",
		Host: "127.0.0.1" + ":" + "8081",
		Path:"/comments/detail",
		RawQuery:"id=" + id,
	}
	resp, err := http.Get(requestUrl.String())
	if err != nil{
		return result, err
	}
	body, _ := ioutil.ReadAll(resp.Body)
	​jsonErr := json.Unmarshal(body, result)
	if jsonErr != nil {
		return result, jsonErr
	}
	return result, nil
}
