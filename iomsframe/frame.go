package iomsframe

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	pb "github.com/golang/protobuf/proto"
	ini "github.com/vaughan0/go-ini"

	"tradebank/logging"
	"tradebank/proto"
	"tradebank/util"
)

type exchMsg struct {
	command int64
	message []byte
	pbmsg   pb.Message
}

type exchConn struct {
	conn      *net.TCPConn
	conStatus bool
	regStatus bool

	recvChan chan *exchMsg
	sendChan chan *exchMsg
}

// Config configre struct
type Config struct {
	FileConf string

	LogPrefix     string
	LogDir        string
	LogSwitchTime int64
	LogFileLevel  int64
	LogTermLevel  int64

	ExchAddr string
	ExchPort int64
	//ControlAddr string
	//ControlPort int64

	TimeOutReconn int64
	TimeOutSess   int64
}

// ExchFrame represents the inout money server
type ExchFrame struct {
	*Config

	fileConf string

	Log *logging.Log

	exch exchConn

	sessChan chan struct{}

	cmdChan  chan string
	exitChan chan struct{}

	Bank MyBank
}

func (m *ExchFrame) exchHandleRecv() {
	err := error(nil)
	for {
		buf := make([]byte, proto.GetHeaderLength())
		_, err = m.exch.conn.Read(buf)
		if err != nil {
			m.Log.Error("read exch packet head failed. ERR=%s\n", err)
			break
		}
		var head *proto.MessageHeader
		head, err = proto.ParseHeader(buf)
		if err != nil {
			m.Log.Error("Parse Header failed. err:%s\n", err)
			break
		}
		buf = make([]byte, head.Length-uint32(proto.GetHeaderLength()))
		_, err = m.exch.conn.Read(buf)
		if err != nil {
			m.Log.Error("read exch packet head failed. ERR=%s\n", err)
			break
		}
		// TODO BUF 解密
		var msg pb.Message
		msg, err = proto.Parse(int64(head.Command), buf)
		if err != nil {
			m.Log.Error("parse message failed discard message. err=%s\n", err)
			continue
		}
		if head.Command != proto.CMD_HEARTBEAT_REQ && head.Command != proto.CMD_HEARTBEAT_RSP {
			m.Log.Info("RECV: %s\n", proto.Debug(int64(head.Command), msg))
		}

		if head.Command == proto.CMD_HEARTBEAT_REQ || head.Command == proto.CMD_SVR_REG_RSP {
			m.HandleDef(int64(head.Command), msg)
			continue
		}

		if head.Command%2 == 0 {
			go m.Bank.ExchRsp(int64(head.Command), msg)
		} else {
			go m.Bank.ExchReq(int64(head.Command), msg)
		}
	}
	if err != nil {
		m.exch.conStatus = false
		m.exch.regStatus = false
		m.Log.Info("try to reconnect to exch, after %d seconds\n", m.TimeOutReconn)
		go m.timerOnce(m.TimeOutReconn, "exchreconnect", m.cmdChan, "reconnect")

	}
}

func (m *ExchFrame) exchHandleSend() {
END:
	for {
		select {
		case msg, isOpen := <-m.exch.sendChan:
			{
				if !isOpen {
					m.Log.Error("exch receive chan is close.\n")
					break END
				}
				if msg.command != proto.CMD_HEARTBEAT_REQ && msg.command != proto.CMD_HEARTBEAT_RSP {
					m.Log.Info("SEND: %s\n", proto.Debug(int64(msg.command), msg.pbmsg))
				}
				n, err := m.exch.conn.Write(msg.message)
				if err != nil {
					m.Log.Error("write msg failed, err=%s\n", err)
					continue
				}
				if msg.command != proto.CMD_HEARTBEAT_REQ && msg.command != proto.CMD_HEARTBEAT_RSP {
					m.Log.Info("write %d byte to exch, cmd=0x%x\n", n, msg.command)
				}
			}
		}
	}
}

func (m *ExchFrame) exchHandleMsg() {
	go m.exchHandleSend()
	go m.exchHandleRecv()
}

