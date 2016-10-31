package logging

import (
	"encoding/json"
	"github.com/fatih/color"
)

var colorLevel = make(map[int64]color.Attribute)

type consoleLogger struct {
	level int64
}

func (c *consoleLogger) Name() string {
	return "console"
}
func (c *consoleLogger) Open(conf string) error {
	err := json.Unmarshal([]byte(conf), *c)
	return err
}
func (c *consoleLogger) Write(msg *Message) (int, error) {
	n, err := 0, error(nil)
	if msg.msgType >= c.level {
		n, err = color.New(colorLevel[msg.msgType]).Print(msg.message)
	}
	return n, err
}
func (c *consoleLogger) Close() error {
	return nil
}
func (c *consoleLogger) Sync() error {
	return nil
}

func init() {

	colorLevel[LevelCritical] = color.FgHiRed
	colorLevel[LevelError] = color.FgRed
	colorLevel[LevelWarning] = color.FgHiCyan
	colorLevel[LevelNotice] = color.FgHiBlue
	colorLevel[LevelInformational] = color.FgBlue
	colorLevel[LevelDebug] = color.FgHiGreen
	colorLevel[LevelTrace] = color.FgGreen
	colorLevel[LevelAll] = color.FgWhite

	c := &consoleLogger{}
	Register(c)
}
