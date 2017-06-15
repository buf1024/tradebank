package main

import (
	"bytes"
	"fmt"
	"strconv"
	"tradebank/iomsframe"

	"tradebank/proto"

	"tradebank/util"

	"net/url"

	pb "github.com/golang/protobuf/proto"
	ini "github.com/vaughan0/go-ini"
)

type YaodeMall struct {
	iomsframe.ExchFrame

	nocard *NocardPay

	BankName string
	BankID   int64
	//no cardpay
	NocardMchNo   string
	NocardMchKey  string
	NocardReqHost string
}

type PayUrlValues struct {
	buf bytes.Buffer
}

func (v *PayUrlValues) Add(key, val string) {
	if v.buf.Len() > 0 {
		v.buf.WriteByte('&')
	}
	v.buf.WriteString(url.QueryEscape(key))
	v.buf.WriteByte('=')
	v.buf.WriteString(url.QueryEscape(val))
}
func (v *PayUrlValues) Encode() string {

	return v.buf.String()
}

func (b *YaodeMall) loadBankConf(path string) error {
	f, err := ini.LoadFile(path)
	if err != nil {
		return err
	}
	ok := false
	str := ""

	// common
	b.BankName, ok = f.Get("YDM", "BANK_NAME")
	if !ok {
		return fmt.Errorf("missing configure, sec=YDM, key=BANK_NAME")
	}

	str, ok = f.Get("YDM", "BANK_ID")
	if !ok {
		return fmt.Errorf("missing configure, sec=YDM, key=BANK_ID")
	}
	i, err := strconv.Atoi(str)
	if err != nil {
		return fmt.Errorf("convert %s to interger failed", str)
	}
	b.BankID = int64(i)

	//nocard pay
	b.NocardMchNo, ok = f.Get("YDM", "NOCARDPAY_MCHNO")
	if !ok {
		return fmt.Errorf("missing configure, sec=YDM, key=NOCARDPAY_MCHNO")
	}
	b.NocardMchKey, ok = f.Get("YDM", "NOCARPAY_MCHKEY")
	if !ok {
		return fmt.Errorf("missing configure, sec=YDM, key=NOCARPAY_MCHKEY")
	}
	b.NocardReqHost, ok = f.Get("YDM", "NOCARPAY_PAYHOST")
	if !ok {
		return fmt.Errorf("missing configure, sec=YDM, key=NOCARPAY_PAYHOST")
	}

	return nil
}

func (b *YaodeMall) Name() string {
	return b.BankName
}

func (b *YaodeMall) ID() int64 {
	return b.BankID
}

func (b *YaodeMall) InitBank(m *iomsframe.ExchFrame) error {
	err := b.loadBankConf(m.FileConf)
	if err != nil {
		return err
	}
	return nil
}
func (b *YaodeMall) StopBank(m *iomsframe.ExchFrame) {

}
func (b *YaodeMall) ExchReq(command int64, msg pb.Message) error {
	switch command {
	case proto.CMD_E2B_IN_MONEY_REQ:
		{
			req := msg.(*proto.E2BInMoneyReq)
			payway := util.GetSplitData(req.GetReversed(), "PAYWAY=")
			if payway == "" {
				b.Log.Warning("req missing payway. use the default one\n")
				return b.nocard.HandleInMoney(req)
			}
			paywayNum, err := strconv.Atoi(payway)
			if err != nil {
				b.Log.Error("unknown payway, payway=%s\n", payway)
				return err
			}
			switch paywayNum {
			case iomsframe.PAYWAY_NOCARD:
				{
					return b.nocard.HandleInMoney(req)
				}
			default:
				{
					b.Log.Error("unknown payway, payway=%s\n", payway)
					return fmt.Errorf("unknown payway, payway=%s", payway)
				}
			}
		}
	case proto.CMD_E2B_OUT_MONEY_REQ:
		{

		}
	default:
		{
			return b.HandleDef(command, msg)
		}
	}
	return nil
}
func (b *YaodeMall) ExchRsp(command int64, msg pb.Message) error {
	switch command {
	case proto.CMD_B2E_IN_MONEY_RSP:
		{

		}
	default:
		{
			return b.HandleDef(command, msg)
		}
	}
	return nil
}

// YaodeMallServer
func YaodeMallServer() *YaodeMall {
	m := &YaodeMall{}

	m.Bank = m
	m.nocard = &NocardPay{}
	m.nocard.mall = m

	return m
}
