package utils

import (
	"testing"
	"fmt"
)

func TestAesEncrypt_Encrypt(t *testing.T) {
	cbc := NewAes("ck-cookie")
	cipher_text,err := cbc.Encrypt([]byte("2"))
	if err != nil {
		panic(err)
	}
	fmt.Println(string(cipher_text))
	text,err := cbc.Decrypt(cipher_text)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(text))
}