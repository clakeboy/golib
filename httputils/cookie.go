package httputils

import (
	"net/http"
	"ck_go_lib/utils"
)

//cookie 操作
type HttpCookie struct {
	request *http.Request
	writer http.ResponseWriter
	aes *utils.AesEncrypt
	options *CookieOptions
}

type CookieOptions struct {
	AesKey  string //加密KEY
	Domain  string  //设置域
	Path    string  //设置路径
	Secure  bool //是否安全传送
	HttpOnly  bool //是否开启只用于HTTP
}

func NewHttpCookie(req *http.Request,writer http.ResponseWriter,options *CookieOptions) *HttpCookie {
	if options == nil {
		options = &CookieOptions{
			AesKey:"ck-cookie",
			Domain:"",
			Path:"/",
			Secure:false,
			HttpOnly:true,
		}
	}

	return &HttpCookie{
		request:req,
		writer:writer,
		options:options,
		aes:utils.NewAes(string(options.AesKey)),
	}
}
//设置一个COOKIE maxAge
func (c *HttpCookie) Set(name string,val string,maxAge int) {
	aesStr,_ := c.aes.EncryptString(val)
	http.SetCookie(c.writer,&http.Cookie{
		Name:     name,
		Value:    aesStr,
		MaxAge:   maxAge,
		Path:     c.options.Path,
		Domain:   c.options.Domain,
		Secure:   c.options.Secure,
		HttpOnly: c.options.HttpOnly,
	})
}
//得到一个COOKIE
func (c *HttpCookie) Get(name string) (string,error) {
	cookie,err := c.request.Cookie(name)
	if err != nil {
		return "",err
	}

	val,_ := c.aes.DecryptString(cookie.Value)
	return val,nil
}
//删除一个COOKIE
func (c *HttpCookie) Delete(name string) {
	http.SetCookie(c.writer,&http.Cookie{
		Name:     name,
		Value:    "",
		MaxAge:   -1,
		Path:     c.options.Path,
		Domain:   c.options.Domain,
		Secure:   c.options.Secure,
		HttpOnly: c.options.HttpOnly,
	})
}