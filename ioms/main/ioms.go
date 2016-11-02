package main

import (
	"flag"
	"fmt"
	"os"

	"tradebank/ioms"
	"tradebank/util"
)

func main() {
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

	m := ioms.NewIomServer()

	err := error(nil)
	// read configure
	m.Config, err = ioms.NewConfig(*file)
	if err != nil {
		fmt.Printf("LoadConfig failed, file = %s, ERR=%s\n", file, err.Error())
		os.Exit(-1)
	}

	// setup logging
	err = m.SetupLog()
	if err != nil {
		fmt.Printf("SetupLog file failed. ERR=%s\n", err.Error())
		os.Exit(-1)
	}

	// connect to exch
	m.Log.Info("connecting to exch\n")
	err = m.ConnectExch()
	if err != nil {
		m.Log.Info("connect to exch failed. ERR=%s\n", err.Error())
		os.Exit(-1)
	}
	m.Log.Info("start exch reconnect go routine\n")
	go m.StartRecon()

	// listen bank
	go m.ListenBank()

	// exch hearbeat
	go m.ExchHeartbeat()

	// timer call
	go m.IomTimer()

	// control and trace
	go m.ControlTrace()

	// handle signal
	go m.HandleSignal()
	//<-m.ExitChan

}

func init() {

}
