package utils

import (
	"fmt"
	"testing"
)

func TestAesEncrypt_Encrypt(t *testing.T) {
	cbc := NewAes("ck-cookie")
	cipher_text, err := cbc.Encrypt([]byte("askdjfh3827349238^sdkfjh219222"))
	if err != nil {
		panic(err)
	}
	fmt.Println(string(cipher_text))
	text, err := cbc.Decrypt(cipher_text)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(text))

	urlText, err := cbc.EncryptUrl([]byte("askdjfh3827349238^sdkfjh219222"))
	if err != nil {
		panic(err)
	}

	fmt.Println(string(urlText))

	deUrlText, err := cbc.DecryptUrl(urlText)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(deUrlText))
}