func (m *ExchFrame) exchConnect() error {
	if m.exch.conn != nil {
		m.exch.conn.Close()
	}

	m.Log.Info("connect to exch, addr = %s:%d\n", m.ExchAddr, m.ExchPort)

	conn, err := net.DialTimeout("tcp",
		fmt.Sprintf("%s:%d", m.ExchAddr, m.ExchPort),
		time.Duration(m.TimeOutReconn*int64(time.Second)))
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
func (m *ExchFrame) exchRegister() error {
	//reg packet
	m.Log.Info("reg to exch\n")
	msg, err := proto.Message(proto.CMD_SVR_REG_REQ)
	if err != nil {
		m.Log.Critical("create message failed, ERR=%s\n", err.Error())
		return err
	}

	req := msg.(*proto.SvrRegReq)
	req.SID = pb.String(util.SID())
	req.SvrType = pb.Int32(int32(m.Bank.ID()))
	req.SvrId = pb.String(m.Bank.Name())

	return m.WriteMsg(proto.CMD_SVR_REG_REQ, req)

}

func (m *ExchFrame) MakeReq(command int64, msg pb.Message) error {
	if command%2 == 0 {
		command = command + 1
	}
	return m.WriteMsg(command, msg)
}
func (m *ExchFrame) MakeRsp(command int64, msg pb.Message) error {
	if command%2 != 0 {
		command = command + 1
	}
	return m.WriteMsg(command, msg)
}

func (m *ExchFrame) WriteMsg(command int64, msg pb.Message) error {

	reqMsg := &exchMsg{}
	reqMsg.command = command
	reqMsg.pbmsg = msg
	var err error
	reqMsg.message, err = proto.SerializeMessage(command, msg, false)
	if err != nil {
		m.Log.Critical("serialize message failed, ERR=%s\n", err.Error())
		return err
	}
	m.exch.sendChan <- reqMsg

	return nil
}

func (m *ExchFrame) timerOnce(to int64, typ string, ch interface{}, cmd interface{}) {
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
func (m *ExchFrame) cmdTask() {
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
						m.ExchAddr, m.ExchPort, m.TimeOutReconn)
					go m.timerOnce(m.TimeOutReconn, "exchreconnect", m.cmdChan, "reconnect")
					continue
				}
				m.cmdChan <- "register"
			}
		case "register":
			{
				err := m.exchRegister()
				if err != nil {
					m.Log.Error("register server failed, try to reconnect.\n")
					m.cmdChan <- "register"
				} else {
					// post
				}
			}
		}
	}
}

func (m *ExchFrame) sigTask() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR2, syscall.SIGUSR1)
	m.Log.Info("setup sig task.\n")
END:
	for {
		select {
		case sig := <-ch:
			switch sig {
			case syscall.SIGHUP:
				m.Log.Info("catch SIGHUP, stop the server.\n")
				m.Stop()
				break END
			case syscall.SIGTERM:
				m.Log.Info("catch SIGTERM, stop the server.\n")
				m.Stop()
				break END
			case syscall.SIGINT:
				m.Log.Info("catch SIGINT, stop the server.\n")
				m.Stop()
				break END
			case syscall.SIGUSR1:
				m.Log.Info("catch SIGUSR1.\n")
			case syscall.SIGUSR2:
				m.Log.Info("catch SIGUSR2.\n")
				m.Log.Sync()

			}
		}

	}

}

// initLog Setup Logger
func (m *ExchFrame) initLog() (err error) {
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
func (m *ExchFrame) initBank() error {
	if m.Bank == nil {
		return fmt.Errorf("bank interface is nil")
	}
	m.Log.Info("init bank, name=%s\n", m.Bank.Name())

	err := m.Bank.InitBank(m)
	if err != nil {
		return err
	}
	return nil
}

// parseArgs parse the commandline arguments
func (m *ExchFrame) parseArgs() {
	file := flag.String("c", "", "configuration file")
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
		fmt.Printf("file %s missiong or not exists\n", *file)
		os.Exit(-1)
	}
	m.fileConf = *file
}

// StartTimer start the timer
/*func (m *ExchFrame) StartTimer(value int64, handler TimerHandler) int64 {
	return 1
}
*/
// StopTimer stop the timer
func (m *ExchFrame) StopTimer(id int64) {

}

