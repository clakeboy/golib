package soap

import (
	"encoding/xml"
	"github.com/clakeboy/golib/utils"
	"strings"
)

type WsdlElement struct {
	XMLName         xml.Name        `xml:"http://schemas.xmlsoap.org/wsdl/ definitions"`
	Xsd             xml.Attr        `xml:"xsd,attr"`
	Wsdl            xml.Attr        `xml:"wsdl,attr"`
	Tns             xml.Attr        `xml:"tns,attr"`
	Soap            xml.Attr        `xml:"soap,attr"`
	Nsl             xml.Attr        `xml:"nsl,attr"`
	Name            string          `xml:"name,attr"`
	TargetNamespace string          `xml:"targetNamespace,attr"`
	Types           *WsdlTypes      `xml:"types"`
	Message         []*WsdlMessage  `xml:"message"`
	PortType        []*WsdlPortType `xml:"portType"`
	Binding         []*WsdlBinding  `xml:"binding"`
	Service         *WsdlService    `xml:"service"`
}

type WsdlTypes struct {
	XMLName xml.Name      `xml:"types"`
	Schema  []*WsdlSchema `xml:"schema"`
}

type WsdlSchema struct {
	XMLName              xml.Name `xml:"schema"`
	AttributeFormDefault string   `xml:"attributeFormDefault,attr"`
	ElementFormDefault   string   `xml:"elementFormDefault,attr"`
	TargetNamespace      string   `xml:"targetNamespace,attr"`
	Import               []struct {
		XMLName   xml.Name `xml:"import"`
		Namespace string   `xml:"namespace,attr"`
	} `xml:"import"`
	ComplexType []*WsdlComplexType   `xml:"complexType"`
	Element     []*WsdlSchemaElement `xml:"element"`
}

type WsdlComplexType struct {
	XMLName  xml.Name `xml:"complexType"`
	Name     string   `xml:"name,attr"`
	Sequence struct {
		Element []*WsdlTypeElement `xml:"element"`
	} `xml:"sequence"`
	ComplexContent  *WsdlComplexContent `xml:"complexContent"`
	TargetNamespace string
}

type WsdlComplexContent struct {
	XMLName   xml.Name `xml:"complexContent"`
	Extension struct {
		Base     string `xml:"base,attr"`
		Sequence struct {
			Element []*WsdlTypeElement `xml:"element"`
		} `xml:"sequence"`
	} `xml:"extension"`
}

type WsdlTypeElement struct {
	XMLName   xml.Name `xml:"element"`
	Name      string   `xml:"name,attr"`
	MinOccurs int      `xml:"minOccurs,attr"`
	MaxOccurs string   `xml:"maxOccurs,attr"`
	Nillable  bool     `xml:"nillable,attr"`
	Type      string   `xml:"type,attr"`
}

type WsdlSchemaElement struct {
	XMLName         xml.Name         `xml:"element"`
	Name            string           `xml:"name,attr"`
	ComplexType     *WsdlComplexType `xml:"complexType"`
	TargetNamespace string
}

//message
type WsdlMessage struct {
	XMLName xml.Name    `xml:"message"`
	Name    string      `xml:"name,attr"`
	Parts   []*WsdlPart `xml:"part"`
}

type WsdlPart struct {
	XMLName xml.Name `xml:"part"`
	Name    string   `xml:"name,attr"`
	Type    string   `xml:"type,attr"`
	Element string   `xml:"element,attr"`
}

type WsdlPortType struct {
	XMLName      xml.Name         `xml:"portType"`
	Name         string           `xml:"name,attr"`
	Operation    []*WsdlOperation `xml:"operation"`
	OperationMap map[string]*WsdlOperation
}

type WsdlOperation struct {
	XMLName       xml.Name       `xml:"operation"`
	Name          string         `xml:"name,attr"`
	Message       string         `xml:"message,attr"`
	Action        string         `xml:"Action,attr"`
	Input         *WsdlInput     `xml:"input"`
	Output        *WsdlOutput    `xml:"output"`
	SoapOperation *SoapOperation `xml:"operation"`
}

