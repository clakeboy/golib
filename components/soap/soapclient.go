package soap

import (
	"ck_go_lib/utils"
	"errors"
	"fmt"
	"strings"

	"regexp"
)

type SoapClient struct {
	ws *Wsdl
	soap_url string
	http_client *utils.HttpClient
	request_url string
}

//新建一个SOAP客户端
func NewSoapClient(soap_url string) (*SoapClient,error) {
	soap := &SoapClient{
		ws:NewWsdl(),
		soap_url:soap_url,
		http_client:utils.NewHttpClient(),
	}

	err := soap.explain()
	if err != nil {
		return nil,err
	}


	return soap,nil
}

//调用一个方法
func (s *SoapClient) Call(func_name string,args ...interface{}) (string,error) {
	funcs := s.ws.GetFunc(func_name)
	if funcs == nil {
		return "",errors.New("not this function!")
	}

	xml_con := "<?xml version='1.0' encoding='UTF-8'?>"
	xml_con += `<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">`
	xml_con += `<soap:Body xmlns:ns1="http://access.vchl.echannel.acic.com/">`
	xml_con += fmt.Sprintf("<ns1:%s>",funcs.Name)
	for i,v := range funcs.Args {
		xml_con += fmt.Sprintf("<ns1:%s><![CDATA[%s]]></ns1:%s>",v.Name,args[i],v.Name)
	}
	xml_con += fmt.Sprintf("</ns1:%s>",funcs.Name)
	xml_con += `</soap:Body>`
	xml_con += `</soap:Envelope>`
	return_str,err := s.httpPost(s.request_url,xml_con)

	if err != nil {
		return "",err
	}

	reg := regexp.MustCompile(fmt.Sprintf("(?si)<%s>(.+?)</%s>", funcs.ReturnArgs[0].Name, funcs.ReturnArgs[0].Name))
	list := reg.FindStringSubmatch(return_str)
	if len(list) <= 0 {
		return "",nil
	}

	return replaceXml(list[1],false),nil
}
//重新设置调用地址
func (s *SoapClient) SetAddress(addr string) {
	s.request_url = addr
}
//发起HTTP请求
func (s *SoapClient) httpPost(url_str ,content string) (string,error) {
	s.http_client.SetHeader("Content-Type","application/soap+xml; charset=utf-8")
	req,err := s.http_client.Request("POST",url_str,strings.NewReader(content))
	if err != nil {
		return "",err
	}

	return string(req.Content),nil
}

//解释WSDL文件
func (s *SoapClient) explain() error {
	wsdl_str,err := s.http_client.Get(s.soap_url)
	if err != nil {
		return err
	}
	err = s.ws.Explain(wsdl_str)
	if err != nil {
		return err
	}

	s.request_url = s.ws.GetAddress()
	return nil
}
//替换XML特殊字符
func replaceXml(str string, toxml bool) string {
	replace_list := map[string]string{
		"&":  "&amp;",
		"<":  "&lt;",
		">":  "&gt;",
		"\"": "&quot;",
		"'":  "&apos;",
	}

	for k, v := range replace_list {
		if toxml {
			str = strings.Replace(str, k, v, -1)
		} else {
			str = strings.Replace(str, v, k, -1)
		}
	}
	return str
}