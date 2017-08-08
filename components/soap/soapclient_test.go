package soap

import (
	"testing"
	"fmt"
	"ck_go_lib/utils"
)

func TestNewSoapClient(t *testing.T) {
	soap,err := NewSoapClient("http://123.147.190.130:28081/vchl-channel/services/accessService?wsdl")
	if err != nil {
		panic(err)
	}

	no := utils.RandStr(10,nil)

	xml_con := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?><Packet><Head><EChannelID>CHPC</EChannelID><RequestType>Q00</RequestType><SerialNo>%s</SerialNo></Head><Body><LicensePlateNo>渝B358U7</LicensePlateNo><VIN>LVSHJCAC3FE215302</VIN><AreaCode>50</AreaCode></Body></Packet>`,no)

	args := []interface{}{
		"CHPC",
		xml_con,
		utils.EncodeMD5(fmt.Sprintf("%s%s%s%s","113.204.136.118","CHPC","e435rfe3dwxd180e5ea7e5f145c4ccb8",xml_con)),
	}

	fmt.Println(args...)
	soap.SetAddress("http://123.147.190.130:28081/vchl-channel/services/accessService")
	soap.Call("vchlRequest",args...)
}