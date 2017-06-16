package iomsframe

import (
	"github.com/golang/protobuf/proto"
)

const (
	PAYWAY_DEFAULT = 0
	PAYWAY_NOCARD  = 1
	PAYWAY_NETBANK = 3
)

const (
	BANK_OUTMONEY = 1
	BANK_INMONEY  = 2
)

// MyBank define the interface of exch
type MyBank interface {
	Name() string
	ID() int64

	ExchReq(command int64, msg proto.Message) error
	ExchRsp(command int64, msg proto.Message) error

	InitBank(m *ExchFrame) error
	StopBank(m *ExchFrame)
}
