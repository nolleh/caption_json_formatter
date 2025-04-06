package caption_json_formatter

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

type Entry struct {
	*logrus.Entry
}

func (entry *Entry) Trace(args ...any) {
	entry.printLog(logrus.TraceLevel, args...)
}

func (entry *Entry) Debug(args ...any) {
	entry.printLog(logrus.DebugLevel, args...)
}

func (entry *Entry) Info(args ...any) {
	entry.printLog(logrus.InfoLevel, args...)
}

func (entry *Entry) Warn(args ...any) {
	entry.printLog(logrus.WarnLevel, args...)
}

func (entry *Entry) Error(args ...any) {
	entry.printLog(logrus.ErrorLevel, args...)
}

func (entry *Entry) Fatal(args ...any) {
	entry.printLog(logrus.FatalLevel, args...)
}

func (entry *Entry) Panic(args ...any) {
	entry.printLog(logrus.PanicLevel, args...)
}

func (entry *Entry) printLog(level logrus.Level, args ...any) {
	datas := make([]any, len(args))
	for i, v := range args {
		datas[i] = Stringify(v)
	}
	entry.Log(level, fmt.Sprint(datas...))
}
