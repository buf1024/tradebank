package main

import (
	"crypto/md5"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"io/ioutil"
	"strconv"
	"strings"
	"tradebank/ioms"
	"tradebank/proto"
	"tradebank/util"

	"fmt"

	"crypto/rsa"

	myrsa "github.com/buf1024/golib/crypt"
	pb "github.com/golang/protobuf/proto"
	ini "github.com/vaughan0/go-ini"
)

type NetOutPay struct {
	mall *YaodeMall

	//no cardpay
	NetOutPublicKey  string
	NetOutPrivateKey string
	NetOutReqHost    string

	pub   *rsa.PublicKey
	privt *rsa.PrivateKey
}
type NetOutQueryResult struct {
	extflow    string
	retryTimes int64
}

type TransReq struct {
	VersionID    string `json:"versionId,omitempty"`
	BusinessType string `json:"businessType,omitempty"`
	MerID        string `json:"merId,omitempty"`
	TransBody    string `json:"transBody,omitempty"`
	SignType     string `json:"signType,omitempty"`
	SignData     string `json:"signData,omitempty"`
}

type TransBody struct {
	OrderId     string `json:"orderId,omitempty"`
	TransDate   string `json:"transDate,omitempty"`
	TransAmount string `json:"transAmount,omitempty"`
	AccNo       string `json:"accNo,omitempty"`
	AccName     string `json:"accName,omitempty"`
}

type RespBody struct {
	OrderId string `json:"orderId,omitempty"`
	RefCode string `json:"refCode,omitempty"`
	RefMsg  string `json:"refMsg,omitempty"`
}

func (p *NetOutPay) Encrypt(req *PayReq) (string, error) {

	t := &TransBody{}
	t.OrderId = req.orderId
	t.TransDate = util.CurrentDate()
	t.TransAmount = req.transAmount
	t.AccNo = req.cardByNo
	t.AccName = req.cardByName

	js, err := json.Marshal(t)
	if err != nil {
		return "", nil
	}
	js = []byte(`{"accName":"%B2%E2%CA%D4","accNo":"621226111111111111111","orderId":"tx20170703051473","transDate":"20170703092105","transAmount":"1"}`)
	encJS, err := myrsa.PrivateEncrypt(p.privt, js)
	if err != nil {
		return "", nil
	}
	encStr := base64.StdEncoding.EncodeToString(encJS)
	fmt.Printf("enc=%s\n", encStr)
	return encStr, nil

}
func (p *NetOutPay) Decrypt(data []byte) (string, error) {

	decBase64 := make([]byte, len(data))
	n, err := base64.RawStdEncoding.Decode(decBase64, data)
	if err != nil {
		return "", err
	}
	decData, err := myrsa.PublicDecrypt(p.pub, decBase64[:n])
	if err != nil {
		return "", err
	}
	return string(decData), nil
}

func (p *NetOutPay) loadConf(f *ini.File) error {
	//netbank pay
	ok := false
	p.NetOutPublicKey, ok = f.Get("YDM", "NETOUTPAY_PUBLIC_KEY")
	if !ok {
		return fmt.Errorf("missing configure, sec=YDM, key=NETOUTPAY_PUBLIC_KEY")
	}
	p.NetOutPrivateKey, ok = f.Get("YDM", "NETOUTPAY_PRIVITE_KEY")
	if !ok {
		return fmt.Errorf("missing configure, sec=YDM, key=NETOUTPAY_PRIVITE_KEY")
	}
	if err := p.loadRSAkeys(p.NetOutPublicKey, p.NetOutPrivateKey); err != nil {
		return err
	}
	p.NetOutReqHost, ok = f.Get("YDM", "NETOUTPAY_PAYHOST")
	if !ok {
		return fmt.Errorf("missing configure, sec=YDM, key=NETOUTPAY_PAYHOST")
	}

	return nil

}
func (p *NetOutPay) loadRSAkeys(public string, private string) error {
	// load private key
	bs, err := ioutil.ReadFile(private)
	if err != nil {
		return err
	}
	block, _ := pem.Decode(bs)
	if block == nil {
		return fmt.Errorf("decode private key failed")
	}
	p.privt, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return err
	}
	// load public key
	bs, err = ioutil.ReadFile(public)
	if err != nil {
		return err
	}
	block, _ = pem.Decode(bs)
	if block == nil {
		return fmt.Errorf("decode certifacte key failed")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return err
	}
	ok := false
	p.pub, ok = cert.PublicKey.(*rsa.PublicKey)
	if !ok {
		return fmt.Errorf("convert to public key failed")
	}
	return nil
}

