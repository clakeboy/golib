package txsms

import (
	"encoding/json"
	"fmt"
	"github.com/clakeboy/golib/utils"
	"reflect"
	"testing"
)

func TestType(t *testing.T) {
	data := NewSendData(nil)
	data.Tel = []*PhoneData{
		{"12312312312", "86"},
	}

	fmt.Println(reflect.TypeOf(data.Tel).Kind())

	switch data.Tel.(type) {
	case []*PhoneData:
		fmt.Println("this is slice")
	}
}

func TestPtr(t *testing.T) {
	tel := &PhoneData{}
	str := "{\"phone\":\"34234123412\"}"

	err := json.Unmarshal([]byte(str), tel)
	if err != nil {
		t.Error(err)
	}
}

func TestTxSms_Send(t *testing.T) {
	data := NewSendData(&SendData{
		Params: []string{"test"},
		Sign:   "优服严选",
		Tel: &PhoneData{
			Mobile:     "18523510321",
			Nationcode: "86",
		},
		TplId: 428111,
	})

	txsms := NewTxSms("1400262921", "c9b81b11d12b99a91c46ee67dcf64691")
	res, err := txsms.Send(data)
	if err != nil {
		t.Error(err)
		return
	}

	utils.PrintAny(res)
}
