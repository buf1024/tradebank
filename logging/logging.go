package logging

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"
)

const (
	LevelAll = iota
	LevelTrace
	LevelDebug
	LevelInformational
	LevelNotice
	LevelWarning
	LevelError
	LevelCritical
)

const (
	defBufferSize = 1024
)
const (
	statusReady = iota
	statusRunning
	statusClosing
	statusClosed
)

type Loger interface {
	Name() string
	Open(conf string) error
	Write(msg *Message) (int64, error)
	Close() error
}

type logConfig struct {
	prefix     string `json:"prefix"`
	switchSize int64  `json:"switchsize"`

	termLevel    int64
	termLevelStr string `json:"termlevel"`

	fileDir      string `json:"filedir"`
	fileLevel    int64
	fileLevelStr string `json:"filelevel"`

	bufSize int64 `json:"buffersize"`
}

type Message struct {
	msgType int64
	message string
}

type Log struct {
	*logConfig

	status int64

	sync   bool
	mutex  sync.Mutex
	logMsg chan *Message
	sigMsg chan string
}

var lvlStr map[string]int64
var lvlHead map[int64]string
var logger map[string]*Loger

func (l *Log) Critical(format ...interface{}) error {
	now := time.Now()
	msg := fmt.Sprintf("[%02d%02d%02d.%06d][%s] ",
		now.Hour(), now.Minute(), now.Second(), now.Nanosecond(),
		lvlHead[LevelCritical])
	l.logMessage(msg, format...)
}

func (l *Log) Error(format ...interface{}) error {
	now := time.Now()
	msg := fmt.Sprintf("[%02d%02d%02d.%06d][%s] ",
		now.Hour(), now.Minute(), now.Second(), now.Nanosecond(),
		lvlHead[LevelError])
	l.logMessage(msg, format...)
}

func (l *Log) Warning(format ...interface{}) error {
	now := time.Now()
	msg := fmt.Sprintf("[%02d%02d%02d.%06d][%s] ",
		now.Hour(), now.Minute(), now.Second(), now.Nanosecond(),
		lvlHead[LevelWarning])
	l.logMessage(msg, format...)
}
func (l *Log) Notice(format ...interface{}) error {
	now := time.Now()
	msg := fmt.Sprintf("[%02d%02d%02d.%06d][%s] ",
		now.Hour(), now.Minute(), now.Second(), now.Nanosecond(),
		lvlHead[LevelNotice])
	l.logMessage(msg, format...)
}

func (l *Log) Info(format ...interface{}) error {
	now := time.Now()
	msg := fmt.Sprintf("[%02d%02d%02d.%06d][%s] ",
		now.Hour(), now.Minute(), now.Second(), now.Nanosecond(),
		lvlHead[LevelInformational])
	l.logMessage(msg, format...)
}

func (l *Log) Debug(format ...interface{}) error {
	now := time.Now()
	msg := fmt.Sprintf("[%02d%02d%02d.%06d][%s] ",
		now.Hour(), now.Minute(), now.Second(), now.Nanosecond(),
		lvlHead[LevelDebug])
	l.logMessage(msg, format...)
}
func (l *Log) Trace(format ...interface{}) error {
	now := time.Now()
	msg := fmt.Sprintf("[%02d%02d%02d.%06d][%s] ",
		now.Hour(), now.Minute(), now.Second(), now.Nanosecond(),
		lvlHead[LevelTrace])

	l.logMessage(msg, format...)
}
func (l *Log) logMessage(msg, format ...interface{}) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	var newMsg Message

	newMsg.message = fmt.SPrintf(format...)
	newMsg.message = fmt.SPrintf("%s%s", msg, newMsg.message)

	if l.sync {
		for _, log := range logger {
			_, err := log.Write(&newMsg)
			if err != nil {
				fmt.Printf("write log message failed. msg = %s\n", newMsg.message)
			}
		}
		return
	}

	l.logMsg <- &newMsg
}

