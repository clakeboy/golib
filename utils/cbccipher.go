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
	AES_CFB = "cfb"
)

type AesEncrypt struct {
	key        []byte
	stringType string
	aesType    string
	iv         []byte
	isBase64   bool
	isPkcs     bool
}

// 创建加密
func NewAes(k string) *AesEncrypt {
	enc := &AesEncrypt{aesType: AES_CBC, isBase64: true, isPkcs: true}
	enc.SetKey(k)
	return enc
}

// 设置是否自动base64返回
func (a *AesEncrypt) SetBase64(chk bool) {
	a.isBase64 = chk
}

// 设置加密类型
func (a *AesEncrypt) SetType(t string) {
	a.aesType = t
}

// 设置加密KEY
func (a *AesEncrypt) SetKey(k string) {
	keyLen := len(k)
	if keyLen < 16 {
		k = StrPad(k, "c", 16, STR_PAD_RIGHT)
	}
	a.key = []byte(k)
}
func (a *AesEncrypt) SetKeyBytes(k []byte) {
	a.key = k
}

// 设置向量IV
func (a *AesEncrypt) SetIV(iv []byte) {
	a.iv = iv
}

// 设置是否启动pkcs
func (a *AesEncrypt) SetPkcs(flag bool) {
	a.isPkcs = flag
}

// 得到加密KEY
func (a *AesEncrypt) GetKey() []byte {
	keyLen := len(a.key)
	//if keyLen < 16 {
	//	a.key = StrPad(a.key, "c", 16, STR_PAD_RIGHT)
	//}
	//arrKey := []byte(a.key)
	if keyLen >= 32 {
		//取前32个字节
		return a.key[:32]
	}
	if keyLen >= 24 {
		//取前24个字节
		return a.key[:24]
	}
	//取前16个字节
	return a.key[:16]
}

// 加密方法
func (a *AesEncrypt) Encrypt(plantText []byte) ([]byte, error) {
	key := a.GetKey()
	block, err := aes.NewCipher(key) //选择加密算法
	if err != nil {
		return nil, err
	}
	if a.isPkcs {
		plantText = PKCS7Padding(plantText, block.BlockSize())
	}

	var blockModel cipher.BlockMode
	var blockSteam cipher.Stream
	switch a.aesType {
	case AES_CBC:
		iv := YN(a.iv == nil, key[:block.BlockSize()], a.iv).([]byte)
		blockModel = cipher.NewCBCEncrypter(block, iv)
	case AES_CFB:
		iv := YN(a.iv == nil, key[:block.BlockSize()], a.iv).([]byte)
		blockSteam = cipher.NewCFBEncrypter(block, iv)
	case AES_ECB:
		blockModel = NewECBEncrypter(block)
	}

	ciphertext := make([]byte, len(plantText))
	switch a.aesType {
	case AES_CFB:
		blockSteam.XORKeyStream(ciphertext, plantText)
	default:
		blockModel.CryptBlocks(ciphertext, plantText)
	}

	if !a.isBase64 {
		return ciphertext, nil
	}
	buf := make([]byte, base64.StdEncoding.EncodedLen(len(ciphertext)))
	base64.StdEncoding.Encode(buf, ciphertext)
	return buf, nil
}

// 加密返回字符串
func (a *AesEncrypt) EncryptString(plantText string) (string, error) {
	res, err := a.Encrypt([]byte(plantText))
	if err != nil {
		return "", err
	}

	return string(res), nil
}

// 解密方法
func (a *AesEncrypt) Decrypt(deStr []byte) ([]byte, error) {
	key := a.GetKey()
	var cipherText []byte
	var txtLen int
	var err error
	if a.isBase64 {
		cipherText = make([]byte, base64.StdEncoding.DecodedLen(len(deStr)))
		txtLen, err = base64.StdEncoding.Decode(cipherText, deStr)
		if err != nil {
			return nil, err
		}
	} else {
		txtLen = len(deStr)
		cipherText = deStr
	}
	block, err := aes.NewCipher(key) //选择加密算法
	if err != nil {
		return nil, err
	}
	var blockModel cipher.BlockMode
	var blockSteam cipher.Stream
	switch a.aesType {
	case AES_ECB:
		blockModel = NewECBDecrypter(block)
	case AES_CBC:
		iv := YN(a.iv == nil, key[:block.BlockSize()], a.iv).([]byte)
		blockModel = cipher.NewCBCDecrypter(block, iv)
	case AES_CFB:
		iv := YN(a.iv == nil, key[:block.BlockSize()], a.iv).([]byte)
		blockSteam = cipher.NewCFBDecrypter(block, iv)
	}

	plantText := make([]byte, len(cipherText[:txtLen]))

	switch a.aesType {
	case AES_CFB:
		blockSteam.XORKeyStream(plantText, cipherText[:txtLen])
	default:
		blockModel.CryptBlocks(plantText, cipherText[:txtLen])
	}

	if a.isPkcs {
		plantText = PKCS7UnPadding(plantText)
	}
	return plantText, nil
}

// 解密字符串
func (a *AesEncrypt) DecryptString(deStr string) (string, error) {
	res, err := a.Decrypt([]byte(deStr))
	if err != nil {
		return "", err
	}

	return string(res), nil
}

// 加密返回BASE64 URL
func (a *AesEncrypt) EncryptUrl(plantText []byte) ([]byte, error) {
	res, err := a.Encrypt(plantText)
	if err != nil {
		return nil, err
	}

	cipherText, err := base64.StdEncoding.DecodeString(string(res))
	if err != nil {
		return nil, err
	}
	buf := make([]byte, base64.URLEncoding.EncodedLen(len(cipherText)))
	base64.URLEncoding.Encode(buf, cipherText)
	return buf, nil
}

// 解密BASE64URL的密码
func (a *AesEncrypt) DecryptUrl(plantText []byte) ([]byte, error) {
	cipherText, err := base64.URLEncoding.DecodeString(string(plantText))
	if err != nil {
		return nil, err
	}
	buf := make([]byte, base64.StdEncoding.EncodedLen(len(cipherText)))
	base64.StdEncoding.Encode(buf, cipherText)

	return a.Decrypt(buf)
}

// PKCS7 处理
func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// PKCS7 反解
func PKCS7UnPadding(plantText []byte) []byte {
	length := len(plantText)
	unpadding := int(plantText[length-1])
	// println(length, (length - unpadding))
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
