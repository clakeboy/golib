package union

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/clakeboy/golib/utils"
	"golang.org/x/crypto/pkcs12"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

//银联支付API
type Pay struct {
	conf        *Config
	urls        *UrlConfig
	callbackUrl string //支付回调地址
}

func NewPay(cfg *Config, urls *UrlConfig) *Pay {
	return &Pay{
		conf:        cfg,
		urls:        urls,
		callbackUrl: urls.CallbackUrl,
	}
}

//设置支付回调地址
func (u *Pay) SetCallbackUrl(url_str string) {
	u.callbackUrl = url_str
}

//绑定卡号
func (u *Pay) BackBind(user *UserInfo, bind *BackBind) error {
	enCrt, err := u.getEncryptCertInfo()
	if err != nil {
		return err
	}

	accno, err := u.EncryptData(bind.AccNo)
	if err != nil {
		return err
	}

	enc_custom, err := u.EncryptCustomerData(user)
	if err != nil {
		return err
	}

	post_data := UMP{
		"version":       "5.1.0",
		"encoding":      "utf-8",
		"signMethod":    "01",
		"txnType":       "72",
		"txnSubType":    "01",
		"bizType":       "000501",
		"accessType":    "0",
		"channelType":   "07",
		"encryptCertId": fmt.Sprintf("%d", enCrt.Cert.SerialNumber),
		"merId":         bind.MerId,
		"orderId":       bind.OrderId,
		"txnTime":       bind.TxnTime,
		"accNo":         accno,
		"customerInfo":  enc_custom,
		"bindId":        bind.BindId,
	}

	err = u.sign(post_data)

	if err != nil {
		return err
	}

	req, err := u.post(u.urls.BackTransUrl, post_data)
	if err != nil {
		return err
	}

	success, err := u.validate(req)
	if err != nil {
		return err
	}

	if success {
		fmt.Println("返回数据验证成功!")
	}

	fmt.Println(req["respMsg"])

	return nil
}

//代收费用
func (u *Pay) BackPay(user *UserInfo, dk *BackDK) (UMP, error) {
	enCrt, err := u.getEncryptCertInfo()
	if err != nil {
		return nil, err
	}

	accno, err := u.EncryptData(dk.AccNo)
	if err != nil {
		return nil, err
	}

	enc_custom, err := u.EncryptCustomerData(user)
	if err != nil {
		return nil, err
	}

	post_data := UMP{
		"version":       "5.1.0",
		"encoding":      "utf-8",
		"signMethod":    "01",
		"txnType":       "11",
		"txnSubType":    "00",
		"bizType":       "000501",
		"accessType":    "0",
		"channelType":   "07",
		"currencyCode":  "156",
		"backUrl":       u.callbackUrl,
		"encryptCertId": fmt.Sprintf("%d", enCrt.Cert.SerialNumber),
		"merId":         dk.MerId,
		"orderId":       dk.OrderId,
		"txnTime":       dk.TxnTime,
		"txnAmt":        dk.TxnAmt,
		"customerInfo":  enc_custom,
		"accNo":         accno,
	}

	err = u.sign(post_data)

	if err != nil {
		return nil, err
	}

	req, err := u.post(u.urls.BackTransUrl, post_data)
	if err != nil {
		return nil, err
	}

	success, err := u.validate(req)
	if err != nil {
		return nil, err
	}

	if !success {
		return nil, errors.New("数据返回验证不成功")
	}

	fmt.Println(req["respMsg"])

	return req, nil
}

//退款
func (u *Pay) UndoPay(undo *OrderUndo) (UMP, error) {
	post_data := UMP{
		"version":     "5.1.0",
		"encoding":    "utf-8",
		"signMethod":  "01",
		"txnType":     "04",
		"txnSubType":  "00",
		"bizType":     "000501",
		"accessType":  "0",
		"channelType": "07",
		"backUrl":     u.urls.UndoCallbackUrl,
		"merId":       undo.MerId,
		"orderId":     undo.OrderId,
		"txnTime":     undo.TxnTime,
		"txnAmt":      undo.TxnAmt,
		"origQryId":   undo.QueryId,
	}

	err := u.sign(post_data)

	if err != nil {
		return post_data, err
	}

	req, err := u.post(u.urls.BackTransUrl, post_data)
	if err != nil {
		return post_data, err
	}

	success, err := u.validate(req)
	if err != nil {
		return req, err
	}

	if !success {
		return req, errors.New("数据返回验证不成功")
	}

	return req, nil
}

//查询订单交易是否成功
func (u *Pay) QueryPay(query *QueryOrder) (UMP, error) {
	post_data := UMP{
		"version":    "5.1.0",
		"encoding":   "UTF-8",
		"signMethod": "01",
		"txnType":    "00",
		"txnSubType": "00",
		"bizType":    "000501",
		"accessType": "0",
		"merId":      query.MerId,
		"orderId":    query.OrderId,
		"txnTime":    query.TxnTime,
	}

	err := u.sign(post_data)

	if err != nil {
		return nil, err
	}

	req, err := u.post(u.urls.BackTransUrl, post_data)
	if err != nil {
		return nil, err
	}

	success, err := u.validate(req)
	if err != nil {
		return nil, err
	}

	if !success {
		return nil, errors.New("数据返回验证不成功")
	}

	fmt.Println(req["respMsg"])

	return req, nil
}

