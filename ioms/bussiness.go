package ioms

import ()

type Handler func(int, interface{}) int

type NetContext struct {
}

type IomServer struct {
	*Config

	ExitChan     chan struct{}
	ExchSendChan chan *NetContext
	ExchRecvChan chan *NetContext

	ConnectStatus bool
}

func (m *IomServer) ListenBank() error {
	m.Config
}

func (m *IomServer) ExchRecv() {

}
func (m *IomServer) ExchSend() {

}

func (m *IomServer) ExchRecon() {

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
