package ioms

import (
	"fmt"
	"strconv"
	"strings"

	ini "github.com/vaughan0/go-ini"
)

const (
	PROTO_TCP_SHORT = iota
	PROTO_TCP_LONG
	PROTO_HTTP
	PROTO_HTTPS
)

type Config struct {
	FileConf string

	LogPrefix     string
	LogDir        string
	LogSwitchTime int64
	LogFileLevel  int64
	LogTermLevel  int64

	ListenAddr  string
	ListenPort  int64
	ExchAddr    string
	ExchPort    int64
	BankAddr    string
	BankPort    int64
	ControlAddr string
	ControlPort int64

	BankProto int64

	Banks []string

	TimerValue   int64
	TimeReconn   int64
	TimeOutValue int64
}

func LoadConfig(path string) (*Config, error) {
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
	c.LogSwitchTime, err = int64(strconv.Atoi(str))
	if err != nil {
		return nil, fmt.Errorf("convert %s to interger failed.", str)
	}

	str, ok = f.Get("COMMON", "LOG_LEVEL_FILE")
	if !ok {
		return nil, fmt.Errorf("missing configure, sec=COMMON, key=LOG_LEVEL_FILE")
	}
	c.LogFileLevel, err = int64(strconv.Atoi(str))
	if err != nil {
		return nil, fmt.Errorf("convert %s to interger failed.", str)
	}

	str, ok = f.Get("COMMON", "LOG_LEVEL_TERM")
	if !ok {
		return nil, fmt.Errorf("missing configure, sec=COMMON, key=LOG_LEVEL_TERM")
	}
	c.LogTermLevel, err = int64(strconv.Atoi(str))
	if err != nil {
		return nil, fmt.Errorf("convert %s to interger failed.", str)
	}

	// net
	c.ExchAddr, ok = f.Get("COMMON", "EXCH_IP")
	if !ok {
		return nil, fmt.Errorf("missing configure, sec=COMMON, key=EXCH_IP")
	}
	str, ok = f.Get("COMMON", "EXCH_PORT")
	if !ok {
		return nil, fmt.Errorf("missing configure, sec=COMMON, key=EXCH_PORT")
	}
	c.ExchPort = int64(strconv.Atoi(str))
	if err != nil {
		return nil, fmt.Errorf("convert %s to interger failed.", str)
	}

	c.ListenAddr, ok = f.Get("COMMON", "LISTEN_IP")
	if !ok {
		return nil, fmt.Errorf("missing configure, sec=COMMON, key=LISTEN_IP")
	}
	str, ok = f.Get("COMMON", "LISTEN_PORT")
	if !ok {
		return nil, fmt.Errorf("missing configure, sec=COMMON, key=LISTEN_PORT")
	}
	c.ListenPort = int64(strconv.Atoi(str))
	if err != nil {
		return nil, fmt.Errorf("convert %s to interger failed.", str)
	}

	c.BankAddr, ok = f.Get("COMMON", "BANK_IP")
	if !ok {
		return nil, fmt.Errorf("missing configure, sec=COMMON, key=BANK_IP")
	}
	str, ok = f.Get("COMMON", "CONTROL_PORT")
	if !ok {
		return nil, fmt.Errorf("missing configure, sec=COMMON, key=BANK_PORT")
	}
	c.BankPort = int64(strconv.Atoi(str))
	if err != nil {
		return nil, fmt.Errorf("convert %s to interger failed.", str)
	}

	c.ControlAddr, ok = f.Get("COMMON", "CONTROL_IP")
	if !ok {
		return nil, fmt.Errorf("missing configure, sec=COMMON, key=CONTROL_IP")
	}
	str, ok = f.Get("COMMON", "CONTROL_PORT")
	if !ok {
		return nil, fmt.Errorf("missing configure, sec=COMMON, key=CONTROL_PORT")
	}
	c.ControlPort = int64(strconv.Atoi(str))
	if err != nil {
		return nil, fmt.Errorf("convert %s to interger failed.", str)
	}

	str, ok = f.Get("COMMON", "BANK_PROTO")
	if !ok {
		return nil, fmt.Errorf("missing configure, sec=COMMON, key=BANK_PROTO")
	}
	str = strings.ToLower(str)
	switch str {
	case "tcp_long":
		c.BankProto = int64(PROTO_TCP_LONG)
	case "tcp_short":
		c.BankProto = int64(PROTO_TCP_SHORT)
	case "http":
		c.BankProto = int64(PROTO_HTTP)
	case "https":
		c.BankProto = int64(PROTO_HTTPS)
	default:
		return nil, fmt.Errorf("unknown proto %s", str)

	}

	// timer
	str, ok = f.Get("COMMON", "TIMEOUT_TIME")
	if !ok {
		return nil, fmt.Errorf("missing configure, sec=COMMON, key=TIMEOUT_TIME")
	}
	c.TimeOutValue = int64(strconv.Atoi(str))
	if err != nil {
		return nil, fmt.Errorf("convert %s to interger failed.", str)
	}
	str, ok = f.Get("COMMON", "TIMER_INTERVAL")
	if !ok {
		return nil, fmt.Errorf("missing configure, sec=COMMON, key=TIMER_INTERVAL")
	}
	c.TimerValue = int64(strconv.Atoi(str))
	if err != nil {
		return nil, fmt.Errorf("convert %s to interger failed.", str)
	}
	str, ok = f.Get("COMMON", "RECONNECT_INTERVAL")
	if !ok {
		return nil, fmt.Errorf("missing configure, sec=COMMON, key=RECONNECT_INTERVAL")
	}
	c.TimeReconn = int64(strconv.Atoi(str))
	if err != nil {
		return nil, fmt.Errorf("convert %s to interger failed.", str)
	}
	// busi
	str, ok = f.Get("COMMON", "BANK_BUSI")
	if !ok {
		return nil, fmt.Errorf("missing configure, sec=COMMON, key=BANK_BUSI")
	}
	c.Banks = strings.Split(str, "|")
	if len(c.Banks) <= 0 {
		return nil, fmt.Errorf("invalid configure, sec=COMMON, key=BANK_BUSI")
	}

	return &c, nil
}
