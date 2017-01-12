package ioms

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"tradebank/ioms/bank"
	"tradebank/logging"
	"tradebank/proto"
	"tradebank/util"
)

//type Handler func(int, interface{}) int

type exchMsg struct {
	conn    net.Conn
	command int64
	message []byte
}

type Conn struct {
	conn      net.Conn
	Status    bool
	regStatus bool

	recvChan chan *exchMsg
	sendChan chan *exchMsg
}

// Server represents the inout money server
type Server struct {
	*Config

	ExitChan chan struct{}

	exchCtx *exchConn
	Log     *logging.Log

	sessChan  chan struct{}
	timerChan chan struct{}
}

func (m *Server) ListenBank() error {
	return nil
}

func (m *Server) ExchRecv() {

}
func (m *Server) ExchSend() {

}

func (m *Server) ExchTimer() {
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
			msg, err := proto.Message(proto.CMD_SVR_REG_REQ)
			if err != nil {
				m.Log.Fatal("create message failed, ERR=%s\n", err.Error())
				time.Sleep(m.TimeOutValue * time.Second)
				continue
			}

			req := (*bankmsg.SvrRegReq)(msg)
			*req.SID = "123"
			*req.SvrType = bankid
			*req.SvrId = bankid

			reqMsg := &exchMsg{}
			reqMsg.conn = m.exchCtx.conn
			reqMsg.command = proto.CMD_SVR_REG_REQ
			reqMsg.message, err = proto.Serialize(req)
			if err != nil {
				m.Log.Fatal("serialize message failed, ERR=%s\n", err.Error())
				time.Sleep(m.TimeOutValue * time.Second)
				continue
			}

			time.Sleep(m.TimeOutValue * time.Second)
			continue
		}

		time.Sleep(m.TimeReconn * time.Second)
	}
}
func (m *Server) ConnectExch() error {
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

func (m *Server) NewBankReq() {

}

func (m *Server) NewExchReq() {

}

func (m *Server) IomTimer() {

}

func (m *Server) HandleSignal() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGUSR2, syscall.SIGUSR1)
END:
	for {
		select {
		case sig := <-ch:
			switch sig {
			case syscall.SIGHUP:
				m.Log.Info("catch SIGHUP, stop the server.")
				m.StopServer()
				break END
			case syscall.SIGTERM:
				m.Log.Info("catch SIGTERM, stop the server.")
				m.StopServer()
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
func (m *Server) Stop() {

}
func (m *Server) InitLog() (err error) {
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
func (m *Server) LoadBankConf() (err error) {
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
		err = myank.LoadConfig(m.FileConf)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Server) Control() {

}

func NewServer() *Server {
	m := &Server{}
	m.exchCtx = make(exchConn)

	return m
}
