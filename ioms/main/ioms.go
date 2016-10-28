package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"tradebank/ioms"
)

func main() {
	file := flag.String("c", "ioms.conf", "configuration file")
	nodaemon := flag.Bool("e", false, "not run in daemon proccess")
	help := flag.Bool("h", false, "show help message")

	flag.Parse()
	if *help {
		flag.PrintDefaults()
		os.Exit(0)
	}
	if *nodaemon {
		// todo
		fmt.Print("force daemo\n")

	}

	// read configure
	conf, err := ioms.NewConfig(*file)
	if err != nil {
		fmt.Printf("LoadConfig failed, file = %s, errors = %s\n", file, err.Error())
		os.Exit(-1)
	}

	m := ioms.NewIomServer()
	m.Config = conf

	// setup logger
	//{"prefix":"filename", "switchsize":1024, "fieldir":"./", "filelevel":5, "termlevel":5}
	logConf := struct {
		prefix     string
		switchsize string
		filedir    string
		filelevel  string
		termlevel  string
	}{m.LogPrefix, m.LogDir, m.LogSwitchSize, m.LogFileLevel, m.LogTermLevel}
	var js, err = json.Marshal(logConf)
	if err != nil {
		fmt.Printf("Json marshal failed. errors = %s\n", err.Error())
		os.Exit(-1)
	}
	m.Log, err = logging.NewLogging(js)

	m.Log.Info("connecting to exch\n")
	// connect to exch
	err = m.ConnectExch()
	if err != nil {
		m.Log.Info("connect to exch failed. errors = %s\n", err.Error())
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

	<-m.ExitChan

}

func init() {

}
