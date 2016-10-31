package logging

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type fileLogger struct {
	prefix     string `json:"prefix"`
	fileDir    string `json:"filedir"`
	level      int64  `json:"level"`
	switchSize int64  `json:"switchsize"`
	switchTime int64  `json:"switchtime"`

	status bool

	file      *os.File
	fileDate  int64
	fileName  string
	fileSize  int64
	fileIndex int64
}

func (f *fileLogger) logSwitch() error {
	n := time.Now()
	if f.file == nil {
		f.fileName = fmt.Sprintf("%s_%d_%04d%02d%02d_%d.log.tmp",
			f.prefix, os.Getpid(), n.Year(), n.Month(), n.Day(), f.fileIndex)
		f.fileDate = time.Date(n.Year(), n.Month(), n.Day(), 0, 0, 0, 0, time.Local).Unix()

		var err error
		f.file, err = os.OpenFile(f.fileName, os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			return err
		}
		f.fileIndex = 0
		f.fileSize = 0

		return nil
	}
	switchFlag := true
	curDate := time.Date(n.Year(), n.Month(), n.Day(), 0, 0, 0, 0, time.Local).Unix()

	if f.switchTime > 0 {
		swt := 86400 + f.switchTime
		cur := int64(n.Hour()*3600+n.Minute()*60+n.Second()) + curDate
		if cur-f.fileDate >= swt {
			switchFlag = true
		}
	}
	if f.switchSize > 0 {
		if f.fileSize >= f.switchSize {
			switchFlag = true
		}
	}

	if switchFlag {
		f.Close()

		if curDate == f.fileDate {
			f.fileIndex++
			return f.logSwitch()
		}
		return f.logSwitch()
	}
	return nil
}

func (f *fileLogger) Name() string {
	return "file"
}

//`{"prefix"="hello", "filedir"="./", "level":"trace"}, "switchsize"=1024, "switchtime"="86400"}`
func (f *fileLogger) Open(conf string) error {
	err := json.Unmarshal([]byte(conf), *f)
	if err != nil {
		return err
	}
	f.fileDir = filepath.Dir(f.fileDir)
	if !strings.HasSuffix(f.fileDir, string(filepath.Separator)) {
		f.fileDir += string(filepath.Separator)
	}
	return nil
}
func (f *fileLogger) Write(msg *Message) (int, error) {
	n, err := 0, error(nil)
	if f.file != nil {
		if msg.msgType >= f.level {
			n, err = f.file.Write([]byte(msg.message))
			if err != nil {
				return n, err
			}
			f.fileSize += int64(n)
			err = f.logSwitch()
			if err != nil {
				return 0, err
			}
			return n, err
		}
	}
	return n, err
}
func (f *fileLogger) Close() error {
	err := f.file.Close()
	if err != nil {
		return err
	}
	f.file = nil
	return nil
}
func (f *fileLogger) Sync() error {
	if f.file != nil {
		return f.file.Sync()
	}
	return nil
}

func init() {

	f := &fileLogger{}
	Register(f)
}
