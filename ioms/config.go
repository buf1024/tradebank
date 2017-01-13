package ioms

import (
	"fmt"
	"strconv"
	"strings"

	ini "github.com/vaughan0/go-ini"
)

const (
	// ProtoTCPShort TCP short connection
	ProtoTCPShort = iota
	// ProtoTCPLong TCP long connection
	ProtoTCPLong
	// ProtoHTTP http connection
	ProtoHTTP
	// ProtoHTTPS https connection
	ProtoHTTPS
)

// Config configre struct
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

	Bank string

	TimeReconn   int64
	TimeOutValue int64
}

// LoadConfig load the configuration
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

	c.ListenAddr, ok = f.Get("COMMON", "LISTEN_IP")
	if !ok {
		return nil, fmt.Errorf("missing configure, sec=COMMON, key=LISTEN_IP")
	}
	str, ok = f.Get("COMMON", "LISTEN_PORT")
	if !ok {
		return nil, fmt.Errorf("missing configure, sec=COMMON, key=LISTEN_PORT")
	}
	i, err = strconv.Atoi(str)
	if err != nil {
		return nil, fmt.Errorf("convert %s to interger failed", str)
	}
	c.ListenPort = int64(i)

	c.BankAddr, ok = f.Get("COMMON", "BANK_IP")
	if !ok {
		return nil, fmt.Errorf("missing configure, sec=COMMON, key=BANK_IP")
	}
	str, ok = f.Get("COMMON", "CONTROL_PORT")
	if !ok {
		return nil, fmt.Errorf("missing configure, sec=COMMON, key=BANK_PORT")
	}
	i, err = strconv.Atoi(str)
	if err != nil {
		return nil, fmt.Errorf("convert %s to interger failed", str)
	}
	c.BankPort = int64(i)

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

	str, ok = f.Get("COMMON", "BANK_PROTO")
	if !ok {
		return nil, fmt.Errorf("missing configure, sec=COMMON, key=BANK_PROTO")
	}
	str = strings.ToLower(str)
	switch str {
	case "tcp-long":
		c.BankProto = int64(ProtoTCPLong)
	case "tcp-short":
		c.BankProto = int64(ProtoTCPShort)
	case "http":
		c.BankProto = int64(ProtoHTTP)
	case "https":
		c.BankProto = int64(ProtoHTTPS)
	default:
		return nil, fmt.Errorf("unknown proto %s", str)

	}

	// time
	str, ok = f.Get("COMMON", "TIMEOUT_TIME")
	if !ok {
		return nil, fmt.Errorf("missing configure, sec=COMMON, key=TIMEOUT_TIME")
	}
	i, err = strconv.Atoi(str)
	if err != nil {
		return nil, fmt.Errorf("convert %s to interger failed", str)
	}
	c.TimeOutValue = int64(i)

	str, ok = f.Get("COMMON", "RECONNECT_INTERVAL")
	if !ok {
		return nil, fmt.Errorf("missing configure, sec=COMMON, key=RECONNECT_INTERVAL")
	}
	i, err = strconv.Atoi(str)
	if err != nil {
		return nil, fmt.Errorf("convert %s to interger failed", str)
	}
	c.TimeReconn = int64(i)
	// busi
	str, ok = f.Get("COMMON", "BANK")
	if !ok {
		return nil, fmt.Errorf("missing configure, sec=COMMON, key=BANK")
	}
	c.Bank = str

	return c, nil
}
