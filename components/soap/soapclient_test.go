package soap

import (
	"fmt"
	"github.com/clakeboy/golib/utils"
	"log"
	"regexp"
	"strings"
	"testing"
)

func TestNewSoapClient(t *testing.T) {
	soap_client, err := NewSoapClient("http://123.147.190.130:28081/vchl-channel/services/accessService?wsdl")
	if err != nil {
		panic(err)
	}

	//no := utils.RandStr(10,nil)

	//xml_con := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?><Packet><Head><EChannelID>CHPC</EChannelID><RequestType>Q00</RequestType><SerialNo>%s</SerialNo></Head><Body><LicensePlateNo>渝B358U7</LicensePlateNo><VIN>LVSHJCAC3FE215302</VIN><AreaCode>50</AreaCode></Body></Packet>`,no)
	xml_con := `
<Packet>
  <Head>
    <EChannelID>CHPC</EChannelID>
    <RequestType>Q00</RequestType>
    <SerialNo>PC1502857381266963</SerialNo>
  </Head>
  <Body>
    <LicensePlateNo>渝B507R8</LicensePlateNo>
    <VIN>LGBH12E04BY194260</VIN>
    <AreaCode>50</AreaCode>
  </Body>
</Packet>`
	args := []interface{}{
		"CHPC",
		xml_con,
		//"4288cecb49821b5a75e8a1d09b3fdda5",
		utils.EncodeMD5Std(fmt.Sprintf("%v%v%v%v", "113.204.136.118", "CHPC", "e435rfe3dwxd180e5ea7e5f145c4ccb8", xml_con)),
	}

	fmt.Println(args...)
	soap_client.SetAddress("http://123.147.190.130:28081/vchl-channel/services/accessService")
	str, err := soap_client.Call("vchlRequest", args...)
	fmt.Println(str, err)
}

func TestSoapClient_Call(t *testing.T) {
	str := `<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/"><soap:Body><ns1:vchlRequestResponse xmlns:ns1="http://access.vchl.echannel.acic.com/"><return>&lt;?xml version=&quot;1.0&quot; encoding=&quot;UTF-8&quot;?&gt;&lt;Packet&gt;  &lt;Head&gt;    &lt;RequestType&gt;Q00&lt;/RequestType&gt;    &lt;SerialNo&gt;PC1502857381266963&lt;/SerialNo&gt;    &lt;ResponseCode&gt;0&lt;/ResponseCode&gt;    &lt;ErrorCode&gt;0000&lt;/ErrorCode&gt;    &lt;ErrorMessage&gt;成功&lt;/ErrorMessage&gt;  &lt;/Head&gt;&lt;Body&gt;&lt;EngineNo&gt;526649E&lt;/EngineNo&gt;&lt;EngineModel&gt;HR16&lt;/EngineModel&gt;&lt;VehicleModel&gt;&lt;Model&gt;RCABSD0039&lt;/Model&gt;&lt;CarName&gt;轩逸DFL7162ACC轿车&lt;/CarName&gt;&lt;BrandName&gt;东风日产&lt;/BrandName&gt;&lt;FamilyName&gt;轩逸&lt;/FamilyName&gt;&lt;DisplaceMent&gt;1598.0&lt;/DisplaceMent&gt;&lt;Power&gt;86.0&lt;/Power&gt;&lt;Whole_Weight&gt;1220&lt;/Whole_Weight&gt;&lt;RatedPassengerCapacity&gt;5&lt;/RatedPassengerCapacity&gt;&lt;Tonnage&gt;0.0&lt;/Tonnage&gt;&lt;Haulage&gt;0.0&lt;/Haulage&gt;&lt;CSeat&gt;5&lt;/CSeat&gt;&lt;ReplacementValue&gt;116800&lt;/ReplacementValue&gt;&lt;IpmORLoc&gt;JV&lt;/IpmORLoc&gt;&lt;PlateVhlCode&gt;BDFJXZUC0020&lt;/PlateVhlCode&gt;&lt;HfCode&gt;0&lt;/HfCode&gt;&lt;HfName&gt;正常&lt;/HfName&gt;&lt;GroupName&gt;轩逸G11(06/08-)&lt;/GroupName&gt;&lt;GroupCode&gt;&lt;/GroupCode&gt;&lt;NewVhlClassName&gt;紧凑型轿车(A)&lt;/NewVhlClassName&gt;&lt;NewVhlClassCode&gt;&lt;/NewVhlClassCode&gt;&lt;FuelType&gt;0&lt;/FuelType&gt;&lt;FuelTypeCode&gt;D1&lt;/FuelTypeCode&gt;&lt;/VehicleModel&gt;&lt;VehicleModel&gt;&lt;Model&gt;RCABSD0040&lt;/Model&gt;&lt;CarName&gt;轩逸DFL7162ACC轿车&lt;/CarName&gt;&lt;BrandName&gt;东风日产&lt;/BrandName&gt;&lt;FamilyName&gt;轩逸&lt;/FamilyName&gt;&lt;DisplaceMent&gt;1598.0&lt;/DisplaceMent&gt;&lt;Power&gt;86.0&lt;/Power&gt;&lt;Whole_Weight&gt;1220&lt;/Whole_Weight&gt;&lt;RatedPassengerCapacity&gt;5&lt;/RatedPassengerCapacity&gt;&lt;Tonnage&gt;0.0&lt;/Tonnage&gt;&lt;Haulage&gt;0.0&lt;/Haulage&gt;&lt;CSeat&gt;5&lt;/CSeat&gt;&lt;ReplacementValue&gt;127800&lt;/ReplacementValue&gt;&lt;IpmORLoc&gt;JV&lt;/IpmORLoc&gt;&lt;PlateVhlCode&gt;BDFJXZUD0025&lt;/PlateVhlCode&gt;&lt;HfCode&gt;0&lt;/HfCode&gt;&lt;HfName&gt;正常&lt;/HfName&gt;&lt;GroupName&gt;轩逸G11(06/08-)&lt;/GroupName&gt;&lt;GroupCode&gt;&lt;/GroupCode&gt;&lt;NewVhlClassName&gt;紧凑型轿车(A)&lt;/NewVhlClassName&gt;&lt;NewVhlClassCode&gt;&lt;/NewVhlClassCode&gt;&lt;FuelType&gt;0&lt;/FuelType&gt;&lt;FuelTypeCode&gt;D1&lt;/
	FuelTypeCode&gt;&lt;/
	VehicleModel&gt;&lt;/Body&gt;&lt;/Packet&gt;</return></ns1:vchlRequestResponse></soap:Body></soap:Envelope>`
	name := "return"
	reg := regexp.MustCompile(fmt.Sprintf("(?si)<%s>(.+?)</%s>", name, name))
	matchs := reg.FindStringSubmatch(str)
	fmt.Println(matchs)
	//fmt.Println(replaceXml(matchs[1],false))
}

