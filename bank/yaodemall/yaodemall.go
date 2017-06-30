package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strconv"
	"tradebank/ioms"

	"tradebank/proto"

	"tradebank/util"

	"sync"

	"io/ioutil"

	pb "github.com/golang/protobuf/proto"
	ini "github.com/vaughan0/go-ini"
)

type YaodeMall struct {
	ioms.ExchFrame

	pay map[string]YaodePay
	db  *YaodeMallDB

	BankName string
	BankID   int64
	DBPath   string

	MchNo  string
	MchKey string

	CheckPath string
	FtpHost   string
	FtpUser   string
	FtpPass   string
	FtpPath   string
}

type CheckContext struct {
	pay        YaodePay
	date       string
	log        *InoutLog
	file       *ioms.CheckFile
	retryTimes int64
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
	b.MchNo, ok = f.Get("YDM", "PAY_MCHNO")
	if !ok {
		return fmt.Errorf("missing configure, sec=YDM, key=PAY_MCHNO")
	}
	b.MchNo, err = util.DBDecrypt(b.MchNo)
	if err != nil {
		return fmt.Errorf("decrypt mch no failed, err=%s", str)
	}
	b.MchKey, ok = f.Get("YDM", "PAY_MCHKEY")
	if !ok {
		return fmt.Errorf("missing configure, sec=YDM, key=PAY_MCHKEY")
	}
	b.MchKey, err = util.DBDecrypt(b.MchKey)
	if err != nil {
		return fmt.Errorf("decrypt mch key failed, err=%s", str)
	}

	b.FtpHost, ok = f.Get("YDM", "FTP_CHECK_HOST")
	if !ok {
		return fmt.Errorf("missing configure, sec=YDM, key=FTP_CHECK_HOST")
	}
	b.FtpHost, err = util.DBDecrypt(b.FtpHost)
	if err != nil {
		return fmt.Errorf("decrypt ftp host failed, err=%s", str)
	}
	b.FtpUser, ok = f.Get("YDM", "FTP_CHECK_USER")
	if !ok {
		return fmt.Errorf("missing configure, sec=YDM, key=FTP_CHECK_USER")
	}
	b.FtpUser, err = util.DBDecrypt(b.FtpUser)
	if err != nil {
		return fmt.Errorf("decrypt ftp user failed, err=%s", str)
	}
	b.FtpPass, ok = f.Get("YDM", "FTP_CHECK_PASS")
	if !ok {
		return fmt.Errorf("missing configure, sec=YDM, key=FTP_CHECK_PASS")
	}
	b.FtpPass, err = util.DBDecrypt(b.FtpPass)
	if err != nil {
		return fmt.Errorf("decrypt ftp pass failed, err=%s", str)
	}
	b.FtpPath, ok = f.Get("YDM", "FTP_CHECK_PATH")
	if !ok {
		return fmt.Errorf("missing configure, sec=YDM, key=FTP_CHECK_PATH")
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

func (b *YaodeMall) InitBank(m *ioms.ExchFrame) error {
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
func (b *YaodeMall) StopBank(m *ioms.ExchFrame) {

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
			pay := b.GetPay(ioms.BANK_INMONEY, paywayNum)
			if pay == nil {
				b.Log.Error("unknown payway, payway=%s\n", payway)
				return fmt.Errorf("unknown payway, payway=%s", payway)
			}
			return pay.InMoneyReq(req)
		}
	case proto.CMD_E2B_OUT_MONEY_REQ:
		{
			pay := b.GetPay(ioms.BANK_OUTMONEY, 0)
			if pay == nil {
				return fmt.Errorf("not support out money req")
			}
			req := msg.(*proto.E2BOutMoneyReq)
			return pay.OutMoneyReq(req)
		}
	case proto.CMD_E2B_CHECK_START_REQ:
		{
			req := msg.(*proto.E2BCheckStartReq)

			rspMsg, err := proto.Message(proto.CMD_E2B_CHECK_START_RSP)
			if err != nil {
				b.Log.Error("create message failed, ERR=%s\n", err.Error())
				return err
			}

			rsp := rspMsg.(*proto.E2BCheckStartRsp)
			rsp.BankID = pb.Int32(req.GetBankID())
			rsp.ExchSID = pb.String(req.GetExchSID())
			rsp.RetCode = pb.Int32(util.E_SUCCESS)

			err = b.MakeRsp(proto.CMD_E2B_CHECK_START_RSP, rsp)
			if err != nil {
				b.Log.Error("MakeRsp failed, ERR=%s\n", err.Error())
				return err
			}

			t, err := util.DateStrToUTCMicroSec(req.GetTradeDate())
			if err != nil {
				b.Log.Error("convert date str to utc micro sec failed, ERR=%s\n", err.Error())
				return err
			}

			logs, err := b.db.QueryCheckLog(t)
			if err != nil {
				b.Log.Error("query database failed, ERR=%s\n", err.Error())
				return err
			}
			chkFile := ioms.NewCheckFile(b.CheckPath, int(b.BankID),
				int(req.GetBatchNo()), req.GetTradeDate(), len(logs))

			if len(logs) == 0 {
				b.Log.Info("inout logs is nil, rsp to exch\n")
				b.CheckFileNotify(chkFile)
			} else {
				for _, log := range logs {
					pay := b.GetPay(log.iotype, log.payway)
					if pay == nil {
						b.Log.Error("unknown payway, payway=%s\n", log.payway)
						return fmt.Errorf("unknown payway, payway=%d", log.payway)
					}
					ctx := &CheckContext{
						pay:        pay,
						log:        &log,
						file:       chkFile,
						date:       req.GetTradeDate(),
						retryTimes: 0,
					}
					b.CheckReq(b.TimeOutSess, ctx)
				}
			}
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
	case proto.CMD_B2E_INOUTNOTIFY_RSP:
		{

		}
	case proto.CMD_B2E_CHECK_FILE_NOTIFICATION_RSP:
		{
			// no need further precess
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
func (b *YaodeMall) CheckFileNotify(f *ioms.CheckFile) {
	if err := util.FtpPut(b.FtpHost, b.FtpUser, b.FtpPass, f.FileName, b.FtpPath, f.FullPath); err != nil {
		b.Log.Error("ftp put failed. err=%s\n", err)
		return
	}

	reqMsg, err := proto.Message(proto.CMD_B2E_CHECK_FILE_NOTIFICATION_REQ)
	if err != nil {
		b.Log.Error("create message failed, ERR=%s\n", err.Error())
		return
	}

	req := reqMsg.(*proto.B2ECheckFileNotificationReq)
	req.BankSID = pb.String(util.SID())
	req.BankID = pb.Int32(int32(b.BankID))
	req.CheckFileName = pb.String(f.FileName)
	req.CheckFileCount = pb.Int32(int32(f.Total))

	signData, err := ioutil.ReadFile(f.FullPath)
	if err != nil {
		b.Log.Error("read file failed, err=%s\n", err)
		return
	}
	h := md5.New()
	md5sum := hex.EncodeToString(h.Sum(signData))
	req.CheckFileNameMD5 = pb.String(md5sum)

	b.MakeReq(proto.CMD_B2E_CHECK_FILE_NOTIFICATION_REQ, req)
}
func (b *YaodeMall) CheckReq(to int64, data interface{}) {
	go func() {
		ctx := data.(*CheckContext)
		st, err := ctx.pay.CheckReq(ctx.log.extflow)
		if err != nil {
			if ctx.retryTimes < QUERYRESULT_RETRY_TIMES {
				ctx.retryTimes = ctx.retryTimes + 1
				b.Log.Error("query state failed, retry later")
				util.CallMeLater(to, b.CheckReq, ctx)
			} else {
				b.Log.Error("query state failed, left for next check")
				ctx.file.Append(nil)
			}
		} else {
			if st == util.E_HALF_SUCCESS {
				b.Log.Info("query state processing, left for next check")
				ctx.file.Append(nil)
			} else {
				log := ctx.log
				if st == util.E_SUCCESS {
					log.status = 1
				} else {
					log.status = 2
				}
				if err := b.db.UpdateLog(log.extflow, log.status, ctx.date); err != nil {
					b.Log.Error("update log failed. extflow = %s\n", log.extflow)
				}
			}
		}

		if ctx.file.CheckDone() {
			b.CheckFileNotify(ctx.file)
		}

	}()
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
	m.AddPay(ioms.BANK_INMONEY, ioms.PAYWAY_NOCARD, nocard)

	netbank := &NetBankPay{mall: m}
	m.AddPay(ioms.BANK_INMONEY, ioms.PAYWAY_NETBANK, netbank)

	netout := &NetOutPay{mall: m}
	m.AddPay(ioms.BANK_OUTMONEY, ioms.PAYWAY_DEFAULT, netout)

	return m
}
