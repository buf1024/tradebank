package main

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"strings"
	"tradebank/iomsframe"
	"tradebank/proto"
	"tradebank/util"

	"fmt"

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
	n.NetOutReqHost, ok = f.Get("YDM", "NETOUTPAY_PAYHOST")
	if !ok {
		return fmt.Errorf("missing configure, sec=YDM, key=NETOUTPAY_PAYHOST")
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
	_, err = n.ParseRsp(bankRsp) // todo
	if err != nil {
		return nil
	}

	rspMsg, err := proto.Message(proto.CMD_E2B_OUT_MONEY_RSP)
	if err != nil {
		return err
	}

	rsp := rspMsg.(*proto.E2BOutMoneyRsp)
	rsp.ExchSID = pb.String(req.GetExchSID())
	rsp.BankID = pb.Int32(req.GetBankID())
	rsp.RetCode = pb.Int32(int32(util.E_SUCCESS))
	rsp.RetMsg = pb.String(util.GetErrMsg(util.E_SUCCESS))

	return n.mall.MakeRsp(proto.CMD_E2B_OUT_MONEY_RSP, rsp)
}
func (n *NetOutPay) VerifyReq(req *proto.E2BVerifyCodeReq) error {
	return fmt.Errorf("not surport verify code")

}

func (n *NetOutPay) CheckReq(orderId string) (int32, error) {
	return 0, nil
}
func (n *NetOutPay) GetExchCode(rsp *PayRsp) int32 {
	return 0
}
func (n *NetOutPay) SignReqData(v *PayUrlValues, key string) (string, error) {

	signStr := v.Encode() + key
	h := md5.New()
	h.Write([]byte(signStr))
	md5sum := strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
	v.Add("signData", md5sum)
	srcData := v.Encode()
	n.mall.Log.Info("REQ:%s\n", srcData)

	return srcData, nil
}
func (n *NetOutPay) PayReq(req *PayReq) (string, error) {
	//v := url.Values{}
	v := PayUrlValues{}
	v.Add("versionId", "001")
	v.Add("businessType", "1100")
	v.Add("insCode", "")
	v.Add("merId", req.merId)
	v.Add("orderId", req.orderId)
	v.Add("transDate", util.CurrentDate())
	v.Add("transAmount", req.transAmount)
	v.Add("transCurrency", "156")
	v.Add("transChanlName", "")
	v.Add("openBankName", "")
	v.Add("pageNotifyUrl", req.pageNotifyUrl)
	v.Add("backNotifyUrl", req.backNotifyUrl)
	v.Add("orderDesc", "")
	v.Add("dev", "")

	return n.SignReqData(&v, req.mchKey)
}
func (n *NetOutPay) ParseRsp(rspStr []byte) (*PayRsp, error) {
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