func TestNewSoapClient2(t *testing.T) {
	wsdl_xml := `
<?xml version="1.0" encoding="UTF-8"?>
<wsdl:definitions xmlns:wsdl="http://schemas.xmlsoap.org/wsdl/" xmlns:ns="http://service.front.sinosoft.com" xmlns:xs="http://www.w3.org/2001/XMLSchema" xmlns:ns1="http://org.apache.axis2/xsd" xmlns:ax273="http://busInterface.front.sinosoft.com/xsd" xmlns:wsaw="http://www.w3.org/2006/05/addressing/wsdl" xmlns:http="http://schemas.xmlsoap.org/wsdl/http/" xmlns:ax276="http://common.dto.front.sinosoft.com/xsd" xmlns:ax275="http://dto.front.sinosoft.com/xsd" xmlns:soap="http://schemas.xmlsoap.org/wsdl/soap/" xmlns:mime="http://schemas.xmlsoap.org/wsdl/mime/" xmlns:ax279="http://out.CarQueryService.dto.front.sinosoft.com/xsd" xmlns:soap12="http://schemas.xmlsoap.org/wsdl/soap12/" targetNamespace="http://service.front.sinosoft.com">
    <wsdl:documentation>CarQueryService</wsdl:documentation>
    <wsdl:types>
        <xs:schema xmlns:ax277="http://common.dto.front.sinosoft.com/xsd" xmlns:ax280="http://out.CarQueryService.dto.front.sinosoft.com/xsd" attributeFormDefault="qualified" elementFormDefault="qualified" targetNamespace="http://dto.front.sinosoft.com/xsd">
            <xs:import namespace="http://common.dto.front.sinosoft.com/xsd"/>
            <xs:import namespace="http://out.CarQueryService.dto.front.sinosoft.com/xsd"/>
            <xs:complexType name="CarQueryRequest">
                <xs:sequence>
                    <xs:element minOccurs="0" name="checkCode" nillable="true" type="xs:string"/>
                    <xs:element minOccurs="0" name="checkNo" nillable="true" type="xs:string"/>
                    <xs:element minOccurs="0" name="cityAreaCode" nillable="true" type="xs:string"/>
                    <xs:element minOccurs="0" name="comCode" nillable="true" type="xs:string"/>
                    <xs:element minOccurs="0" name="engineNo" nillable="true" type="xs:string"/>
                    <xs:element minOccurs="0" name="frameNo" nillable="true" type="xs:string"/>
                    <xs:element minOccurs="0" name="franchiserCode" nillable="true" type="xs:string"/>
                    <xs:element minOccurs="0" name="handler1Code" nillable="true" type="xs:string"/>
                    <xs:element minOccurs="0" name="licenseNo" nillable="true" type="xs:string"/>
                    <xs:element minOccurs="0" name="licenseNoQueryFlag" nillable="true" type="xs:string"/>
                    <xs:element minOccurs="0" name="licenseType" nillable="true" type="xs:string"/>
                    <xs:element minOccurs="0" name="policySort" nillable="true" type="xs:string"/>
                    <xs:element minOccurs="0" name="txInsuranceRequestEhm" nillable="true" type="ax277:TxInsuranceRequestEhm"/>
                    <xs:element minOccurs="0" name="txInsuranceRequestExtensionEhm" nillable="true" type="ax277:TxInsuranceRequestExtensionEhm"/>
                </xs:sequence>
            </xs:complexType>
            <xs:complexType name="CarQueryResponse">
                <xs:sequence>
                    <xs:element maxOccurs="unbounded" minOccurs="0" name="carQueryInfoArr" nillable="true" type="ax280:CarQueryInfo"/>
                    <xs:element minOccurs="0" name="checkCode" nillable="true" type="xs:string"/>
                    <xs:element minOccurs="0" name="checkNo" nillable="true" type="xs:string"/>
                    <xs:element minOccurs="0" name="policySort" nillable="true" type="xs:string"/>
                    <xs:element minOccurs="0" name="txInsuranceResponseEhm" nillable="true" type="ax277:TxInsuranceResponseEhm"/>
                    <xs:element minOccurs="0" name="txInsuranceResponseExtensionEhm" nillable="true" type="ax277:TxInsuranceResponseExtensionEhm"/>
                </xs:sequence>
            </xs:complexType>
        </xs:schema>
        <xs:schema attributeFormDefault="qualified" elementFormDefault="qualified" targetNamespace="http://busInterface.front.sinosoft.com/xsd">
            <xs:complexType name="CarQueryTpInterface">
                <xs:sequence/>
            </xs:complexType>
            <xs:complexType name="CarQueryTpCheckInterface">
                <xs:sequence/>
            </xs:complexType>
        </xs:schema>
        <xs:schema xmlns:ax274="http://busInterface.front.sinosoft.com/xsd" xmlns:ax278="http://dto.front.sinosoft.com/xsd" attributeFormDefault="qualified" elementFormDefault="qualified" targetNamespace="http://service.front.sinosoft.com">
            <xs:import namespace="http://busInterface.front.sinosoft.com/xsd"/>
            <xs:import namespace="http://dto.front.sinosoft.com/xsd"/>
            <xs:element name="setCarQueryTpInterface">
                <xs:complexType>
                    <xs:sequence>
                        <xs:element minOccurs="0" name="carQueryTpInterface" nillable="true" type="ax273:CarQueryTpInterface"/>
                    </xs:sequence>
                </xs:complexType>
            </xs:element>
            <xs:element name="setCarQueryTpCheckInterface">
                <xs:complexType>
                    <xs:sequence>
                        <xs:element minOccurs="0" name="carQueryTpCheckInterface" nillable="true" type="ax273:CarQueryTpCheckInterface"/>
                    </xs:sequence>
                </xs:complexType>
            </xs:element>
            <xs:element name="getCarQueryTpInterface">
                <xs:complexType>
                    <xs:sequence/>
                </xs:complexType>
            </xs:element>
            <xs:element name="getCarQueryTpInterfaceResponse">
                <xs:complexType>
                    <xs:sequence>
                        <xs:element minOccurs="0" name="return" nillable="true" type="ax273:CarQueryTpInterface"/>
                    </xs:sequence>
                </xs:complexType>
            </xs:element>
            <xs:element name="getCarQueryTpCheckInterface">
                <xs:complexType>
                    <xs:sequence/>
                </xs:complexType>
            </xs:element>
            <xs:element name="getCarQueryTpCheckInterfaceResponse">
                <xs:complexType>
                    <xs:sequence>
                        <xs:element minOccurs="0" name="return" nillable="true" type="ax273:CarQueryTpCheckInterface"/>
                    </xs:sequence>
                </xs:complexType>
            </xs:element>
            <xs:element name="carQuery">
                <xs:complexType>
                    <xs:sequence>
                        <xs:element minOccurs="0" name="request" nillable="true" type="ax275:CarQueryRequest"/>
                    </xs:sequence>
                </xs:complexType>
            </xs:element>
            <xs:element name="carQueryResponse">
                <xs:complexType>
                    <xs:sequence>
                        <xs:element minOccurs="0" name="return" nillable="true" type="ax275:CarQueryResponse"/>
                    </xs:sequence>
                </xs:complexType>
            </xs:element>
        </xs:schema>
        <xs:schema attributeFormDefault="qualified" elementFormDefault="qualified" targetNamespace="http://out.CarQueryService.dto.front.sinosoft.com/xsd">
            <xs:complexType name="CarQueryInfo">
                <xs:sequence>
                    <xs:element minOccurs="0" name="brandName" nillable="true" type="xs:string"/>
                    <xs:element minOccurs="0" name="carKindCode" nillable="true" type="xs:string"/>
                    <xs:element minOccurs="0" name="carOwner" nillable="true" type="xs:string"/>
                    <xs:element minOccurs="0" name="color" nillable="true" type="xs:string"/>
                    <xs:element minOccurs="0" name="displacement" nillable="true" type="xs:string"/>
                    <xs:element minOccurs="0" name="engineNo" nillable="true" type="xs:string"/>
                    <xs:element minOccurs="0" name="frameNo" nillable="true" type="xs:string"/>
                    <xs:element minOccurs="0" name="haulage" nillable="true" type="xs:string"/>
                    <xs:element minOccurs="0" name="ineffectualDate" nillable="true" type="xs:string"/>
                    <xs:element minOccurs="0" name="lastCheckDate" nillable="true" type="xs:string"/>
                    <xs:element minOccurs="0" name="licenseNo" nillable="true" type="xs:string"/>
                    <xs:element minOccurs="0" name="licenseType" nillable="true" type="xs:string"/>
                    <xs:element minOccurs="0" name="limitLoad" nillable="true" type="xs:string"/>
                    <xs:element minOccurs="0" name="limitLoadPerson" nillable="true" type="xs:string"/>
                    <xs:element minOccurs="0" name="madeDate" nillable="true" type="xs:string"/>
                    <xs:element minOccurs="0" name="madeFactory" nillable="true" type="xs:string"/>
                    <xs:element minOccurs="0" name="modelCode" nillable="true" type="xs:string"/>
                    <xs:element minOccurs="0" name="motorTypeCode" nillable="true" type="xs:string"/>
                    <xs:element minOccurs="0" name="pmFuelType" nillable="true" type="xs:string"/>
                    <xs:element minOccurs="0" name="producerType" nillable="true" type="xs:string"/>
                    <xs:element minOccurs="0" name="registDate" nillable="true" type="xs:string"/>
                    <xs:element minOccurs="0" name="rejectDate" nillable="true" type="xs:string"/>
                    <xs:element minOccurs="0" name="salePrice" nillable="true" type="xs:string"/>
                    <xs:element minOccurs="0" name="status" nillable="true" type="xs:string"/>
                    <xs:element minOccurs="0" name="transferDate" nillable="true" type="xs:string"/>
                    <xs:element minOccurs="0" name="useType" nillable="true" type="xs:string"/>
                    <xs:element minOccurs="0" name="vehicleBrand1" nillable="true" type="xs:string"/>
                    <xs:element minOccurs="0" name="vehicleBrand2" nillable="true" type="xs:string"/>
                    <xs:element minOccurs="0" name="vehicleStyle" nillable="true" type="xs:string"/>
                    <xs:element minOccurs="0" name="wholeWeight" nillable="true" type="xs:string"/>
                </xs:sequence>
            </xs:complexType>
        </xs:schema>
        <xs:schema attributeFormDefault="qualified" elementFormDefault="qualified" targetNamespace="http://common.dto.front.sinosoft.com/xsd">
            <xs:complexType name="TxInsuranceEhm">
                <xs:sequence>
                    <xs:element minOccurs="0" name="transExeDate" nillable="true" type="xs:string"/>
                    <xs:element minOccurs="0" name="transExeTime" nillable="true" type="xs:string"/>
                    <xs:element minOccurs="0" name="transRefGUID" nillable="true" type="xs:string"/>
                    <xs:element minOccurs="0" name="transSubType" nillable="true" type="xs:string"/>
                    <xs:element minOccurs="0" name="transType" nillable="true" type="xs:string"/>
                </xs:sequence>
            </xs:complexType>
            <xs:complexType name="TxInsuranceRequestEhm">
                <xs:complexContent>
                    <xs:extension base="ax276:TxInsuranceEhm">
                        <xs:sequence>
                            <xs:element minOccurs="0" name="iinsuranceExtensionEhm" nillable="true" type="ax276:IinsuranceExtensionEhm"/>
                        </xs:sequence>
                    </xs:extension>
                </xs:complexContent>
            </xs:complexType>
            <xs:complexType name="IinsuranceExtensionEhm">
                <xs:sequence>
                    <xs:element minOccurs="0" name="maxRecords" nillable="true" type="xs:string"/>
                    <xs:element minOccurs="0" name="orderField" nillable="true" type="xs:string"/>
                    <xs:element minOccurs="0" name="orderFlag" nillable="true" type="xs:string"/>
                    <xs:element minOccurs="0" name="pageFlag" nillable="true" type="xs:string"/>
                    <xs:element minOccurs="0" name="pageRowNum" nillable="true" type="xs:string"/>
                    <xs:element minOccurs="0" name="rowNumStart" nillable="true" type="xs:string"/>
                </xs:sequence>
            </xs:complexType>
            <xs:complexType name="TxInsuranceExtensionEhm">
                <xs:sequence>
                    <xs:element minOccurs="0" name="operator" nillable="true" type="xs:string"/>
                    <xs:element minOccurs="0" name="operatorKey" nillable="true" type="xs:string"/>
                </xs:sequence>
            </xs:complexType>
            <xs:complexType name="TxInsuranceRequestExtensionEhm">
                <xs:complexContent>
                    <xs:extension base="ax276:TxInsuranceExtensionEhm">
                        <xs:sequence/>
                    </xs:extension>
                </xs:complexContent>
            </xs:complexType>
            <xs:complexType name="TxInsuranceResponseEhm">
                <xs:complexContent>
                    <xs:extension base="ax276:TxInsuranceEhm">
                        <xs:sequence>
                            <xs:element minOccurs="0" name="oinsuranceExtensionEhm" nillable="true" type="ax276:OinsuranceExtensionEhm"/>
                            <xs:element minOccurs="0" name="transResultEhm" nillable="true" type="ax276:TransResultEhm"/>
                        </xs:sequence>
                    </xs:extension>
                </xs:complexContent>
            </xs:complexType>
            <xs:complexType name="OinsuranceExtensionEhm">
                <xs:sequence>
                    <xs:element minOccurs="0" name="maxRecords" nillable="true" type="xs:string"/>
                </xs:sequence>
            </xs:complexType>
            <xs:complexType name="TransResultEhm">
                <xs:sequence>
                    <xs:element minOccurs="0" name="errorNo" nillable="true" type="xs:string"/>
                    <xs:element minOccurs="0" name="errorType" nillable="true" type="xs:string"/>
                    <xs:element minOccurs="0" name="resultCode" nillable="true" type="xs:string"/>
                    <xs:element minOccurs="0" name="resultInfoDesc" nillable="true" type="xs:string"/>
                    <xs:element minOccurs="0" name="stackTrace" nillable="true" type="xs:string"/>
                </xs:sequence>
            </xs:complexType>
            <xs:complexType name="TxInsuranceResponseExtensionEhm">
                <xs:complexContent>
                    <xs:extension base="ax276:TxInsuranceExtensionEhm">
                        <xs:sequence/>
                    </xs:extension>
                </xs:complexContent>
            </xs:complexType>
        </xs:schema>
    </wsdl:types>
    <wsdl:message name="carQueryRequest">
        <wsdl:part name="parameters" element="ns:carQuery"/>
    </wsdl:message>
    <wsdl:message name="carQueryResponse">
        <wsdl:part name="parameters" element="ns:carQueryResponse"/>
    </wsdl:message>
    <wsdl:message name="setCarQueryTpInterfaceRequest">
        <wsdl:part name="parameters" element="ns:setCarQueryTpInterface"/>
    </wsdl:message>
    <wsdl:message name="getCarQueryTpInterfaceRequest">
        <wsdl:part name="parameters" element="ns:getCarQueryTpInterface"/>
    </wsdl:message>
    <wsdl:message name="getCarQueryTpInterfaceResponse">
        <wsdl:part name="parameters" element="ns:getCarQueryTpInterfaceResponse"/>
    </wsdl:message>
    <wsdl:message name="getCarQueryTpCheckInterfaceRequest">
        <wsdl:part name="parameters" element="ns:getCarQueryTpCheckInterface"/>
    </wsdl:message>
    <wsdl:message name="getCarQueryTpCheckInterfaceResponse">
        <wsdl:part name="parameters" element="ns:getCarQueryTpCheckInterfaceResponse"/>
    </wsdl:message>
    <wsdl:message name="setCarQueryTpCheckInterfaceRequest">
        <wsdl:part name="parameters" element="ns:setCarQueryTpCheckInterface"/>
    </wsdl:message>
    <wsdl:portType name="CarQueryServicePortType">
        <wsdl:operation name="carQuery">
            <wsdl:input message="ns:carQueryRequest" wsaw:Action="urn:carQuery"/>
            <wsdl:output message="ns:carQueryResponse" wsaw:Action="urn:carQueryResponse"/>
        </wsdl:operation>
        <wsdl:operation name="setCarQueryTpInterface">
            <wsdl:input message="ns:setCarQueryTpInterfaceRequest" wsaw:Action="urn:setCarQueryTpInterface"/>
        </wsdl:operation>
        <wsdl:operation name="getCarQueryTpInterface">
            <wsdl:input message="ns:getCarQueryTpInterfaceRequest" wsaw:Action="urn:getCarQueryTpInterface"/>
            <wsdl:output message="ns:getCarQueryTpInterfaceResponse" wsaw:Action="urn:getCarQueryTpInterfaceResponse"/>
        </wsdl:operation>
        <wsdl:operation name="getCarQueryTpCheckInterface">
            <wsdl:input message="ns:getCarQueryTpCheckInterfaceRequest" wsaw:Action="urn:getCarQueryTpCheckInterface"/>
            <wsdl:output message="ns:getCarQueryTpCheckInterfaceResponse" wsaw:Action="urn:getCarQueryTpCheckInterfaceResponse"/>
        </wsdl:operation>
        <wsdl:operation name="setCarQueryTpCheckInterface">
            <wsdl:input message="ns:setCarQueryTpCheckInterfaceRequest" wsaw:Action="urn:setCarQueryTpCheckInterface"/>
        </wsdl:operation>
    </wsdl:portType>
    <wsdl:binding name="CarQueryServiceSoap11Binding" type="ns:CarQueryServicePortType">
        <soap:binding transport="http://schemas.xmlsoap.org/soap/http" style="document"/>
        <wsdl:operation name="carQuery">
            <soap:operation soapAction="urn:carQuery" style="document"/>
            <wsdl:input>
                <soap:body use="literal"/>
            </wsdl:input>
            <wsdl:output>
                <soap:body use="literal"/>
            </wsdl:output>
        </wsdl:operation>
        <wsdl:operation name="setCarQueryTpInterface">
            <soap:operation soapAction="urn:setCarQueryTpInterface" style="document"/>
            <wsdl:input>
                <soap:body use="literal"/>
            </wsdl:input>
        </wsdl:operation>
        <wsdl:operation name="getCarQueryTpInterface">
            <soap:operation soapAction="urn:getCarQueryTpInterface" style="document"/>
            <wsdl:input>
                <soap:body use="literal"/>
            </wsdl:input>
            <wsdl:output>
                <soap:body use="literal"/>
            </wsdl:output>
        </wsdl:operation>
        <wsdl:operation name="getCarQueryTpCheckInterface">
            <soap:operation soapAction="urn:getCarQueryTpCheckInterface" style="document"/>
            <wsdl:input>
                <soap:body use="literal"/>
            </wsdl:input>
            <wsdl:output>
                <soap:body use="literal"/>
            </wsdl:output>
        </wsdl:operation>
        <wsdl:operation name="setCarQueryTpCheckInterface">
            <soap:operation soapAction="urn:setCarQueryTpCheckInterface" style="document"/>
            <wsdl:input>
                <soap:body use="literal"/>
            </wsdl:input>
        </wsdl:operation>
    </wsdl:binding>
    <wsdl:binding name="CarQueryServiceSoap12Binding" type="ns:CarQueryServicePortType">
        <soap12:binding transport="http://schemas.xmlsoap.org/soap/http" style="document"/>
        <wsdl:operation name="carQuery">
            <soap12:operation soapAction="urn:carQuery" style="document"/>
            <wsdl:input>
                <soap12:body use="literal"/>
            </wsdl:input>
            <wsdl:output>
                <soap12:body use="literal"/>
            </wsdl:output>
        </wsdl:operation>
        <wsdl:operation name="setCarQueryTpInterface">
            <soap12:operation soapAction="urn:setCarQueryTpInterface" style="document"/>
            <wsdl:input>
                <soap12:body use="literal"/>
            </wsdl:input>
        </wsdl:operation>
        <wsdl:operation name="getCarQueryTpInterface">
            <soap12:operation soapAction="urn:getCarQueryTpInterface" style="document"/>
            <wsdl:input>
                <soap12:body use="literal"/>
            </wsdl:input>
            <wsdl:output>
                <soap12:body use="literal"/>
            </wsdl:output>
        </wsdl:operation>
        <wsdl:operation name="getCarQueryTpCheckInterface">
            <soap12:operation soapAction="urn:getCarQueryTpCheckInterface" style="document"/>
            <wsdl:input>
                <soap12:body use="literal"/>
            </wsdl:input>
            <wsdl:output>
                <soap12:body use="literal"/>
            </wsdl:output>
        </wsdl:operation>
        <wsdl:operation name="setCarQueryTpCheckInterface">
            <soap12:operation soapAction="urn:setCarQueryTpCheckInterface" style="document"/>
            <wsdl:input>
                <soap12:body use="literal"/>
            </wsdl:input>
        </wsdl:operation>
    </wsdl:binding>
    <wsdl:binding name="CarQueryServiceHttpBinding" type="ns:CarQueryServicePortType">
        <http:binding verb="POST"/>
        <wsdl:operation name="carQuery">
            <http:operation location="carQuery"/>
            <wsdl:input>
                <mime:content type="application/xml" part="parameters"/>
            </wsdl:input>
            <wsdl:output>
                <mime:content type="application/xml" part="parameters"/>
            </wsdl:output>
        </wsdl:operation>
        <wsdl:operation name="setCarQueryTpInterface">
            <http:operation location="setCarQueryTpInterface"/>
            <wsdl:input>
                <mime:content type="application/xml" part="parameters"/>
            </wsdl:input>
        </wsdl:operation>
        <wsdl:operation name="getCarQueryTpInterface">
            <http:operation location="getCarQueryTpInterface"/>
            <wsdl:input>
                <mime:content type="application/xml" part="parameters"/>
            </wsdl:input>
            <wsdl:output>
                <mime:content type="application/xml" part="parameters"/>
            </wsdl:output>
        </wsdl:operation>
        <wsdl:operation name="getCarQueryTpCheckInterface">
            <http:operation location="getCarQueryTpCheckInterface"/>
            <wsdl:input>
                <mime:content type="application/xml" part="parameters"/>
            </wsdl:input>
            <wsdl:output>
                <mime:content type="application/xml" part="parameters"/>
            </wsdl:output>
        </wsdl:operation>
        <wsdl:operation name="setCarQueryTpCheckInterface">
            <http:operation location="setCarQueryTpCheckInterface"/>
            <wsdl:input>
                <mime:content type="application/xml" part="parameters"/>
            </wsdl:input>
        </wsdl:operation>
    </wsdl:binding>
    <wsdl:service name="CarQueryService">
        <wsdl:port name="CarQueryServiceHttpSoap11Endpoint" binding="ns:CarQueryServiceSoap11Binding">
            <soap:address location="http://114.255.29.204:8010/frontServiceCenter/services/CarQueryService.CarQueryServiceHttpSoap11Endpoint/"/>
        </wsdl:port>
        <wsdl:port name="CarQueryServiceHttpSoap12Endpoint" binding="ns:CarQueryServiceSoap12Binding">
            <soap12:address location="http://114.255.29.204:8010/frontServiceCenter/services/CarQueryService.CarQueryServiceHttpSoap12Endpoint/"/>
        </wsdl:port>
        <wsdl:port name="CarQueryServiceHttpEndpoint" binding="ns:CarQueryServiceHttpBinding">
            <http:address location="http://114.255.29.204:8010/frontServiceCenter/services/CarQueryService.CarQueryServiceHttpEndpoint/"/>
        </wsdl:port>
    </wsdl:service>
</wsdl:definitions>`
	soap_client, err := NewSoapClientForContent(wsdl_xml)
	if err != nil {
		panic(err)
	}
	fmt.Println(soap_client.requestUrl)

	res, err := soap_client.Call("carQuery", utils.M{
		"policySort":   "CQ0",
		"licenseNo":    "渝B6S919",
		"licenseType":  "02",
		"cityAreaCode": "03",
		"txInsuranceRequestEhm": utils.M{
			"transExeDate": "",
			"transExeTime": "",
		},
		"txInsuranceRequestExtensionEhm": utils.M{
			"operator":    "CQ0_Test",
			"operatorKey": "123456",
		},
	})
	if err != nil {
		log.Println(err)
	}

	fmt.Println(res)
}

