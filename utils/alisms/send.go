package alisms

import (
	"../../utils"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"
)

//短信配置
type AliSmsConfig struct {
	AccessID  string `json:"access_id" yaml:"access_id"`
	AccessKey string `json:"access_key" yaml:"access_key"`
	RegionID  string `json:"region_id" yaml:"region_id"`
	Template  string `json:"template" yaml:"template"`
	SignName  string `json:"sign_name" yaml:"sign_name"`
}

//短信发送回复消息
type AliSmsRes struct {
	RequestId string `json:"RequestId" bson:"request_id"`
	Message   string `json:"Message" bson:"message"`
	Code      string `json:"Code" bson:"code"`
	BizId     string `json:"BizId" bson:"biz_id"`
}

//阿里短信发送
type AliSms struct {
	signName  string
	accessID  string
	accessKey string
	regionID  string
}

//阿里短信
func NewAliSms(access_id, access_key, region_id, sign_name string) *AliSms {
	return &AliSms{
		accessID:  access_id,
		accessKey: access_key,
		signName:  sign_name,
		regionID:  region_id,
	}
}

//设置短信签名
func (a *AliSms) SetSignName(sign_name string) {
	a.signName = sign_name
}

//发送短信
func (a *AliSms) SendSms(phone_num string, template_name string, params utils.M) (*AliSmsRes, error) {
	data := utils.M{
		"SignatureMethod":  "HMAC-SHA1",
		"SignatureNonce":   utils.CreateUUID(true),
		"AccessKeyId":      a.accessID,
		"SignatureVersion": "1.0",
		"Timestamp":        time.Now().UTC().Format("2006-01-02T15:04:05Z"),
		"Format":           "JSON",
		"Action":           "SendSms",
		"Version":          "2017-05-25",
		"RegionId":         a.regionID,
		"PhoneNumbers":     phone_num,
		"SignName":         a.signName,
		"TemplateParam":    params.ToJsonString(),
		"TemplateCode":     template_name,
		"OutId":            "123",
	}

	sign_str, query_str := a.Sign(data)

	url_str := fmt.Sprintf("http://dysmsapi.aliyuncs.com/?Signature=%s&%s", a.specialUrlEncode(sign_str), query_str)

	return a.send(url_str)
}

//发送请求
func (a *AliSms) send(url_str string) (*AliSmsRes, error) {
	res_str, err := utils.HttpGet(url_str)
	if err != nil {
		return nil, err
	}

	res_data := AliSmsRes{}
	err = json.Unmarshal([]byte(res_str), &res_data)
	if err != nil {
		return nil, err
	}

	return &res_data, err
}

//签名
func (a *AliSms) Sign(data utils.M) (sign string, query_str string) {
	buf := []string{}
	keys := utils.MapKeys(data)
	sort.Strings(keys)
	for _, key := range keys {
		buf = append(buf, fmt.Sprintf("%s=%s", a.specialUrlEncode(key), a.specialUrlEncode(data[key].(string))))
	}

	sortedQueryString := strings.Join(buf, "&")
	sign_str := fmt.Sprintf("%s&%s&%s", "GET", a.specialUrlEncode("/"), a.specialUrlEncode(sortedQueryString))
	mac := hmac.New(sha1.New, []byte(a.accessKey+"&"))
	mac.Write([]byte(sign_str))

	sign = base64.StdEncoding.EncodeToString(mac.Sum(nil))
	query_str = sortedQueryString
	return
}

//特殊编码
func (a *AliSms) specialUrlEncode(str string) string {
	en_str := url.QueryEscape(str)
	en_str = strings.Replace(en_str, "+", "%20", -1)
	en_str = strings.Replace(en_str, "*", "%2A", -1)
	en_str = strings.Replace(en_str, "%7E", "~", -1)
	return en_str
}