func (p *NetOutPay) Init(f *ini.File) error {
	err := p.loadConf(f)
	if err != nil {
		p.mall.Log.Error("netout loadConf failed, err=%s\n", err)
		return err
	}

	return nil
}

func (p *NetOutPay) InMoneyReq(req *proto.E2BInMoneyReq) error {
	return fmt.Errorf("not surport inmoney")

}
func (p *NetOutPay) OutMoneyReq(req *proto.E2BOutMoneyReq) error {
	bankReq := &PayReq{}
	bankReq.cardByName = base64.StdEncoding.EncodeToString([]byte(req.GetCustName()))
	bankReq.cardByNo = req.GetBankAcct()
	bankReq.cerNumber = req.GetCustCID()
	bankReq.mchKey = p.mall.MchKey
	bankReq.merId = p.mall.MchNo
	bankReq.orderId = req.GetExchSID()
	bankReq.transAmount = fmt.Sprintf("%.2f", req.GetAmount())

	bankMsg, err := p.PayReq(bankReq)
	if err != nil {
		return err
	}
	dbData := InoutLog{
		extflow: bankReq.orderId,
		iotype:  ioms.BANK_OUTMONEY,
		amount:  req.GetAmount(),
	}
	err = p.mall.db.InsertLog(dbData)
	if err != nil {
		p.mall.Log.Error("InsertLog failed. extflow=%s\n", dbData.extflow)
		return err
	}

	p.mall.Log.Info("POST REQ:\nURL=%s\n, DATA=%s\n", p.NetOutReqHost, bankMsg)
	bankRsp, err := util.PostData(p.NetOutReqHost, []byte(bankMsg))
	if err != nil {
		return err
	}
	p.mall.Log.Info("POST RSP:%s\n", string(bankRsp))
	payRsp, err := p.ParseRsp(bankRsp) // todo
	if err != nil {
		p.mall.Log.Error("parse resp message failed.")
		return nil
	}

	rspMsg, err := proto.Message(proto.CMD_E2B_OUT_MONEY_RSP)
	if err != nil {
		return err
	}
	ctx := &QueryResultContext{
		extflow:    dbData.extflow,
		orderId:    payRsp.orderId,
		amount:     bankReq.transAmount,
		retryTimes: 0,
	}
	util.CallMeLater(p.mall.TimeOutSess, p.CheckResult, ctx)

	rsp := rspMsg.(*proto.E2BOutMoneyRsp)
	rsp.ExchSID = pb.String(req.GetExchSID())
	rsp.BankID = pb.Int32(req.GetBankID())
	rsp.RetCode = pb.Int32(int32(util.E_SUCCESS)) // 默认成功
	rsp.RetMsg = pb.String(payRsp.refMsg)

	return p.mall.MakeRsp(proto.CMD_E2B_OUT_MONEY_RSP, rsp)
}
func (p *NetOutPay) VerifyReq(req *proto.E2BVerifyCodeReq) error {
	return fmt.Errorf("not surport verify code")
}

func (p *NetOutPay) CheckReq(orderId string) (int32, error) {
	bankReq := &PayReq{}
	bankReq.merId = p.mall.MchNo
	bankReq.mchKey = p.mall.MchKey
	bankReq.orderId = orderId

	bankMsg, err := p.QueryReq(bankReq)
	if err != nil {
		return 0, err
	}
	p.mall.Log.Info("POST REQ:\nURL=%s\n, DATA=%s\n", p.NetOutPublicKey, bankMsg)
	bankRsp, err := util.PostData(p.NetOutReqHost, []byte(bankMsg))
	if err != nil {
		return 0, err
	}
	p.mall.Log.Info("POST RSP:%s\n", string(bankRsp))
	rspData, err := p.ParseRsp(bankRsp)
	if err != nil {
		return 0, err
	}
	packSt, buzSt := p.GetExchCode(rspData)
	if packSt == util.E_SUCCESS {
		return buzSt, nil
	}
	return 0, fmt.Errorf("packet state not success")

}
func (p *NetOutPay) GetExchCode(rsp *PayRsp) (int32, int32) {
	if rsp.status != "00" {
		return util.E_BANK_ERR, 0
	}
	if rsp.refCode == "1" {
		return util.E_SUCCESS, util.E_SUCCESS
	}

	return util.E_SUCCESS, util.E_HALF_SUCCESS
}
func (p *NetOutPay) EncryptReqData(v *TransReq, key string) (string, error) {

	signStr := fmt.Sprintf("businessType=%s&merId=%s&transBody=%s&versionId=%s&key=%s",
		v.BusinessType, v.MerID, v.TransBody, v.VersionID, key)
	h := md5.New()
	h.Write([]byte(signStr))
	md5sum := strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
	v.SignType = "MD5"
	v.SignData = md5sum
	jsData, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	jsStr := string(jsData)
	p.mall.Log.Info("REQ:%s\n", jsStr)

	return jsStr, nil
}
func (p *NetOutPay) PayReq(req *PayReq) (string, error) {
	tb, err := p.Encrypt(req)
	if err != nil {
		return "", err
	}

	tranReq := &TransReq{
		VersionID:    "001",
		BusinessType: "470000",
		MerID:        req.merId,
		TransBody:    tb,
	}

	return p.EncryptReqData(tranReq, req.mchKey)
}
func (p *NetOutPay) QueryReq(req *PayReq) (string, error) {
	tb, err := p.Encrypt(req)
	if err != nil {
		return "", err
	}

	tranReq := &TransReq{
		VersionID:    "001",
		BusinessType: "460000",
		MerID:        req.merId,
		TransBody:    tb,
	}

	return p.EncryptReqData(tranReq, req.mchKey)
}

