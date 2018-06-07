package utils

import "encoding/json"

type JsonParse struct {
	JsonObject interface{}
}

func (m *JsonParse) ToJson() []byte {
	data ,err := json.Marshal(m.JsonObject)
	if err != nil {
		return nil
	}
	return data
}

func (m *JsonParse) ToJsonString() string {
	data := m.ToJson()
	if data == nil {
		return ""
	}
	return string(data)
}

func (m *JsonParse) ParseJson(data []byte) error {
	err := json.Unmarshal(data,m.JsonObject)
	if err != nil{
		return err
	}
	return nil
}

func (m *JsonParse) ParseJsonString(data string) error {
	return m.ParseJson([]byte(data))
}