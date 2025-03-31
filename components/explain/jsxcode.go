package explain

import (
	"bytes"
	"fmt"
	"github.com/clakeboy/golib/utils"
	"io"
	"regexp"
	"strings"
)

type JSXNode struct {
	Original   []byte
	Len        int
	StartIndex int
	EndIndex   int

	NodeName   string
	EndNode    bool
	SingleNode bool

	properties []*JSXNodeProperty
}

type JSXNodeProperty struct {
	Len        int
	StartIndex int
	EndIndex   int

	Name  string
	Value string
}

var regNodeName = regexp.MustCompile(`<(/?)([\w.]+)([^>]*?)(/?)>`)
var regPropertyName = regexp.MustCompile(`\s+(\w+)=`)

func (s *JSXNode) String() string {
	label :=
		`Node:    %s
Len:     %d
Start:   %d
End:     %d
Name:    %s
EndNode: %v
Single:  %v
Property:
%v`
	return fmt.Sprintf(label,
		string(s.Original),
		s.Len,
		s.StartIndex,
		s.EndIndex,
		s.NodeName,
		s.EndNode,
		s.SingleNode,
		s.properties)
}

func (s *JSXNode) Compile() string {
	var list []string
	if s.EndNode {
		return fmt.Sprintf("</%s>", s.NodeName)
	}
	for _, v := range s.properties {
		list = append(list, fmt.Sprintf("%s=%s", v.Name, v.Value))
	}
	return fmt.Sprintf("<%s %s%s>",
		s.NodeName,
		strings.Join(list, " "),
		utils.YN(s.SingleNode, "/", ""))
}

func (s *JSXNodeProperty) String() string {
	label := "Name:%s--Value:%s\n"
	return fmt.Sprintf(label, s.Name, s.Value)
}

//explain XML or HTML code
func JSX(content []byte) []*JSXNode {
	buf := bytes.NewBuffer(content)
	start, _ := buf.ReadBytes('<')
	var nodeList []*JSXNode
	idx := len(start)
	for {
		node, err := nextNode(idx, buf)
		if err != nil {
			panic(err)
		}
		nodeList = append(nodeList, node)
		cache, err := buf.ReadBytes('<')
		if err != nil {
			break
		}
		idx = node.EndIndex + len(cache) - 1
	}

	return nodeList
}

func JSXString(content string) []*JSXNode {
	return JSX([]byte(content))
}

//explain property
func JSXProperty(rd io.Reader) {

}

func nextNode(startIndex int, bf *bytes.Buffer) (*JSXNode, error) {
	var buf bytes.Buffer
	//_,err := bf.ReadBytes('<')
	//if err != nil {
	//	return nil,err
	//}
	buf.WriteByte('<')
	flag := 1
	for {
		cache, err := bf.ReadBytes('>')
		if err != nil {
			return nil, err
		}
		buf.Write(cache)
		if bytes.Equal(cache[len(cache)-2:], []byte("=>")) {
			continue
		}
		flag -= 1
		idx := regNodeName.FindAllIndex(cache, -1)
		flag += len(idx)
		if flag == 0 {
			break
		}
	}

	rs := regNodeName.FindSubmatch(buf.Bytes())

	node := &JSXNode{
		Original:   buf.Bytes(),
		Len:        buf.Len(),
		StartIndex: startIndex,
		EndIndex:   startIndex + buf.Len(),
		NodeName:   utils.YN(len(rs) == 0, "", string(rs[2])).(string),
		EndNode:    utils.YN(len(rs[1]) == 0, false, true).(bool),
		SingleNode: utils.YN(len(rs[4]) == 0, false, true).(bool),
	}

	allIndex := regPropertyName.FindAllSubmatchIndex(node.Original, -1)

	if len(allIndex) > 0 {
		var propertyList []*JSXNodeProperty
		pptBuf := bytes.NewBuffer(node.Original)
		for _, v := range allIndex {
			fmt.Println(allIndex)
			//fmt.Println(string(bytes.TrimSpace(buf.Bytes()[v[0]:v[1]])))
			property, err := nextProperty(v, pptBuf, node.Original)
			if err != nil {
				fmt.Println(node.NodeName)
				panic(err)
			}
			propertyList = append(propertyList, property)
		}
		node.properties = propertyList
	}

	return node, nil
}

func nextProperty(idx []int, srcBuf *bytes.Buffer, src []byte) (*JSXNodeProperty, error) {
	name := src[idx[2]:idx[3]]
	var buf bytes.Buffer
	_, err := srcBuf.ReadBytes('{')
	if err != nil {
		return nil, err
	}
	buf.WriteByte('{')
	flag := 1
	for {
		cache, err := srcBuf.ReadBytes('}')
		if err != nil {
			return nil, err
		}
		flag -= 1
		buf.Write(cache)
		count := bytes.Count(cache, []byte("{"))
		flag += count
		if flag == 0 {
			break
		}
	}
	property := &JSXNodeProperty{
		Name:  string(name),
		Value: buf.String(),
	}
	return property, nil
}
