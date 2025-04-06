package main

import (
	"log"
	"os"

	"github.com/nolleh/caption_json_formatter"
	"github.com/sirupsen/logrus"
)

var (
	logger  = NewLogger()
	logger2 = NewLogger2()
	logger3 = NewLogger3()
	logger4 = NewJsonLogger()
)

type (
	JO map[string]any
)

func NewLogger() *logrus.Logger {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetOutput(os.Stdout)

	logger := logrus.New()
	logger.Level = logrus.TraceLevel

	consoleLogger := caption_json_formatter.Console()
	logger.SetFormatter(consoleLogger)
	return logger
}

func NewLogger2() *logrus.Logger {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetOutput(os.Stdout)

	logger := logrus.New()
	logger.Level = logrus.TraceLevel

	consoleLogger := caption_json_formatter.Console()
	// if you want to change default value, modify it.
	consoleLogger.CustomCaption = JO{"name": "nolleh", "say": "hello"}
	logger.SetFormatter(consoleLogger)
	return logger
}

func NewLogger3() *logrus.Logger {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetOutput(os.Stdout)

	logger := logrus.New()
	logger.Level = logrus.TraceLevel

	// if you no need syntatic sugar to generate formatter, you can do by yourself.
	logger.SetFormatter(&caption_json_formatter.Formatter{PrettyPrint: true, Colorize: true})
	return logger
}

func NewJsonLogger() *logrus.Logger {

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetOutput(os.Stdout)

	logger := logrus.New()
	logger.Level = logrus.TraceLevel

	json := caption_json_formatter.Json()
	// you can use this
	// json.Colorize = true
	// json.PrettyPrint = true
	logger.SetFormatter(json)
	return logger
}

func Log() *caption_json_formatter.Entry {
	return &caption_json_formatter.Entry{Entry: logrus.NewEntry(logger)}
}

func Log2() *caption_json_formatter.Entry {
	return &caption_json_formatter.Entry{Entry: logrus.NewEntry(logger2)}
}

func Log3() *caption_json_formatter.Entry {
	return &caption_json_formatter.Entry{Entry: logrus.NewEntry(logger3)}
}

func JsonLog() *caption_json_formatter.Entry {
	return &caption_json_formatter.Entry{Entry: logrus.NewEntry(logger4)}
}
func main() {
	type Request struct {
		Url    string `json:"url"`
		Method string `json:"method"`
	}
	type Response struct {
		User    int `json:"user"`
		Balance int `json:"balance"`
	}
	type Message struct {
		Request  Request  `json:"request"`
		Response Response `json:"response"`
	}

	message := Message{Request{"/user/123456/balance", "GET"},
		Response{123456, 1000}}

	// without entry, you can't get json formatted message...
	// output: 2020-01-20T18:08:12.6077798+09:00 [DEBUG] [nollehLog] {{/user/123456/balance GET} {123456 1000}}
	logger.Debug(message)

	/* in current logrus implementation, there isn't way for set Custom Entry.
	 * hook or formatter, doesn't have opportunity for marshaling message.
	 * so before pull request was allowed, you need to use custom function to use extended entry.
	 */
	// 2020-01-20T18:08:12.811776+09:00 [TRACE] [nollehLog] trace
	Log().Trace("trace")

	//  2020-01-20T18:08:12.8127822+09:00 [DEBUG] [{"name":"nolleh","say":"hello"}] {"key":"value"}
	Log2().Debug(map[string]string{"key": "value"})

	// 2020-01-20T18:08:12.8127822+09:00 [INFO] [nollehLog] {
	// "request": {
	//  "method": "GET",
	//  "url": "/user/123456/balance"
	// },
	// "response": {
	//  "balance": 1000,
	//  "user": 123456
	// }
	//}
	Log().Info(message)

	// 2020-01-20T18:08:12.8127822+09:00 [WARNING] [{"name":"nolleh","say":"hello"}] {"request":{"method":"GET","url":"/user/123456/balance"},"response":{"balance":1000,"user":123456}}
	Log2().Warn(message)

	// 2020-01-20T18:08:12.8127822+09:00 [ERROR] [nollehLog] {
	// "request": {
	//  "method": "GET",
	//  "url": "/user/123456/balance"
	// },
	// "response": {
	//  "balance": 1000,
	//  "user": 123456
	// }
	//}
	Log().Error(message)

	// 2020-01-20T18:08:12.8127822+09:00 [FATAL] [{"name":"nolleh","say":"hello"}] {"request":{"method":"GET","url":"/user/123456/balance"},"response":{"balance":1000,"user":123456}}
	Log2().Fatal(message)

	// 2020-01-20T18:08:12.8127822+09:00 [INFO] {
	// "request": {
	// 		"method": "GET",
	//  	"url": "/user/123456/balance"
	// },
	// "response": {
	// 		"balance": 1000,
	// 		"user": 123456
	// }
	//}
	Log3().Info(message)

	//
	JsonLog().Info(message)
}
