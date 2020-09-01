package components

import (
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/clakeboy/golib/utils"
	"regexp"
	"sort"
	"strings"
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

//用户ACCESS
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
	Subscribe      int    `json:"subscribe"`
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

//文件名
var fileReg = regexp.MustCompile(`filename="(.+)"`)

func NewWeixin(app_id, app_secret string) *Weixin {
	return &Weixin{
		appId:     app_id,
		appSecret: app_secret,
	}
}

//发起HTTP连接
func (w *Weixin) Http(uri string) (string, error) {
	s, err := utils.HttpGet(uri)
	if err != nil {
		return "", err
	}
	var info WxError
	err = json.Unmarshal([]byte(s), &info)
	if err != nil {
		return "", err
	}

	if info.ErrorCode != 0 {
		err := errors.New(fmt.Sprintf("error_code:%d,errmsg:%s", info.ErrorCode, info.ErrorMsg))
		return "", err
	}

	return s, nil
}

//发起HTTP POSTJSON
func (w *Weixin) post(uri string, data utils.M) (string, error) {
	res, err := utils.HttpPostJsonString(uri, data)
	if err != nil {
		return "", err
	}
	var info WxError
	err = json.Unmarshal([]byte(res), &info)
	if err != nil {
		return res, err
	}
	if info.ErrorCode != 0 {
		err := errors.New(fmt.Sprintf("error_code:%d,errmsg:%s", info.ErrorCode, info.ErrorMsg))
		return "", err
	}

	return res, nil
}

//得到公众号全局 ACCESS_TOKEN
func (w *Weixin) GetAccessToken() (*WxAccessToken, error) {
	url := fmt.Sprintf(
		"https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%v&secret=%v",
		w.appId,
		w.appSecret)
	raw, err := w.Http(url)
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
func (w *Weixin) WebAuth(redirect_uri string, mpid string, base bool) string {
	scope := "snsapi_userinfo"
	if base {
		scope = "snsapi_base"
	}

	u := fmt.Sprintf("https://open.weixin.qq.com/connect/oauth2/authorize?appid=%v&redirect_uri=%v&response_type=code&scope=%s&state=%v%v",
		w.appId,
		redirect_uri,
		scope,
		mpid,
		"#wechat_redirect",
	)

	return u
}

//使用code获取用户access token
func (w *Weixin) GetUserAccessToken(code string) (*WxUserAccessToken, error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/oauth2/access_token?%v&%v&%v&%v",
		"appid="+w.appId,
		"secret="+w.appSecret,
		"code="+code,
		"grant_type=authorization_code",
	)
	raw, err := w.Http(url)
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
func (w *Weixin) GetUserInfoAccessToken(access_token string, openid string) (*WxUser, error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/userinfo?%v&%v&%v",
		"access_token="+access_token,
		"openid="+openid,
		"lang=zh_CN",
	)
	raw, err := w.Http(url)
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
func (w *Weixin) GetUserInfo(access_token string, openid string) (*WxUserInfo, error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/user/info?access_token=%v&openid=%v&lang=zh_CN",
		access_token,
		openid,
	)
	raw, err := w.Http(url)
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
func (w *Weixin) GetJsapiTicket(access_token string) (*WxJsapiTicket, error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=%v&type=jsapi", access_token)
	raw, err := w.Http(url)
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
func (w *Weixin) SignJsTicket(data map[string]interface{}) string {
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
func (w *Weixin) SendTemplateMessage(access_token string, data utils.M) error {
	url_str := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=%s", access_token)
	_, err := w.post(url_str, data)
	if err != nil {
		return err
	}

	return nil
}

type MediaData struct {
	FileType string
	FileName string
	Content  []byte
}

//得到临时素材
func (w *Weixin) GetMedia(access_token string, media_id string) (*MediaData, error) {
	url_str := fmt.Sprintf(
		"https://api.weixin.qq.com/cgi-bin/media/get?access_token=%s&media_id=%s",
		access_token,
		media_id,
	)
	ck_http := utils.NewHttpClient()
	res, err := ck_http.Request("GET", url_str, nil)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("request error code: %d", res.StatusCode))
	}

	if res.Headers.Get("Content-Type") != "text/plain" || res.Headers.Get("Content-Type") != "application/json" {
		dis := res.Headers.Get("Content-disposition")
		//reg := regexp.MustCompile(`filename="(.+)"`)
		list := fileReg.FindStringSubmatch(dis)
		media := &MediaData{
			FileType: res.Headers.Get("Content-Type"),
			FileName: list[1],
			Content:  res.Content,
		}
		return media, nil
	}
	err_msg := utils.M{}
	json.Unmarshal(res.Content, &err_msg)
	return nil, errors.New(fmt.Sprintf("code:%v,msg:%v", err_msg["errcode"], err_msg["errmsg"]))
}

//微信二维码结构
type WxQrCode struct {
	Ticket        string `json:"ticket"`
	ExpireSeconds int    `json:"expire_seconds"`
	Url           string `json:"url"`
}

//创建临时二维码
func (w *Weixin) CreateTempQrCode(token, content string, time_s int) (*WxQrCode, error) {
	data := utils.M{
		"expire_seconds": time_s,
		"action_name":    "QR_STR_SCENE",
		"action_info": utils.M{
			"scene": utils.M{
				"scene_str": content,
			},
		},
	}

	url_str := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/qrcode/create?access_token=%v", token)

	res, err := w.post(url_str, data)
	if err != nil {
		return nil, err
	}

	res_data := &WxQrCode{}
	json.Unmarshal([]byte(res), res_data)

	return res_data, nil
}

//创建永久二维码
func (w *Weixin) CreateQrCode(token, content string) (*WxQrCode, error) {
	data := utils.M{
		"action_name": "QR_LIMIT_STR_SCENE",
		"action_info": utils.M{
			"scene": utils.M{
				"scene_str": content,
			},
		},
	}

	url_str := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/qrcode/create?access_token=%v", token)

	res, err := w.post(url_str, data)
	if err != nil {
		return nil, err
	}

	res_data := &WxQrCode{}
	json.Unmarshal([]byte(res), res_data)

	return res_data, nil
}

//创建微信小程序二维码
func (w *Weixin) CreateWxAppQrCode(token, scene, page string, width int) (string, error) {
	data := utils.M{
		"scene": scene,
		"page":  page,
		"width": width,
	}

	url_str := fmt.Sprintf("https://api.weixin.qq.com/wxa/getwxacodeunlimit?access_token=%s", token)
	res, err := utils.HttpPostJsonString(url_str, data)
	if err != nil {
		return "", err
	}
	return res, err
}

//向用户发送客户消息
func (w *Weixin) SendCustomMessage(token, openid, msg_type string, data utils.M) (string, error) {
	msg := utils.M{
		"touser":  openid,
		"msgtype": msg_type,
		msg_type:  data,
	}

	url_str := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/message/custom/send?access_token=%s", token)
	res, err := utils.HttpPostJsonString(url_str, msg)
	if err != nil {
		return "", err
	}
	return res, err
}

///微信标签管理
//用户标签JSON
type UserTag struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Count int    `json:"count"`
}

//创建用户标签
func (w *Weixin) CreateTag(token, name string) (string, error) {
	data := utils.M{
		"tag": UserTag{
			Name: name,
		},
	}

	urlStr := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/tags/create?access_token=%s", token)
	//res, err := utils.HttpPostJsonString(urlStr, data)
	res, err := w.post(urlStr, data)
	if err != nil && res == "" {
		return "", err
	}

	return res, err
}

type TagList struct {
	Tags []*UserTag `json:"tags"`
}

//获取已创建标签
func (w *Weixin) GetTags(token string) (*TagList, error) {
	urlStr := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/tags/get?access_token=%s", token)
	res, err := utils.HttpGet(urlStr)
	if err != nil {
		return nil, err
	}
	tags := &TagList{}
	err = json.Unmarshal([]byte(res), tags)
	if err != nil {
		return nil, err
	}
	return tags, err
}

//修改标签名
func (w *Weixin) UpdateTag(token string, id int, name string) (string, error) {
	data := utils.M{
		"tag": utils.M{
			"id":   id,
			"name": name,
		},
	}
	urlStr := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/tags/update?access_token=%s", token)
	res, err := utils.HttpPostJsonString(urlStr, data)
	if err != nil {
		return "", err
	}
	return res, err
}

//删除标签
func (w *Weixin) DeleteTag(token string, id int) (string, error) {
	data := utils.M{
		"tag": utils.M{
			"id": id,
		},
	}
	urlStr := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/tags/delete?access_token=%s", token)
	res, err := utils.HttpPostJsonString(urlStr, data)
	if err != nil {
		return "", err
	}
	return res, err
}

// 获取标签下粉丝列表
func (w *Weixin) GetTagUsers(token string, tagId int, nextOpenid string) (string, error) {
	data := utils.M{
		"tagid":       tagId,
		"next_openid": nextOpenid,
	}
	urlStr := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/user/tag/get?access_token=%s", token)
	res, err := utils.HttpPostJsonString(urlStr, data)
	if err != nil {
		return "", err
	}
	return res, err
}

//批量为用户打标签
func (w *Weixin) SetUserTag(token string, tagId int, openIdList []string) (string, error) {
	data := utils.M{
		"tagid":       tagId,
		"openid_list": openIdList,
	}
	urlStr := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/tags/members/batchtagging?access_token=%s", token)
	//res, err := utils.HttpPostJsonString(urlStr, data)
	res, err := w.post(urlStr, data)
	if err != nil {
		return "", err
	}
	return res, err
}

//批量为用户取消标签
func (w *Weixin) CancelUserTag(token string, tagId int, openIdList []string) (string, error) {
	data := utils.M{
		"tagid":       tagId,
		"openid_list": openIdList,
	}
	urlStr := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/tags/members/batchuntagging?access_token=%s", token)
	//res, err := utils.HttpPostJsonString(urlStr, data)
	res, err := w.post(urlStr, data)
	if err != nil {
		return "", err
	}
	return res, err
}

//获取用户身上的标签列表
func (w *Weixin) GetUserTags(token, openId string) (string, error) {
	data := utils.M{
		"openid": openId,
	}
	urlStr := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/tags/getidlist?access_token=%s", token)
	res, err := utils.HttpPostJsonString(urlStr, data)
	if err != nil {
		return "", err
	}
	return res, err
}
