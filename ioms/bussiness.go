package ioms

import (
	"fmt"
	"net"
	"os"
	"os/sinal"
	"syscall"
	"time"

	"tradebank/ioms/bank"
	"tradebank/logging"
	"tradebank/util"
)

type Handler func(int, interface{}) int

type exchMsg struct {
	conn    net.Conn
	command int64
	message []byte
}

type exchConn struct {
	conn      net.Conn
	conStatus bool
	regStatus bool

	recvChan chan *exchMsg
	sendChan chan *exchMsg
}

type IomServer struct {
	*Config

	ExitChan chan struct{}

	exchCtx *exchConn

	Log *logging.Log
}

func (m *IomServer) ListenBank() error {
	return nil
}

func (m *IomServer) ExchRecv() {

}
func (m *IomServer) ExchSend() {

}

func (m *IomServer) ExchTimer() {
	for {
		if !m.exchCtx.conStatus {
			m.Log.Info("reconnect to exch, ADDR=%s:%d\n", m.ExchAddr, m.ExchPort)
			conn, err := net.DialTimeout("tcp",
				fmt.Sprint("%s:%d", m.ExchAddr, m.ExchPort), m.TimeOutValue*time.Second)
			if err != nil {
				m.Log.Info("reconnect to exch failed, ADDR=%s:%d, ERR=%s\n",
					m.ExchAddr, m.ExchPort, err.Error())
				time.Sleep(m.TimeReconn * time.Second)
				continue
			}

			conn.SetDeadline(m.TimeOutValue, time.Second)
			(*net.TCPConn)(conn).SetNoDelay(true)
			(*net.TCPConn)(conn).SetKeepAlive(true)

			m.exchConn.conn = conn
			m.exchCtx.conStatus = true
			m.exchCtx.regStatus = false

			m.exchCtx.recvChan = make(chan *exchMsg, 1024)
			m.exchCtx.sendChan = make(chan *exchMsg, 1024)
			m.Log.Info("exch connected\n")

		}
		if !m.exchCtx.regStatus {
			// reg packet
			m.Log.Info("reg to exch\n")
			time.Sleep(m.TimeOutValue * time.Second)
			continue
		}

		time.Sleep(m.TimeReconn * time.Second)
	}
}
func (m *IomServer) ConnectExch() error {
	m.Log.Info("connect to exch, ADDR=%s:%d\n", m.ExchAddr, m.ExchPort)

	conn, err := net.DialTimeout("tcp",
		fmt.Sprint("%s:%d", m.ExchAddr, m.ExchPort), m.TimeOutValue*time.Second)
	if err != nil {
		return err
	}
	conn.SetDeadline(m.TimeOutValue, time.Second)
	(*net.TCPConn)(conn).SetNoDelay(true)
	(*net.TCPConn)(conn).SetKeepAlive(true)

	m.exchConn.conn = conn
	m.exchCtx.conStatus = true
	m.exchCtx.regStatus = false

	m.exchCtx.recvChan = make(chan *exchMsg, 1024)
	m.exchCtx.sendChan = make(chan *exchMsg, 1024)

	m.Log.Info("exch connected\n")

	go m.ExchTimer()
	go m.ExchRecv()
	go m.ExchSend()

	return nil
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
				m.Log.Info("catch SIGUSR1.")
			case syscall.SIGUSR2:
				m.Log.Info("catch SIGUSR2.")
				m.Log.Sync()

			}
		}

	}

}
func (m *IomServer) Stop() {

}
func (m *IomServer) InitLog() (err error) {
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
func (m *IomServer) LoadBankConf() (err error) {
	for _, name := range m.Banks {
		m.Log.Info("init bank, name=%s", name)

		id, err := util.ID(name)
		if err != nil {
			return err
		}
		myank := bank.MyBank(nil)
		mybank, err = bank.Bank(id)
		if err != nil {
			return nil
		}
		err = myank.Init(m.FileConf)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *IomServer) Control() {

}

func NewIomServer() *IomServer {
	m := &IomServer{}
	m.exchCtx = make(exchConn)

	return m
}
