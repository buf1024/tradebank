package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/url"
	"strings"
	"time"
)

type NocardPay struct {
	mall *YaodeMall
}

type NocarPayReq struct {
	mchKey      string
	merId       string
	orderId     string
	transAmount string
	cardByName  string
	cardByNo    string
	cerNumber   string
	mobile      string
}
type NocardVerifyCodeReq struct {
	mchKey       string
	merId        string
	yzm          string
	ksPayOrderId string
}
type NocardQueryReq struct {
	mchKey    string
	merId     string
	orderId   string
	transDate string
}

type NocardRsp struct {
	status        string // 00：成功 01：失败 02：系统错误
	orderId       string
	ksPayOrderId  string
	chanelRefcode string // 89   要求手机验证码
	bankOrderId   string
	refCode       string // ‘00’交易成功 01’预交易成功 ‘02’交易失败 03  交易处理中
	refMsg        string
}

func currentDate() string {
	t := time.Now()
	return fmt.Sprintf("%04d%02d%02d%02d%02d%02d",
		t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
}

func (m *NocardPay) signReqData(v *url.Values, key string) (string, error) {

	signStr := v.Encode() + key
	h := md5.New()
	md5sum := strings.ToUpper(hex.EncodeToString(h.Sum([]byte(signStr))))
	v.Add("signType", "MD5")
	v.Add("signData", md5sum)

	srcData := v.Encode()
	m.mall.Log.Info("REQ(NO ENC):%s\n", srcData)
	dstData, err := NoCardPayEncrypt(srcData, key)
	if err != nil {
		m.mall.Log.Error("encrypt req failed, err=%s\n", err)
		return "", nil
	}
	m.mall.Log.Info("REQ(ENC):%s\n", dstData)
	return dstData, err
}
func (m *NocardPay) signCheckData(strData string, key string, md5val string) bool {
	signStr := strData + key
	h := md5.New()
	md5sum := strings.ToUpper(hex.EncodeToString(h.Sum([]byte(signStr))))
	m.mall.Log.Info("compute md5:%s, receive md5: %s\n", md5sum, md5val)
	if md5sum != md5val {
		return false
	}
	return true
}

func (m *NocardPay) MakePayReq(req *NocarPayReq) (string, error) {
	v := url.Values{}
	v.Add("versionId", "001")
	v.Add("businessType", "1401")
	//v.Add("insCode", "")
	v.Add("merId", req.merId)
	v.Add("orderId", req.orderId)
	v.Add("transDate", currentDate())
	v.Add("transAmount", req.transAmount)
	//v.Add("transCurrency", "156")
	v.Add("cardByName", req.cardByName)
	v.Add("cardByNo", req.cardByNo)
	v.Add("cardType", "01")
	//v.Add("expireDate", "")
	//v.Add("CVV", "")
	//v.Add("bankCode", "")
	//v.Add("openBankName", "")
	v.Add("cerType", "01")
	v.Add("cerNumber", req.cerNumber)
	v.Add("mobile", req.mobile)
	v.Add("isAcceptYzm", "00")
	//v.Add("pageNotifyUrl", "")
	//v.Add("backNotifyUrl", "")
	//v.Add("orderDesc", "")
	v.Add("instalTransFlag", "01")
	//v.Add("instalTransNums", "")
	//v.Add("dev", "")
	//v.Add("fee", "")
	return m.signReqData(&v, req.mchKey)
}

func (m *NocardPay) MakeVerifyCodeReq(req *NocardVerifyCodeReq) (string, error) {
	v := url.Values{}
	v.Add("versionId", "001")
	v.Add("businessType", "1411")
	//v.Add("insCode", "")
	v.Add("merId", req.merId)
	v.Add("yzm", req.yzm)
	v.Add("ksPayOrderId", req.ksPayOrderId)

	return m.signReqData(&v, req.mchKey)
}
func (m *NocardPay) MakeQueryReq(req *NocardQueryReq) (string, error) {
	v := url.Values{}
	v.Add("versionId", "001")
	v.Add("businessType", "1421")
	//v.Add("insCode", "")
	v.Add("merId", req.merId)
	v.Add("orderId", req.orderId)
	v.Add("transDate", currentDate())

	return m.signReqData(&v, req.mchKey)
}

func (m *NocardPay) ParseRsp(rspStr string, key string) (*NocardRsp, error) {
	m.mall.Log.Info("RSP(ENC):%s\n", rspStr)
	decStr, err := NoCardPayEncrypt(rspStr, key)
	if err != nil {
		return nil, err
	}
	m.mall.Log.Info("RSP(DEC):%s\n", decStr)
	idx := strings.Index(decStr, "&signType")
	if idx < 0 {
		return nil, fmt.Errorf("&signType not found")
	}
	strData := decStr[0:idx]
	signData := decStr[idx+1:]
	idx = strings.Index(signData, "Data=")
	md5val := signData[idx+len("Data="):]

	if m.signCheckData(strData, key, md5val) {
		return nil, fmt.Errorf("md5 not equal")
	}
	v, err := url.ParseQuery(strData)
	if err != nil {
		return nil, err
	}

	rsp := &NocardRsp{
		status:        v.Get("status"),
		orderId:       v.Get("orderId"),
		ksPayOrderId:  v.Get("ksPayOrderId"),
		chanelRefcode: v.Get("chanelRefcode"),
		bankOrderId:   v.Get("bankOrderId"),
		refCode:       v.Get("refCode"),
		refMsg:        v.Get("refMsg"),
	}
	return rsp, nil
}
