package prov

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"text/template"
)

type goProtoData struct {
	command     int64
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
	proto    goProto
	template string
}

func (p *goRrovider) GenCmdFile(protoFile string) (string, error) {
	err := p.parseProto(protoFile)
	if err != nil {
		return "", err
	}

	report, err := template.New("genfile").Parse(p.template)
	if err != nil {
		return "", fmt.Errorf("parse template failed, err = %s", err)
	}
	var output bytes.Buffer
	err = report.Execute(&output, p.proto)
	if err != nil {
		return "", fmt.Errorf("rend template failed, err = %s", err)
	}
	return output.String(), nil
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
	case strings.HasPrefix(line, "package"):
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
		idx := strings.LastIndex(lastComment, "0x")
		if idx < 0 {
			return fmt.Errorf("comment style not right, line = %s, comment = %s", line, lastComment)
		}
		strCommand := strings.Trim(lastComment[idx:len(lastComment)], " ")
		command, err := strconv.ParseInt(strCommand, 0, 64)
		if err != nil {
			return fmt.Errorf("comment style not right, line = %s, str = %s, strconv error = %s", line, strCommand, err)
		}
		if command == 0x0 {
			p.proto.curData = goProtoData{}
			return nil
		}
		p.proto.curData.command = command
		msgs := strings.Split(line, " ")
		if len(msgs) < 2 {
			return fmt.Errorf("message style not right, line = %s", line)
		}
		idx = strings.LastIndex(msgs[1], "{")
		if idx < 0 {
			idx = len(msgs[1])
		}
		msg := strings.Trim(msgs[1][0:idx], " ")
		p.proto.curData.msg = msg

		p.proto.data = append(p.proto.data, p.proto.curData)

		p.proto.curData = goProtoData{}
	}
	return nil

}

func init() {
	var template = `
package {{.pack}}

import (
	"fmt"

	pb "github.com/golang/protobuf/proto"	
)

type CommandCode int64

const (
	{{rang .data}}
	{{.comments}}
	CMD_{{strings.Upper(.msg)}} CommandCode = {{.command}} {{.lastComment}}
	{{end}}
)

var message map[int64]pb.Message

func Message(command int64) (pb.Message, error) {
	if m, ok := message[command]; ok {
		pb.Clone(m)
		return m, nil
	}
	return nil, fmt.Errorf("command %d not found", command)
}

func Parse(command int64, buf []byte) (pb.Message, error) {
	msg, err := Message(command)
	if err != nil {
		return nil, err
	}
	err = pb.Unmarshal(buf, msg)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func Serialize(msg pb.Message) ([]byte, error) {
	return pb.Marshal(msg)
}

func Debug(command int64, msg pb.Message) string {
	return fmt.Sprintf("command:0x%x %s", command, pb.CompactTextString(msg))
}

func init() {
	message = make(map[int64]pb.Message)

	{{rang .data}}
	message[CMD_{{strings.Upper(.msg)}}] = &{{.msg}}{} {{.lastComment}}
	{{end}}
}
`
	provider := &goRrovider{}
	provider.template = template
	Register("go", provider)
}
