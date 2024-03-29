package utils

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
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

func TestEncryptPrng(t *testing.T) {
	demo := `F9F6B74DCBB0075DDAC9CEEF9CADB96278F408D5635A9C5D17D7C714A3F771A4D368D12732B36FE716B99FDABA594CCA012EE791F8E5D8A60B7F7518D7E6EECAC64C80E3F7700BD22FDA1150626CA9201DEBC82A1C405E48124BFE7C5618DEE9546B134718ED41E069BC66904B873E552FDA8EEC1D4032644A111BA0763C655D74DFC19E9E07CC65CED476206FD34D59AA868FCC68DA7028714F7D5D3132F3D83C35B97549297D3E78C9357F5E899770162718CA2AD7D062F8A812E3F514E911C5C2BD2EB524661EA4C15D3B913772DB0BA37F0A964337BA27FB8AEC5D6294D5482C435A82CE9B51AD20501CF03BB4A26847B5138E397DB1E5B0CD0B4782E5586786CFDD322F179A3987B3FAB00FE30B36063453C6F17A7A6ED681AC2EDBA7EFBECB97857CF7CC3D0F44558D09A86B9B9E968857A9637BDDD0D405FCDD4F0703DCD6A85C26E65710622CF3ECA8CD178DD1ED7A4B042754CDC6465511C2C0E37C245A302E910F1D3144913EEBC410E6C6E9DD953B9066C5688E59A38BEF430DB0EA157E7ABFB150CC5DE54CDE33612E65BD34C1E20E424FC46F29FADE2A687429BFBDF3334EBA6698E071C9203E0DD2B1754AA50B1D658322864ABEE53AB2815419D236A5AA6394AA6EC9D3FAE81FBFCE26E78C41A910FFB4B3F25074EF6E03548F2573DD7D271D140B3EBF7ECAD9C87E51D0103756F74F61F4DC4974A2A4DCF93568AE814E0C807E4E4A8508AC0B3D06DA60916BCA0496385DAE010A9DE9496C0B8AD03BDC695EF84251BA9C64BE6CB2411AFE627662E570EAB56E1335337E1ABD5F3600B271D9A5AA5FCC82475A8DABAD459C3A6FDED6637BE04D76D4C89A2EF1B30ECFF420C895C970A213CB4EA8383040FDD816E611F10CA6EFC633253DC6FFBE40EE1BBAE3037FCB956414E0426818C748807B15407BE77DB791C2F2BC0345C6F40BAB1A186D264648BB2143C60F7B64290F16757B0B6639912245A4108466FE21EB9B5F401C187F1E462F541BEDE41823A609FCAFE926B7FDFAA08B8C6663DA1455098CF0B1E0BA96FECE4440D688520D5786BB61A13E7ABCCA685FE0D2BC8C43AC6A28F3418821C55B92DF308B0B7CA1CCB91B6C1AF5AE6A05863EC668FA40503DA5F147F4E9A3A1D4694517B574C48B0173D8C241B933CC7A7AE5C09502AC06C9C449D2F9268AFC0AA05B1E00F13B7396D67BB8F6EB0AAEC0B15F6EA4396E6ED6330D582DAF1114FBA331B199D6C56CD6D945630015BE76A5ABC122F243423DC53DB2CB585F0070BE6D4FDCA695EA2B93B8D41677F518053458E89421B0432DD70A0F448A3FCAE1A9D6D70679F951169D6642C7454B79476EDBC5BD0B4AB845724D8A091B93E3C5A535991191`
	org := []byte("11234afee4中方")
	//encode
	enStr, err := encrypt(org)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println("加密串:", enStr)
	//decode
	deStr, err := decrypt(demo)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println("解密串:", string(deStr))
}

func encrypt(in []byte) (string, error) {
	gb, _ := UTF82GBK(in)
	buf := make([]byte, base64.StdEncoding.EncodedLen(len(gb)))
	base64.StdEncoding.Encode(buf, gb)
	key := AesKeySecureRandom([]byte("MTX123456789"))
	aes := NewAes(string(key))
	aes.SetType(AES_ECB)
	aes.SetBase64(false)
	enstr, err := aes.Encrypt(buf)
	if err != nil {
		return "", err
	}
	enHexStr := strings.ToUpper(hex.EncodeToString(enstr))
	return enHexStr, nil
}

func decrypt(de string) ([]byte, error) {
	deHexStr, err := hex.DecodeString(de)
	if err != nil {
		return nil, err
	}
	key := AesKeySecureRandom([]byte("MTX123456789"))
	aes := NewAes(string(key))
	aes.SetType(AES_ECB)
	aes.SetBase64(false)
	decryptStr, err := aes.Decrypt(deHexStr)
	if err != nil {
		return nil, err
	}
	decryptText := make([]byte, base64.StdEncoding.DecodedLen(len(decryptStr)))
	n, err := base64.StdEncoding.Decode(decryptText, decryptStr)
	if err != nil {
		return nil, err
	}
	ut, err := GBK2UTF8(decryptText[:n])
	if err != nil {
		return nil, err
	}
	return ut, nil
}

func TestNonePkcs(t *testing.T) {
	key := EncodeMD5Std("mOzTwjD1o6Q0OhLu")
	// key := "mOzTwjD1o6Q0OhLu"
	// h := md5.New()
	// h.Write([]byte(key))
	// md5Char := h.Sum(nil)
	println(len([]byte(key)))
	str := "zkUpAGQlqyOzR+t8JtFrow=="
	aes := NewAes(key)
	// aes.SetKeyBytes(md5Char)
	// aes.SetBase64(false)
	// aes.SetPkcs(false)
	aes.SetType(AES_ECB)
	deStr, err := aes.DecryptString(str)
	println(deStr, err)

	enOrgStr := `Hello`

	// aes.SetPkcs(true)
	enStr, err := aes.EncryptString(enOrgStr)
	println(enStr, err)
}