type WsdlInput struct {
	XMLName xml.Name  `xml:"input"`
	Message string    `xml:"message,attr"`
	Name    string    `xml:"name,attr"`
	Action  string    `xml:"Action,attr"`
	Body    *SoapBody `xml:"body"`
}

type WsdlOutput struct {
	XMLName xml.Name `xml:"output"`
	Message string   `xml:"message,attr"`
	Name    string   `xml:"name,attr"`
}

type WsdlBinding struct {
	XMLName      xml.Name         `xml:"binding"`
	Name         string           `xml:"name,attr"`
	Type         string           `xml:"type,attr"`
	Binding      *SoapBinding     `xml:"binding"`
	Operation    []*WsdlOperation `xml:"operation"`
	OperationMap map[string]*WsdlOperation
}

type WsdlService struct {
	XMLName xml.Name `xml:"service"`
	Name    string   `xml:"name,attr"`
	Port    []struct {
		XMLName xml.Name `xml:"port"`
		Binding string   `xml:"binding,attr"`
		Name    string   `xml:"name,attr"`
		Address struct {
			XMLName  xml.Name `xml:"address"`
			Location string   `xml:"location,attr"`
		} `xml:"address"`
	} `xml:"port"`
}

type SoapBinding struct {
	XMLName   xml.Name `xml:"binding"`
	Style     string   `xml:"style,attr"`
	Transport string   `xml:"transport,attr"`
}

type SoapOperation struct {
	XMLName    xml.Name `xml:"operation"`
	SoapAction string   `xml:"soapAction,attr"`
	Style      string   `xml:"style,attr"`
}

type SoapBody struct {
	XMLName   xml.Name `xml:"body"`
	Namespace string   `xml:"namespace,attr"`
	Use       string   `xml:"use,attr"`
}

type Wsdl struct {
	ws           *WsdlElement
	messages     map[string]*WsdlMessage
	portTypes    map[string]*WsdlPortType
	buildings    map[string]*WsdlBinding
	elements     map[string]*WsdlSchemaElement
	complexTypes map[string]*WsdlComplexType
}

//新建一个WSDL解释器
func NewWsdl() *Wsdl {
	return &Wsdl{
		messages:     make(map[string]*WsdlMessage),
		portTypes:    make(map[string]*WsdlPortType),
		buildings:    make(map[string]*WsdlBinding),
		elements:     make(map[string]*WsdlSchemaElement),
		complexTypes: make(map[string]*WsdlComplexType),
	}
}

//解释WSDL文件得到解构
func (w *Wsdl) Explain(body []byte) error {
	wsdl := &WsdlElement{}
	err := xml.Unmarshal(body, wsdl)
	if err != nil {
		return err
	}

	w.ws = wsdl
	w.explainPortType()
	w.explainMessage()
	w.explainBuilding()
	w.explainComplexType()
	return nil
}

//解释 port type
func (w *Wsdl) explainPortType() {
	for _, v := range w.ws.PortType {
		v.OperationMap = make(map[string]*WsdlOperation)
		for _, n := range v.Operation {
			v.OperationMap[n.Name] = n
		}
		w.portTypes[v.Name] = v
	}
}

//解释 message
func (w *Wsdl) explainMessage() {
	for _, v := range w.ws.Message {
		w.messages[v.Name] = v
	}
}

//解释 building
func (w *Wsdl) explainBuilding() {
	for _, v := range w.ws.Binding {
		v.OperationMap = make(map[string]*WsdlOperation)
		for _, n := range v.Operation {
			v.OperationMap[n.Name] = n
		}
		w.buildings[v.Name] = v
	}
}

//解释 ComplexType,Element
func (w *Wsdl) explainComplexType() {
	if w.ws.Types == nil {
		return
	}
	for _, v := range w.ws.Types.Schema {
		namespace := utils.YN(v.TargetNamespace == "", w.ws.TargetNamespace, v.TargetNamespace).(string)

		for _, t := range v.ComplexType {
			t.TargetNamespace = namespace
			w.complexTypes[t.Name] = t
		}

		for _, e := range v.Element {
			e.TargetNamespace = namespace
			e.ComplexType.TargetNamespace = namespace
			w.elements[e.Name] = e
		}
	}
}

