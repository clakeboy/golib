package utils

import (
	"net/http"
	"net/url"
	"strings"
	"io/ioutil"
	"time"
)

type HttpRequestData struct {
	Status  string
	StatusCode int
	Headers http.Header
	Content []byte
}

type HttpClient struct {
	client  *http.Client
	headers HD
}

func NewHttpClient() *HttpClient {
	client := &http.Client{}
	client.Timeout = time.Second * 30
	return &HttpClient{client: client, headers: HD{}}
}

func (h *HttpClient) Post(url_str string, data M) ([]byte, error) {
	post_data := &url.Values{}

	for k, v := range data {
		post_data.Add(k, v.(string))
	}

	req, err := http.NewRequest("POST", url_str, strings.NewReader(post_data.Encode()))
	if err != nil {
		return nil, err
	}

	if len(h.headers) > 0 {
		for k, v := range h.headers {
			req.Header.Set(k, v.(string))
		}
	} else {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	resp,err := h.client.Do(req)
	if err != nil {
		return nil,err
	}
	defer resp.Body.Close()

	body,err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil,err
	}

	return body, nil
}

func (h *HttpClient) Get(url_str string) ([]byte, error) {
	req,err := http.NewRequest("GET",url_str,nil)
	if err != nil {
		return nil,err
	}
	resp,err :=  h.client.Do(req)
	if err != nil {
		return nil,err
	}
	defer resp.Body.Close()

	body,err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil,err
	}
	return body,nil
}

func (h *HttpClient) Request(method string,url_str string,data M) (*HttpRequestData,error) {
	post_data := &url.Values{}

	for k, v := range data {
		post_data.Add(k, v.(string))
	}

	req,err := http.NewRequest(method,url_str,strings.NewReader(post_data.Encode()))
	if err != nil {
		return nil,err
	}

	if method == "POST" {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	if len(h.headers) > 0 {
		for k, v := range h.headers {
			req.Header.Set(k, v.(string))
		}
	}

	resp,err := h.client.Do(req)
	if err != nil {
		return nil,err
	}
	defer resp.Body.Close()

	res := &HttpRequestData{}

	body,err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil,err
	}
	res.Content = body
	res.Status = resp.Status
	res.StatusCode = resp.StatusCode
	res.Headers = resp.Header

	return res,nil
}

func (h *HttpClient) SetHeader(key, val string) {
	h.headers[key] = val
}

func (h *HttpClient) SetTimeout(sc time.Duration) {
	h.client.Timeout = sc
}
