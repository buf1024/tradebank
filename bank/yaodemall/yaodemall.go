package main

import (
	"fmt"
	"strconv"
	"tradebank/iomsframe"

	"github.com/golang/protobuf/proto"
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
func (b *YaodeMall) ExchReq(command int64, msg proto.Message) error {
	return b.HandleDef(command, msg)
}
func (b *YaodeMall) ExchRsp(command int64, msg proto.Message) error {
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