func (l *Log) StartSync() {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.status = statusRunning
	l.sync = true
}
func (l *Log) StartAsync() {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.status = statusRunning
	l.sync = false

	log.logMsg = make(chan *Message, log.bufSize)
	log.sigMsg = make(chan string)
}

func (l *Log) Start() {
	l.StartAsync()
}

func (l *Log) Stop() {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if l.status == statusRunning {
		if l.sync == false {
			l.sigMsg <- "closing"

			// wait for closed
			sig <- l.sigMsg
			if sig == "closed" {
				l.status = statusClosed

				close(l.logMsg)
				close(l.sigMsg)
			}
		} else {
			l.status = statusClosed
		}
	}
}

func (l *Log) waitMsg() {
END:
	for {
		select {
		case msg := <-l.logMsg:
			for _, log := range logger {
				_, err := log.Write(msg)
				if err != nil {
					fmt.Printf("log write message failed, logger = %s, type = %ld, message = %s\n",
						log.Name(), msg.msgType, msg.message)
				}
			}
			if l.status == statusClosing {
				if len(l.logMsg) == 0 {
					break END
				}
			}
		case sig := <-l.sigMsg:
			if sig == "closing" {
				if len(l.logMsg) == 0 {
					break END
				}
				l.status = statusClosing
			}

		}
	}

	// exit logger
	for _, log := range logger {
		err := log.Close()
		if err != nil {
			fmt.Printf("log write message failed, logger = %s, type = %ld, message = %s\n",
				log.Name(), msg.msgType, msg.message)
		}
	}
	l.sigMsg <- "closed"
}

//{"buffersize": 2048, "prefix":"filename", "switchsize":1024, "fieldir":"./", "filelevel":"debug", "termlevel":"info"}
func NewLogging(conf string) (*Log, error) {

	if len(logger) == 0 {
		return nil, fmt.Errorf("no logger available")
	}

	var log Log
	err := json.Unmarshal(conf, *log.logConfig)
	if err != nil {
		return nil, err
	}

	lvl, ok := lvlStr[log.termLevelStr]
	if !ok {
		return nil, fmt.Errorf("termlevel value not valid, value = %s", log.termLevelStr)
	}
	log.termLevel = lvl

	lvl, ok = lvlStr[log.fileLevelStr]
	if !ok {
		return nil, fmt.Errorf("filelevel value not valid, value = %s", log.fileLevelStr)
	}
	log.fileLevel = lvl

	if log.bufSize <= 0 {
		log.bufSize = defBufferSize
	}

	log.mutex = sync.Mutex

	for name, l := range logger {
		fmt.Printf("open logger %s\n", name)
		//{"prefix":"filename", "switchsize":1024, "fieldir":"./"}
		s := struct {
			prefix     string
			switchsize string
			filedir    string
		}{log.prefix, log.switchSize, log.fileDir}

		conf, err := json.Marshal(s)
		if err != nil {
			return nil, err
		}

		err = l.Open(conf)
		if err != nil {
			return nil, err
		}
	}

	return &log, nil
}

func Register(log *Loger) error {
	name = log.Name()
	if _, ok := logger[name]; !ok {
		return fmt.Errorf("logger %s exists", name)
	}
	logger[name] = log

	return nil
}

func init() {
	lvlStr["all"] = LevelAll
	lvlStr["trace"] = LevelTrace
	lvlStr["debug"] = LevelDebug
	lvlStr["info"] = LevelInformational
	lvlStr["notice"] = LevelNotice
	lvlStr["warn"] = LevelWarning
	lvlStr["error"] = LevelError
	lvlStr["critical"] = LevelCritical

	lvlHead[LevelAll] = "[A]"
	lvlHead[LevelTrace] = "[T]"
	lvlHead[LevelDebug] = "[D]"
	lvlHead[LevelInformational] = "[I]"
	lvlHead[LevelNotice] = "[N]"
	lvlHead[LevelWarning] = "[W]"
	lvlHead[LevelError] = "[E]"
	lvlHead[LevelCritical] = "[C]"

}
