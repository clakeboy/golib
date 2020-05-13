package utils

import (
	"crypto/cipher"
	"crypto/des"
	"encoding/base64"
)

type DesCipher struct {
	key     []byte
	DesType string
	iv      []byte
}

func NewDesCipher(key []byte) *DesCipher {
	return &DesCipher{
		key:     key[:8],
		DesType: AES_ECB,
		iv:      nil,
	}
}

func (d *DesCipher) SetKey(key []byte) {
	d.key = key[:8]
}

func (d *DesCipher) SetIV(iv []byte) {
	d.iv = iv
}

func (d *DesCipher) Encrypt(data []byte) ([]byte, error) {
	block, err := des.NewCipher(d.key) //选择加密算法
	if err != nil {
		return nil, err
	}

	data = PKCS7Padding(data, block.BlockSize())

	var blockModel cipher.BlockMode

	switch d.DesType {
	case AES_CBC:
		iv := YN(d.iv == nil, d.key[:block.BlockSize()], d.iv).([]byte)
		blockModel = cipher.NewCBCEncrypter(block, iv)
	case AES_ECB:
		blockModel = NewECBEncrypter(block)
	default:
		blockModel = NewECBEncrypter(block)
	}
	cipherData := make([]byte, len(data))
	blockModel.CryptBlocks(cipherData, data)
	buf := make([]byte, base64.StdEncoding.EncodedLen(len(cipherData)))
	base64.StdEncoding.Encode(buf, cipherData)
	return buf, nil
}

func (d *DesCipher) EncryptString(data string) (string, error) {
	res, err := d.Encrypt([]byte(data))
	if err != nil {
		return "", err
	}

	return string(res), nil
}

func (d *DesCipher) Decrypt(data []byte) ([]byte, error) {
	cipherData := make([]byte, base64.StdEncoding.DecodedLen(len(data)))
	n, err := base64.StdEncoding.Decode(cipherData, data)
	if err != nil {
		return nil, err
	}
	block, err := des.NewCipher(d.key) //选择加密算法
	if err != nil {
		return nil, err
	}

	var blockModel cipher.BlockMode

	switch d.DesType {
	case AES_ECB:
		blockModel = NewECBDecrypter(block)
	case AES_CBC:
		iv := YN(d.iv == nil, d.key[:block.BlockSize()], d.iv).([]byte)
		blockModel = cipher.NewCBCDecrypter(block, iv)
	default:
		blockModel = NewECBDecrypter(block)
	}

	plantText := make([]byte, len(cipherData[:n]))
	blockModel.CryptBlocks(plantText, cipherData[:n])
	plantText = PKCS7UnPadding(plantText)
	return plantText, nil
}

func (d *DesCipher) DecryptString(data string) (string, error) {
	res, err := d.Decrypt([]byte(data))
	if err != nil {
		return "", err
	}

	return string(res), nil
}
