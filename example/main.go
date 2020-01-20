package main

import (
	"caption_json_formatter"
	"github.com/sirupsen/logrus"
	"log"
	"os"
)
var (
	logger = NewLogger()
	logger2 = NewLogger2()
)

type (
	JO map[string]interface{}
)

func NewLogger() *logrus.Logger {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetOutput(os.Stdout)

	logger := logrus.New()
	logger.Level = logrus.TraceLevel

	logger.SetFormatter(&caption_json_formatter.Formatter{ PrettyPrint: true, CustomCaption: "nollehLog" })
	return logger
}

func NewLogger2() *logrus.Logger {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetOutput(os.Stdout)

	logger := logrus.New()
	logger.Level = logrus.TraceLevel

	logger.SetFormatter(&caption_json_formatter.Formatter{ CustomCaption: JO{"name":"nolleh", "say":"hello"} })
	return logger
}

func Log() *caption_json_formatter.Entry {
	return &caption_json_formatter.Entry{ Entry: logrus.NewEntry(logger) }
}

func Log2() *caption_json_formatter.Entry {
	return &caption_json_formatter.Entry{ Entry: logrus.NewEntry(logger2) }
}

func main() {
	type Request struct {
		Url string `json:"url"`
		Method string `json:"method"`
	}
	type Response struct {
		User int `json:"user"`
		Balance int `json:"balance"`
	}
	type Message struct {
		Request Request `json:"request"`
		Response Response `json:"response"`
	}

	message := Message { Request{ "/user/123456/balance", "GET" },
		Response{123456, 1000} }

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
	Log2().Debug(map[string]interface{}{ "key": "value"})

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
}