//签名数据
func (u *Pay) sign(data UMP) error {
	crt, err := u.getSignCertInfo()
	if err != nil {
		return err
	}

	data["certId"] = fmt.Sprintf("%d", crt.Cert.SerialNumber)

	enc_str := u.createLinkString(data, true)

	rng := rand.Reader

	hashed := sha256.Sum256([]byte(fmt.Sprintf("%x", sha256.Sum256([]byte(enc_str)))))

	signer, err := rsa.SignPKCS1v15(rng, crt.PrivateKey, crypto.SHA256, hashed[:])

	if err != nil {
		return err
	}

	data["signature"] = base64.StdEncoding.EncodeToString(signer)
	return nil
}

//签证回传数据签名
func (u *Pay) validate(req UMP) (bool, error) {
	//fmt.Println(req)
	signature_base64 := req["signature"].(string)
	delete(req, "signature")
	link_str := u.createLinkString(req, true)

	pubkey, err := u.getVerifyCertPublicKey(req["signPubKeyCert"].(string))
	if err != nil {
		return false, err
	}

	hashed := sha256.Sum256([]byte(fmt.Sprintf("%x", sha256.Sum256([]byte(link_str)))))

	signature_str, err := base64.StdEncoding.DecodeString(signature_base64)
	if err != nil {
		return false, err
	}

	err = rsa.VerifyPKCS1v15(pubkey, crypto.SHA256, hashed[:], signature_str)
	if err != nil {
		return false, err
	}

	return true, nil
}

//得到用户信息加密
func (u *Pay) EncryptCustomerData(user *UserInfo) (string, error) {
	base := UMP{
		"certifTp":   user.CertifTp,
		"certifId":   user.CertifId,
		"customerNm": user.CustomerNm,
	}

	enc := UMP{
		"phoneNo": user.PhotoNo,
		"cvn2":    user.Cvn2,
		"expired": user.Expired,
	}

	enc_link := u.createLinkString(enc, false)

	enc_str, err := u.EncryptData(enc_link)
	if err != nil {
		return "", err
	}

	base["encryptedInfo"] = enc_str

	all_enc := fmt.Sprintf("{%s}", u.createLinkString(base, false))

	return base64.StdEncoding.EncodeToString([]byte(all_enc)), nil
}

//加密一般数据
func (u *Pay) EncryptData(data string) (string, error) {
	crt, err := u.getEncryptCertInfo()
	if err != nil {
		return "", err
	}
	rng := rand.Reader

	enc, err := rsa.EncryptPKCS1v15(rng, crt.PublicKey, []byte(data))

	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(enc), nil
}

//得到签名证书信息
func (u *Pay) getSignCertInfo() (*SignCertInfo, error) {
	certInfo := &SignCertInfo{}

	data, err := ioutil.ReadFile(u.conf.SignCertPath)
	if err != nil {
		return nil, err
	}

	privateKey, crt, err := pkcs12.Decode(data, u.conf.SignCertPassword)

	if err != nil {
		return nil, err
	}
	certInfo.Cert = crt
	certInfo.PrivateKey = privateKey.(*rsa.PrivateKey)
	return certInfo, nil
}

//得到内容加密证书信息
func (u *Pay) getEncryptCertInfo() (*EncryptCertInfo, error) {
	data, err := ioutil.ReadFile(u.conf.EncryptCertPath)
	if err != nil {
		return nil, err
	}

	p, data := pem.Decode(data)

	crt, err := x509.ParseCertificate(p.Bytes)

	if err != nil {
		return nil, err
	}

	certInfo := &EncryptCertInfo{
		Cert:      crt,
		PublicKey: crt.PublicKey.(*rsa.PublicKey),
	}
	return certInfo, nil
}

//得到验证签名证书
func (u *Pay) getVerifyCertPublicKey(sign_data string) (*rsa.PublicKey, error) {
	p, _ := pem.Decode([]byte(sign_data))

	crt, err := x509.ParseCertificate(p.Bytes)
	if err != nil {
		return nil, err
	}

	//fmt.Println()

	return crt.PublicKey.(*rsa.PublicKey), nil
}

//把MAP解为key=value&....串
func (u *Pay) createLinkString(data UMP, is_sort bool) string {
	var link_list []string

	if is_sort {
		keys := utils.MapKeys(data)
		sort.Strings(keys)

		for _, k := range keys {
			link_list = append(link_list, fmt.Sprintf("%s=%s", k, data[k]))
		}
	} else {
		for k, v := range data {
			link_list = append(link_list, fmt.Sprintf("%s=%s", k, v))
		}
	}

	return strings.Join(link_list, "&")
}

//POST数据
func (u *Pay) post(url_str string, post_data UMP) (UMP, error) {
	post := url.Values{}

	for k, v := range post_data {
		post.Add(k, v.(string))
	}

	req, err := http.PostForm(u.urls.BackTransUrl, post)

	if err != nil {
		return nil, err
	}
	defer req.Body.Close()

	r, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	return u.explainResponse(r), nil
}

//解释POST返回数据
func (u *Pay) explainResponse(data []byte) UMP {
	//fmt.Println(string(data))
	req_map := UMP{}
	req_str := string(data)
	req_str_list := strings.Split(req_str, "&")
	for _, v := range req_str_list {
		res := strings.SplitN(v, "=", 2)
		if len(res) > 1 {
			req_map[res[0]] = res[1]
		}
	}

	return req_map
}

//创建订单号
func (u *Pay) CreateOrderNo() string {
	prefix := "PCB"
	ti := time.Now()
	date := ti.Format("060102150405")
	return fmt.Sprintf("%s%s%d", prefix, date, ti.Nanosecond())
}

//得到推送返回数据
func (u *Pay) Receive(post_data []byte) (UMP, error) {
	data := u.explainResponse(post_data)
	ok, err := u.validate(data)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, nil
	}

	return data, nil
}
