package alisms

import (
	"fmt"
	"github.com/clakeboy/golib/utils"
	"log"
	"net/url"
	"testing"
)

func TestAliSms_Sign(t *testing.T) {
	str := url.QueryEscape("+*~")
	pstr := url.PathEscape("+*~")
	fmt.Println(str, pstr)
	ss := "http://dysmsapi.aliyuncs.com/?Signature=zJDF%2BLrzhj%2FThnlvIToysFRq6t4%3D&AccessKeyId=testId&Action=SendSms&Format=XML&OutId=123&PhoneNumbers=15300000001&RegionId=cn-hangzhou&SignName=%E9%98%BF%E9%87%8C%E4%BA%91%E7%9F%AD%E4%BF%A1%E6%B5%8B%E8%AF%95%E4%B8%93%E7%94%A8&SignatureMethod=HMAC-SHA1&SignatureNonce=45e25e9b-0a6f-4070-8c85-2956eda1b466&SignatureVersion=1.0&TemplateCode=SMS_71390007&TemplateParam=%7B%22customer%22%3A%22test%22%7D&Timestamp=2017-07-12T02%3A42%3A19Z&Version=2017-05-25"
	ss1 := "http://dysmsapi.aliyuncs.com/?Signature=zJDF%2BLrzhj%2FThnlvIToysFRq6t4%3D&AccessKeyId=testId&Action=SendSms&Format=XML&OutId=123&PhoneNumbers=15300000001&RegionId=cn-hangzhou&SignName=%E9%98%BF%E9%87%8C%E4%BA%91%E7%9F%AD%E4%BF%A1%E6%B5%8B%E8%AF%95%E4%B8%93%E7%94%A8&SignatureMethod=HMAC-SHA1&SignatureNonce=45e25e9b-0a6f-4070-8c85-2956eda1b466&SignatureVersion=1.0&TemplateCode=SMS_71390007&TemplateParam=%7B%22customer%22%3A%22test%22%7D&Timestamp=2017-07-12T02%3A42%3A19Z&Version=2017-05-25"
	fmt.Println(utils.EncodeMD5Std(ss))
	fmt.Println(utils.EncodeMD5Std(ss1))
}

func TestAliSms_SendSms(t *testing.T) {
	sms := NewAliSms("LTAIqNfX49aao0c3", "bL2UY9915xuLPha3evWLYmlYKtyrSl", "cn-hangzhou", "兔保长")
	res, err := sms.SendSms("18623451781", "SMS_80285147", utils.M{
		"code": "test",
	})
	if err != nil {
		log.Println(err)
	}

	fmt.Println(res)
}
