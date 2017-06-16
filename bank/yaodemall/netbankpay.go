package main

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"strings"
	"tradebank/proto"
	"tradebank/util"

	"fmt"

	"strconv"

	"net/http"

	"net"

	pb "github.com/golang/protobuf/proto"
	ini "github.com/vaughan0/go-ini"
)

type NetBankPay struct {
	mall       *YaodeMall
	listenPort int
	notifyUrl  string

	//no cardpay
	NetBankMchNo   string
	NetBankMchKey  string
	NetBankReqHost string

	PageNotifyUrl string
	BackNotifyUrl string
}

func (n *NetBankPay) loadConf(f *ini.File) error {
	//netbank pay
	ok := false
	n.NetBankMchNo, ok = f.Get("YDM", "NETBANKPAY_MCHNO")
	if !ok {
		return fmt.Errorf("missing configure, sec=YDM, key=NETBANKPAY_MCHNO")
	}
	n.NetBankMchKey, ok = f.Get("YDM", "NETBANKPAY_MCHKEY")
	if !ok {
		return fmt.Errorf("missing configure, sec=YDM, key=NETBANKPAY_MCHKEY")
	}
	n.NetBankReqHost, ok = f.Get("YDM", "NETBANKPAY_PAYHOST")
	if !ok {
		return fmt.Errorf("missing configure, sec=YDM, key=NETBANKPAY_PAYHOST")
	}
	n.PageNotifyUrl, ok = f.Get("YDM", "NETBANKPAGE_NOTIFY_URL")
	if !ok {
		return fmt.Errorf("missing configure, sec=YDM, key=NETBANKPAGE_NOTIFY_URL")
	}
	n.BackNotifyUrl, ok = f.Get("YDM", "NETBANKBACK_NOTIFY_URL")
	if !ok {
		return fmt.Errorf("missing configure, sec=YDM, key=NETBANKBACK_NOTIFY_URL")
	}
	iStart := strings.LastIndex(n.BackNotifyUrl, ":")
	if iStart == -1 || iStart == len(n.BackNotifyUrl)-1 {
		return fmt.Errorf("missing listen port")
	}
	iEnd := strings.Index(n.BackNotifyUrl[iStart+1:], "/")
	if iEnd == -1 {
		var err error
		n.listenPort, err = strconv.Atoi(n.BackNotifyUrl[iStart+1:])
		if err != nil {
			return err
		}
		n.notifyUrl = "/"
	} else {
		iEnd = iStart + 1 + iEnd
		var err error
		n.listenPort, err = strconv.Atoi(n.BackNotifyUrl[iStart+1 : iEnd])
		if err != nil {
			return err
		}
		n.notifyUrl = n.BackNotifyUrl[iEnd:]
	}
	n.mall.Log.Debug("listen=%d, url=%s\n", n.listenPort, n.notifyUrl)
	return nil

}
func (n *NetBankPay) NotifyHandler(rsp http.ResponseWriter, req *http.Request) {

}
func (n *NetBankPay) TestHttpListen() bool {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", n.listenPort))
	if err != nil {
		n.mall.Log.Error("try to listen port %d failed.\n", n.listenPort)
		return false
	}
	l.Close()
	return true
}
func (n *NetBankPay) HttpListen() {
	n.mall.Log.Info("http start listen %d \n", n.listenPort)
	http.HandleFunc(n.notifyUrl, n.NotifyHandler)
	err := http.ListenAndServe(fmt.Sprintf(":%d", n.listenPort), nil)
	if err != nil {
		n.mall.Log.Error("http listen %s failed\n", n.listenPort)

	}
}
func (n *NetBankPay) Init(f *ini.File) error {
	err := n.loadConf(f)
	if err != nil {
		n.mall.Log.Error("netbank loadConf failed, err=%s\n", err)
		return err
	}
	if n.TestHttpListen() {
		go n.HttpListen()
	} else {
		return fmt.Errorf("listen netbank addr failed")
	}

	return nil
}

func (n *NetBankPay) InMoneyReq(req *proto.E2BInMoneyReq) error {
	bankReq := &PayReq{}
	bankReq.cardByName = base64.StdEncoding.EncodeToString([]byte(req.GetCustName()))
	bankReq.cardByNo = req.GetBankAcct()
	bankReq.cerNumber = req.GetCustCID()
	bankReq.mchKey = n.NetBankMchKey
	bankReq.merId = n.NetBankMchNo
	bankReq.orderId = req.GetExchSID()
	bankReq.transAmount = fmt.Sprintf("%.2f", req.GetAmount())
	bankReq.backNotifyUrl = n.BackNotifyUrl
	bankReq.pageNotifyUrl = n.PageNotifyUrl

	bankMsg, err := n.PayReq(bankReq)
	if err != nil {
		return err
	}

	rspMsg, err := proto.Message(proto.CMD_E2B_IN_MONEY_RSP)
	if err != nil {
		return err
	}

	rsp := rspMsg.(*proto.E2BInMoneyRsp)
	rsp.ExchSID = pb.String(req.GetExchSID())
	rsp.BankID = pb.Int32(req.GetBankID())
	rsp.RetCode = pb.Int32(int32(util.E_SUCCESS))
	rsp.RetMsg = pb.String(util.GetErrMsg(util.E_SUCCESS))
	rsp.PostUrl = pb.String(n.NetBankReqHost)
	rsp.PostData = pb.String(bankMsg)

	return n.mall.MakeRsp(proto.CMD_E2B_IN_MONEY_RSP, rsp)
}
func (n *NetBankPay) OutMoneyReq(req *proto.E2BOutMoneyReq) error {
	return fmt.Errorf("not surport outmoney")
}
func (n *NetBankPay) VerifyReq(req *proto.E2BVerifyCodeReq) error {
	return fmt.Errorf("not surport verify code")

}

func (n *NetBankPay) CheckReq(orderId string) error {
	return nil
}
func (n *NetBankPay) SignReqData(v *PayUrlValues, key string) (string, error) {

	signStr := v.Encode() + key
	h := md5.New()
	h.Write([]byte(signStr))
	md5sum := strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
	v.Add("signData", md5sum)
	srcData := v.Encode()
	n.mall.Log.Info("REQ:%s\n", srcData)

	return srcData, nil
}
func (n *NetBankPay) PayReq(req *PayReq) (string, error) {
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
