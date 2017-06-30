package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"bytes"
	"encoding/base64"
	"strconv"
	"strings"
)

type AesEncrypt struct {
	key string
	stringType string
}

func NewAes(k string) *AesEncrypt{
	enc := &AesEncrypt{}
	enc.SetKey(k)
	return enc
}

func (a *AesEncrypt) SetKey(k string) {
	keyLen := len(k)
	if keyLen < 16 {
		k = StrPad(k,"c",16,STR_PAD_RIGHT)
	}
	a.key = k
}

func (a *AesEncrypt) GetKey() []byte {
	keyLen := len(a.key)
	if keyLen < 16 {
		a.key = StrPad(a.key,"c",16,STR_PAD_RIGHT)
	}
	arrKey := []byte(a.key)
	if keyLen >= 32 {
		//取前32个字节
		return arrKey[:32]
	}
	if keyLen >= 24 {
		//取前24个字节
		return arrKey[:24]
	}
	//取前16个字节
	return arrKey[:16]
}

func (a *AesEncrypt) Encrypt(plantText []byte) ([]byte, error) {
	key := a.GetKey()
	block, err := aes.NewCipher(key) //选择加密算法
	if err != nil {
		return nil, err
	}
	plantText = a.PKCS7Padding(plantText, block.BlockSize())

	blockModel := cipher.NewCBCEncrypter(block, key[:block.BlockSize()])

	ciphertext := make([]byte, len(plantText))

	blockModel.CryptBlocks(ciphertext, plantText)
	//return base64.StdEncoding.EncodeToString(ciphertext), nil
	//fmt.Println(hex.EncodeToString(ciphertext))
	buf := make([]byte, base64.StdEncoding.EncodedLen(len(ciphertext)))
	base64.StdEncoding.Encode(buf, ciphertext)
	return buf, nil
}

func (a *AesEncrypt) Decrypt(deStr []byte) ([]byte, error) {
	key := a.GetKey()
	ciphertext := make([]byte, base64.StdEncoding.DecodedLen(len(deStr)))
	n,err := base64.StdEncoding.Decode(ciphertext,deStr)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key) //选择加密算法
	if err != nil {
		return nil, err
	}
	blockModel := cipher.NewCBCDecrypter(block, key[:block.BlockSize()])
	plantText := make([]byte, len(ciphertext[:n]))
	blockModel.CryptBlocks(plantText, ciphertext[:n])
	plantText = a.PKCS7UnPadding(plantText)
	return plantText, nil
}

func (a *AesEncrypt) PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func (a *AesEncrypt) PKCS7UnPadding(plantText []byte) []byte {
	length := len(plantText)
	unpadding := int(plantText[length-1])
	return plantText[:(length - unpadding)]
}

func (a *AesEncrypt) convert( b []byte ) string {
	s := make([]string,len(b))
	for i := range b {
		s[i] = strconv.Itoa(int(b[i]))
	}
	return strings.Join(s,"")
}