package ioms

import ()

const (
	PROTO_TCP_SHORT = itoa
	PROTO_TCP_LONG
	PROTO_HTTP
	PROTO_HTTPS
)

type Config struct {
	ListenAddress string
	ListenPort    int64

	UniqioAddress string
	UniqioPort    int64

	TimerValue   int64
	TimeOutValue int64
}

func LoadConfig(file string) (*Config, error) {

}