func (p *NetOutPay) ParseRsp(rspStr []byte) (*PayRsp, error) {
	v := make(map[string]interface{})
	err := json.Unmarshal(rspStr, &v)
	if err != nil {
		return nil, err
	}

	status := ""
	ok := false

	rspJs := &RespBody{}

	if status, ok = v["status"].(string); !ok {
		return nil, fmt.Errorf("status not found")
	}

	if status == "00" {
		rspBody, ok := v["resBody"].(string)
		if !ok {
			return nil, fmt.Errorf("resBody not found")
		}
		rspStr, err := p.Decrypt([]byte(rspBody))
		if err != nil {
			return nil, err
		}
		if err = json.Unmarshal([]byte(rspStr), rspJs); err != nil {
			return nil, err
		}
	}

	rsp := &PayRsp{
		status:  v["status"].(string),
		orderId: rspJs.OrderId,
		refCode: rspJs.RefCode,
		refMsg:  rspJs.RefMsg,
	}
	return rsp, nil
}

func (p *NetOutPay) CheckResult(to int64, data interface{}) {
	ctx := data.(*QueryResultContext)
	if ctx.retryTimes < QUERYRESULT_RETRY_TIMES {
		p.mall.Log.Debug("query out money result, extflow=%s\n", ctx.extflow)
		ctx.retryTimes++
		p.mall.Log.Info("query out money result for sid=%s\n", ctx.extflow)
		ret, err := p.CheckReq(ctx.extflow)
		if err != nil {
			p.mall.Log.Info("check req out money failed, err=%s\n", err)
		} else {
			if ret != util.E_HALF_SUCCESS {
				err = p.mall.db.UpdateLog(ctx.extflow, int(ret), "")
				if err != nil {
					p.mall.Log.Error("update database error, err=%s\n", ctx.extflow)
					return
				}

				reqMsg, err := proto.Message(proto.CMD_B2E_INOUTNOTIFY_REQ)
				if err != nil {
					p.mall.Log.Info("proto.message error: %s\n", err)
					return
				}
				req := reqMsg.(*proto.B2EInOutNotifyReq)
				req.TransType = pb.Int32(ioms.BANK_OUTMONEY)
				req.BankAcct = pb.String(ctx.bankacct)
				req.BankId = pb.Int32(int32(p.mall.BankID))
				req.BankSID = pb.String(ctx.orderId)
				req.Currency = pb.Int32(1)
				req.ExchSID = pb.String(ctx.extflow)
				req.Status = pb.Int32(ret)
				req.RetMsg = pb.String(util.GetErrMsg(int64(ret)))
				amt, err := strconv.ParseFloat(ctx.amount, 64)
				if err != nil {
					p.mall.Log.Info("parse float error: %s\n", err)
					return
				}
				req.Amount = pb.Float64(amt)
				p.mall.Log.Info("query out money result, notfiy status req : %s\n", proto.Debug(proto.CMD_B2E_INOUTNOTIFY_REQ, req))

				p.mall.MakeRsp(proto.CMD_B2E_INOUTNOTIFY_REQ, req)

				// insert session

				return
			}
		}
		util.CallMeLater(ctx.retryTimes*p.mall.TimeOutReconn, p.CheckResult, ctx)

		return
	}
	p.mall.Log.Info("no result for sid = %s, query result in daily check.\n", ctx.extflow)
}
