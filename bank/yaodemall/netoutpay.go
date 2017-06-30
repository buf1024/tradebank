package main

import (
	"crypto/md5"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"io/ioutil"
	"strings"
	"tradebank/iomsframe"
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
	NetOutMchNo      string
	NetOutMchKey     string
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

func (n *NetOutPay) Encrypt(req *PayReq) (string, error) {

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
	encJS, err := myrsa.PrivateEncrypt(n.privt, js)
	if err != nil {
		return "", nil
	}
	return string(encJS), nil

}
func (n *NetOutPay) Decrypt(data []byte) (string, error) {
	decData, err := myrsa.PublicDecrypt(n.pub, data)
	if err != nil {
		return "", err
	}
	return string(decData), nil
}

func (n *NetOutPay) loadConf(f *ini.File) error {
	//netbank pay
	ok := false
	n.NetOutMchNo, ok = f.Get("YDM", "NETOUTPAY_MCHNO")
	if !ok {
		return fmt.Errorf("missing configure, sec=YDM, key=NETOUTPAY_MCHNO")
	}
	n.NetOutMchKey, ok = f.Get("YDM", "NETOUTPAY_MCHKEY")
	if !ok {
		return fmt.Errorf("missing configure, sec=YDM, key=NETOUTPAY_MCHKEY")
	}
	n.NetOutPublicKey, ok = f.Get("YDM", "NETOUTPAY_PUBLIC_KEY")
	if !ok {
		return fmt.Errorf("missing configure, sec=YDM, key=NETOUTPAY_PUBLIC_KEY")
	}
	n.NetOutPrivateKey, ok = f.Get("YDM", "NETOUTPAY_PRIVITE_KEY")
	if !ok {
		return fmt.Errorf("missing configure, sec=YDM, key=NETOUTPAY_PRIVITE_KEY")
	}
	if err := n.loadRSAkeys(n.NetOutPublicKey, n.NetOutPrivateKey); err != nil {
		return err
	}
	n.NetOutReqHost, ok = f.Get("YDM", "NETOUTPAY_PAYHOST")
	if !ok {
		return fmt.Errorf("missing configure, sec=YDM, key=NETOUTPAY_PAYHOST")
	}

	return nil

}
func (n *NetOutPay) loadRSAkeys(public string, private string) error {
	// load private key
	bs, err := ioutil.ReadFile(private)
	if err != nil {
		return err
	}
	block, _ := pem.Decode(bs)
	if block == nil {
		return fmt.Errorf("decode private key failed")
	}
	n.privt, err = x509.ParsePKCS1PrivateKey(block.Bytes)
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
	n.pub, ok = cert.PublicKey.(*rsa.PublicKey)
	if !ok {
		return fmt.Errorf("convert to public key failed")
	}
	return nil
}

func (n *NetOutPay) Init(f *ini.File) error {
	err := n.loadConf(f)
	if err != nil {
		n.mall.Log.Error("netout loadConf failed, err=%s\n", err)
		return err
	}

	return nil
}

func (n *NetOutPay) InMoneyReq(req *proto.E2BInMoneyReq) error {
	return fmt.Errorf("not surport inmoney")

}
func (n *NetOutPay) OutMoneyReq(req *proto.E2BOutMoneyReq) error {
	bankReq := &PayReq{}
	bankReq.cardByName = base64.StdEncoding.EncodeToString([]byte(req.GetCustName()))
	bankReq.cardByNo = req.GetBankAcct()
	bankReq.cerNumber = req.GetCustCID()
	bankReq.mchKey = n.NetOutMchKey
	bankReq.merId = n.NetOutMchNo
	bankReq.orderId = req.GetExchSID()
	bankReq.transAmount = fmt.Sprintf("%.2f", req.GetAmount())

	bankMsg, err := n.PayReq(bankReq)
	if err != nil {
		return err
	}
	dbData := InoutLog{
		extflow: bankReq.orderId,
		iotype:  iomsframe.BANK_OUTMONEY,
		amount:  req.GetAmount(),
	}
	err = n.mall.db.InsertLog(dbData)
	if err != nil {
		n.mall.Log.Error("InsertLog failed. extflow=%s\n", dbData.extflow)
		return err
	}

	n.mall.Log.Info("POST REQ:\nURL=%s\n, DATA=%s\n", n.NetOutReqHost, bankMsg)
	bankRsp, err := util.PostData(n.NetOutReqHost, []byte(bankMsg))
	if err != nil {
		return err
	}
	n.mall.Log.Info("POST RSP:%s\n", string(bankRsp))
	payRsp, err := n.ParseRsp(bankRsp) // todo
	if err != nil {
		n.mall.Log.Error("parse resp message failed.")
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
	util.CallMeLater(n.mall.TimeOutSess, n.CheckResult, ctx)

	rsp := rspMsg.(*proto.E2BOutMoneyRsp)
	rsp.ExchSID = pb.String(req.GetExchSID())
	rsp.BankID = pb.Int32(req.GetBankID())
	rsp.RetCode = pb.Int32(int32(util.E_SUCCESS)) // 默认成功
	rsp.RetMsg = pb.String(payRsp.refMsg)

	return n.mall.MakeRsp(proto.CMD_E2B_OUT_MONEY_RSP, rsp)
}
func (n *NetOutPay) VerifyReq(req *proto.E2BVerifyCodeReq) error {
	return fmt.Errorf("not surport verify code")
}

func (n *NetOutPay) CheckReq(orderId string) (int32, error) {
	bankReq := &PayReq{}
	bankReq.merId = n.NetOutMchNo
	bankReq.mchKey = n.NetOutMchKey
	bankReq.orderId = orderId

	bankMsg, err := n.QueryReq(bankReq)
	if err != nil {
		return 0, err
	}
	n.mall.Log.Info("POST REQ:\nURL=%s\n, DATA=%s\n", n.NetOutPublicKey, bankMsg)
	bankRsp, err := util.PostData(n.NetOutReqHost, []byte(bankMsg))
	if err != nil {
		return 0, err
	}
	n.mall.Log.Info("POST RSP:%s\n", string(bankRsp))
	rspData, err := n.ParseRsp(bankRsp)
	if err != nil {
		return 0, err
	}
	packSt, buzSt := n.GetExchCode(rspData)
	if packSt == util.E_SUCCESS {
		return buzSt, nil
	}
	return 0, fmt.Errorf("packet state not success")

}
func (n *NetOutPay) GetExchCode(rsp *PayRsp) (int32, int32) {
	if rsp.status != "00" {
		return util.E_BANK_ERR, 0
	}
	if rsp.refCode == "1" {
		return util.E_SUCCESS, util.E_SUCCESS
	}

	return util.E_SUCCESS, util.E_HALF_SUCCESS
}
func (n *NetOutPay) EncryptReqData(v *TransReq, key string) (string, error) {

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
	n.mall.Log.Info("REQ:%s\n", jsStr)

	return jsStr, nil
}
func (n *NetOutPay) PayReq(req *PayReq) (string, error) {
	tb, err := n.Encrypt(req)
	if err != nil {
		return "", err
	}

	tranReq := &TransReq{
		VersionID:    "001",
		BusinessType: "470000",
		MerID:        req.merId,
		TransBody:    tb,
	}

	return n.EncryptReqData(tranReq, req.mchKey)
}
func (n *NetOutPay) QueryReq(req *PayReq) (string, error) {
	tb, err := n.Encrypt(req)
	if err != nil {
		return "", err
	}

	tranReq := &TransReq{
		VersionID:    "001",
		BusinessType: "460000",
		MerID:        req.merId,
		TransBody:    tb,
	}

	return n.EncryptReqData(tranReq, req.mchKey)
}

func (n *NetOutPay) ParseRsp(rspStr []byte) (*PayRsp, error) {
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
		rspStr, err := n.Decrypt([]byte(rspBody))
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

func (n *NetOutPay) CheckResult(to int64, data interface{}) {
	ctx := data.(*QueryResultContext)
	if ctx.retryTimes < QUERYRESULT_RETRY_TIMES {
		n.mall.Log.Debug("query out money result, extflow=%s\n", ctx.extflow)
		ctx.retryTimes++
		n.mall.Log.Info("query out money result for sid=%s\n", ctx.extflow)
		ret, err := n.CheckReq(ctx.extflow)
		if err != nil {
			n.mall.Log.Info("check req out money failed, err=%s\n", err)
		} else {
			if ret != util.E_HALF_SUCCESS {
				err = n.mall.db.UpdateLog(ctx.extflow, int(ret), "")
				if err != nil {
					n.mall.Log.Error("update database error, err=%s\n", ctx.extflow)
					return
				}

				reqMsg, err := proto.Message(proto.CMD_B2E_INOUTNOTIFY_REQ)
				if err != nil {
					n.mall.Log.Info("proto.message error: %s\n", err)
					return
				}
				req := reqMsg.(*proto.B2EInOutNotifyReq)
				req.TransType = pb.Int32(iomsframe.BANK_OUTMONEY)
				req.BankAcct = pb.String(ctx.bankacct)
				req.BankId = pb.Int32(int32(n.mall.BankID))
				req.BankSID = pb.String(ctx.orderId)
				req.Currency = pb.Int32(1)
				req.ExchSID = pb.String(ctx.extflow)
				req.Status = pb.Int32(ret)
				req.RetMsg = pb.String(util.GetErrMsg(int64(ret)))
				n.mall.Log.Info("query out money result, notfiy status req : %s\n", proto.Debug(proto.CMD_B2E_INOUTNOTIFY_REQ, req))

				n.mall.MakeRsp(proto.CMD_B2E_INOUTNOTIFY_REQ, req)

				// insert session

				return
			}
		}
		util.CallMeLater(ctx.retryTimes*n.mall.TimeOutReconn, n.CheckResult, ctx)

		return
	}
	n.mall.Log.Info("no result for sid = %s, query result in daily check.\n", ctx.extflow)
}
