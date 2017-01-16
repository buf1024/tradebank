package ioms

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"tradebank/logging"
	"tradebank/util"
)

type exchMsg struct {
	command int64
	message []byte
}

type exchConn struct {
	conn      *net.TCPConn
	conStatus bool
	regStatus bool

	recvChan chan *exchMsg
	sendChan chan *exchMsg
}

// Server represents the inout money server
type Server struct {
	*Config

	Log *logging.Log

	exch exchConn

	sessChan  chan struct{}
	timerChan chan struct{}

	cmdChan  chan string
	exitChan chan struct{}
}

const (
	statusInited = iota
	statusStarted
	statusReady
	statusStoped
)

func (m *Server) bankListen() error {
	return nil
}

func (m *Server) exchHandleMsg() {
END:
	for {
		select {
		// case msg := <-m.exch.recvChan:
		// 	{
		// 	}
		// case msg := <-m.exch.sendChan:
		// 	{

		// 	}
		default:
			// chan close
			break END

		}
	}
}

func (m *Server) exchConnect() error {
	m.Log.Info("connect to exch, addr = %s:%d\n", m.ExchAddr, m.ExchPort)
	conn, err := net.DialTimeout("tcp",
		fmt.Sprintf("%s:%d", m.ExchAddr, m.ExchPort),
		time.Duration(m.TimeOutValue*int64(time.Second)))
	if err != nil {
		m.Log.Error("connect to exch failed, add = %s:%d, err=%s\n",
			m.ExchAddr, m.ExchPort, err)
		return err
	}

	(conn).(*net.TCPConn).SetNoDelay(true)
	(conn).(*net.TCPConn).SetKeepAlive(true)

	m.exch.conn = conn.(*net.TCPConn)
	m.exch.conStatus = true
	m.exch.regStatus = false

	if m.exch.recvChan != nil {
		close(m.exch.recvChan)
	}
	if m.exch.sendChan != nil {
		close(m.exch.sendChan)
	}

	m.exch.recvChan = make(chan *exchMsg, 1024)
	m.exch.sendChan = make(chan *exchMsg, 1024)

	go m.exchHandleMsg()

	m.Log.Info("exch connected\n")

	return nil
}
func (m *Server) exchRegister() error {
	// reg packet
	// m.Log.Info("reg to exch\n")
	// msg, err := proto.Message(proto.CMD_SVR_REG_REQ)
	// if err != nil {
	// 	m.Log.Critical("create message failed, ERR=%s\n", err.Error())
	// 	return err
	// }

	// req := (*bankmsg.SvrRegReq)(msg)
	// *req.SID = "123"
	// *req.SvrType = bankid
	// *req.SvrId = bankid

	// reqMsg := &exchMsg{}
	// reqMsg.conn = m.exchCtx.conn
	// reqMsg.command = proto.CMD_SVR_REG_REQ
	// reqMsg.message, err = proto.Serialize(req)
	// if err != nil {
	// 	m.Log.Critical("serialize message failed, ERR=%s\n", err.Error())
	// 	return err
	// }

	// return err
	return nil
}

func (m *Server) timerOnce(to int64, typ string, ch interface{}, cmd interface{}) {
	t := time.NewTimer((time.Duration)((int64)(time.Second) * to))
	<-t.C

	switch {
	case typ == "exchreconnect":
		exchChan := (ch).(chan string)
		exchMsg := (cmd).(string)

		exchChan <- exchMsg
	}
}

// task represents exchange connect go routine
func (m *Server) cmdTask() {
	for {
		cmd := <-m.cmdChan
		switch cmd {
		case "connect":
			{
				err := m.exchConnect()
				if err != nil {
					m.Stop()
					return
				}
				m.cmdChan <- "register"
			}
		case "reconnect":
			{
				err := m.exchConnect()
				if err != nil {
					m.Log.Info("reconnect to addr=%s:%d\n after %d second",
						m.ExchAddr, m.ExchPort, m.TimeReconn)
					go m.timerOnce(m.TimeReconn, "exchreconnect", m.cmdChan, "reconnect")
					continue
				}
				m.cmdChan <- "register"
			}
		case "register":
			{
				err := m.exchRegister()
				if err != nil {
					// err
				} else {
					// post
				}
			}
		case "listen":
			{
				m.Log.Info("start to listen, addr = %s:%d\n", m.BankAddr, m.BankPort)
				err := m.bankListen()
				if err != nil {
					m.Log.Critical("listen failed, err = %s\n", err)
					m.Stop()
				} else {
					m.Log.Info("listen success")
				}
			}
		case "stop":
			{
				fmt.Printf("stop")
			}
		}
	}
}

func (m *Server) sigTask() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGUSR2, syscall.SIGUSR1)
END:
	for {
		select {
		case sig := <-ch:
			switch sig {
			case syscall.SIGHUP:
				m.Log.Info("catch SIGHUP, stop the server.")
				m.Stop()
				break END
			case syscall.SIGTERM:
				m.Log.Info("catch SIGTERM, stop the server.")
				m.Stop()
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

// initLog Setup Logger
func (m *Server) initLog() (err error) {
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

// // initBank setup the bank configure
func (m *Server) initBank() error {
	m.Log.Info("init bank, name=%s", m.Bank)

	myank := GetBank(m.Bank)
	if myank == nil {
		return fmt.Errorf("bank %s not found", m.Bank)
	}
	err := myank.Init(m)
	if err != nil {
		return err
	}
	return nil
}

// parseArgs parse the commandline arguments
func (m *Server) parseArgs() {
	file := flag.String("c", "ioms.conf", "configuration file")
	trace := flag.Int64("e", 0, "print error message")
	help := flag.Bool("h", false, "show help message")

	flag.Parse()
	if *help {
		flag.PrintDefaults()
		os.Exit(0)
	}
	if *trace > 0 {
		fmt.Printf("%s\n", util.NewError(*trace))
		os.Exit(0)
	}
	_, err := os.Stat(*file)
	if err != nil && os.IsNotExist(err) {
		fmt.Printf("file %s not exists\n", *file)
		os.Exit(-1)
	}
	m.FileConf = *file
}

// StartTimer start the timer
func (m *Server) StartTimer(value int64, handler TimerHandler) int64 {
	return 1
}

// StopTimer stop the timer
func (m *Server) StopTimer(id int64) {

}

// InitServer init the io money server
func (m *Server) InitServer() {
	m.parseArgs()

	err := error(nil)

	// read configure
	m.Config, err = LoadConfig(m.FileConf)
	if err != nil {
		fmt.Printf("LoadConfig failed, file = %s, ERR=%s\n", m.FileConf, err.Error())
		os.Exit(-1)
	}

	// init logging
	err = m.initLog()
	if err != nil {
		fmt.Printf("init log failed. ERR=%s\n", err.Error())
		os.Exit(-1)
	}

	// load bank configure
	err = m.initBank()
	if err != nil {
		m.Log.Critical("load bank configure failed. ERR=%s\n", err.Error())
		m.Stop()
		os.Exit(-1)
	}

	go m.sigTask()

	m.cmdChan = make(chan string, 1024)
	m.Log.Info("start task go routine\n")
	go m.cmdTask()

	m.exitChan = make(chan struct{})
}

// Start the server
func (m *Server) Start() {
	m.cmdChan <- "listen"
	m.cmdChan <- "connect"
}

// Stop the server
func (m *Server) Stop() {

}

// Wait wait the server to stop
func (m *Server) Wait() {
	<-m.exitChan

}
