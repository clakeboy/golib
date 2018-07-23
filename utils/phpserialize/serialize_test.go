package phpserialize

import (
	"testing"
	"fmt"
)

func TestNewSerializer(t *testing.T) {
	data := PhpArray{
		"asdfasdf":"123",
		"asddfasdf":"111",
	}

	str ,err:= NewSerializer().Encode(data)
	fmt.Println(err)
	fmt.Println(str)
}
