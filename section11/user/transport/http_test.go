package transport

import (
	"encoding/json"
	"flag"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"
)

/**
 * @Author : dixuanhuang
 * @File : http_test.go
 * @Date : 2020/8/4 11:26 上午
 * @Description:
**/

func TestRegister(t *testing.T)  {

	if !flag.Parsed(){
		flag.Parse()
	}

	args := flag.Args()
	postUrl := "http://127.0.0.1:30036/register"
	if len(args) > 0{
		postUrl = args[0]
	}

	body := map[string]string{
		"email":"aoho1@mail.com",
		"password":"aoho",
		"username": "aoho",
	}
	result, err := httpPost(postUrl, body)

	if err != nil{
		t.Errorf("http post err %s", err)
		t.FailNow()
	}

	t.Logf("result is %v", result)

}

func TestLogin(t *testing.T)  {
	if !flag.Parsed(){
		flag.Parse()
	}

	args := flag.Args()
	postUrl := "http://127.0.0.1:30036/login"
	if len(args) > 0{
		postUrl = args[0]
	}
	body := map[string]string{
		"email":"aoho@mail.com",
		"password":"aoho1",

	}
	result, err := httpPost(postUrl, body)

	if err != nil{
		t.Errorf("http post err %s", err)
		t.FailNow()
	}
	t.Logf("result is %v", result)
}

func httpPost(postUrl string, body map[string]string) (interface{}, error) {

	// 超时时间：5秒
	client := &http.Client{Timeout: 5 * time.Second}

	dataUrlVal := url.Values{}

	for k, v := range body {
		dataUrlVal.Add(k, v)
	}

	req,err := http.NewRequest("POST", postUrl, strings.NewReader(dataUrlVal.Encode()))
	if err != nil{
		return nil, err
	}

	req.Header.Add("Content-Type","application/x-www-form-urlencoded")

	//提交请求
	resp, err := client.Do(req)
	if err != nil{
		return nil, err
	}
	defer resp.Body.Close()
	//读取返回值
	decode := json.NewDecoder(resp.Body)
	var result interface{}

	err = decode.Decode(&result)
	if err != nil{
		return nil, err
	}
	return result, nil

}