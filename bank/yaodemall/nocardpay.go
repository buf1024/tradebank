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

	"tradebank/ioms"

	pb "github.com/golang/protobuf/proto"
	ini "github.com/vaughan0/go-ini"
)

type NoCardPay struct {
	mall *YaodeMall
	//no cardpay
	NocardReqHost string
}

func (p *NoCardPay) PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func (p *NoCardPay) PKCS5Unpadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
func (p *NoCardPay) NoCardPayEncrypt(src string, key string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	blockSize := block.BlockSize()
	srcPadding := p.PKCS5Padding([]byte(src), blockSize)
	dstEnc := make([]byte, len(srcPadding))
	for bs, be := 0, blockSize; bs < len(srcPadding); bs, be = bs+blockSize, be+blockSize {
		block.Encrypt(dstEnc[bs:be], srcPadding[bs:be])
	}
	dstStr := hex.EncodeToString(dstEnc)
	return strings.ToUpper(dstStr), nil
}
func (p *NoCardPay) NoCardPayDecrypt(src string, key string) (string, error) {
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
	decStr := string(p.PKCS5Unpadding(decBytes))
	return decStr, nil
}

func (p *NoCardPay) SignReqData(v *PayUrlValues, key string) (string, error) {

	signStr := v.Encode() + key
	h := md5.New()
	h.Write([]byte(signStr))
	md5sum := strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
	v.Add("signType", "MD5")
	v.Add("signData", md5sum)
	srcData := v.Encode()
	p.mall.Log.Info("REQ(NO ENC):%s\n", srcData)
	dstData, err := p.NoCardPayEncrypt(srcData, key)
	if err != nil {
		p.mall.Log.Error("encrypt req failed, err=%s\n", err)
		return "", nil
	}
	p.mall.Log.Info("REQ(ENC):%s\n", dstData)

	post := &PayUrlValues{}
	post.Add("merId", p.mall.MchNo)
	post.Add("transData", dstData)

	dstData = post.Encode()

	return dstData, err
}
func (p *NoCardPay) SignCheckData(strData string, key string, md5val string) bool {
	signStr := strData + key
	h := md5.New()
	md5sum := strings.ToUpper(hex.EncodeToString(h.Sum([]byte(signStr))))
	p.mall.Log.Info("compute md5:%s, receive md5: %s\n", md5sum, md5val)
	if md5sum != md5val {
		return false
	}
	return true
}

func (p *NoCardPay) Init(f *ini.File) error {
	//nocard pay
	ok := false
	p.NocardReqHost, ok = f.Get("YDM", "NOCARPAY_PAYHOST")
	if !ok {
		return fmt.Errorf("missing configure, sec=YDM, key=NOCARPAY_PAYHOST")
	}
	return nil
}
func (p *NoCardPay) InMoneyReq(req *proto.E2BInMoneyReq) error {
	bankReq := &PayReq{}
	bankReq.cardByName = base64.StdEncoding.EncodeToString([]byte(req.GetCustName()))
	bankReq.cardByNo = req.GetBankAcct()
	bankReq.cerNumber = req.GetCustCID()
	bankReq.mchKey = p.mall.MchKey
	bankReq.merId = p.mall.MchNo
	mobile := util.GetSplitData(req.GetReversed(), "PHONE=")
	if mobile == "" {
		p.mall.Log.Error("missing required field phone no\n")
		return fmt.Errorf("missing required field phone no")
	}
	bankReq.mobile = mobile
	bankReq.orderId = req.GetExchSID()
	bankReq.transAmount = fmt.Sprintf("%.2f", req.GetAmount())

	bankMsg, err := p.PayReq(bankReq)
	if err != nil {
		return err
	}

	dbData := InoutLog{
		extflow: bankReq.orderId,
		iotype:  ioms.BANK_INMONEY,
		amount:  req.GetAmount(),
	}
	err = p.mall.db.InsertLog(dbData)
	if err != nil {
		p.mall.Log.Error("InsertLog failed. extflow=%s\n", dbData.extflow)
		return err
	}

	p.mall.Log.Info("POST REQ:\nURL=%s\n, DATA=%s\n", p.NocardReqHost, bankMsg)
	bankRsp, err := util.PostData(p.NocardReqHost, []byte(bankMsg))
	if err != nil {
		return err
	}
	p.mall.Log.Info("POST RSP:%s\n", string(bankRsp))
	rspData, err := p.ParseRsp(bankRsp)
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
	rsp.RetCode = pb.Int32(util.E_BANK_ERR)
	packSt, buzSt := p.GetExchCode(rspData)
	if packSt == util.E_SUCCESS {
		rsp.RetCode = pb.Int32(buzSt)
	}
	rsp.RetMsg = pb.String(rspData.refMsg)

	if rsp.GetRetCode() == util.E_SUCCESS {
		ctx := &QueryResultContext{
			extflow:    dbData.extflow,
			orderId:    rspData.orderId,
			amount:     bankReq.transAmount,
			retryTimes: 0,
		}
		util.CallMeLater(p.mall.TimeOutReconn, p.CheckResult, ctx)
	}

	return p.mall.MakeRsp(proto.CMD_E2B_IN_MONEY_RSP, rsp)

}
func (p *NoCardPay) OutMoneyReq(req *proto.E2BOutMoneyReq) error {
	return fmt.Errorf("not support outmoney")
}
func (p *NoCardPay) VerifyReq(req *proto.E2BVerifyCodeReq) error {
	return nil
}

