package prov

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type goProtoData struct {
	cmommand    int64
	comments    string
	lastComment string
	msg         string
}

type goProto struct {
	pack    string
	data    []goProtoData
	curData goProtoData
}

type goRrovider struct {
	proto goProto
}

func (p *goRrovider) GenCmdFile(protoFile string) (string, error) {
	return "Hello", nil
}

func (p *goRrovider) parseProto(protoFile string) error {
	file, err := os.Open(protoFile)
	if err != nil {
		return err
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	lines := strings.Split(string(b), "\n")
	for _, line := range lines {
		if err = p.parseLine(line); err != nil {
			return err
		}
	}
	return nil
}
func (p *goRrovider) parseLine(line string) error {
	line = strings.Trim(line, "\r\n\t ;")
	switch {
	case line == "package":
		packString := strings.Split(line, " ")
		if len(packString) != 2 {
			return fmt.Errorf("package statement missing package name, line = %s", line)
		}
		p.proto.pack = strings.Trim(packString[1], " ")
	case strings.HasPrefix(line, "//"):
		p.proto.curData.comments = fmt.Sprintf("%s\n%s", p.proto.curData.comments, line)
		p.proto.curData.lastComment = line
	case strings.LastIndex(line, "/*") != -1 || strings.LastIndex(line, "*/") != -1:
		return fmt.Errorf("not support /**/ style comments, please use // comment style, line = %s", line)
	case strings.HasPrefix(line, "message"):
		lastComment := p.proto.curData.lastComment
		if len(lastComment) <= 0 {
			return fmt.Errorf("missing message command comment info, line = %s", line)
		}
		comment := strings.Split(line, " ")
		if len(comment) != 2 {
			return fmt.Errorf("format style not right, line = %s", line)
		}
		command, err := strconv.ParseInt(strings.Trim(comment[1], " "), 16, 64)
		if err != nil {
			return fmt.Errorf("format style not right, line = %s, strconv error = %s", line, err)
		}
		if command == 0x0 {
			p.proto.curData = goProtoData{}
		}

	}

}

func init() {
	Register("go", &goRrovider{})
}
