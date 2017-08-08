package soap

import (
	"ck_go_lib/utils"
	"errors"
	"fmt"
	"strings"
	"encoding/xml"
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
func (s *SoapClient) Call(func_name string,args ...interface{}) (error) {
	funcs := s.ws.GetFunc(func_name)
	if funcs == nil {
		return errors.New("not this function!")
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
	fmt.Println(xml_con)
	s.httpPost(s.request_url,xml_con)
	return nil
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
	fmt.Printf("%+v\n",req)
	fmt.Println(string(req.Content))
	xml_map := utils.XMLMap{}
	err = xml.Unmarshal(req.Content,&xml_map)
	if err != nil {
		return "",err
	}
	fmt.Println(xml_map)
	return "",nil
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