func TestNewSoapClient4(t *testing.T) {
	soap_client, err := NewSoapClient("http://114.255.29.204:8010/frontServiceCenter/services/PremiumCaculateService?wsdl")
	if err != nil {
		panic(err)
	}
	fmt.Println(soap_client.requestUrl)

	res, err := soap_client.Call("premiumCaculate", getSoapData())
	if err != nil {
		log.Println(err)
	}

	fmt.Println(res)
}

func TestNewSoapClient3(t *testing.T) {
	content := `<soap:Envelope xmlns:soap="http://www.w3.org/2003/05/soap-envelope" xmlns:ser="http://service.front.sinosoft.com" xmlns:xsd="http://dto.front.sinosoft.com/xsd" xmlns:xsd1="http://common.dto.front.sinosoft.com/xsd">
   <soap:Header/>
   <soap:Body>
      <ser:carQuery>
      </ser:carQuery>
   </soap:Body>
</soap:Envelope>`

	url_str := "http://114.255.29.204:8010/frontServiceCenter/services/CarQueryService.CarQueryServiceHttpSoap12Endpoint/"
	client := utils.NewHttpClient()
	client.SetHeader("Content-Type", "application/soap+xml; charset=utf-8")
	client.SetHeader("Content-Length", fmt.Sprintf("%v", len(content)))
	req, err := client.Request("POST", url_str, strings.NewReader(content))
	if err != nil {
		log.Println(err)
	}

	fmt.Println(string(req.Content))
}

