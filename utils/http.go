package utils

import (
	"net/http"
	"encoding/json"
	"bytes"
	"io/ioutil"
	"net/url"
)

type HD map[string]interface{}

func (h *HD) ToJson() string {
	data ,err := json.Marshal(h)
	if err != nil {
		return ""
	}
	return string(data)
}

func HttpPost(url_str string, post_data *url.Values) (M,error){
	req,err := http.PostForm(url_str,*post_data)
	if err != nil {
		return nil,err
	}

	return getRequestData(req)
}

func HttpPostJson(url_str string,post_data M) (M,error){
	b,err := json.Marshal(post_data)
	if err != nil {
		return nil,err
	}
	body := bytes.NewBuffer(b)
	req,err := http.Post(url_str,"application/json;charset=utf-8",body)
	if err != nil {
		return nil,err
	}
	return getRequestData(req)
}

func HttpPostJsonString(url_str string,post_data M) (string,error){
	b,err := json.Marshal(post_data)
	if err != nil {
		return "",err
	}
	body := bytes.NewBuffer(b)
	req,err := http.Post(url_str,"application/json;charset=utf-8",body)
	if err != nil {
		return "",err
	}
	return getRequestString(req)
}

func HttpGet(url_str string) (string,error){
	req,err := http.Get(url_str)
	if err != nil {
		return "",err
	}

	return getRequestString(req)
}

func getRequestData(req *http.Response) (M,error){
	r, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil,err
	}
	req.Body.Close()
	var res M
	err = json.Unmarshal(r,&res)

	return res,err
}

func getRequestString(req *http.Response) (string,error){
	r, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return "",err
	}
	req.Body.Close()
	return string(r),nil
}