package logging

import ()

type Log struct {
	mutex     *sync.Mutex
	logs      map[int]*Appender
	chanMsg   chan LogMessage
	signalMsg chan string
}

func (l *Log) Critical(format ...interface{}) error {
	a.writeMessage(LevelCritical, message)
}

func (l *Log) Error(format ...interface{}) error {

}

func (l *Log) Warning(format ...interface{}) error {

}
func (l *Log) Notice(format ...interface{}) error {

}

func (l *Log) Info(format ...interface{}) error {

}

func (l *Log) Debug(format ...interface{}) error {

}
func (l *Log) Trace(format ...interface{}) error {

}