func TestExplainXml(t *testing.T) {
	content := `<?xml version='1.0' encoding='UTF-8'?>
<wsdl:definitions name="AppServiceService" targetNamespace="http://service.webservice.yaic.com/" 
    xmlns:ns1="http://cxf.apache.org/bindings/xformat" 
    xmlns:soap="http://schemas.xmlsoap.org/wsdl/soap/" 
    xmlns:tns="http://service.webservice.yaic.com/" 
    xmlns:wsdl="http://schemas.xmlsoap.org/wsdl/" 
    xmlns:xsd="http://www.w3.org/2001/XMLSchema">
    <wsdl:types>
        <xs:schema attributeFormDefault="unqualified" elementFormDefault="unqualified" targetNamespace="http://service.webservice.yaic.com/" 
            xmlns:tns="http://service.webservice.yaic.com/" 
            xmlns:xs="http://www.w3.org/2001/XMLSchema">
            <xs:complexType name="appRequestBean">
                <xs:sequence>
                    <xs:element maxOccurs="unbounded" minOccurs="0" name="appInfos" nillable="true" type="tns:appInfo" />
                    <xs:element minOccurs="0" name="baseInfo" type="tns:baseInfo" />
                    <xs:element maxOccurs="unbounded" minOccurs="0" name="kindInfos" nillable="true" type="tns:kindInfo" />
                </xs:sequence>
            </xs:complexType>
            <xs:complexType name="appInfo">
                <xs:sequence>
                    <xs:element minOccurs="0" name="appBaseInfo" type="tns:appBaseInfo" />
                    <xs:element maxOccurs="unbounded" minOccurs="0" name="cbnfcInfos" nillable="true" type="tns:cbnfcInfo" />
                    <xs:element minOccurs="0" name="extraInfo" type="tns:extraInfo" />
                    <xs:element maxOccurs="unbounded" minOccurs="0" name="insuredInfos" nillable="true" type="tns:insuredInfo" />
                    <xs:element minOccurs="0" name="payInfo" type="tns:payInfo" />
                    <xs:element minOccurs="0" name="payerInfo" type="tns:payerInfo" />
                    <xs:element maxOccurs="unbounded" minOccurs="0" name="paymentInfos" nillable="true" type="tns:paymentInfo" />
                    <xs:element maxOccurs="unbounded" minOccurs="0" name="propertyInfos" nillable="true" type="tns:propertyInfo" />
                </xs:sequence>
            </xs:complexType>
            <xs:complexType name="appBaseInfo">
                <xs:sequence>
                    <xs:element minOccurs="0" name="appBirthday" type="xs:dateTime" />
                    <xs:element minOccurs="0" name="appSex" type="xs:string" />
                    <xs:element minOccurs="0" name="appaddress" type="xs:string" />
                    <xs:element minOccurs="0" name="appemail" type="xs:string" />
                    <xs:element minOccurs="0" name="appid" type="xs:string" />
                    <xs:element minOccurs="0" name="appmrk" type="xs:string" />
                    <xs:element minOccurs="0" name="appname" type="xs:string" />
                    <xs:element minOccurs="0" name="appphone" type="xs:string" />
                    <xs:element minOccurs="0" name="apppost" type="xs:string" />
                    <xs:element minOccurs="0" name="apprel" type="xs:string" />
                    <xs:element minOccurs="0" name="apptel" type="xs:string" />
                    <xs:element minOccurs="0" name="apptype" type="xs:string" />
                    <xs:element minOccurs="0" name="appworkdpt" type="xs:string" />
                    <xs:element minOccurs="0" name="businessid" type="xs:int" />
                    <xs:element minOccurs="0" name="cntrcertfcde" type="xs:string" />
                    <xs:element minOccurs="0" name="cntrcertfendtm" type="xs:dateTime" />
                    <xs:element minOccurs="0" name="cntrnme" type="xs:string" />
                    <xs:element minOccurs="0" name="legalcertfcde" type="xs:string" />
                    <xs:element minOccurs="0" name="legalcertfendtm" type="xs:dateTime" />
                    <xs:element minOccurs="0" name="legalnme" type="xs:string" />
                    <xs:element minOccurs="0" name="resvtxt1" type="xs:string" />
                    <xs:element minOccurs="0" name="resvtxt2" type="xs:string" />
                    <xs:element minOccurs="0" name="resvtxt3" type="xs:string" />
                </xs:sequence>
            </xs:complexType>
            <xs:complexType name="cbnfcInfo">
                <xs:sequence>
                    <xs:element minOccurs="0" name="cbenford" type="xs:string" />
                    <xs:element minOccurs="0" name="cbnfcnme" type="xs:string" />
                    <xs:element minOccurs="0" name="ccertfcde" type="xs:string" />
                    <xs:element minOccurs="0" name="ccertfcls" type="xs:string" />
                    <xs:element minOccurs="0" name="crelcde" type="xs:string" />
                    <xs:element minOccurs="0" name="csex" type="xs:string" />
                    <xs:element minOccurs="0" name="nbenfprop" type="xs:double" />
                    <xs:element minOccurs="0" name="tbirthday" type="xs:dateTime" />
                </xs:sequence>
            </xs:complexType>
            <xs:complexType name="extraInfo">
                <xs:sequence>
                    <xs:element minOccurs="0" name="agencycode" type="xs:string" />
                    <xs:element minOccurs="0" name="agreementno" type="xs:string" />
                    <xs:element minOccurs="0" name="businessnature" type="xs:string" />
                    <xs:element minOccurs="0" name="businessnature2" type="xs:string" />
                    <xs:element minOccurs="0" name="businessnature3" type="xs:string" />
                    <xs:element minOccurs="0" name="cbrkslscde" type="xs:string" />
                    <xs:element minOccurs="0" name="gscomcode" type="xs:string" />
                    <xs:element minOccurs="0" name="handlercode" type="xs:string" />
                    <xs:element minOccurs="0" name="handlername" type="xs:string" />
                    <xs:element minOccurs="0" name="hxcdopercode" type="xs:string" />
                    <xs:element minOccurs="0" name="prjctgmidtyp" type="xs:string" />
                    <xs:element minOccurs="0" name="prjctgsubtyp" type="xs:string" />
                    <xs:element minOccurs="0" name="prjctgtyp" type="xs:string" />
                </xs:sequence>
            </xs:complexType>
            <xs:complexType name="insuredInfo">
                <xs:sequence>
                    <xs:element minOccurs="0" name="accountBank" type="xs:string" />
                    <xs:element minOccurs="0" name="accountNo" type="xs:string" />
                    <xs:element minOccurs="0" name="allinesflag" type="xs:string" />
                    <xs:element minOccurs="0" name="apprel" type="xs:string" />
                    <xs:element minOccurs="0" name="businessPlace" type="xs:string" />
                    <xs:element minOccurs="0" name="CResvTxt27" type="xs:string" />
                    <xs:element minOccurs="0" name="CResvTxt5" type="xs:string" />
                    <xs:element minOccurs="0" name="cduty" type="xs:string" />
                    <xs:element minOccurs="0" name="clntaddr" type="xs:string" />
                    <xs:element minOccurs="0" name="czipcde" type="xs:string" />
                    <xs:element minOccurs="0" name="nperamt" type="xs:double" />
                    <xs:element minOccurs="0" name="nperprm" type="xs:double" />
                    <xs:element minOccurs="0" name="occupcde" type="xs:string" />
                    <xs:element minOccurs="0" name="page" type="xs:int" />
                    <xs:element minOccurs="0" name="pbirthday" type="xs:dateTime" />
                    <xs:element minOccurs="0" name="pemail" type="xs:string" />
                    <xs:element minOccurs="0" name="pid" type="xs:string" />
                    <xs:element minOccurs="0" name="pmrk" type="xs:string" />
                    <xs:element minOccurs="0" name="pname" type="xs:string" />
                    <xs:element minOccurs="0" name="pphone" type="xs:string" />
                    <xs:element minOccurs="0" name="psex" type="xs:string" />
                    <xs:element minOccurs="0" name="ptel" type="xs:string" />
                    <xs:element minOccurs="0" name="ptype" type="xs:string" />
                    <xs:element minOccurs="0" name="pworkdpt" type="xs:string" />
                    <xs:element minOccurs="0" name="resvtxt7" type="xs:string" />
                    <xs:element minOccurs="0" name="resvtxt8" type="xs:string" />
                    <xs:element minOccurs="0" name="socialSecurityNo" type="xs:string" />
                    <xs:element minOccurs="0" name="socialSecurityNoPlace" type="xs:string" />
                </xs:sequence>
            </xs:complexType>
            <xs:complexType name="payInfo">
                <xs:sequence>
                    <xs:element minOccurs="0" name="amount" type="xs:double" />
                    <xs:element minOccurs="0" name="currencyType" type="xs:string" />
                    <xs:element minOccurs="0" name="dataType" type="xs:string" />
                    <xs:element minOccurs="0" name="departmentCode" type="xs:string" />
                    <xs:element minOccurs="0" name="inpaymentDate" type="xs:string" />
                    <xs:element minOccurs="0" name="inpaymentTime" type="xs:string" />
                    <xs:element minOccurs="0" name="insumidNo" type="xs:string" />
                    <xs:element minOccurs="0" name="payAppNo" type="xs:string" />
                    <xs:element minOccurs="0" name="payChannel" type="xs:string" />
                    <xs:element minOccurs="0" name="payWay" type="xs:string" />
                    <xs:element minOccurs="0" name="subCompany" type="xs:string" />
                </xs:sequence>
            </xs:complexType>
            <xs:complexType name="payerInfo">
                <xs:sequence>
                    <xs:element minOccurs="0" name="cardFlag" type="xs:string" />
                    <xs:element minOccurs="0" name="custaccountNo" type="xs:string" />
                    <xs:element minOccurs="0" name="cvv2" type="xs:string" />
                    <xs:element minOccurs="0" name="expiredDate" type="xs:string" />
                    <xs:element minOccurs="0" name="payerId" type="xs:string" />
                    <xs:element minOccurs="0" name="payerName" type="xs:string" />
                    <xs:element minOccurs="0" name="payerTel" type="xs:string" />
                    <xs:element minOccurs="0" name="payerType" type="xs:string" />
                    <xs:element minOccurs="0" name="cBankArea" type="xs:string" />
                    <xs:element minOccurs="0" name="cBankCde" type="xs:string" />
                    <xs:element minOccurs="0" name="cBankPro" type="xs:string" />
                    <xs:element minOccurs="0" name="cPubPri" type="xs:string" />
                </xs:sequence>
            </xs:complexType>
            <xs:complexType name="paymentInfo">
                <xs:sequence>
                    <xs:element minOccurs="0" name="bsnum" type="xs:long" />
                    <xs:element minOccurs="0" name="bsnums" type="xs:long" />
                    <xs:element minOccurs="0" name="seqno" type="xs:long" />
                    <xs:element minOccurs="0" name="serialno" type="xs:string" />
                    <xs:element minOccurs="0" name="nPayablePrm" type="xs:double" />
                    <xs:element minOccurs="0" name="nPrmVar" type="xs:double" />
                    <xs:element minOccurs="0" name="tPayBgnTm" type="xs:dateTime" />
                    <xs:element minOccurs="0" name="tPayEndTm" type="xs:dateTime" />
                </xs:sequence>
            </xs:complexType>
            <xs:complexType name="propertyInfo">
                <xs:sequence>
                    <xs:element minOccurs="0" name="blname" type="xs:string" />
                    <xs:element minOccurs="0" name="blno" type="xs:string" />
                    <xs:element minOccurs="0" name="carrybillno" type="xs:string" />
                    <xs:element minOccurs="0" name="carshipno" type="xs:string" />
                    <xs:element minOccurs="0" name="checicode" type="xs:string" />
                    <xs:element minOccurs="0" name="classname" type="xs:string" />
                    <xs:element minOccurs="0" name="conveyance" type="xs:string" />
                    <xs:element minOccurs="0" name="country" type="xs:string" />
                    <xs:element minOccurs="0" name="cresvtxt2" type="xs:string" />
                    <xs:element minOccurs="0" name="cresvtxt4" type="xs:string" />
                    <xs:element minOccurs="0" name="cresvtxt5" type="xs:string" />
                    <xs:element minOccurs="0" name="cresvtxt6" type="xs:string" />
                    <xs:element minOccurs="0" name="cresvtxt7" type="xs:string" />
                    <xs:element minOccurs="0" name="ctgtaddr" type="xs:string" />
                    <xs:element minOccurs="0" name="ctgttxtfld5" type="xs:string" />
                    <xs:element minOccurs="0" name="ctgtzip" type="xs:string" />
                    <xs:element minOccurs="0" name="delaydate" type="xs:dateTime" />
                    <xs:element minOccurs="0" name="driverseatnum" type="xs:int" />
                    <xs:element minOccurs="0" name="endsitename" type="xs:string" />
                    <xs:element minOccurs="0" name="endtime" type="xs:dateTime" />
                    <xs:element minOccurs="0" name="engineNo" type="xs:string" />
                    <xs:element minOccurs="0" name="financialInstitution" type="xs:string" />
                    <xs:element minOccurs="0" name="flyOrderNo" type="xs:string" />
                    <xs:element minOccurs="0" name="flydate" type="xs:dateTime" />
                    <xs:element minOccurs="0" name="flylanddate" type="xs:dateTime" />
                    <xs:element minOccurs="0" name="flyno" type="xs:string" />
                    <xs:element minOccurs="0" name="frameNo" type="xs:string" />
                    <xs:element minOccurs="0" name="fromname" type="xs:string" />
                    <xs:element minOccurs="0" name="hospital" type="xs:string" />
                    <xs:element minOccurs="0" name="hospitallevel" type="xs:string" />
                    <xs:element minOccurs="0" name="initLoanAmount" type="xs:double" />
                    <xs:element minOccurs="0" name="licenseNo" type="xs:string" />
                    <xs:element minOccurs="0" name="loanApprovalNo" type="xs:string" />
                    <xs:element minOccurs="0" name="loanContractNo" type="xs:string" />
                    <xs:element minOccurs="0" name="loanamount" type="xs:string" />
                    <xs:element minOccurs="0" name="loanbankcode" type="xs:string" />
                    <xs:element minOccurs="0" name="loanbegDate" type="xs:dateTime" />
                    <xs:element minOccurs="0" name="loancontractno" type="xs:string" />
                    <xs:element minOccurs="0" name="loandate" type="xs:dateTime" />
                    <xs:element minOccurs="0" name="loanmonth" type="xs:int" />
                    <xs:element minOccurs="0" name="loanstopDate" type="xs:dateTime" />
                    <xs:element minOccurs="0" name="matchplace" type="xs:string" />
                    <xs:element minOccurs="0" name="nresvnum1" type="xs:int" />
                    <xs:element minOccurs="0" name="nresvnum2" type="xs:int" />
                    <xs:element minOccurs="0" name="opertime" type="xs:dateTime" />
                    <xs:element minOccurs="0" name="passengerseatnum" type="xs:int" />
                    <xs:element minOccurs="0" name="payNumber" type="xs:string" />
                    <xs:element minOccurs="0" name="power" type="xs:double" />
                    <xs:element minOccurs="0" name="reflydate" type="xs:dateTime" />
                    <xs:element minOccurs="0" name="reflyno" type="xs:string" />
                    <xs:element minOccurs="0" name="repaiddate" type="xs:dateTime" />
                    <xs:element minOccurs="0" name="repaidtype" type="xs:string" />
                    <xs:element minOccurs="0" name="schoolname" type="xs:string" />
                    <xs:element minOccurs="0" name="seatno" type="xs:string" />
                    <xs:element minOccurs="0" name="seatnum" type="xs:int" />
                    <xs:element minOccurs="0" name="seqno" type="xs:int" />
                    <xs:element minOccurs="0" name="startsitename" type="xs:string" />
                    <xs:element minOccurs="0" name="starttime" type="xs:dateTime" />
                    <xs:element minOccurs="0" name="theBorrower" type="xs:string" />
                    <xs:element minOccurs="0" name="ticketAmt" type="xs:double" />
                    <xs:element minOccurs="0" name="ticketEndDate" type="xs:dateTime" />
                    <xs:element minOccurs="0" name="ticketdate" type="xs:dateTime" />
                    <xs:element minOccurs="0" name="ticketno" type="xs:string" />
                    <xs:element minOccurs="0" name="toname" type="xs:string" />
                    <xs:element minOccurs="0" name="university" type="xs:string" />
                    <xs:element minOccurs="0" name="userquality" type="xs:string" />
                    <xs:element minOccurs="0" name="yafcProDate1" type="xs:dateTime" />
                    <xs:element minOccurs="0" name="yafcProDate2" type="xs:dateTime" />
                    <xs:element minOccurs="0" name="yafcProNum1" type="xs:double" />
                    <xs:element minOccurs="0" name="yafcProNum2" type="xs:double" />
                    <xs:element minOccurs="0" name="yafcProNum3" type="xs:double" />
                    <xs:element minOccurs="0" name="yafcProNum4" type="xs:int" />
                    <xs:element minOccurs="0" name="yafcProText1" type="xs:string" />
                    <xs:element minOccurs="0" name="yafcProText2" type="xs:string" />
                    <xs:element minOccurs="0" name="yafcProText3" type="xs:string" />
                    <xs:element minOccurs="0" name="yafcProText4" type="xs:string" />
                    <xs:element minOccurs="0" name="yafcProText5" type="xs:string" />
                </xs:sequence>
            </xs:complexType>
            <xs:complexType name="baseInfo">
                <xs:sequence>
                    <xs:element minOccurs="0" name="amount" type="xs:double" />
                    <xs:element minOccurs="0" name="amt" type="xs:double" />
                    <xs:element minOccurs="0" name="arguesolution" type="xs:string" />
                    <xs:element minOccurs="0" name="backurl" type="xs:string" />
                    <xs:element minOccurs="0" name="bsnum" type="xs:int" />
                    <xs:element minOccurs="0" name="businessno" type="xs:string" />
                    <xs:element minOccurs="0" name="cagreeType" type="xs:string" />
                    <xs:element minOccurs="0" name="ccardbsnstyp" type="xs:string" />
                    <xs:element minOccurs="0" name="checkcode" type="xs:string" />
                    <xs:element minOccurs="0" name="isgrp" type="xs:string" />
                    <xs:element minOccurs="0" name="isticket" type="xs:string" />
                    <xs:element minOccurs="0" name="opercode" type="xs:string" />
                    <xs:element minOccurs="0" name="password" type="xs:string" />
                    <xs:element minOccurs="0" name="payappno" type="xs:string" />
                    <xs:element minOccurs="0" name="payflag" type="xs:string" />
                    <xs:element minOccurs="0" name="rationcode" type="xs:string" />
                    <xs:element minOccurs="0" name="realflag" type="xs:string" />
                    <xs:element minOccurs="0" name="serialno" type="xs:string" />
                    <xs:element minOccurs="0" name="tapptm" type="xs:dateTime" />
                    <xs:element minOccurs="0" name="tinsrncbgntm" type="xs:dateTime" />
                    <xs:element minOccurs="0" name="tinsrncendtm" type="xs:dateTime" />
                    <xs:element minOccurs="0" name="tissuetm" type="xs:dateTime" />
                    <xs:element minOccurs="0" name="user" type="xs:string" />
                </xs:sequence>
            </xs:complexType>
            <xs:complexType name="kindInfo">
                <xs:sequence>
                    <xs:element minOccurs="0" name="amount" type="xs:double" />
                    <xs:element minOccurs="0" name="cresvtxt8" type="xs:string" />
                    <xs:element minOccurs="0" name="itemcode" type="xs:string" />
                    <xs:element minOccurs="0" name="kindcode" type="xs:string" />
                    <xs:element minOccurs="0" name="ndductamt" type="xs:double" />
                    <xs:element minOccurs="0" name="ndductrate" type="xs:double" />
                    <xs:element minOccurs="0" name="nfloatrate" type="xs:double" />
                    <xs:element minOccurs="0" name="nindemlmt" type="xs:double" />
                    <xs:element minOccurs="0" name="nonceindemlmt" type="xs:double" />
                    <xs:element minOccurs="0" name="nperamt" type="xs:double" />
                    <xs:element minOccurs="0" name="nperprm" type="xs:double" />
                    <xs:element minOccurs="0" name="nresvnum3" type="xs:int" />
                    <xs:element minOccurs="0" name="premium" type="xs:double" />
                    <xs:element minOccurs="0" name="rate" type="xs:double" />
                    <xs:element minOccurs="0" name="yafcProNum1" type="xs:double" />
                    <xs:element minOccurs="0" name="yafcProNum2" type="xs:double" />
                    <xs:element minOccurs="0" name="yafcProNum3" type="xs:double" />
                    <xs:element minOccurs="0" name="yafcProText1" type="xs:string" />
                </xs:sequence>
            </xs:complexType>
            <xs:complexType name="appResponseBean">
                <xs:sequence>
                    <xs:element maxOccurs="unbounded" minOccurs="0" name="appInfoRes" nillable="true" type="tns:appInfoRes" />
                    <xs:element minOccurs="0" name="flag" type="xs:string" />
                    <xs:element minOccurs="0" name="reason" type="xs:string" />
                </xs:sequence>
            </xs:complexType>
            <xs:complexType name="appInfoRes">
                <xs:sequence>
                    <xs:element minOccurs="0" name="businessno" type="xs:string" />
                    <xs:element minOccurs="0" name="payappno" type="xs:string" />
                    <xs:element minOccurs="0" name="paymenturl" type="xs:string" />
                    <xs:element minOccurs="0" name="pdfurl" type="xs:string" />
                    <xs:element minOccurs="0" name="policyno" type="xs:string" />
                    <xs:element minOccurs="0" name="seqno" type="xs:int" />
                    <xs:element minOccurs="0" name="serialno" type="xs:string" />
                    <xs:element minOccurs="0" name="nAddedTax" type="xs:double" />
                    <xs:element minOccurs="0" name="nNoTaxPrm" type="xs:double" />
                </xs:sequence>
            </xs:complexType>
            <xs:element name="appBatchRequest" type="tns:appBatchRequest" />
            <xs:complexType name="appBatchRequest">
                <xs:sequence>
                    <xs:element minOccurs="0" name="arg0" type="tns:appRequestBean" />
                </xs:sequence>
            </xs:complexType>
            <xs:element name="appBatchRequestResponse" type="tns:appBatchRequestResponse" />
            <xs:complexType name="appBatchRequestResponse">
                <xs:sequence>
                    <xs:element minOccurs="0" name="return" type="tns:appResponseBean" />
                </xs:sequence>
            </xs:complexType>
            <xs:element name="appRequest" type="tns:appRequest" />
            <xs:complexType name="appRequest">
                <xs:sequence>
                    <xs:element minOccurs="0" name="arg0" type="tns:appRequestBean" />
                </xs:sequence>
            </xs:complexType>
            <xs:element name="appRequestResponse" type="tns:appRequestResponse" />
            <xs:complexType name="appRequestResponse">
                <xs:sequence>
                    <xs:element minOccurs="0" name="return" type="tns:appResponseBean" />
                </xs:sequence>
            </xs:complexType>
        </xs:schema>
    </wsdl:types>
    <wsdl:message name="appRequestResponse">
        <wsdl:part element="tns:appRequestResponse" name="parameters">
        </wsdl:part>
    </wsdl:message>
    <wsdl:message name="appBatchRequest">
        <wsdl:part element="tns:appBatchRequest" name="parameters">
        </wsdl:part>
    </wsdl:message>
    <wsdl:message name="appRequest">
        <wsdl:part element="tns:appRequest" name="parameters">
        </wsdl:part>
    </wsdl:message>
    <wsdl:message name="appBatchRequestResponse">
        <wsdl:part element="tns:appBatchRequestResponse" name="parameters">
        </wsdl:part>
    </wsdl:message>
    <wsdl:portType name="AppService">
        <wsdl:operation name="appBatchRequest">
            <wsdl:input message="tns:appBatchRequest" name="appBatchRequest">
            </wsdl:input>
            <wsdl:output message="tns:appBatchRequestResponse" name="appBatchRequestResponse">
            </wsdl:output>
        </wsdl:operation>
        <wsdl:operation name="appRequest">
            <wsdl:input message="tns:appRequest" name="appRequest">
            </wsdl:input>
            <wsdl:output message="tns:appRequestResponse" name="appRequestResponse">
            </wsdl:output>
        </wsdl:operation>
    </wsdl:portType>
    <wsdl:binding name="AppServiceServiceSoapBinding" type="tns:AppService">
        <soap:binding style="document" transport="http://schemas.xmlsoap.org/soap/http" />
        <wsdl:operation name="appBatchRequest">
            <soap:operation soapAction="" style="document" />
            <wsdl:input name="appBatchRequest">
                <soap:body use="literal" />
            </wsdl:input>
            <wsdl:output name="appBatchRequestResponse">
                <soap:body use="literal" />
            </wsdl:output>
        </wsdl:operation>
        <wsdl:operation name="appRequest">
            <soap:operation soapAction="" style="document" />
            <wsdl:input name="appRequest">
                <soap:body use="literal" />
            </wsdl:input>
            <wsdl:output name="appRequestResponse">
                <soap:body use="literal" />
            </wsdl:output>
        </wsdl:operation>
    </wsdl:binding>
    <wsdl:service name="AppServiceService">
        <wsdl:port binding="tns:AppServiceServiceSoapBinding" name="AppServicePort">
            <soap:address location="http://tapi.yaic.com.cn:80/yaicservice/api/appservice" />
        </wsdl:port>
    </wsdl:service>
</wsdl:definitions>`
	soap_client, err := NewSoapClientForContent(content)
	if err != nil {
		panic(err)
	}

	soap_client.Call("")
}

