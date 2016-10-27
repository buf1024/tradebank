package main

import (
	"os"

	"tradebank/bankmsg"
	//	"tradebank/ioms"
)

var eArgs = os.Args

var m IomServer

func main() {
	// parse args

	// read configure
	conf, err := LoadConfig()

	// connect to exch
	go m.ExchRecon()

	// listen bank
	go m.ListenBank()

	// exch hearbeat
	go m.ExchHeartbeat()

	// timer call
	go m.IomTimer()

	// control and trace
	go m.ControlTrace()

	<-m.ExitChan
}

func init() {

}
