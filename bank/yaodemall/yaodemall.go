package main

import (
	"tradebank/iomsframe"

	"github.com/golang/protobuf/proto"
)

type YaodeMall struct {
	iomsframe.ExchFrame
}

func (b *YaodeMall) Name() string {
	return "YaodeMall"
}

func (b *YaodeMall) ID() int64 {
	return 56
}

func (b *YaodeMall) InitBank(m *iomsframe.ExchFrame) error {
	return nil
}
func (b *YaodeMall) StopBank(m *iomsframe.ExchFrame) {

}
func (b *YaodeMall) ExchReq(command int64, msg proto.Message) error {
	return nil
}
func (b *YaodeMall) ExchRsp(command int64, msg proto.Message) error {
	return nil
}

// YaodeMallServer
func YaodeMallServer() *YaodeMall {
	m := &YaodeMall{}
	m.Bank = m

	return m
}
