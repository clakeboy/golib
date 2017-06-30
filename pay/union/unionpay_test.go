package union

import (
	"testing"
	"time"
)

var cfg = &Config{
	SignCertPath:     "../../certs/700000000000001_acp.pfx",
	SignCertPassword: "000000",
	EncryptCertPath:  "../../certs/acp_test_enc.cer",
	MiddleCertPath:   "../../certs/acp_test_middle.cer",
	RootCertPath:     "../../certs/acp_test_root.cer",
}

var urls = &UrlConfig{
	FrontTransUrl:  "https://gateway.test.95516.com/gateway/api/frontTransReq.do",
	BackTransUrl:   "https://gateway.test.95516.com/gateway/api/backTransReq.do",
	SingleQueryUrl: "https://gateway.test.95516.com/gateway/api/queryTrans.do",
	BatchTransUrl:  "https://gateway.test.95516.com/gateway/api/batchTrans.do",
	FileTransUrl:   "https://filedownload.test.95516.com/",
	AppTransUrl:    "https://gateway.test.95516.com/gateway/api/appTransReq.do",
	CardTransUrl:   "https://gateway.test.95516.com/gateway/api/cardTransReq.do",
	CallbackUrl:    "http://pay.b.clake.cc:7908/pay/union/recv",
}

func TestUnionPay_BackBind(t *testing.T) {

	union_pay := NewPay(cfg, urls)

	user := &UserInfo{
		PhotoNo:    "13552535506",
		CertifId:   "341126197709218366",
		CertifTp:   "01",
		CustomerNm: "互联网",
		Cvn2:       "123",
		Expired:    "1711",
	}

	bind := &BackBind{
		MerId:   "777290058147175",
		OrderId: "TEST000000002",
		//TxnTime: time.Now().Format("20060102150405"),
		TxnTime:"20170615190823",
		AccNo:   "6221558812340000",
		BindId:  "UN000000002",
	}

	err := union_pay.BackBind(user, bind)
	if err != nil {
		panic(err)
	}
}

func TestPay_BackPay(t *testing.T) {
	union_pay := NewPay(cfg, urls)

	user := &UserInfo{
		PhotoNo:    "13552535506",
		CertifId:   "341126197709218366",
		CertifTp:   "01",
		CustomerNm: "互联网",
		Cvn2:       "123",
		Expired:    "1711",
	}

	dk := &BackDK{
		MerId:   "777290058147175",
		OrderId: "TP000000001",
		TxnTime: time.Now().Format("20060102150405"),
		TxnAmt:  "1000",
		BindId:  "UN000000001",
	}

	err := union_pay.BackPay(user,dk)
	if err != nil {
		panic(err)
	}
}