func (p *NoCardPay) CheckReq(orderId string) (int32, error) {
	bankReq := &QueryReq{}
	bankReq.merId = p.mall.MchNo
	bankReq.mchKey = p.mall.MchKey
	bankReq.orderId = orderId

	bankMsg, err := p.QueryReq(bankReq)
	if err != nil {
		return 0, err
	}
	p.mall.Log.Info("POST REQ:\nURL=%s\n, DATA=%s\n", p.NocardReqHost, bankMsg)
	bankRsp, err := util.PostData(p.NocardReqHost, []byte(bankMsg))
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

func (p *NoCardPay) PayReq(req *PayReq) (string, error) {
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
	return p.SignReqData(&v, req.mchKey)
}

func (p *NoCardPay) VerifyCodeReq(req *VerifyReq) (string, error) {
	//v := url.Values{}
	v := PayUrlValues{}
	v.Add("versionId", "001")
	v.Add("businessType", "1411")
	//v.Add("insCode", "")
	v.Add("merId", req.merId)
	v.Add("yzm", req.yzm)
	v.Add("ksPayOrderId", req.ksPayOrderId)

	return p.SignReqData(&v, req.mchKey)
}
func (p *NoCardPay) QueryReq(req *QueryReq) (string, error) {
	//v := url.Values{}
	v := PayUrlValues{}
	v.Add("versionId", "001")
	v.Add("businessType", "1421")
	//v.Add("insCode", "")
	v.Add("merId", req.merId)
	v.Add("orderId", req.orderId)
	v.Add("transDate", util.CurrentDate())

	return p.SignReqData(&v, req.mchKey)
}

func (p *NoCardPay) ParseRsp(rspStr []byte) (*PayRsp, error) {
	v := make(map[string]interface{})
	err := json.Unmarshal(rspStr, &v)
	if err != nil {
		return nil, err
	}
	rsp := &PayRsp{}
	if t, exists := v["status"]; exists {
		rsp.status, _ = t.(string)
	}
	if t, exists := v["orderId"]; exists {
		rsp.orderId, _ = t.(string)
	}
	if t, exists := v["ksPayOrderId"]; exists {
		rsp.ksPayOrderId, _ = t.(string)
	}
	if t, exists := v["chanelRefcode"]; exists {
		rsp.chanelRefcode, _ = t.(string)
	}
	if t, exists := v["bankOrderId"]; exists {
		rsp.bankOrderId, _ = t.(string)
	}
	if t, exists := v["refCode"]; exists {
		rsp.refCode, _ = t.(string)
	}
	if t, exists := v["refMsg"]; exists {
		rsp.refMsg, _ = t.(string)
	}
	return rsp, nil
}

func (p *NoCardPay) GetExchCode(rsp *PayRsp) (packStatus int32, buzStatus int32) {
	packStatus = util.E_SUCCESS
	buzStatus = 0
	if rsp.status != "00" {
		packStatus = util.E_BANK_ERR
		return
	}
	buzStatus = util.E_BANK_ERR
	if rsp.refCode == "00" || rsp.refCode == "01" {
		if rsp.chanelRefcode == "89" {
			buzStatus = util.E_HALF_SUCCESS
		} else {
			buzStatus = util.E_SUCCESS
		}
	} else if rsp.refCode == "03" {
		buzStatus = util.E_HALF_SUCCESS
	}
	return
}

func (p *NoCardPay) CheckResult(to int64, data interface{}) {
	ctx := data.(*QueryResultContext)
	if ctx.retryTimes < QUERYRESULT_RETRY_TIMES {
		p.mall.Log.Debug("query in money result, extflow=%s\n", ctx.extflow)
		ctx.retryTimes++
		p.mall.Log.Info("query in money result for sid=%s\n", ctx.extflow)
		ret, err := p.CheckReq(ctx.extflow)
		if err != nil {
			p.mall.Log.Info("check req failed, err=%s\n", err)
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
				p.mall.Log.Info("query in money result, notfiy status req : %s\n", proto.Debug(proto.CMD_B2E_INOUTNOTIFY_REQ, req))

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
