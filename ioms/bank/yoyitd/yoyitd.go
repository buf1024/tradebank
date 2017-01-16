package yoyitd

import (
	"tradebank/ioms"

	"github.com/golang/protobuf/proto"
)

type YoyitdServer struct {
	*ioms.Server
}

func (b *YoyitdServer) Name() string {
	return "YoyitdServer"
}

func (b *YoyitdServer) ID() int64 {
	return 26
}

func (b *YoyitdServer) Init(s *ioms.Server) error {
	return nil
}

func (b *YoyitdServer) ExchReq(command int64, msg proto.Message) ([]ioms.Action, error) {
	return nil, nil
}
func (b *YoyitdServer) ExchRsp(command int64, msg proto.Message) ([]ioms.Action, error) {
	return nil, nil
}

func (b *YoyitdServer) BankReq([]byte) ([]ioms.Action, error) {
	return nil, nil
}
func (b *YoyitdServer) BankRsp([]byte) ([]ioms.Action, error) {
	return nil, nil
}

// NewServer
func NewServer() *YoyitdServer {
	m := &YoyitdServer{}

	return m
}
