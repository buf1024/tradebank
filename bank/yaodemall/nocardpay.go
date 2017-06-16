package main

import (
	"bytes"
	"crypto/aes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
	"tradebank/proto"
	"tradebank/util"

	"encoding/json"

	"tradebank/iomsframe"

	pb "github.com/golang/protobuf/proto"
	ini "github.com/vaughan0/go-ini"
)

type NoCardPay struct {
	mall *YaodeMall
	//no cardpay
	NocardMchNo   string
	NocardMchKey  string
	NocardReqHost string
}

func (m *NoCardPay) PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func (m *NoCardPay) PKCS5Unpadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
func (m *NoCardPay) NoCardPayEncrypt(src string, key string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	blockSize := block.BlockSize()
	srcPadding := m.PKCS5Padding([]byte(src), blockSize)
	dstEnc := make([]byte, len(srcPadding))
	for bs, be := 0, blockSize; bs < len(srcPadding); bs, be = bs+blockSize, be+blockSize {
		block.Encrypt(dstEnc[bs:be], srcPadding[bs:be])
	}
	dstStr := hex.EncodeToString(dstEnc)
	return strings.ToUpper(dstStr), nil
}
func (m *NoCardPay) NoCardPayDecrypt(src string, key string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	data, err := hex.DecodeString(src)
	if err != nil {
		return "", err
	}
	blockSize := block.BlockSize()
	decBytes := make([]byte, len(data))
	srcBytes := []byte(data)
	for bs, be := 0, blockSize; bs < len(data); bs, be = bs+blockSize, be+blockSize {
		block.Decrypt(decBytes[bs:be], srcBytes[bs:be])
	}
	decStr := string(m.PKCS5Unpadding(decBytes))
	return decStr, nil
}

func (m *NoCardPay) SignReqData(v *PayUrlValues, key string) (string, error) {

	signStr := v.Encode() + key
	h := md5.New()
	h.Write([]byte(signStr))
	md5sum := strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
	v.Add("signType", "MD5")
	v.Add("signData", md5sum)
	srcData := v.Encode()
	m.mall.Log.Info("REQ(NO ENC):%s\n", srcData)
	dstData, err := m.NoCardPayEncrypt(srcData, key)
	if err != nil {
		m.mall.Log.Error("encrypt req failed, err=%s\n", err)
		return "", nil
	}
	m.mall.Log.Info("REQ(ENC):%s\n", dstData)

	post := &PayUrlValues{}
	post.Add("merId", m.NocardMchNo)
	post.Add("transData", dstData)

	dstData = post.Encode()

	return dstData, err
}
func (m *NoCardPay) SignCheckData(strData string, key string, md5val string) bool {
	signStr := strData + key
	h := md5.New()
	md5sum := strings.ToUpper(hex.EncodeToString(h.Sum([]byte(signStr))))
	m.mall.Log.Info("compute md5:%s, receive md5: %s\n", md5sum, md5val)
	if md5sum != md5val {
		return false
	}
	return true
}

