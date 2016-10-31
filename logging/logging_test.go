package logging

import (
	"testing"
	"time"
)

func TestLogging(t *testing.T) {
	log, err := NewLogging()
	if err != nil {
		t.Errorf("NewLogging failed. err = %s\n", err.Error())
	}
	f := AddLog("file")
	if f == nil {
		t.Errorf("get file logger failed.")
	}
	c := AddLog("console")
	if c == nil {
		t.Errorf("get console logger failed.")
	}
	err = f.Open(`{"prefix"="hello", "filedir"="./", "level":"trace"}, "switchsize"=1024, "switchtime"="86400"}`)
	if err != nil {
		t.Errorf("open file logger failed. err = %s\n", err.Error())
	}
	err = c.Open(`{"level":"debug"}`)
	if err != nil {
		t.Errorf("open console logger failed. err = %s\n", err.Error())
	}

	running := 86400

	for {
		log.Trace("trace\n")
		log.Debug("debug\n")
		log.Info("info\n")
		log.Notice("notice\n")
		log.Warning("warning\n")
		log.Error("error\n")
		log.Critical("critical\n")

		time.Sleep(1 * time.Second)
		running--
		if running <= 0 {
			break
		}
	}
	f.Close()
	c.Close()
}
