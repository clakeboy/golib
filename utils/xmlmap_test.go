package utils

import (
	"testing"
	"encoding/xml"
	"fmt"
	"strings"
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
				xml_map[root_token.Name.Local] = parse_token(decoder,token)
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
