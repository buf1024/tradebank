package ioms

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"tradebank/logging"
	"tradebank/util"
)

type exchMsg struct {
	conn    net.Conn
	command int64
	message []byte
}

type ExchConn struct {
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

	Log *logging.Log

	sessChan  chan struct{}
	timerChan chan struct{}

	cmdChan chan string
}

func (m *Server) exchConnect() error {
	m.Log.Info("reconnect to exch, ADDR=%s:%d\n", m.ExchAddr, m.ExchPort)
	// conn, err := net.DialTimeout("tcp",
	// 	fmt.Sprintf("%s:%d", m.ExchAddr, m.ExchPort), time.Duration(m.TimeOutValue*int64(time.Second)))
	// if err != nil {
	// 	m.Log.Info("reconnect to exch failed, ADDR=%s:%d, ERR=%s\n",
	// 		m.ExchAddr, m.ExchPort, err.Error())
	// 	return err
	// }

	// conn.SetDeadline(m.TimeOutValue, time.Second)
	// (*net.TCPConn)(conn).SetNoDelay(true)
	// (*net.TCPConn)(conn).SetKeepAlive(true)

	// m.exchConn.conn = conn
	// m.exchCtx.conStatus = true
	// m.exchCtx.regStatus = false

	// m.exchCtx.recvChan = make(chan *exchMsg, 1024)
	// m.exchCtx.sendChan = make(chan *exchMsg, 1024)
	// m.Log.Info("exch connected\n")

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

// task represents exchange connect go routine
func (m *Server) task() {
	for {
		cmd := <-m.cmdChan
		switch cmd {
		case "connect":
			{
				err := m.exchConnect()
				if err != nil {
					// Sleep emit connect
				} else {
					// emit register
				}
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
		}
	}
}

func (m *Server) handleSignal() {
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

// // loadBankConf setup the bank configure
// func (m *Server) loadBankConf() (err error) {
// 	for _, name := range m.Banks {
// 		m.Log.Info("init bank, name=%s", name)

// 		id, err := util.ID(name)
// 		if err != nil {
// 			return err
// 		}
// 		myank := bank.MyBank(nil)
// 		mybank, err = bank.Bank(id)
// 		if err != nil {
// 			return nil
// 		}
// 		err = myank.LoadConfig(m.FileConf)
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

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
	// err = m.LoadBankConf()
	// if err != nil {
	// 	m.Log.Critical("load bank configure failed. ERR=%s\n", err.Error())
	// 	m.Stop()
	// 	os.Exit(-1)
	// }

	// connect to exch
	m.Log.Info("connecting to exch\n")
	// err = m.ConnectExch()
	// if err != nil {
	// 	m.Log.Critical("connect to exch failed. ERR=%s\n", err.Error())
	// 	m.Stop()
	// 	os.Exit(-1)
	// }
	m.Log.Info("start exch reconnect go routine\n")

}

// Start the server
func (m *Server) Start() {

	// go m.StartRecon()

	// // listen bank
	// go m.ListenBank()

	// // exch hearbeat
	// go m.ExchHeartbeat()

	// // timer call
	// go m.IomTimer()

	// // control and trace
	// go m.ControlTrace()

	// // handle signal
	// go m.HandleSignal()
	// //<-m.ExitChan
}

// Stop the server
func (m *Server) Stop() {

}
