package union

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
)

type UMP map[string]interface{}

func (m *UMP) ToJson() []byte {
	data, err := json.Marshal(m)
	if err != nil {
		return nil
	}
	return data
}

func (m *UMP) ToJsonString() string {
	data := m.ToJson()
	if data == nil {
		return ""
	}
	return string(data)
}

func (m *UMP) ParseJson(data []byte) error {
	err := json.Unmarshal(data, m)
	if err != nil {
		return err
	}
	return nil
}

func (m *UMP) ParseJsonString(data string) error {
	return m.ParseJson([]byte(data))
}

//银联证书配置
type Config struct {
	SignCertPath     string `json:"sign_cert_path"`     //加密证书路径
	SignCertPassword string `json:"sign_cert_password"` //加密证书密码
	EncryptCertPath  string `json:"encrypt_cert_path"`  //敏感信息加密证书路径
	MiddleCertPath   string `json:"middle_cert_path"`   //验签中级证书路径
	RootCertPath     string `json:"root_cert_path"`     //验签根证书路径
	MerId            string `json:"mer_id"`             //商户ID
}

//银联调用地址配置
type UrlConfig struct {
	FrontTransUrl   string `json:"front_trans_url"`  //前台调用地址
	BackTransUrl    string `json:"back_trans_url"`   //后台调用地址
	SingleQueryUrl  string `json:"single_query_url"` //单次查询调用地址
	BatchTransUrl   string `json:"batch_trans_url"`  //
	FileTransUrl    string `json:"file_trans_url"`
	AppTransUrl     string `json:"app_trans_url"`
	CardTransUrl    string `json:"card_trans_url"`
	CallbackUrl     string `json:"callback_url"`      //后台通知调用地址
	UndoCallbackUrl string `json:"undo_callback_url"` //退款通知调用地址
}

//后台绑定用户银行卡数据
type UserInfo struct {
	PhotoNo    string `json:"phoneNo"`    //手机号
	CertifTp   string `json:"certifTp"`   //证件类型
	CertifId   string `json:"certifId"`   //证件号
	CustomerNm string `json:"customerNm"` //姓名
	Cvn2       string `json:"cvn2"`       //CVN2 码
	Expired    string `json:"expired"`    //银行卡有效期
}

//绑定用户卡数据
type BackBind struct {
	MerId   string `json:"merId"`   //商户代码
	OrderId string `json:"orderId"` //商户订单号
	TxnTime string `json:"txnTime"` //订单发送时间
	AccNo   string `json:"accNo"`   //银行卡号
	BindId  string `json:"bindId"`  //绑定标识号
}

//绑定用户扣费数据
type BackDK struct {
	MerId   string `json:"merId"`   //商户代码
	OrderId string `json:"orderId"` //商户订单号
	TxnTime string `json:"txnTime"` //订单发送时间
	TxnAmt  string `json:"txnAmt"`  //订单金额
	BindId  string `json:"bindId"`  //绑定标识号
	AccNo   string `json:"accNo"`   //银行卡号
}

//退款数据
type OrderUndo struct {
	MerId   string `json:"merId"`    //商户代码
	OrderId string `json:"orderId"`  //商户订单号
	TxnTime string `json:"txnTime"`  //订单发送时间
	TxnAmt  string `json:"txnAmt"`   //原订单扣费金额
	QueryId string `json:"query_id"` //原订单query_id
}

//查询数据
type QueryOrder struct {
	MerId   string `json:"merId"`    //商户代码
	OrderId string `json:"orderId"`  //商户订单号
	TxnTime string `json:"txnTime"`  //订单发送时间
}

//pkcs 证书信息
type SignCertInfo struct {
	Cert       *x509.Certificate
	PrivateKey *rsa.PrivateKey
}

//enc 加密证书
type EncryptCertInfo struct {
	Cert      *x509.Certificate
	PublicKey *rsa.PublicKey
}
