package main

import (
	"bytes"
	"fmt"
	"strconv"
	"tradebank/iomsframe"

	"tradebank/proto"

	"tradebank/util"

	"sync"

	pb "github.com/golang/protobuf/proto"
	ini "github.com/vaughan0/go-ini"
)

type YaodeMall struct {
	iomsframe.ExchFrame

	pay map[string]YaodePay
	db  *YaodeMallDB

	BankName string
	BankID   int64
	DBPath   string
}

type PayUrlValues struct {
	buf bytes.Buffer
}

func (v *PayUrlValues) Add(key, val string) {
	if v.buf.Len() > 0 {
		v.buf.WriteByte('&')
	}
	v.buf.WriteString(key)
	v.buf.WriteByte('=')
	v.buf.WriteString(val)
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

	b.DBPath, ok = f.Get("YDM", "DB_PATH")
	if !ok {
		return fmt.Errorf("missing configure, sec=YDM, key=DB_PATH")
	}

	for _, v := range b.pay {
		if err = v.Init(&f); err != nil {
			return err
		}
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
	err = b.db.Init(b.DBPath)
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
				payway = "0"
			}
			paywayNum, err := strconv.Atoi(payway)
			if err != nil {
				b.Log.Error("unknown payway, payway=%s\n", payway)
				return err
			}
			pay := b.GetPay(iomsframe.BANK_INMONEY, paywayNum)
			if pay == nil {
				b.Log.Error("unknown payway, payway=%s\n", payway)
				return fmt.Errorf("unknown payway, payway=%s", payway)
			}
			return pay.InMoneyReq(req)
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
func (b *YaodeMall) AddPay(inout int, payway int, pay YaodePay) {
	key := fmt.Sprintf("%d-%d", inout, payway)
	b.pay[key] = pay
}
func (b *YaodeMall) GetPay(inout int, payway int) YaodePay {
	key := fmt.Sprintf("%d-%d", inout, payway)

	if v, exist := b.pay[key]; exist {
		return v
	}
	return nil
}

// YaodeMallServer
func YaodeMallServer() *YaodeMall {
	m := &YaodeMall{
		pay: make(map[string]YaodePay),
		db:  &YaodeMallDB{lock: &sync.Mutex{}},
	}
	m.Bank = m
	m.db.mall = m

	nocard := &NoCardPay{mall: m}
	m.AddPay(iomsframe.BANK_INMONEY, iomsframe.PAYWAY_NOCARD, nocard)

	netbank := &NetBankPay{mall: m}
	m.AddPay(iomsframe.BANK_INMONEY, iomsframe.PAYWAY_NETBANK, netbank)

	return m
}