//得到一个方法
func (w *Wsdl) GetFunction(build_name, func_name string) *WsdlFunction {
	build_name = formatPrefixNs(build_name)
	wsdl_fun := &WsdlFunction{}
	build, ok := w.buildings[build_name]
	if !ok {
		return nil
	}

	ports, ok := w.portTypes[formatPrefixNs(build.Type)]
	if !ok {
		return nil
	}

	fun, ok := ports.OperationMap[func_name]
	if !ok {
		return nil
	}

	wsdl_fun.Action = fun.Input.Action
	reqmsg, ok := w.messages[formatPrefixNs(fun.Input.Message)]
	if !ok {
		return nil
	}

	resmsg, ok := w.messages[formatPrefixNs(fun.Output.Message)]
	if !ok {
		return nil
	}

	wsdl_fun.Name = func_name
	wsdl_fun.RequestArgs = &WsdlFunctionArgs{
		Name:      fun.Name,
		Namespace: w.ws.TargetNamespace,
	}
	for _, part := range reqmsg.Parts {
		if part.Name == "parameters" && part.Element != "" {
			elm, ok := w.elements[formatPrefixNs(part.Element)]
			if !ok {
				continue
			}
			wsdl_fun.RequestArgs.Elements = w.getWsdlArgs(elm.ComplexType)
		} else {
			args := &WsdlFunctionArgs{}
			args.Name = part.Name
			args.Type = formatPrefixNs(part.Type)
			args.Namespace = w.ws.TargetNamespace
			wsdl_fun.RequestArgs.Elements = append(wsdl_fun.RequestArgs.Elements, args)
		}
	}
	wsdl_fun.ResponseName = resmsg.Name
	wsdl_fun.ResponseArgs = &WsdlFunctionArgs{
		Name:      resmsg.Name,
		Namespace: w.ws.TargetNamespace,
	}
	for _, part := range resmsg.Parts {
		if part.Name == "parameters" && part.Element != "" {
			elm, ok := w.elements[formatPrefixNs(part.Element)]
			if !ok {
				continue
			}
			wsdl_fun.ResponseArgs.Elements = w.getWsdlArgs(elm.ComplexType)
		} else {
			args := &WsdlFunctionArgs{}
			args.Name = part.Name
			wsdl_fun.ResponseArgs.Elements = append(wsdl_fun.ResponseArgs.Elements, args)
		}
	}

	return wsdl_fun
}

//得到参数和子参数递归
func (w *Wsdl) getWsdlArgs(ct *WsdlComplexType) ArgsMap {
	args_map := ArgsMap{}
	var seq []*WsdlTypeElement
	if ct.ComplexContent != nil {
		args_map = append(args_map, w.getWsdlArgs(w.complexTypes[formatPrefixNs(ct.ComplexContent.Extension.Base)])...)

		seq = ct.ComplexContent.Extension.Sequence.Element
	} else {
		seq = ct.Sequence.Element
	}

	for _, v := range seq {
		args := &WsdlFunctionArgs{}
		args.Name = v.Name
		args.Type = formatPrefixNs(v.Type)
		args.Namespace = ct.TargetNamespace
		args.MinOccurs = v.MinOccurs
		args.MaxOccurs = v.MaxOccurs
		args.Nillable = v.Nillable

		child, ok := w.complexTypes[args.Type]
		if ok {
			args.Elements = w.getWsdlArgs(child)
		}
		args_map = append(args_map, args)
	}
	return args_map
}

func formatPrefixNs(str string) string {
	strarr := strings.Split(str, ":")
	if len(strarr) > 1 {
		return strarr[1]
	}
	return strarr[0]
}

//得到WSDL服务调用地址
func (w *Wsdl) GetAddress(idx int) string {
	return w.ws.Service.Port[idx].Address.Location
}
