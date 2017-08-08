package soap

import (
	"encoding/xml"
)

type WsdlElement struct {
	XMLName  xml.Name `xml:"http://schemas.xmlsoap.org/wsdl/ definitions"`
	Xsd      xml.Attr   `xml:"xsd,attr"`
	Wsdl     xml.Attr    `xml:"wsdl,attr"`
	Tns      xml.Attr    `xml:"tns,attr"`
	Soap     xml.Attr    `xml:"soap,attr"`
	Nsl      xml.Attr     `xml:"nsl,attr"`
	Name     string     `xml:"name,attr"`
	TargetNamespace  string  `xml:"targetNamespace,attr"`
	Message  []*WsdlMessage  `xml:"message"`
	PortType  *WsdlPortType  `xml:"portType"`
	Binding   *WsdlBinding    `xml:"binding"`
	Service   *WsdlService    `xml:"service"`
}

type WsdlMessage struct {
	XMLName  xml.Name `xml:"message"`
	Name     string  `xml:"name,attr"`
	Parts    []*WsdlPart `xml:"part"`
}

type WsdlPart struct {
	XMLName  xml.Name `xml:"part"`
	Name    string  `xml:"name,attr"`
	Type    string   `xml:"type,attr"`
}

type WsdlPortType struct {
	XMLName xml.Name `xml:"portType"`
	Name string  `xml:"name,attr"`
	Operation []*WsdlOperation `xml:"operation"`
}

type WsdlOperation struct {
	XMLName xml.Name `xml:"operation"`
	Name  string   `xml:"name,attr"`
	Input *WsdlInput  `xml:"input"`
	Output *WsdlOutput  `xml:"output"`
	SoapOperation *SoapOperation `xml:"operation"`
}

type WsdlInput  struct {
	XMLName xml.Name  `xml:"input"`
	Message string   `xml:"message,attr"`
	Name    string    `xml:"name,attr"`
	Body    *SoapBody  `xml:"body"`
}

type WsdlOutput  struct {
	XMLName xml.Name  `xml:"output"`
	Message string   `xml:"message,attr"`
	Name    string    `xml:"name,attr"`
}

type WsdlBinding struct {
	XMLName xml.Name  `xml:"binding"`
	Name    string    `xml:"name,attr"`
	Type    string   `xml:"type,attr"`
	Binding *SoapBinding `xml:"binding"`
	Operation *WsdlOperation `xml:"operation"`
}

type WsdlService struct {
	XMLName xml.Name  `xml:"service"`
	Name    string    `xml:"name,attr"`
	Port    struct{
		XMLName xml.Name  `xml:"port"`
		Binding string `xml:"binding,attr"`
		Name    string    `xml:"name,attr"`
		Address struct{
			XMLName xml.Name  `xml:"address"`
			Location string `xml:"location,attr"`
		} `xml:"address"`
	} `xml:"port"`
}

type SoapBinding struct {
	XMLName xml.Name  `xml:"binding"`
	Style string `xml:"style,attr"`
	Transport string  `xml:"transport,attr"`
}

type SoapOperation struct {
	XMLName xml.Name  `xml:"operation"`
	SoapAction string `xml:"soapAction,attr"`
	Style string `xml:"style,attr"`
}

type SoapBody struct {
	XMLName xml.Name  `xml:"body"`
	Namespace string `xml:"namespace,attr"`
	Use string `xml:"use,attr"`
}

type Wsdl struct {
	ws *WsdlElement
}

//新建一个WSDL解释器
func NewWsdl() *Wsdl {
	return &Wsdl{}
}

//解释WSDL文件得到解构
func (w *Wsdl) Explain(body []byte) error {
	wsdl := &WsdlElement{}
	err := xml.Unmarshal(body,wsdl)
	if err != nil {
		return err
	}

	w.ws = wsdl
	return nil
}

//得到WSDL服务调用地址
func (w *Wsdl) GetAddress() string {
	return w.ws.Service.Port.Address.Location
}

type SoapFunc struct {
	Name string
	Args []*SoapFuncArgs
	Namespace string
	ReturnName string
	ReturnArgs []*SoapFuncArgs
}

type SoapFuncArgs struct {
	Name string
	Type string
}

//得到可以调用的方法
func (w *Wsdl) GetFunc(name string) *SoapFunc {
	var funcs *SoapFunc
	for _,v := range w.ws.PortType.Operation {
		if v.Name == name {
			funcs = &SoapFunc{
				Name:v.Name,
				Args:w.GetArgs(name),
				Namespace:w.GetNamespace(name),
				ReturnName:v.Output.Name,
				ReturnArgs:w.GetArgs(v.Output.Name),
			}
		}
	}

	return funcs
}
//得到方法参数集
func (w *Wsdl) GetArgs(name string) []*SoapFuncArgs {
	var args []*SoapFuncArgs
	for _,v := range w.ws.Message {
		if v.Name == name {
			for _,a := range v.Parts {
				args = append(args,&SoapFuncArgs{
					Name:a.Name,
					Type:a.Type,
				})
			}
			break
		}
	}
	return args
}

func (w *Wsdl) GetNamespace(name string) string {
	if w.ws.Binding.Operation.Name == name {
		return w.ws.Binding.Operation.Input.Body.Namespace
	}
	return ""
}


