package soap

import (
	"testing"
	"encoding/xml"
	"fmt"
)

func TestWsdl_Explain(t *testing.T) {
	str := `<?xml version='1.0' encoding='UTF-8'?><wsdl:definitions xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:wsdl="http://schemas.xmlsoap.org/wsdl/" xmlns:tns="http://access.vchl.echannel.acic.com/" xmlns:soap="http://schemas.xmlsoap.org/wsdl/soap/" xmlns:ns1="http://schemas.xmlsoap.org/soap/http" name="IAccessServiceImplService" targetNamespace="http://access.vchl.echannel.acic.com/">
  <wsdl:message name="vchlRequestResponse">
    <wsdl:part name="return" type="xsd:string">
    </wsdl:part>
  </wsdl:message>
  <wsdl:message name="vchlRequest">
    <wsdl:part name="arg0" type="xsd:string">
    </wsdl:part>
    <wsdl:part name="arg1" type="xsd:string">
    </wsdl:part>
    <wsdl:part name="arg2" type="xsd:string">
    </wsdl:part>
  </wsdl:message>
  <wsdl:portType name="IAccessServiceImpl">
    <wsdl:operation name="vchlRequest">
      <wsdl:input message="tns:vchlRequest" name="vchlRequest">
    </wsdl:input>
      <wsdl:output message="tns:vchlRequestResponse" name="vchlRequestResponse">
    </wsdl:output>
    </wsdl:operation>
  </wsdl:portType>
  <wsdl:binding name="IAccessServiceImplServiceSoapBinding" type="tns:IAccessServiceImpl">
    <soap:binding style="rpc" transport="http://schemas.xmlsoap.org/soap/http"/>
    <wsdl:operation name="vchlRequest">
      <soap:operation soapAction="" style="rpc"/>
      <wsdl:input name="vchlRequest">
        <soap:body namespace="http://access.vchl.echannel.acic.com/" use="literal"/>
      </wsdl:input>
      <wsdl:output name="vchlRequestResponse">
        <soap:body namespace="http://access.vchl.echannel.acic.com/" use="literal"/>
      </wsdl:output>
    </wsdl:operation>
  </wsdl:binding>
  <wsdl:service name="IAccessServiceImplService">
    <wsdl:port binding="tns:IAccessServiceImplServiceSoapBinding" name="IAccessServiceImplPort">
      <soap:address location="http://123.147.190.130/vchl-channel/services/accessService"/>
    </wsdl:port>
  </wsdl:service>
</wsdl:definitions>`
	wsdl := NewWsdl()
	wsdl.Explain([]byte(str))

}

func TestWsdl_GetFunc(t *testing.T) {
	var Envelope struct{
		XMLName xml.Name
		Soap    xml.Attr  `xml:"soap,attr"`
		Body struct{
			XMLName xml.Name `xml:"body"`
		} `xml:"body"`
	}
	xml_str := `
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/"><soap:Body><soap:Fault><faultcode>soap:Server</faultcode><faultstring>No binding operation info while invoking unknown method with params unknown.</faultstring></soap:Fault></soap:Body></soap:Envelope>`
	xml.Unmarshal([]byte(xml_str),&Envelope)
	fmt.Printf("%+v",Envelope.Soap)
	xml_con,_ := xml.Marshal(&Envelope)
	fmt.Println(string(xml_con))
}
