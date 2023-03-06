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
	Func      string
	Level     logrus.Level
	Fields    interface{}
}

type Formatter struct {
	/** if enabled, the whole log will be marshaled to json (including captions)
	  work as json Formatter. */
	TransportToJson bool
	/** if enabled, log prefixed with 'timestamp, level' */
	UseDefaultCaption bool
	/** if enabled, CustomCaption will be marshaled to json */
	CustomCaptionPrettyPrint bool
	/** if has value, it attached right before message(object). custom caption can be struct, string, whatever */
	CustomCaption interface{}
	/** do PrettyPrint for message(object) */
	PrettyPrint bool
	/** if enabled, the message(object) will be colorized by predefined color code, along with logLevel */
	Colorize bool
}

/** syntatic sugar for using formatter with predefined option (for console) */
func Console() *Formatter {
	return &Formatter{UseDefaultCaption: true, PrettyPrint: true, Colorize: true}
}

/** syntatic sugar for using formatter with predefined option (for server, json) */
func Json() *Formatter {
	return &Formatter{TransportToJson: true, UseDefaultCaption: true}
}

type (
	JO map[string]interface{}
)

func (f *Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	b := &bytes.Buffer{}

	root := RootFields{Fields: encode(entry.Message)}

	if f.UseDefaultCaption {
		if f.TransportToJson {
			root = RootFields{
				Fields: JO{
					"time_stamp": entry.Time.Format(time.RFC3339Nano),
					"level":      entry.Level,
					"message":    encode(entry.Message),
				},
			}
		} else {
			root = RootFields{Timestamp: entry.Time.Format(time.RFC3339Nano), Level: entry.Level,
				//CustomCaption: entry.CustomCaption, // not possible in logrus...
				Fields: encode(entry.Message)}
			b.WriteString(root.Timestamp)
			b.WriteString(" [")

			if f.Colorize {
				_, _ = fmt.Fprintf(b, "\x1b[%dm", colorDarkBlue)
			}

			b.WriteString(strings.ToUpper(root.Level.String()))

			if f.Colorize {
				b.WriteString("\x1b[0m")
			}

			b.WriteString("] ")
		}
	}

	//if entry.HasCaller() {
	//	caller := getCaller(entry.Caller)
	//	fc := caller.Function
	//	file := fmt.Sprintf("%s:%d", caller.File, caller.Line)
	//	b.WriteString(prettierCaller(file, fc))
	//}

	if f.CustomCaption != nil {

		b.WriteString("[")
		captionStr := ""
		if f.CustomCaptionPrettyPrint {
			captionStr = marshalIndent(f.CustomCaption)
		} else {
			captionStr = marshal(f.CustomCaption)
		}
		b.WriteString(captionStr)
		b.WriteString("] ")
	}

	if f.Colorize {
		levelColor := getColorByLevel(entry.Level)
		_, _ = fmt.Fprintf(b, "\x1b[%dm", levelColor)
	}

	var data string
	if f.PrettyPrint {
		data = marshalIndent(root.Fields)
		// data = marshalIndent(root)
	} else {
		data = marshal(root.Fields)
		// data = marshal(root)
	}

	b.WriteString(data)
	if f.Colorize {
		b.WriteString("\x1b[0m")
	}

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
