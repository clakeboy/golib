package utils

import (
	"fmt"
	"testing"
)

func TestNewDesCipher(t *testing.T) {
	//key := AesKeySecureRandom([]byte("1234ABCD"))
	desCipher := NewDesCipher([]byte("1234ABCD"))
	res, err := desCipher.EncryptString("00010001120000000测试报文加解密！")
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println("encrypt:", res)
	deRes, err := desCipher.DecryptString(res)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println("decrypt:", deRes)
}
