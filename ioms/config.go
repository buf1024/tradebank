package ioms

import ()

const (
	PROTO_TCP_SHORT = iota
	PROTO_TCP_LONG
	PROTO_HTTP
	PROTO_HTTPS
)

//{"prefix":"filename", "switchsize":1024, "fieldir":"./", "filelevel":5, "termlevel":5}

type Config struct {
	LogPrefix     string
	LogDir        string
	LogSwitchSize string
	LogFileLevel  string
	LogTermLevel  string

	ListenAddress string
	ListenPort    int64

	UniqioAddress string
	UniqioPort    int64

	TimerValue   int64
	TimeOutValue int64
}

func NewConfig(file string) (*Config, error) {
	var c Config
	c.ListenAddress = "123"

	return &c, nil
}
