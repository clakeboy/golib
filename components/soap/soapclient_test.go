package soap

import (
	"ck_go_lib/utils"
	"fmt"
	"regexp"
	"testing"
)

func TestNewSoapClient(t *testing.T) {
	soap_client, err := NewSoapClient("http://123.147.190.130:28081/vchl-channel/services/accessService?wsdl")
	if err != nil {
		panic(err)
	}

	//no := utils.RandStr(10,nil)

	//xml_con := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?><Packet><Head><EChannelID>CHPC</EChannelID><RequestType>Q00</RequestType><SerialNo>%s</SerialNo></Head><Body><LicensePlateNo>渝B358U7</LicensePlateNo><VIN>LVSHJCAC3FE215302</VIN><AreaCode>50</AreaCode></Body></Packet>`,no)
	xml_con := `<?xml version="1.0" encoding="UTF-8"?>
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
