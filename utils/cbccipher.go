package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"strconv"
	"strings"
)

const (
	AES_CBC = "cbc"
	AES_ECB = "ecb"
)

type AesEncrypt struct {
	key        string
	stringType string
	aesType    string
	iv         []byte
}

//创建加密
func NewAes(k string) *AesEncrypt {
	enc := &AesEncrypt{aesType: AES_CBC}
	enc.SetKey(k)
	return enc
}

//设置加密类型
func (a *AesEncrypt) SetType(t string) {
	a.aesType = t
}

//设置加密KEY
func (a *AesEncrypt) SetKey(k string) {
	keyLen := len(k)
	if keyLen < 16 {
		k = StrPad(k, "c", 16, STR_PAD_RIGHT)
	}
	a.key = k
}

//设置向量IV
func (a *AesEncrypt) SetIV(iv []byte) {
	a.iv = iv
}

//得到加密KEY
func (a *AesEncrypt) GetKey() []byte {
	keyLen := len(a.key)
	if keyLen < 16 {
		a.key = StrPad(a.key, "c", 16, STR_PAD_RIGHT)
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

//加密方法
func (a *AesEncrypt) Encrypt(plantText []byte) ([]byte, error) {
	key := a.GetKey()
	block, err := aes.NewCipher(key) //选择加密算法
	if err != nil {
		return nil, err
	}

	plantText = PKCS7Padding(plantText, block.BlockSize())

	var blockModel cipher.BlockMode

	switch a.aesType {
	case AES_CBC:
		iv := YN(a.iv == nil, key[:block.BlockSize()], a.iv).([]byte)
		blockModel = cipher.NewCBCEncrypter(block, iv)
	case AES_ECB:
		blockModel = NewECBEncrypter(block)
	}

	ciphertext := make([]byte, len(plantText))

	blockModel.CryptBlocks(ciphertext, plantText)
	//return base64.StdEncoding.EncodeToString(ciphertext), nil
	//fmt.Println(hex.EncodeToString(ciphertext))
	buf := make([]byte, base64.StdEncoding.EncodedLen(len(ciphertext)))
	base64.StdEncoding.Encode(buf, ciphertext)
	return buf, nil
}

//加密返回字符串
func (a *AesEncrypt) EncryptString(plantText string) (string, error) {
	res, err := a.Encrypt([]byte(plantText))
	if err != nil {
		return "", err
	}

	return string(res), nil
}

//解密方法
func (a *AesEncrypt) Decrypt(deStr []byte) ([]byte, error) {
	key := a.GetKey()
	ciphertext := make([]byte, base64.StdEncoding.DecodedLen(len(deStr)))
	n, err := base64.StdEncoding.Decode(ciphertext, deStr)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key) //选择加密算法
	if err != nil {
		return nil, err
	}

	var blockModel cipher.BlockMode

	switch a.aesType {
	case AES_ECB:
		blockModel = NewECBDecrypter(block)
	case AES_CBC:
		iv := YN(a.iv == nil, key[:block.BlockSize()], a.iv).([]byte)
		blockModel = cipher.NewCBCDecrypter(block, iv)
	}

	plantText := make([]byte, len(ciphertext[:n]))
	blockModel.CryptBlocks(plantText, ciphertext[:n])
	plantText = PKCS7UnPadding(plantText)
	return plantText, nil
}

//解密字符串
func (a *AesEncrypt) DecryptString(deStr string) (string, error) {
	res, err := a.Decrypt([]byte(deStr))
	if err != nil {
		return "", err
	}

	return string(res), nil
}

//PKCS7 处理
func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

//PKCS7 反解
func PKCS7UnPadding(plantText []byte) []byte {
	length := len(plantText)
	unpadding := int(plantText[length-1])
	return plantText[:(length - unpadding)]
}

func (a *AesEncrypt) convert(b []byte) string {
	s := make([]string, len(b))
	for i := range b {
		s[i] = strconv.Itoa(int(b[i]))
	}
	return strings.Join(s, "")
}

/**
 * ECB 加密算法
 */
type ecb struct {
	b         cipher.Block
	blockSize int
}

type ecbDecrypter ecb

func NewECBEncrypter(b cipher.Block) cipher.BlockMode {
	return &ecb{
		b:         b,
		blockSize: b.BlockSize(),
	}
}

func (e *ecb) BlockSize() int { return e.blockSize }

func (e *ecb) CryptBlocks(dst, src []byte) {
	//分组分块加密
	for index := 0; index < len(src); index += e.blockSize {
		e.b.Encrypt(dst[index:index+e.blockSize], src[index:index+e.blockSize])
	}
}

func NewECBDecrypter(b cipher.Block) cipher.BlockMode {
	return &ecbDecrypter{
		b:         b,
		blockSize: b.BlockSize(),
	}
}

func (e *ecbDecrypter) BlockSize() int { return e.blockSize }

func (e *ecbDecrypter) CryptBlocks(dst, src []byte) {
	//分组分块加密
	for index := 0; index < len(src); index += e.blockSize {
		e.b.Decrypt(dst[index:index+e.blockSize], src[index:index+e.blockSize])
	}
}
