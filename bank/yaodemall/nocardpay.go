package main

import (
	"encoding/base64"
	"fmt"
	"tradebank/proto"
	"tradebank/util"
)

func (m *NocardPay) HandleInMoney(req *proto.E2BInMoneyReq) error {
	bankReq := &NocarPayReq{}
	bankReq.cardByName = base64.StdEncoding.EncodeToString([]byte(req.GetCustName()))
	bankReq.cardByNo = req.GetBankAcct()
	bankReq.cerNumber = req.GetCustCID()
	bankReq.mchKey = m.mall.NocardMchKey
	bankReq.merId = m.mall.NocardMchNo
	mobile := util.GetSplitData(req.GetReversed(), "PHONE=")
	if mobile == "" {
		m.mall.Log.Error("missing required field phone no\n")
		return fmt.Errorf("missing required field phone no")
	}
	bankReq.mobile = mobile
	bankReq.orderId = req.GetExchSID()
	bankReq.transAmount = fmt.Sprintf("%.2f", req.GetAmount())

	bankMsg, err := m.MakePayReq(bankReq)
	if err != nil {
		return err
	}
	m.mall.Log.Info("POST REQ:\nURL=%s\n, DATA=%s\n", m.mall.NocardReqHost, bankMsg)
	bankRsp, err := util.PostData(m.mall.NocardReqHost, []byte(bankMsg))
	if err != nil {
		return err
	}
	m.mall.Log.Info("POST RSP:%s\n", string(bankRsp))

	return nil

}