// InitServer init the io money server
func (m *ExchFrame) InitServer() {
	m.parseArgs()

	err := error(nil)

	// read configure
	m.Config, err = m.loadConfig(m.fileConf)
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
		m.Log.Critical("init bank failed. ERR=%s\n", err.Error())
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
func (m *ExchFrame) Start() {
	m.cmdChan <- "listen"
	m.cmdChan <- "connect"
	m.Wait()
}

// Stop the server
func (m *ExchFrame) Stop() {

	m.Log.Info("stop the server!\n")
	m.Log.Info("stop bank, cleanup!\n")
	m.Bank.StopBank(m)
	m.Log.Stop()

	m.exitChan <- struct{}{}

}

// Wait wait the server to stop
func (m *ExchFrame) Wait() {
	<-m.exitChan

}

// loadConfig load the configuration
func (m *ExchFrame) loadConfig(path string) (*Config, error) {
	c := &Config{}

	c.FileConf = path

	f, err := ini.LoadFile(path)
	if err != nil {
		return nil, err
	}

	ok := false
	str := ""

	// logging
	c.LogDir, ok = f.Get("COMMON", "LOG_DIR")
	if !ok {
		return nil, fmt.Errorf("missing configure, sec=COMMON, key=LOG_DIR")
	}

	c.LogPrefix, ok = f.Get("COMMON", "LOG_HEADER")
	if !ok {
		return nil, fmt.Errorf("missing configure, sec=COMMON, key=LOG_HEADER")
	}

	str, ok = f.Get("COMMON", "LOG_SWITCH_TIME")
	if !ok {
		return nil, fmt.Errorf("missing configure, sec=COMMON, key=LOG_SWITCH_TIME")
	}
	i, err := strconv.Atoi(str)
	if err != nil {
		return nil, fmt.Errorf("convert %s to interger failed", str)
	}
	c.LogSwitchTime = int64(i)

	str, ok = f.Get("COMMON", "LOG_LEVEL_FILE")
	if !ok {
		return nil, fmt.Errorf("missing configure, sec=COMMON, key=LOG_LEVEL_FILE")
	}
	i, err = strconv.Atoi(str)
	if err != nil {
		return nil, fmt.Errorf("convert %s to interger failed", str)
	}
	c.LogFileLevel = int64(i)

	str, ok = f.Get("COMMON", "LOG_LEVEL_TERM")
	if !ok {
		return nil, fmt.Errorf("missing configure, sec=COMMON, key=LOG_LEVEL_TERM")
	}
	i, err = strconv.Atoi(str)
	if err != nil {
		return nil, fmt.Errorf("convert %s to interger failed", str)
	}
	c.LogTermLevel = int64(i)

	// net
	c.ExchAddr, ok = f.Get("COMMON", "EXCH_IP")
	if !ok {
		return nil, fmt.Errorf("missing configure, sec=COMMON, key=EXCH_IP")
	}
	str, ok = f.Get("COMMON", "EXCH_PORT")
	if !ok {
		return nil, fmt.Errorf("missing configure, sec=COMMON, key=EXCH_PORT")
	}
	i, err = strconv.Atoi(str)
	if err != nil {
		return nil, fmt.Errorf("convert %s to interger failed", str)
	}
	c.ExchPort = int64(i)
	/*
		c.ControlAddr, ok = f.Get("COMMON", "CONTROL_IP")
		if !ok {
			return nil, fmt.Errorf("missing configure, sec=COMMON, key=CONTROL_IP")
		}
		str, ok = f.Get("COMMON", "CONTROL_PORT")
		if !ok {
			return nil, fmt.Errorf("missing configure, sec=COMMON, key=CONTROL_PORT")
		}
		i, err = strconv.Atoi(str)
		if err != nil {
			return nil, fmt.Errorf("convert %s to interger failed", str)
		}
		c.ControlPort = int64(i)
	*/
	// time
	str, ok = f.Get("COMMON", "TIMEOUT_TIME")
	if !ok {
		return nil, fmt.Errorf("missing configure, sec=COMMON, key=TIMEOUT_TIME")
	}
	i, err = strconv.Atoi(str)
	if err != nil {
		return nil, fmt.Errorf("convert %s to interger failed", str)
	}
	c.TimeOutSess = int64(i)

	str, ok = f.Get("COMMON", "RECONNECT_INTERVAL")
	if !ok {
		return nil, fmt.Errorf("missing configure, sec=COMMON, key=RECONNECT_INTERVAL")
	}
	i, err = strconv.Atoi(str)
	if err != nil {
		return nil, fmt.Errorf("convert %s to interger failed", str)
	}
	c.TimeOutReconn = int64(i)

	return c, nil
}

func (m *ExchFrame) HandleDef(command int64, msg pb.Message) error {
	switch command {
	case proto.CMD_HEARTBEAT_REQ:
		{
			req := msg.(*proto.HeartBeatReq)
			rspMsg, err := proto.Message(proto.CMD_HEARTBEAT_RSP)
			if err != nil {
				m.Log.Critical("create message failed, ERR=%s\n", err.Error())
				return err
			}

			rsp := rspMsg.(*proto.HeartBeatRsp)
			rsp.SID = pb.String(req.GetSID())

			m.MakeRsp(proto.CMD_HEARTBEAT_RSP, rsp)
		}
	case proto.CMD_SVR_REG_RSP:
		{
			m.exch.regStatus = true
		}
	case proto.CMD_E2B_SIGNINOUT_REQ:
		{
			req := msg.(*proto.E2BSignInOutReq)
			rspMsg, err := proto.Message(proto.CMD_E2B_SIGNINOUT_RSP)
			if err != nil {
				m.Log.Critical("create message failed, ERR=%s\n", err.Error())
				return err
			}

			rsp := rspMsg.(*proto.E2BSignInOutRsp)
			rsp.ExchSID = pb.String(req.GetExchSID())
			rsp.BankID = pb.Int32(req.GetBankID())
			rsp.Type = pb.Int32(req.GetType())
			rsp.RetCode = pb.Int32(int32(util.E_SUCCESS))
			rsp.RetMsg = pb.String(util.GetErrMsg(util.E_SUCCESS))

			m.MakeRsp(proto.CMD_E2B_SIGNINOUT_RSP, rsp)
		}
	default:
		{
			m.Log.Warning("unknown command 0x%x, discard packet\n", command)
		}
	}
	return nil
}

func (m *ExchFrame) Ready() bool {
	return m.exch.conStatus && m.exch.regStatus
}