func getSoapData() utils.M {
	str := "{\"txInsuranceRequestEhm\":{\"transExeDate\":\"20170925\",\"transExeTime\":\"201300\"},\"txInsuranceRequestExtensionEhm\":{\"operator\":\"CQ0_Test\",\"operatorKey\":\"123456\"},\"main\":{\"policySort\":\"CQ0\",\"relationFlag\":\"1\",\"classCode\":\"05\",\"businessNature\":\"1\",\"agentCode\":\"500102199010215230\",\"startDate\":\"2017-09-25\",\"startHour\":\"0\",\"endDate\":\"2017-09-25\",\"endHour\":\"24\",\"inputDate\":\"2017-09-25\",\"operateDate\":\"2017-09-25\",\"cityAreaCode\":\"500100\",\"comCode\":\"5001028001\",\"handler1Code\":\"500102199010215230\",\"ipAddress\":\"113.204.136.118\"},\"applicant\":{\"appliCode\":\"1\",\"appliName\":\"鹏诚保险代理有限公司\",\"insuredType\":\"1\",\"insuredNature\":\"3\",\"identifyType\":\"07\",\"identifyNumber\":\"915001035590277575\"},\"insuredDataArr\":[{\"insuredCode\":\"3\",\"insuredNature\":\"3\",\"insuredType\":\"7\",\"identifyType\":\"01\",\"identifyNumber\":\"511202197512286847\",\"insuredFlag\":\"5\"}],\"carInfo\":{\"actualValue\":\"71400.00\",\"areaCode\":\"04\",\"areaName\":\"中国境内(不含港澳台)\",\"standardName\":\"长安SC7169B轿车\",\"carInsureRelation\":\"1\",\"carOwnerNature\":\"7\",\"carOwnerIdentifyType\":\"01\",\"carOwnerIdentifyNumber\":\"511202197512286847\",\"carKindCode\":\"A0\",\"carType\":\"K33\",\"colorCode\":\"99\",\"completeKerbMass\":1323,\"engineNo\":\"13983505598\",\"frameNo\":\"LS5A2ABE9DA201695\",\"licenseNo\":\"渝F389G1\",\"enrollDate\":\"2013-09-26\",\"exhaustScale\":1.598,\"importFlag\":\"B\",\"licenseType\":\"02\",\"licenseColorCode\":\"99\",\"purchasePrice\":71400,\"modelCode\":\"CAADMD0004\",\"runMileRate\":\"0\",\"runMilers\":\"0\",\"seatCount\":\"5\",\"tonCount\":\"0.0\",\"useNatureCode\":\"8A\",\"useYears\":\"3\",\"vin\":\"LS5A2ABE9DA201695\",\"wholeWeight\":\"0.0\",\"fuleType\":\"0\",\"vehicleStyleDesc\":\"长安SC7169B轿车\",\"chgowerFlag\":0,\"platmodelCode\":\"CAADMD0004\",\"platmodelname\":\"长安SC7169B轿车\",\"carBuyDate\":\"2013-09-26\",\"fairMarketValue\":\"\",\"carPriceType\":\"1\",\"projectCode\":\"1\"},\"carShipTaxInfo\":{\"taxFlag\":\"1N\",\"taxPayerName\":\"王建华\"},\"combosDataArr\":[{\"serialNo\":\"1\",\"riskCode\":\"0507\",\"itemKindArr\":[{\"amount\":\"0.0\",\"startDate\":\"2017-09-27\",\"endDate\":\"2018-09-27\",\"deductableFlag\":0,\"kindName\":\"机动车交通事故责任强制保险\",\"kindCode\":\"BZ\",\"riskCode\":\"0507\"}]},{\"serialNo\":\"2\",\"riskCode\":\"0511\",\"itemKindArr\":[{\"amount\":\"71400.00\",\"startDate\":\"2017-09-27\",\"endDate\":\"2018-09-27\",\"deductableFlag\":1,\"kindName\":\"机动车损失保险\",\"kindCode\":\"001\",\"riskCode\":\"0511\"},{\"amount\":\"0.0\",\"startDate\":\"2017-09-27\",\"endDate\":\"2018-09-27\",\"deductableFlag\":\"0\",\"kindName\":\"不计免赔率险（机动车损失保险）\",\"kindCode\":\"301\",\"riskCode\":\"0511\"},{\"amount\":\"500000\",\"startDate\":\"2017-09-27\",\"endDate\":\"2018-09-27\",\"deductableFlag\":1,\"kindName\":\"第三者责任保险\",\"kindCode\":\"002\",\"riskCode\":\"0511\"},{\"amount\":\"0.0\",\"startDate\":\"2017-09-27\",\"endDate\":\"2018-09-27\",\"deductableFlag\":\"0\",\"kindName\":\"不计免赔率险（第三者责任保险）\",\"kindCode\":\"302\",\"riskCode\":\"0511\"},{\"amount\":\"10000\",\"startDate\":\"2017-09-27\",\"endDate\":\"2018-09-27\",\"deductableFlag\":1,\"kindName\":\"车上人员责任保险(驾驶人)\",\"kindCode\":\"003\",\"riskCode\":\"0511\"},{\"amount\":\"0.0\",\"startDate\":\"2017-09-27\",\"endDate\":\"2018-09-27\",\"deductableFlag\":\"0\",\"kindName\":\"不计免赔率险（车上人员责任保险-驾驶人）\",\"kindCode\":\"303\",\"riskCode\":\"0511\"},{\"amount\":40000,\"startDate\":\"2017-09-27\",\"endDate\":\"2018-09-27\",\"deductableFlag\":1,\"kindName\":\"车上人员责任保险(乘客)\",\"quantity\":4,\"kindCode\":\"006\",\"riskCode\":\"0511\"},{\"amount\":\"0.0\",\"startDate\":\"2017-09-27\",\"endDate\":\"2018-09-27\",\"deductableFlag\":\"0\",\"kindName\":\"不计免赔率险（车上人员责任保险-乘客）\",\"kindCode\":\"305\",\"riskCode\":\"0511\"},{\"amount\":\"0.0\",\"startDate\":\"2017-09-27\",\"endDate\":\"2018-09-27\",\"deductableFlag\":1,\"kindName\":\"发动机涉水损失险\",\"kindCode\":\"206\",\"riskCode\":\"0511\"},{\"amount\":\"0.0\",\"startDate\":\"2017-09-27\",\"endDate\":\"2018-09-27\",\"deductableFlag\":\"0\",\"kindName\":\"不计免赔率险（发动机涉水损失险）\",\"kindCode\":\"310\",\"riskCode\":\"0511\"}]}],\"bzRelationMain\":{\"startDate_bz\":\"2018-09-27\",\"endDate_bz\":\"2018-09-27\",\"startHour_bz\":\"0\",\"endHour_bz\":\"24\"}}"
	data := utils.M{}
	data.ParseJsonString(str)
	return data
}
