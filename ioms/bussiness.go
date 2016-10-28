package ioms

import (
	"tradebank/logging"
)

type Handler func(int, interface{}) int

type NetContext struct {
}

type IomServer struct {
	*Config

	ExitChan     chan struct{}
	ExchSendChan chan *NetContext
	ExchRecvChan chan *NetContext

	ConnectStatus bool

	Log *Log
}

func (m *IomServer) ListenBank() error {
	return nil
}

func (m *IomServer) ExchRecv() {

}
func (m *IomServer) ExchSend() {

}

func (m *IomServer) StartRecon() {

}
func (m *IomServer) ConnectExch() {

}

func (m *IomServer) ExchHeartbeat() {

}

func (m *IomServer) NewBankReq() {

}

func (m *IomServer) NewExchReq() {

}

func (m *IomServer) IomTimer() {

}

func (m *IomServer) SetupSignal() {

}

func (m *IomServer) ControlTrace() {

}

func NewIomServer() *IomServer {
	return &IomServer{}
}
