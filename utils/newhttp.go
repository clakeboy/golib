package utils

import (
	"net/http"
	"net/url"
	"strings"
	"io/ioutil"
	"time"
	"io"
	"fmt"
	"net/http/cookiejar"
)
//http 请求返回数据
type HttpRequestData struct {
	Status  string
	StatusCode int
	Headers http.Header
	Content []byte
	Cookie *HttpCookies
}
//http 请求客户端
type HttpClient struct {
	client  *http.Client
	headers M
	cookies []*http.Cookie
	lastRequest *HttpRequestData
}
//cookie集合处理
type HttpCookies struct {
	Cookies []*http.Cookie
}

func NewHttpCookies(cookies ...*http.Cookie) *HttpCookies {
	return &HttpCookies{
		Cookies:YN(cookies == nil,[]*http.Cookie{},cookies).([]*http.Cookie),
	}
}

func (hc *HttpCookies) GetCookieString() string {
	arr := []string{}
	for _,v := range hc.Cookies {
		arr = append(arr,fmt.Sprintf("%s=%s",v.Name,v.Value))
	}

	return strings.Join(arr,"; ")
}
//创建一个新的HTTP连接客户端
func NewHttpClient() *HttpClient {
	jar, _ := cookiejar.New(nil)
	client := &http.Client{}
	client.Timeout = time.Second * 30
	client.Jar = jar
	return &HttpClient{client: client, headers: M{},cookies:[]*http.Cookie{}}
}
//发起一个POST请求
func (h *HttpClient) Post(url_str string, data M) ([]byte, error) {
	h.SetHeader("Content-Type","application/x-www-form-urlencoded")
	post_data := &url.Values{}

	for k, v := range data {
		post_data.Add(k, v.(string))
	}

	req, err := h.Request("POST",url_str,strings.NewReader(post_data.Encode()))

	if err != nil {
		return nil, err
	}

	return req.Content, nil
}
//发起一个POST JSON请求
func (h *HttpClient) PostJson(url_str string,data M) ([]byte,error) {
	h.SetHeader("Content-Type","application/json")
	req, err := h.Request("POST",url_str,strings.NewReader(data.ToJsonString()))
	if err != nil {
		return nil, err
	}

	return req.Content, nil
}
//发起一个POST XML请求
func (h *HttpClient) PostXml(url_str string,data string) ([]byte,error) {
	h.SetHeader("Content-Type","text/xml;charset=utf-8")
	req, err := h.Request("POST",url_str,strings.NewReader(data))
	if err != nil {
		return nil, err
	}

	return req.Content, nil
}
//发起一个GET请求
func (h *HttpClient) Get(url_str string) ([]byte, error) {
	resp,err := h.Request("GET",url_str,nil)
	if err != nil {
		return nil,err
	}
	return resp.Content,nil
}
//发起一个http request 请求
func (h *HttpClient) Request(method string,url_str string,content io.Reader) (*HttpRequestData,error) {
	req,err := http.NewRequest(method,url_str,content)
	if err != nil {
		return nil,err
	}

	if len(h.headers) > 0 {
		for k, v := range h.headers {
			req.Header.Set(k, v.(string))
		}
	}

	if len(h.cookies) > 0 {
		for _,v := range h.cookies {
			req.AddCookie(v)
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
	res.Cookie = NewHttpCookies(resp.Cookies()...)
	h.lastRequest = res
	return res,nil
}
//设置请求头信息
func (h *HttpClient) SetHeader(key, val string) {
	h.headers[key] = val
}
//设置请求超时时间
func (h *HttpClient) SetTimeout(sc time.Duration) {
	h.client.Timeout = sc
}
//设置请求带的COOKIE信息
func (h *HttpClient) SetCookie(cookie *http.Cookie) {
	if cookie == nil {
		return
	}
	h.cookies = append(h.cookies,cookie)
}
//清除HTTP当前请求的数据
func (h *HttpClient) Clear() {
	h.cookies = []*http.Cookie{}
	h.headers = M{}
}
//得到最后一次请求的返回数据
func (h *HttpClient) GetLastResponse() *HttpRequestData {
	return h.lastRequest
}