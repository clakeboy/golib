package soap

import (
	"ck_go_lib/utils"
	"encoding/xml"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type SoapMessageEnvelope struct {
	XMLName     xml.Name         `xml:"Envelope"`
	Body        *SoapMessageBody `xml:"Body"`
	BodyContent string           `xml:"Body"`
}

type SoapMessageBody struct {
	XMLName xml.Name          `xml:"Body"`
	Fault   *SoapMessageFault `xml:"Fault"`
}

type SoapMessageFault struct {
	XMLName     xml.Name `xml:"Fault"`
	Faultcode   string   `xml:"faultcode"`
	Faultstring string   `xml:"faultstring"`
	Detail      string   `xml:"detail"`
}

type SoapClient struct {
	ws            *Wsdl             //wsdl 解释对像
	soapUrl       string            //wsdl 文档URL地址
	httpClient    *utils.HttpClient //HTTP 调用对像
	requestUrl    string            //SOAP 请求URL地址
	soapVersion   string            //SOAP 结构XML版本
	soapAction    string            //SOAP 当前执行Action
	soapService   int               //SOAP 默认使用SERVICE
	soapBuilding  string            //SOAP 默认使用building
	soapCache     bool              //SOAP 是否开启缓存
	soapCacheTime int               //SOAP 缓存时间
	returnAll     bool              //SOAP 返回时是否返回全格式
	options       utils.M           //配置项
}

//新建一个SOAP客户端
func NewSoapClient(soap_url string) (*SoapClient, error) {
	soap := newDefaultSoap()
	soap.soapUrl = soap_url
	wsdl_con, err := soap.getRemoteWsdl()
	if err != nil {
		return nil, err
	}
	err = soap.explain(wsdl_con)
	if err != nil {
		return nil, err
	}

	return soap, nil
}

//用WSDL内容新建一个SOAP客户端
func NewSoapClientForContent(content string) (*SoapClient, error) {
	soap := newDefaultSoap()

	err := soap.explain([]byte(content))
	if err != nil {
		return nil, err
	}

	return soap, nil
}

func newDefaultSoap() *SoapClient {
	soap := &SoapClient{
		ws:          NewWsdl(),
		soapUrl:     "",
		httpClient:  utils.NewHttpClient(),
		soapVersion: "1.1",
		soapService: 0,
		returnAll:   false,
	}

	return soap
}
//设置接口返回时是否返回全格式
func (s *SoapClient) SetReturnAll(yes bool) {
	s.returnAll = yes
}

//调用一个方法
//func (s *SoapClient) __Call(func_name string, args ...interface{}) (string, error) {
//	funcs := s.ws.GetFunc(func_name)
//	if funcs == nil {
//		return "", errors.New("Not this function!")
//	}
//	body_ns := s.ws.GetNamespace(func_name)
//	xml_con := "<?xml version='1.0' encoding='UTF-8'?>"
//	xml_con += `<soap:Envelope xmlns:soap="http://www.w3.org/2003/05/soap-envelope/">`
//	xml_con += fmt.Sprintf(`<soap:Body xmlns:ns1="%s">`, body_ns)
//	xml_con += fmt.Sprintf("<ns1:%s>", funcs.Name)
//	for i, v := range funcs.Args {
//		xml_con += fmt.Sprintf("<ns1:%s><![CDATA[%s]]></ns1:%s>", v.Name, args[i], v.Name)
//	}
//	xml_con += fmt.Sprintf("</ns1:%s>", funcs.Name)
//	xml_con += `</soap:Body>`
//	xml_con += `</soap:Envelope>`
//	return_str, err := s.httpPost(s.requestUrl, xml_con)
//
//	if err != nil {
//		return "", err
//	}
//	fmt.Println(return_str)
//	reg := regexp.MustCompile(fmt.Sprintf("(?si)<%s>(.+?)</%s>", funcs.ReturnArgs[0].Name, funcs.ReturnArgs[0].Name))
//	list := reg.FindStringSubmatch(return_str)
//	if len(list) <= 0 {
//		return "", nil
//	}
//
//	return replaceXml(list[1], false), nil
//}

func (s *SoapClient) Call(func_name string, args ...interface{}) (string, error) {
	fun := s.ws.GetFunction(s.ws.ws.Service.Port[s.soapService].Binding, func_name)
	if fun == nil {
		return "", errors.New("Not found function name")
	}
	var xml_str string

	if len(args) > 1 {
		xml_str = s.buildSoapXML(fun, args)
	} else {
		xml_str = s.buildSoapXML(fun, args[0])
	}

	s.soapAction = fun.Action
	res, err := s.httpPost(s.requestUrl, xml_str)
	if err != nil {
		return "", err
	}

	res_msg := s.explainResponse(res, fun)

	return res_msg, nil
}

//重新设置调用地址
func (s *SoapClient) SetAddress(addr string) {
	s.requestUrl = addr
}

/*	设置SOAP配置项
 *  配置项名称           说明                     示例值
 *  SOAP-VERSION        [SOAP 构建版本]          1.1|1.2
 *  SOAP-DEF-SERVICE    [SOAP 默认使用SERVICE]   0
 */
func (s *SoapClient) SetOptions(key, val string) {
	s.options[key] = val
}

//发起HTTP请求
func (s *SoapClient) httpPost(url_str, content string) (string, error) {
	if s.soapVersion == "1.2" {
		s.httpClient.SetHeader("Content-Type", "application/soap+xml; charset=utf-8")
	} else {
		s.httpClient.SetHeader("Content-Type", "text/xml; charset=utf-8")
		s.httpClient.SetHeader("SOAPAction", s.soapAction)
	}

	req, err := s.httpClient.Request("POST", url_str, strings.NewReader(content))
	if err != nil {
		return "", err
	}

	return string(req.Content), nil
}

//得到远程WSDL文件
func (s *SoapClient) getRemoteWsdl() ([]byte, error) {
	wsdl_str, err := s.httpClient.Get(s.soapUrl)
	if err != nil {
		return nil, err
	}
	return wsdl_str, nil
}

//解释WSDL文件
func (s *SoapClient) explain(content []byte) error {
	err := s.ws.Explain(content)
	if err != nil {
		return err
	}

	s.requestUrl = s.ws.GetAddress(0)
	return nil
}

//构建SOAP XML 文件
func (s *SoapClient) buildSoapXML(fun *WsdlFunction, params interface{}) string {
	ns := map[string]string{}
	xml_con := "<?xml version='1.0' encoding='UTF-8'?>"
	xml_body := s.buildSoapBody(ArgsMap{fun.RequestArgs}, params, ns)
	body_ns := ""
	for k, v := range ns {
		body_ns += fmt.Sprintf(` xmlns:%s="%s"`, v, k)
	}
	if s.soapVersion == "1.1" {
		xml_con += `<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">`
	} else {
		xml_con += `<soap:Envelope xmlns:soap="http://www.w3.org/2003/05/soap-envelope/">`
	}
	xml_con += fmt.Sprintf(`<soap:Body %s>`, body_ns)
	xml_con += xml_body
	xml_con += `</soap:Body>`
	xml_con += `</soap:Envelope>`

	return xml_con
}

//构建SOAP BODY XML 文件
func (s *SoapClient) buildSoapBody(elms ArgsMap, params interface{}, ns map[string]string) string {
	var xml_con string

	for idx, v := range elms {
		ns_el, ok := ns[v.Namespace]
		if !ok {
			ns_el = "ck" + utils.RandStr(3, "123456789")
			ns[v.Namespace] = ns_el
		}

		xml_con += fmt.Sprintf("<%s:%s>", ns_el, v.Name)
		if len(v.Elements) > 0 {
			switch params.(type) {
			case []interface{}:
				xml_con += s.buildSoapBody(v.Elements,
					params,
					ns)
			case utils.M:
				child_params, ok := params.(utils.M)[v.Name]
				xml_con += s.buildSoapBody(v.Elements, utils.YN(ok, child_params, params), ns)
			case map[string]interface{}:
				child_params, ok := params.(map[string]interface{})[v.Name]
				xml_con += s.buildSoapBody(v.Elements, utils.YN(ok, child_params, params), ns)
			}
		} else {
			var val interface{}
			switch params.(type) {
			case []interface{}:
				val = params.([]interface{})[idx]
			case utils.M:
				val = params.(utils.M)[v.Name]
			case map[string]interface{}:
				val = params.(map[string]interface{})[v.Name]
			}
			if v.Type == "string" {
				xml_con += fmt.Sprintf("<![CDATA[%v]]>", utils.YN(val == nil, "", val))
			} else {
				xml_con += fmt.Sprintf("%v", val)
			}
		}
		xml_con += fmt.Sprintf("</%s:%s>", ns_el, v.Name)
	}

	return xml_con
}

//解释接口返回数据
func (s *SoapClient) explainResponse(res string, fun *WsdlFunction) string {
	//处理SOAP错误返回
	reg_fault := regexp.MustCompile(`(?si)<(\w+:)?Fault([^>]+)?>(.+?)</(\w+:)?Fault>`)
	if reg_fault.MatchString(res) {
		sub := reg_fault.FindString(res)
		fault := &SoapMessageFault{}
		err := xml.Unmarshal([]byte(sub), fault)
		if err != nil {
			return err.Error()
		}
		return fault.Faultstring
	}
	if s.returnAll {
		reg := regexp.MustCompile(fmt.Sprintf(`(?si)<(\w+:)?%s([^>]+)?>(.+?)</(\w+:)?%s>`, fun.ResponseArgs.Name, fun.ResponseArgs.Name))
		res_sub := reg.FindString(res)
		return res_sub
	} else {
		reg := regexp.MustCompile(fmt.Sprintf(`(?si)<(\w+:)?%s([^>]+)?>(.+?)</(\w+:)?%s>`, fun.ResponseArgs.Elements[0].Name, fun.ResponseArgs.Elements[0].Name))
		res_sub := reg.FindString(res)
		return reg.ReplaceAllString(res_sub, "$3")
	}
}

//替换XML特殊字符
func (s *SoapClient) ReplaceXml(str string, toxml bool) string {
	replace_list := map[string]string{
		"<":  "&lt;",
		">":  "&gt;",
		"\"": "&quot;",
		"'":  "&apos;",
		"&":  "&amp;",
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
