package main

import (
	"tradebank/proto"

	ini "github.com/vaughan0/go-ini"
)

type YaodePay interface {
	Init(*ini.File) error

	InMoneyReq(req *proto.E2BInMoneyReq) error
	OutMoneyReq(req *proto.E2BOutMoneyReq) error
	VerifyReq(req *proto.E2BVerifyCodeReq) error

	CheckReq(orderId string) (int32, error)
}

const (
	QUERYRESULT_RETRY_TIMES = 5
)

type QueryResultContext struct {
	extflow    string
	orderId    string
	amount     string
	bankacct   string
	retryTimes int64
}

type PayReq struct {
	mchKey      string
	merId       string
	orderId     string
	transAmount string
	cardByName  string
	cardByNo    string
	cerNumber   string
	mobile      string

	transChanlName string
	pageNotifyUrl  string
	backNotifyUrl  string

	transBody string
}
type VerifyReq struct {
	mchKey       string
	merId        string
	yzm          string
	ksPayOrderId string
}

type QueryReq struct {
	mchKey    string
	merId     string
	orderId   string
	transDate string
}
type NotifyReq struct {
	orderId     string
	transAmount string
	payStatus   string
	payMsg      string
}
type PayRsp struct {
	status        string // 00：成功 01：失败 02：系统错误
	orderId       string
	ksPayOrderId  string
	chanelRefcode string // 89   要求手机验证码
	bankOrderId   string
	refCode       string // ‘00’交易成功 01’预交易成功 ‘02’交易失败 03  交易处理中
	refMsg        string

	respBody string
}
