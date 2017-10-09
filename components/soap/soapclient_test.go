package soap

import (
	"ck_go_lib/utils"
	"fmt"
	"regexp"
	"testing"
	"log"
	"strings"
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
	fmt.Println(str,err)
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
</wsdl:definitions>
`
	soap_client, err := NewSoapClientForContent(wsdl_xml)
	if err != nil {
		panic(err)
	}
	fmt.Println(soap_client.requestUrl)

	res,err := soap_client.Call("carQuery",utils.M{
		"policySort":"CQ0",
		"licenseNo":"渝B6S919",
		"licenseType":"02",
		"cityAreaCode":"03",
		"txInsuranceRequestEhm":utils.M{
			"transExeDate":"",
			"transExeTime":"",
		},
		"txInsuranceRequestExtensionEhm":utils.M{
			"operator":"CQ0_Test",
			"operatorKey":"123456",
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

	res,err := soap_client.Call("premiumCaculate",getSoapData())
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
	client.SetHeader("Content-Type","application/soap+xml; charset=utf-8")
	client.SetHeader("Content-Length",fmt.Sprintf("%v",len(content)))
	req,err := client.Request("POST",url_str,strings.NewReader(content))
	if err != nil {
		log.Println(err)
	}

	fmt.Println(string(req.Content))
}

func getSoapData() utils.M {
	str := "{\"txInsuranceRequestEhm\":{\"transExeDate\":\"20170925\",\"transExeTime\":\"201300\"},\"txInsuranceRequestExtensionEhm\":{\"operator\":\"CQ0_Test\",\"operatorKey\":\"123456\"},\"main\":{\"policySort\":\"CQ0\",\"relationFlag\":\"1\",\"classCode\":\"05\",\"businessNature\":\"1\",\"agentCode\":\"500102199010215230\",\"startDate\":\"2017-09-25\",\"startHour\":\"0\",\"endDate\":\"2017-09-25\",\"endHour\":\"24\",\"inputDate\":\"2017-09-25\",\"operateDate\":\"2017-09-25\",\"cityAreaCode\":\"500100\",\"comCode\":\"5001028001\",\"handler1Code\":\"500102199010215230\",\"ipAddress\":\"113.204.136.118\"},\"applicant\":{\"appliCode\":\"1\",\"appliName\":\"鹏诚保险代理有限公司\",\"insuredType\":\"1\",\"insuredNature\":\"3\",\"identifyType\":\"07\",\"identifyNumber\":\"915001035590277575\"},\"insuredDataArr\":[{\"insuredCode\":\"3\",\"insuredNature\":\"3\",\"insuredType\":\"7\",\"identifyType\":\"01\",\"identifyNumber\":\"511202197512286847\",\"insuredFlag\":\"5\"}],\"carInfo\":{\"actualValue\":\"71400.00\",\"areaCode\":\"04\",\"areaName\":\"中国境内(不含港澳台)\",\"standardName\":\"长安SC7169B轿车\",\"carInsureRelation\":\"1\",\"carOwnerNature\":\"7\",\"carOwnerIdentifyType\":\"01\",\"carOwnerIdentifyNumber\":\"511202197512286847\",\"carKindCode\":\"A0\",\"carType\":\"K33\",\"colorCode\":\"99\",\"completeKerbMass\":1323,\"engineNo\":\"13983505598\",\"frameNo\":\"LS5A2ABE9DA201695\",\"licenseNo\":\"渝F389G1\",\"enrollDate\":\"2013-09-26\",\"exhaustScale\":1.598,\"importFlag\":\"B\",\"licenseType\":\"02\",\"licenseColorCode\":\"99\",\"purchasePrice\":71400,\"modelCode\":\"CAADMD0004\",\"runMileRate\":\"0\",\"runMilers\":\"0\",\"seatCount\":\"5\",\"tonCount\":\"0.0\",\"useNatureCode\":\"8A\",\"useYears\":\"3\",\"vin\":\"LS5A2ABE9DA201695\",\"wholeWeight\":\"0.0\",\"fuleType\":\"0\",\"vehicleStyleDesc\":\"长安SC7169B轿车\",\"chgowerFlag\":0,\"platmodelCode\":\"CAADMD0004\",\"platmodelname\":\"长安SC7169B轿车\",\"carBuyDate\":\"2013-09-26\",\"fairMarketValue\":\"\",\"carPriceType\":\"1\",\"projectCode\":\"1\"},\"carShipTaxInfo\":{\"taxFlag\":\"1N\",\"taxPayerName\":\"王建华\"},\"combosDataArr\":[{\"serialNo\":\"1\",\"riskCode\":\"0507\",\"itemKindArr\":[{\"amount\":\"0.0\",\"startDate\":\"2017-09-27\",\"endDate\":\"2018-09-27\",\"deductableFlag\":0,\"kindName\":\"机动车交通事故责任强制保险\",\"kindCode\":\"BZ\",\"riskCode\":\"0507\"}]},{\"serialNo\":\"2\",\"riskCode\":\"0511\",\"itemKindArr\":[{\"amount\":\"71400.00\",\"startDate\":\"2017-09-27\",\"endDate\":\"2018-09-27\",\"deductableFlag\":1,\"kindName\":\"机动车损失保险\",\"kindCode\":\"001\",\"riskCode\":\"0511\"},{\"amount\":\"0.0\",\"startDate\":\"2017-09-27\",\"endDate\":\"2018-09-27\",\"deductableFlag\":\"0\",\"kindName\":\"不计免赔率险（机动车损失保险）\",\"kindCode\":\"301\",\"riskCode\":\"0511\"},{\"amount\":\"500000\",\"startDate\":\"2017-09-27\",\"endDate\":\"2018-09-27\",\"deductableFlag\":1,\"kindName\":\"第三者责任保险\",\"kindCode\":\"002\",\"riskCode\":\"0511\"},{\"amount\":\"0.0\",\"startDate\":\"2017-09-27\",\"endDate\":\"2018-09-27\",\"deductableFlag\":\"0\",\"kindName\":\"不计免赔率险（第三者责任保险）\",\"kindCode\":\"302\",\"riskCode\":\"0511\"},{\"amount\":\"10000\",\"startDate\":\"2017-09-27\",\"endDate\":\"2018-09-27\",\"deductableFlag\":1,\"kindName\":\"车上人员责任保险(驾驶人)\",\"kindCode\":\"003\",\"riskCode\":\"0511\"},{\"amount\":\"0.0\",\"startDate\":\"2017-09-27\",\"endDate\":\"2018-09-27\",\"deductableFlag\":\"0\",\"kindName\":\"不计免赔率险（车上人员责任保险-驾驶人）\",\"kindCode\":\"303\",\"riskCode\":\"0511\"},{\"amount\":40000,\"startDate\":\"2017-09-27\",\"endDate\":\"2018-09-27\",\"deductableFlag\":1,\"kindName\":\"车上人员责任保险(乘客)\",\"quantity\":4,\"kindCode\":\"006\",\"riskCode\":\"0511\"},{\"amount\":\"0.0\",\"startDate\":\"2017-09-27\",\"endDate\":\"2018-09-27\",\"deductableFlag\":\"0\",\"kindName\":\"不计免赔率险（车上人员责任保险-乘客）\",\"kindCode\":\"305\",\"riskCode\":\"0511\"},{\"amount\":\"0.0\",\"startDate\":\"2017-09-27\",\"endDate\":\"2018-09-27\",\"deductableFlag\":1,\"kindName\":\"发动机涉水损失险\",\"kindCode\":\"206\",\"riskCode\":\"0511\"},{\"amount\":\"0.0\",\"startDate\":\"2017-09-27\",\"endDate\":\"2018-09-27\",\"deductableFlag\":\"0\",\"kindName\":\"不计免赔率险（发动机涉水损失险）\",\"kindCode\":\"310\",\"riskCode\":\"0511\"}]}],\"bzRelationMain\":{\"startDate_bz\":\"2018-09-27\",\"endDate_bz\":\"2018-09-27\",\"startHour_bz\":\"0\",\"endHour_bz\":\"24\"}}"
	data := utils.M{}
	data.ParseJsonString(str)
	return data
}