package utils

import (
	"testing"
	"encoding/xml"
	"fmt"
	"strings"
	"log"
)

func TestXMLMap_UnmarshalXML(t *testing.T) {
	//xml_con := `<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/"><soap:Body><ns1:vchlRequestResponse xmlns:ns1="http://access.vchl.echannel.acic.com/"><return>&lt;?xml version=&quot;1.0&quot; encoding=&quot;UTF-8&quot;?&gt;&lt;Packet&gt;  &lt;Head&gt;    &lt;RequestType&gt;Q00&lt;/RequestType&gt;    &lt;SerialNo&gt;9f2A7nhnlR&lt;/SerialNo&gt;    &lt;ResponseCode&gt;-1&lt;/ResponseCode&gt;    &lt;ErrorCode&gt;1201&lt;/ErrorCode&gt;    &lt;ErrorMessage&gt;报文摘要信息不匹配&lt;/ErrorMessage&gt;  &lt;/Head&gt;&lt;/Packet&gt;</return></ns1:vchlRequestResponse></soap:Body></soap:Envelope>`
	xml_con := `<Person>
        <FullName>Grace R. Emlin</FullName>
        <Company>Example Inc.</Company>
        <Email where="home">
            <Addr>gre@example.com</Addr>
        </Email>
        <Email where='work'>
            <Addr>gre@work.com</Addr>
        </Email>
        <Group>
            <Value>Friends</Value>
            <Value>Squash</Value>
        </Group>
        <City>Hanga Roa</City>
        <State>Easter Island</State>
    </Person>`
	//xml_map := XMLMap{}
	//
	//err := xml.Unmarshal([]byte(xml_con),&xml_map)
	//if err != nil {
	//	panic(err)
	//}

	inputReader := strings.NewReader(xml_con)
	decoder := xml.NewDecoder(inputReader)

	////xml_map := M{}
	//for tt, err = decoder.Token(); err == nil; tt, err = decoder.Token() {
	//
	//	switch token := tt.(type) {
	//	case xml.StartElement:
	//		name := token.Name.Local
	//		fmt.Printf("Token name: %s\n", name)
	//		for _, attr := range token.Attr {
	//			attrName := attr.Name.Local
	//			attrValue := attr.Value
	//			fmt.Printf("An attribute is: %s %s\n", attrName, attrValue)
	//		}
	//	case xml.EndElement:
	//		fmt.Printf("Token of '%s' end\n", token.Name.Local)
	//		// 处理字符数据（这里就是元素的文本）
	//	case xml.CharData:
	//		content := string([]byte(token))
	//		fmt.Printf("This is the content: %v\n", content)
	//	default:
	//
	//	}
	//}
	tt, err := decoder.Token()
	if err != nil {
		fmt.Println(err)
	}

	xml_map := parse_token(decoder,nil,tt)

	fmt.Println(xml_map)
}

func parse_token(decoder *xml.Decoder, current_token xml.Token, parent_token xml.Token) interface{} {
	xml_map := M{}
	root_token := parent_token.(xml.StartElement)
	var err error
	var t xml.Token
	end := false
	for t, err = decoder.Token(); err == nil; t, err = decoder.Token() {
		switch token := t.(type) {
		case xml.StartElement:
			fmt.Println("start ",token.Name.Local)
			if !end {
				end = false
				xml_map[root_token.Name.Local] = parse_token(decoder,token,token)
			}
		case xml.EndElement:
			fmt.Println("end ",token.Name.Local)
			if token.Name.Local == root_token.Name.Local {
				return xml_map
			}

			end = true
		case xml.CharData:
			if end {
				break
			}
			content := strings.Trim(string([]byte(token))," ")
			content = strings.Trim(content,"\n")
			if content != "" {
				xml_map[root_token.Name.Local] = content
				fmt.Println("content",content)
			}
		default:

		}
	}
	return xml_map
}

func TestXMLMap_UnmarshalXML2(t *testing.T) {
	str := `<xml><return_code><![CDATA[SUCCESS]]></return_code>
<return_msg><![CDATA[OK]]></return_msg>
<appid><![CDATA[wx6b1dafd90d1e2db1]]></appid>
<mch_id><![CDATA[1486942202]]></mch_id>
<device_info><![CDATA[WEB]]></device_info>
<nonce_str><![CDATA[HMYC5hq2yDYSelOi]]></nonce_str>
<sign><![CDATA[3B090F9A60FDD4F6187A72C5A18BC4D1]]></sign>
<result_code><![CDATA[SUCCESS]]></result_code>
<openid><![CDATA[owS5q0n3sxQc4l_oV5D9cxfbqGWo]]></openid>
<is_subscribe><![CDATA[N]]></is_subscribe>
<trade_type><![CDATA[JSAPI]]></trade_type>
<bank_type><![CDATA[ICBC_DEBIT]]></bank_type>
<total_fee>232601</total_fee>
<coupon_fee>37</coupon_fee>
<fee_type><![CDATA[CNY]]></fee_type>
<transaction_id><![CDATA[4200000060201801292945278230]]></transaction_id>
<out_trade_no><![CDATA[PC150E37D7F95E0C4C6C0EE8]]></out_trade_no>
<attach><![CDATA[]]></attach>
<time_end><![CDATA[20180129153858]]></time_end>
<trade_state><![CDATA[SUCCESS]]></trade_state>
<coupon_id_0><![CDATA[2000000006103513115]]></coupon_id_0>
<coupon_fee_0>37</coupon_fee_0>
<coupon_count>1</coupon_count>
<cash_fee>232564</cash_fee>
</xml>`

	data := XMLMap{}
	err := xml.Unmarshal([]byte(str),&data)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println(data)
}