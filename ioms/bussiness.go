package ioms

import (
	"fmt"
	"os"
	"os/sinal"
	"syscall"

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

	Log *logging.Log
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
func (m *IomServer) ConnectExch() error {

	return nil
}

func (m *IomServer) ExchHeartbeat() {

}

func (m *IomServer) NewBankReq() {

}

func (m *IomServer) NewExchReq() {

}

func (m *IomServer) IomTimer() {

}

func (m *IomServer) HandleSignal() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGUSR2, syscall.SIGUSR1)
END:
	for {
		select {
		case sig := <-ch:
			switch sig {
			case syscall.SIGHUP:
				m.Log.Info("catch SIGHUP, stop the server.")
				m.StopIomServer()
				break END
			case syscall.SIGTERM:
				m.Log.Info("catch SIGTERM, stop the server.")
				m.StopIomServer()
				break END
			case syscall.SIGUSR1:
			case syscall.SIGUSR2:

			}
		}

	}

}
func (m *IomServer) StopIomServer() {

}
func (m *IomServer) SetupLog() (err error) {
	m.Log, err = logging.NewLogging()
	if err != nil {
		return err
	}
	_, err = logging.SetupLog("file",
		fmt.Sprintf(`{"prefix":"%s", "filedir":"%s", "level":%d, "switchsize":%d, "switchtime":%d}`,
			m.LogPrefix, m.LogDir, m.LogFileLevel, -1, m.LogSwitchTime))
	if err != nil {
		return err
	}
	_, err = logging.SetupLog("console",
		fmt.Sprintf(`{"level":%d}`, m.LogTermLevel))
	if err != nil {
		return err
	}
	m.Log.Start()

	m.Log.Info("log started.\n")

	return nil
}

func (m *IomServer) Control() {

}

func NewIomServer() *IomServer {
	return &IomServer{}
}
