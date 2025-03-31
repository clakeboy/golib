package utils

import (
	"crypto/sha1"
)

func SHAPRNGEncode(plantText []byte, kenCode []byte) ([]byte, error) {
	key := AesKeySecureRandom(kenCode)
	aesen := &AesEncrypt{aesType: AES_ECB}
	aesen.SetKey(string(key))
	//block, err := aes.NewCipher(key) //选择加密算法
	//if err != nil {
	//	return nil, err
	//}
	//cipherText := make([]byte, len(plantText))
	////blockModel := NewECBEncrypter(block)
	////blockModel.CryptBlocks(cipherText,plantText)
	//block.Encrypt(cipherText,plantText)
	//buf := make([]byte, base64.StdEncoding.EncodedLen(len(cipherText)))
	//base64.StdEncoding.Encode(buf, cipherText)
	//return buf, nil
	return aesen.Encrypt(plantText)
}

func SHAPRNGEncodeString(encryptStr string, kenCode string) (string, error) {
	data, err := SHAPRNGEncode([]byte(encryptStr), []byte(kenCode))
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func SHAPRNGDecode(deText []byte, kenCode []byte) ([]byte, error) {
	//plantText := make([]byte, base64.StdEncoding.DecodedLen(len(deText)))
	//n,err := base64.StdEncoding.Decode(plantText,deText)
	//if err != nil {
	//	return nil, err
	//}
	//key := AesKeySecureRandom(kenCode)
	//block, err := aes.NewCipher(key) //选择加密算法
	//if err != nil {
	//	return nil, err
	//}
	//cipherText := make([]byte, len(plantText[:n]))
	////blockModel := NewECBDecrypter(block)
	////blockModel.CryptBlocks(cipherText,plantText[:n])
	//block.Decrypt(cipherText,plantText[:n])
	//return cipherText, nil
	//---------------------------
	key := AesKeySecureRandom(kenCode)
	aesen := &AesEncrypt{aesType: AES_ECB, key: key}
	aesen.SetPkcs(false)
	return aesen.Decrypt(deText)
}

func SHAPRNGDecodeString(decryptStr string, kenCode string) (string, error) {
	data, err := SHAPRNGDecode([]byte(decryptStr), []byte(kenCode))
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func SHA1(data []byte) []byte {
	h := sha1.New()
	h.Write(data)
	return h.Sum(nil)
}

func AesKeySecureRandom(keyword []byte) (key []byte) {
	hashs := SHA1(SHA1(keyword))
	key = hashs[0:16]
	return key
}
