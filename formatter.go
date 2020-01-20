package caption_json_formatter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

type RootFields struct {
	Timestamp string
	Func   string
	Level  logrus.Level
	Fields interface{}
}

type Formatter struct {
	CustomCaptionPrettyPrint bool
	// custom caption can be struct, string, whatever
	CustomCaption interface{}
	PrettyPrint bool
}

func (f *Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	b := &bytes.Buffer{}

	root := RootFields{ Timestamp: entry.Time.Format(time.RFC3339Nano), Level: entry.Level,
		//CustomCaption: entry.CustomCaption, // not possible in logrus...
		Fields: encode(entry.Message) }

	b.WriteString(root.Timestamp)

	b.WriteString(" [")
	_, _ = fmt.Fprintf(b, "\x1b[%dm", colorDarkBlue)
	b.WriteString(strings.ToUpper(root.Level.String()))
	b.WriteString("\x1b[0m")
	b.WriteString("] ")

	//if entry.HasCaller() {
	//	caller := getCaller(entry.Caller)
	//	fc := caller.Function
	//	file := fmt.Sprintf("%s:%d", caller.File, caller.Line)
	//	b.WriteString(prettierCaller(file, fc))
	//}

	b.WriteString("[")
	captionStr := ""
	if f.CustomCaptionPrettyPrint {
		captionStr = marshalIndent(f.CustomCaption)
	} else {
		captionStr = marshal(f.CustomCaption)
	}
	b.WriteString(captionStr)
	b.WriteString("] ")

	levelColor := getColorByLevel(entry.Level)
	_, _ = fmt.Fprintf(b, "\x1b[%dm", levelColor)

	var data string
	if f.PrettyPrint {
		data = marshalIndent(root.Fields)
		//data = root.marshalIndent()
	} else {
		data = marshal(root.Fields)
		//data = root.marshal()
	}

	b.WriteString(data)
	b.WriteString("\x1b[0m")

	b.WriteByte('\n')
	return b.Bytes(), nil
}

func prettierCaller(file string, function string) string {
	dirs := strings.Split(file, "/")
	fileDesc := strings.Join(dirs[len(dirs)-2:], "/")

	funcs := strings.Split(function, ".")
	funcDesc := strings.Join(funcs[len(funcs)-2:], ".")

	return "[" + fileDesc + ":" + funcDesc + "]"
}

func encode(message string) interface{} {
	if data := encodeForJsonString(message); data != nil {
		return data
	} else {
		return message
	}
}

func encodeForJsonString(message string) map[string]interface{} {
	// jsonstring
	inInterface := make(map[string]interface{})
	if err := json.Unmarshal([]byte(message), &inInterface); err != nil {
		//fmt.Print("err !!!! " , err.Error())
		return nil
	}
	return inInterface
}
const (
	colorRed      = 31
	colorGreen    = 32
	colorYellow   = 33
	colorDarkBlue = 34
	colorBlue     = 36
	colorGray     = 37
)

func getColorByLevel(level logrus.Level) int {
	switch level {
	case logrus.TraceLevel:
		return colorGray
	case logrus.DebugLevel:
		return colorBlue
	case logrus.InfoLevel:
		return colorGreen
	case logrus.WarnLevel:
		return colorYellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		return colorRed
	default:
		return colorDarkBlue
	}
}
