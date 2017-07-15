package components

import (
	"ck_go_lib/utils"
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"regexp"
)

type WxError struct {
	ErrorCode int    `json:"errcode"`
	ErrorMsg  string `json:"errmsg"`
}

type WxAccessToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type WxJsapiTicket struct {
	Ticket    string `json:"ticket"`
	ExpiresIn int    `json:"expires_in"`
}

type Weixin struct {
	appId     string
	appSecret string
}

type WxUserAccessToken struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	OpenId       string `json:"openid"`
	Scope        string `json:"scope"`
	UnionId      string `json:"unionid"`
}

//网页受权接口获取用户数据
type WxUser struct {
	NickName string `json:"nickname"`
	OpenId   string `json:"openid"`
	Sex      int    `json:"sex"`
	Province string `json:"province"`
	City     string `json:"city"`
	Country  string `json:"country"`
	HeadImg  string `json:"headimgurl"`
	UnionId  string `json:"unionid"`
}

//用公众号接口获取的用户数据
type WxUserInfo struct {
	Subscribe      int `json:"subscribe"`
	Subscribe_time int    `json:"subscribe_time"`
	Nickname       string `json:"nickname"`
	OpenId         string `json:"openid"`
	Sex            int    `json:"sex"`
	Province       string `json:"province"`
	City           string `json:"city"`
	Country        string `json:"country"`
	HeadImg        string `json:"headimgurl"`
	Language       string `json:"language"`
	UnionId        string `json:"unionid"`
	Remark         string `json:"remark"`
	Groupid        int    `json:"groupid"`
	Tagid_list     []int  `json:"tagid_list"`
}

func NewWeixin(app_id, app_secret string) *Weixin {
	return &Weixin{
		appId:     app_id,
		appSecret: app_secret,
	}
}

//发起HTTP连接
func (this *Weixin) Http(uri string) (string, error) {
	s,err := utils.HttpGet(uri)
	if err != nil {
		return "", err
	}
	var info WxError
	err = json.Unmarshal([]byte(s), &info)
	if err != nil {
		return "", err
	}

	if info.ErrorCode != 0 {
		err := errors.New(fmt.Sprintf("error_code:%s,errmsg:%s", info.ErrorCode, info.ErrorMsg))
		return "", err
	}

	return s, nil
}

//发起HTTP POSTJSON
func (w *Weixin) post(uri string,data utils.M) (string,error) {
	res,err := utils.HttpPostJsonString(uri,data)
	if err != nil {
		return "",err
	}
	var info WxError
	err = json.Unmarshal([]byte(res), &info)
	if err != nil {
		return "", err
	}
	if info.ErrorCode != 0 {
		err := errors.New(fmt.Sprintf("error_code:%s,errmsg:%s", info.ErrorCode, info.ErrorMsg))
		return "", err
	}

	return res,nil
}

//得到公众号全局 ACCESS_TOKEN
func (this *Weixin) GetAccessToken() (*WxAccessToken, error) {
	url := fmt.Sprintf(
		"https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%v&secret=%v",
		this.appId,
		this.appSecret)
	raw, err := this.Http(url)
	if err != nil {
		return nil, err
	}
	var info WxAccessToken

	err = json.Unmarshal([]byte(raw), &info)
	if err != nil {
		return nil, err
	}

	return &info, nil
}

//发起 WEB 验证受权
func (this *Weixin) WebAuth(redirect_uri string,mpid string) string {
	u := fmt.Sprintf("https://open.weixin.qq.com/connect/oauth2/authorize?appid=%v&redirect_uri=%v&response_type=code&scope=snsapi_userinfo&state=%v%v",
		this.appId,
		redirect_uri,
		mpid,
		"#wechat_redirect",
	)
	return u
}

//使用code获取用户access token
func (this *Weixin) GetUserAccessToken(code string) (*WxUserAccessToken, error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/oauth2/access_token?%v&%v&%v&%v",
		"appid="+this.appId,
		"secret="+this.appSecret,
		"code="+code,
		"grant_type=authorization_code",
	)
	raw, err := this.Http(url)
	if err != nil {
		return nil, err
	}

	var user_access WxUserAccessToken
	err = json.Unmarshal([]byte(raw), &user_access)
	if err != nil {
		return nil, err
	}
	return &user_access, nil
}

//使用用户access token 获取用户信息
func (this *Weixin) GetUserInfoAccessToken(access_token string, openid string) (*WxUser, error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/userinfo?%v&%v&%v",
		"access_token="+access_token,
		"openid="+openid,
		"lang=zh_CN",
	)
	raw, err := this.Http(url)
	if err != nil {
		return nil, err
	}

	var wxuser WxUser
	err = json.Unmarshal([]byte(raw), &wxuser)
	if err != nil {
		return nil, err
	}

	return &wxuser, nil
}

//拉取用户信息
func (this *Weixin) GetUserInfo(access_token string, openid string) (*WxUserInfo,error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/user/info?access_token=%v&openid=%v&lang=zh_CN",
		access_token,
		openid,
	)
	raw, err := this.Http(url)
	if err != nil {
		return nil, err
	}

	var wxuser WxUserInfo
	err = json.Unmarshal([]byte(raw), &wxuser)
	if err != nil {
		return nil, err
	}

	return &wxuser, nil
}

//获取JsTicket
func (this *Weixin) GetJsapiTicket(access_token string) (*WxJsapiTicket, error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=%v&type=jsapi", access_token)
	raw, err := this.Http(url)
	if err != nil {
		return nil, err
	}

	var ticket WxJsapiTicket
	err = json.Unmarshal([]byte(raw), &ticket)
	if err != nil {
		return nil, err
	}

	return &ticket, nil
}

//进行JSTICKET签名
func (this *Weixin) SignJsTicket(data map[string]interface{}) string {
	keys := utils.MapKeys(data)
	sort.Strings(keys)
	val := []string{}
	for _, k := range keys {
		val = append(val, fmt.Sprintf("%v=%v", k, data[k]))
	}

	sign_str := strings.Join(val, "&")

	h := sha1.New()
	h.Write([]byte(sign_str))

	return fmt.Sprintf("%x", h.Sum(nil))
}

//发送模板消息
func (w *Weixin) SendTemplateMessage(access_token string,data utils.M) error {
	url_str := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=%s",access_token)
	_,err := w.post(url_str,data)
	if err != nil {
		return err
	}

	return nil
}

type MediaData struct {
	FileType string
	FileName string
	Content []byte
}

//得到临时素材
func (w *Weixin) GetMedia(access_token string,media_id string) (*MediaData,error) {
	url_str := fmt.Sprintf(
		"https://api.weixin.qq.com/cgi-bin/media/get?access_token=%s&media_id=%s",
		access_token,
		media_id,
	)

	ck_http := utils.NewHttpClient()
	res,err := ck_http.Request("GET",url_str,nil)
	if err != nil {
		return nil,err
	}

	if res.StatusCode != 200 {
		return nil,errors.New(fmt.Sprintf("request error code: %d",res.StatusCode))
	}

	if res.Headers.Get("Content-Type") != "text/plain" {
		dis := res.Headers.Get("Content-disposition")
		reg := regexp.MustCompile(`filename="(.+)"`)
		list := reg.FindStringSubmatch(dis)
		media := &MediaData{
			FileType:res.Headers.Get("Content-Type"),
			FileName:list[1],
			Content:res.Content,
		}
		return media,nil
	}
	err_msg := utils.M{}
	json.Unmarshal(res.Content,&err_msg)
	return nil,errors.New(fmt.Sprintf("code:%v,msg:%v",err_msg["errcode"],err_msg["errmsg"]))
}