func (m *NoCardPay) Init(f *ini.File) error {
	//nocard pay
	ok := false
	m.NocardMchNo, ok = f.Get("YDM", "NOCARDPAY_MCHNO")
	if !ok {
		return fmt.Errorf("missing configure, sec=YDM, key=NOCARDPAY_MCHNO")
	}
	m.NocardMchKey, ok = f.Get("YDM", "NOCARPAY_MCHKEY")
	if !ok {
		return fmt.Errorf("missing configure, sec=YDM, key=NOCARPAY_MCHKEY")
	}
	m.NocardReqHost, ok = f.Get("YDM", "NOCARPAY_PAYHOST")
	if !ok {
		return fmt.Errorf("missing configure, sec=YDM, key=NOCARPAY_PAYHOST")
	}
	return nil
}
func (m *NoCardPay) InMoneyReq(req *proto.E2BInMoneyReq) error {
	bankReq := &PayReq{}
	bankReq.cardByName = base64.StdEncoding.EncodeToString([]byte(req.GetCustName()))
	bankReq.cardByNo = req.GetBankAcct()
	bankReq.cerNumber = req.GetCustCID()
	bankReq.mchKey = m.NocardMchKey
	bankReq.merId = m.NocardMchNo
	mobile := util.GetSplitData(req.GetReversed(), "PHONE=")
	if mobile == "" {
		m.mall.Log.Error("missing required field phone no\n")
		return fmt.Errorf("missing required field phone no")
	}
	bankReq.mobile = mobile
	bankReq.orderId = req.GetExchSID()
	bankReq.transAmount = fmt.Sprintf("%.2f", req.GetAmount())

	bankMsg, err := m.PayReq(bankReq)
	if err != nil {
		return err
	}

	dbData := InoutLog{
		extflow: bankReq.orderId,
		iotype:  iomsframe.BANK_INMONEY,
		amount:  req.GetAmount(),
	}
	err = m.mall.db.InsertLog(dbData)
	if err != nil {
		m.mall.Log.Error("InsertLog failed. extflow=%s\n", dbData.extflow)
		return err
	}

	m.mall.Log.Info("POST REQ:\nURL=%s\n, DATA=%s\n", m.NocardReqHost, bankMsg)
	bankRsp, err := util.PostData(m.NocardReqHost, []byte(bankMsg))
	if err != nil {
		return err
	}
	m.mall.Log.Info("POST RSP:%s\n", string(bankRsp))
	rspData, err := m.ParseRsp(bankRsp)
	if err != nil {
		return nil
	}
	rspMsg, err := proto.Message(proto.CMD_E2B_IN_MONEY_RSP)
	if err != nil {
		return err
	}
	rsp := rspMsg.(*proto.E2BInMoneyRsp)
	rsp.ExchSID = pb.String(req.GetExchSID())
	rsp.BankID = pb.Int32(req.GetBankID())
	rsp.RetCode = pb.Int32(m.GetExchCode(rspData))
	rsp.RetMsg = pb.String(rspData.refMsg)

	if rsp.GetRetCode() == util.E_SUCCESS {
		util.CallMeLater(m.mall.TimeOutReconn, m.CheckResult, dbData.extflow)
	}

	return m.mall.MakeRsp(proto.CMD_E2B_IN_MONEY_RSP, rsp)

}
func (m *NoCardPay) OutMoneyReq(req *proto.E2BOutMoneyReq) error {
	return fmt.Errorf("not support outmoney")
}
func (m *NoCardPay) VerifyReq(req *proto.E2BVerifyCodeReq) error {
	return nil
}

func (m *NoCardPay) CheckReq(orderId string) error {
	return nil

}

func (m *NoCardPay) PayReq(req *PayReq) (string, error) {
	//v := url.Values{}
	v := PayUrlValues{}
	v.Add("versionId", "001")
	v.Add("businessType", "1401")
	//v.Add("insCode", "")
	v.Add("merId", req.merId)
	v.Add("orderId", req.orderId)
	v.Add("transDate", util.CurrentDate())
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
	return m.SignReqData(&v, req.mchKey)
}

func (m *NoCardPay) VerifyCodeReq(req *VerifyReq) (string, error) {
	//v := url.Values{}
	v := PayUrlValues{}
	v.Add("versionId", "001")
	v.Add("businessType", "1411")
	//v.Add("insCode", "")
	v.Add("merId", req.merId)
	v.Add("yzm", req.yzm)
	v.Add("ksPayOrderId", req.ksPayOrderId)

	return m.SignReqData(&v, req.mchKey)
}
func (m *NoCardPay) QueryReq(req *QueryReq) (string, error) {
	//v := url.Values{}
	v := PayUrlValues{}
	v.Add("versionId", "001")
	v.Add("businessType", "1421")
	//v.Add("insCode", "")
	v.Add("merId", req.merId)
	v.Add("orderId", req.orderId)
	v.Add("transDate", util.CurrentDate())

	return m.SignReqData(&v, req.mchKey)
}

func (m *NoCardPay) ParseRsp(rspStr []byte) (*PayRsp, error) {
	v := make(map[string]interface{})
	err := json.Unmarshal(rspStr, &v)
	if err != nil {
		return nil, err
	}
	rsp := &PayRsp{
		status:        v["status"].(string),
		orderId:       v["orderId"].(string),
		ksPayOrderId:  v["ksPayOrderId"].(string),
		chanelRefcode: v["chanelRefcode"].(string),
		bankOrderId:   v["bankOrderId"].(string),
		refCode:       v["refCode"].(string),
		refMsg:        v["refMsg"].(string),
	}
	return rsp, nil
}
func (m *NoCardPay) GetExchCode(rsp *PayRsp) int32 {
	if rsp.status != "00" || rsp.refCode != "01" {
		// 处理失败
		return util.E_BANK_ERR
	}
	if rsp.chanelRefcode == "89" {
		// 需发送验证码
		return util.E_HALF_SUCCESS
	}
	return util.E_SUCCESS
}

func (m *NoCardPay) CheckResult(data interface{}) {
}
