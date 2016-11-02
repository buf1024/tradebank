package ioms

import (
//ini "github.com/vaughan0/go-ini"
)

const (
	PROTO_TCP_SHORT = iota
	PROTO_TCP_LONG
	PROTO_HTTP
	PROTO_HTTPS
)

type Config struct {
	LogPrefix     string
	LogDir        string
	LogSwitchTime int64
	LogFileLevel  int64
	LogTermLevel  int64

	ListenAddress string
	ListenPort    int64

	ExchAddress string
	ExchPort    int64

	TimerValue   int64
	TimeOutValue int64
}

func NewConfig(file string) (*Config, error) {
	c := &Config{}

	c.ListenAddress = "123"

	return &c, nil
}
