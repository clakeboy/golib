///腾讯云短信发送
package txsms

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/clakeboy/golib/utils"
	"strings"
	"time"
)

const (
	SMS_SEND_URL = "https://yun.tim.qq.com/v5/tlssmssvr/sendsms?sdkappid=%s&random=%s"
)

//电话号码信息
type PhoneData struct {
	Mobile     string `json:"mobile"`     //手机号码
	Nationcode string `json:"nationcode"` //国家（或地区）码
}

//短信发送数据
type SendData struct {
	Ext              string      `json:"ext"`    //用户的 session 内容，腾讯 server 回包中会原样返回，可选字段，不需要时设置为空
	Extend           string      `json:"extend"` //短信码号扩展号，格式为纯数字串，其他格式无效。默认没有开通
	Params           []string    `json:"params"` //模板参数，具体使用方法请参见下方说明。若模板没有参数，请设置为空数组
	Sig              string      `json:"sig"`    //App 凭证，具体计算方式请参见下方说明
	Sign             string      `json:"sign"`   //短信签名内容，使用 UTF-8 编码，必须填写已审核通过的签名
	Tel              interface{} `json:"tel"`    //电话号码信息 ,数据为群发,单个为单发 对应数据 PhoneData
	Time             int         `json:"time"`   //请求发起时间，UNIX 时间戳（单位：秒），如果和系统时间相差超过 10 分钟则会返回失败
	TplId            int         `json:"tpl_id"` //模板 ID，必须填写已审核通过的模板 ID
	*utils.JsonParse `json:"-"`
}

func NewSendData(data *SendData) *SendData {
	if data == nil {
		data = &SendData{}
	}
	data.JsonParse = utils.NewJsonParse(data)
	return data
}

//发送返回信息
type SentResponse struct {
	Result int    `json:"result"` //错误码，0表示成功（计费依据），非0表示失败
	Errmsg string `json:"errmsg"` //错误消息，result 非0时的具体错误信息
	Ext    string `json:"ext"`    //用户的 session 内容，腾讯 server 回包中会原样返回
	Fee    int    `json:"fee"`    //短信计费的条数
	Sid    string `json:"sid"`    //本次发送标识 ID，标识一次短信下发记录
}

//腾讯云短信发送类
type TxSms struct {
	appKey string //业务 appid 对应的 appkey
	appId  string //业务 appid
}

func NewTxSms(appId, appKey string) *TxSms {
	return &TxSms{
		appId:  appId,
		appKey: appKey,
	}
}

//发送信息
func (t *TxSms) Send(data *SendData) (*SentResponse, error) {
	ranStr := utils.RandStr(16, nil)
	timeStamp := int(time.Now().Unix())
	data.Sig = t.sig(data, ranStr, timeStamp)
	data.Time = timeStamp
	urlStr := fmt.Sprintf(SMS_SEND_URL, t.appId, ranStr)
	res, err := utils.HttpPostJsonBytes(urlStr, data.ToJson())
	if err != nil {
		return nil, err
	}
	resp := new(SentResponse)
	err = json.Unmarshal(res, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

//sig 伪代码
func (t *TxSms) sig(data *SendData, ranStr string, time int) string {
	var mobile string
	switch data.Tel.(type) {
	case []*PhoneData:
		tel := data.Tel.([]*PhoneData)
		var tels []string
		for _, v := range tel {
			tels = append(tels, v.Mobile)
		}
		mobile = strings.Join(tels, ",")
	case *PhoneData:
		tel := data.Tel.(*PhoneData)
		mobile = tel.Mobile
	}
	sig := fmt.Sprintf("appkey=%s&random=%s&time=%d&mobile=%s", t.appKey, ranStr, time, mobile)
	hash := sha256.New()
	hash.Write([]byte(sig))
	hashSig := hash.Sum(nil)
	return hex.EncodeToString(hashSig)
}
