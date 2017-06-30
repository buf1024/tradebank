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
	NetBankReqHost string

	PageNotifyUrl string
	BackNotifyUrl string
}

func (p *NetBankPay) loadConf(f *ini.File) error {
	//netbank pay
	ok := false
	p.NetBankReqHost, ok = f.Get("YDM", "NETBANKPAY_PAYHOST")
	if !ok {
		return fmt.Errorf("missing configure, sec=YDM, key=NETBANKPAY_PAYHOST")
	}
	p.PageNotifyUrl, ok = f.Get("YDM", "NETBANKPAGE_NOTIFY_URL")
	if !ok {
		return fmt.Errorf("missing configure, sec=YDM, key=NETBANKPAGE_NOTIFY_URL")
	}
	p.BackNotifyUrl, ok = f.Get("YDM", "NETBANKBACK_NOTIFY_URL")
	if !ok {
		return fmt.Errorf("missing configure, sec=YDM, key=NETBANKBACK_NOTIFY_URL")
	}
	iStart := strings.LastIndex(p.BackNotifyUrl, ":")
	if iStart == -1 || iStart == len(p.BackNotifyUrl)-1 {
		return fmt.Errorf("missing listen port")
	}
	iEnd := strings.Index(p.BackNotifyUrl[iStart+1:], "/")
	if iEnd == -1 {
		var err error
		p.listenPort, err = strconv.Atoi(p.BackNotifyUrl[iStart+1:])
		if err != nil {
			return err
		}
		p.notifyUrl = "/"
	} else {
		iEnd = iStart + 1 + iEnd
		var err error
		p.listenPort, err = strconv.Atoi(p.BackNotifyUrl[iStart+1 : iEnd])
		if err != nil {
			return err
		}
		p.notifyUrl = p.BackNotifyUrl[iEnd:]
	}
	p.mall.Log.Debug("listen=%d, url=%s\n", p.listenPort, p.notifyUrl)
	return nil

}
func (p *NetBankPay) NotifyHandler(rsp http.ResponseWriter, req *http.Request) {

}
func (p *NetBankPay) TestHttpListen() bool {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", p.listenPort))
	if err != nil {
		p.mall.Log.Error("try to listen port %d failed.\n", p.listenPort)
		return false
	}
	l.Close()
	return true
}
func (p *NetBankPay) HttpListen() {
	p.mall.Log.Info("http start listen %d \n", p.listenPort)
	http.HandleFunc(p.notifyUrl, p.NotifyHandler)
	err := http.ListenAndServe(fmt.Sprintf(":%d", p.listenPort), nil)
	if err != nil {
		p.mall.Log.Error("http listen %s failed\n", p.listenPort)

	}
}
func (p *NetBankPay) Init(f *ini.File) error {
	err := p.loadConf(f)
	if err != nil {
		p.mall.Log.Error("netbank loadConf failed, err=%s\n", err)
		return err
	}
	if p.TestHttpListen() {
		go p.HttpListen()
	} else {
		return fmt.Errorf("listen netbank addr failed")
	}

	return nil
}

func (p *NetBankPay) InMoneyReq(req *proto.E2BInMoneyReq) error {
	bankReq := &PayReq{}
	bankReq.cardByName = base64.StdEncoding.EncodeToString([]byte(req.GetCustName()))
	bankReq.cardByNo = req.GetBankAcct()
	bankReq.cerNumber = req.GetCustCID()
	bankReq.mchKey = p.mall.MchKey
	bankReq.merId = p.mall.MchNo
	bankReq.orderId = req.GetExchSID()
	bankReq.transAmount = fmt.Sprintf("%.2f", req.GetAmount())
	bankReq.backNotifyUrl = p.BackNotifyUrl
	bankReq.pageNotifyUrl = p.PageNotifyUrl

	bankMsg, err := p.PayReq(bankReq)
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
	rsp.PostUrl = pb.String(p.NetBankReqHost)
	rsp.PostData = pb.String(bankMsg)

	return p.mall.MakeRsp(proto.CMD_E2B_IN_MONEY_RSP, rsp)
}
func (p *NetBankPay) OutMoneyReq(req *proto.E2BOutMoneyReq) error {
	return fmt.Errorf("not surport outmoney")
}
func (p *NetBankPay) VerifyReq(req *proto.E2BVerifyCodeReq) error {
	return fmt.Errorf("not surport verify code")
}

func (p *NetBankPay) CheckReq(orderId string) (int32, error) {
	return 0, nil
}
func (p *NetBankPay) GetExchCode(rsp *PayRsp) int32 {
	return 0
}
func (p *NetBankPay) SignReqData(v *PayUrlValues, key string) (string, error) {

	signStr := v.Encode() + key
	h := md5.New()
	h.Write([]byte(signStr))
	md5sum := strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
	v.Add("signData", md5sum)
	srcData := v.Encode()
	p.mall.Log.Info("REQ:%s\n", srcData)

	return srcData, nil
}
func (p *NetBankPay) PayReq(req *PayReq) (string, error) {
	//v := url.Values{}
	v := PayUrlValues{}
	v.Add("versionId", "001")
	v.Add("businessType", "1100")
	//v.Add("insCode", "")
	v.Add("merId", req.merId)
	v.Add("orderId", req.orderId)
	v.Add("transDate", util.CurrentDate())
	v.Add("transAmount", req.transAmount)
	v.Add("transCurrency", "156")
	//v.Add("transChanlName", "")
	//v.Add("openBankName", "")
	v.Add("pageNotifyUrl", req.pageNotifyUrl)
	v.Add("backNotifyUrl", req.backNotifyUrl)
	//v.Add("orderDesc", "")
	//v.Add("dev", "")

	return p.SignReqData(&v, req.mchKey)
}
