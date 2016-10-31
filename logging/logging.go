package logging

import (
	"fmt"
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
	statusReady = iota
	statusRunning
	statusClosing
	statusClosed
)

const (
	defSyncSize = 1024
)

type Loger interface {
	Name() string
	Open(conf string) error
	Write(msg *Message) (int, error)
	Close() error
	Sync() error
}

type Message struct {
	msgType int64
	message string
}

type Log struct {
	status int64
	sync   bool
	mutex  sync.Mutex
	logMsg chan *Message
	sigMsg chan string
}

var levelString = make(map[string]int64)
var levelHeadString = make(map[int64]string)
var loggerRegistered = make(map[string]Loger)
var loggerTraced = make(map[string]Loger)

func (l *Log) Critical(format string, a ...interface{}) {
	now := time.Now()
	logMsg := fmt.Sprintf("[%02d%02d%02d.%06d][%s] ",
		now.Hour(), now.Minute(), now.Second(), now.Nanosecond(),
		levelHeadString[LevelCritical])

	chanMsg := &Message{}
	chanMsg.msgType = LevelCritical

	l.logMessage(chanMsg, logMsg, format, a...)
}

func (l *Log) Error(format string, a ...interface{}) {
	now := time.Now()
	logMsg := fmt.Sprintf("[%02d%02d%02d.%06d][%s] ",
		now.Hour(), now.Minute(), now.Second(), now.Nanosecond(),
		levelHeadString[LevelError])

	chanMsg := &Message{}
	chanMsg.msgType = LevelError

	l.logMessage(chanMsg, logMsg, format, a...)
}

func (l *Log) Warning(format string, a ...interface{}) {
	now := time.Now()
	logMsg := fmt.Sprintf("[%02d%02d%02d.%06d][%s] ",
		now.Hour(), now.Minute(), now.Second(), now.Nanosecond(),
		levelHeadString[LevelWarning])

	chanMsg := &Message{}
	chanMsg.msgType = LevelWarning

	l.logMessage(chanMsg, logMsg, format, a...)
}
func (l *Log) Notice(format string, a ...interface{}) {
	now := time.Now()
	logMsg := fmt.Sprintf("[%02d%02d%02d.%06d][%s] ",
		now.Hour(), now.Minute(), now.Second(), now.Nanosecond(),
		levelHeadString[LevelNotice])

	chanMsg := &Message{}
	chanMsg.msgType = LevelNotice

	l.logMessage(chanMsg, logMsg, format, a...)
}

func (l *Log) Info(format string, a ...interface{}) {
	now := time.Now()
	logMsg := fmt.Sprintf("[%02d%02d%02d.%06d][%s] ",
		now.Hour(), now.Minute(), now.Second(), now.Nanosecond(),
		levelHeadString[LevelInformational])

	chanMsg := &Message{}
	chanMsg.msgType = LevelInformational

	l.logMessage(chanMsg, logMsg, format, a...)
}

func (l *Log) Debug(format string, a ...interface{}) {
	now := time.Now()
	logMsg := fmt.Sprintf("[%02d%02d%02d.%06d][%s] ",
		now.Hour(), now.Minute(), now.Second(), now.Nanosecond(),
		levelHeadString[LevelDebug])

	chanMsg := &Message{}
	chanMsg.msgType = LevelDebug

	l.logMessage(chanMsg, logMsg, format, a...)
}
func (l *Log) Trace(format string, a ...interface{}) {
	now := time.Now()
	logMsg := fmt.Sprintf("[%02d%02d%02d.%06d][%s] ",
		now.Hour(), now.Minute(), now.Second(), now.Nanosecond(),
		levelHeadString[LevelTrace])

	chanMsg := &Message{}
	chanMsg.msgType = LevelTrace

	l.logMessage(chanMsg, logMsg, format, a...)
}
func (l *Log) logMessage(chanMsg *Message, logMsg string, format string, a ...interface{}) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	chanMsg.message = fmt.Sprintf(format, a...)
	chanMsg.message = fmt.Sprintf("%s%s", logMsg, chanMsg.message)

	if l.sync {
		for _, log := range loggerTraced {
			_, err := log.Write(chanMsg)
			if err != nil {
				fmt.Printf("write log message failed. msg = %s\n", chanMsg.message)
			}
		}
		return
	}

	l.logMsg <- chanMsg
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

	l.logMsg = make(chan *Message, defSyncSize)
	l.sigMsg = make(chan string)
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

			sig := <-l.sigMsg
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
			for _, log := range loggerTraced {
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
	for _, log := range loggerTraced {
		err := log.Close()
		if err != nil {
			fmt.Printf("log close failed.\n")
		}
	}
	l.sigMsg <- "closed"
}

func NewLogging() (*Log, error) {
	log := &Log{}
	return log, nil
}

func AddLog(name string) Loger {
	if lg, ok := loggerRegistered[name]; ok {
		loggerTraced[name] = lg
		return lg
	}
	return nil
}

func Register(log Loger) error {
	name := log.Name()
	if _, ok := loggerRegistered[name]; ok {
		return fmt.Errorf("logger %s exists", name)
	}
	loggerRegistered[name] = log

	return nil
}

func init() {
	levelString["all"] = LevelAll
	levelString["trace"] = LevelTrace
	levelString["debug"] = LevelDebug
	levelString["info"] = LevelInformational
	levelString["notice"] = LevelNotice
	levelString["warn"] = LevelWarning
	levelString["error"] = LevelError
	levelString["critical"] = LevelCritical

	levelHeadString[LevelAll] = "[A]"
	levelHeadString[LevelTrace] = "[T]"
	levelHeadString[LevelDebug] = "[D]"
	levelHeadString[LevelInformational] = "[I]"
	levelHeadString[LevelNotice] = "[N]"
	levelHeadString[LevelWarning] = "[W]"
	levelHeadString[LevelError] = "[E]"
	levelHeadString[LevelCritical] = "[C]"

}
