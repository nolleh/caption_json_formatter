[![Go Reference](https://pkg.go.dev/badge/github.com/nolleh/caption_json_formatter.svg)](https://pkg.go.dev/github.com/nolleh/caption_json_formatter)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## caption-json-formatter

### Feature

logrus's 'message(object)' json formatter with human-readable caption added.  
default available caption is (it can be disabled)

> - timestamp
> - logLevel

and also you can use this formatter for using simple json formatter (@since v0.2.0).

### Example

#### OutPut

1. with custom caption (nollehLog)  
   <img src="docs/images/example_01.png" alt="with caption"/>

2. set logger option with unpretty print (off pritty print)  
   <img src="docs/images/example_02.png" alt="unpretty"/>

3. set logger without custom caption and set as pretty  
   <img src="docs/images/example_03.png" alt="no caption"/>

4. use caption_json_formatter as simple json logger, to handled by other program. (like kibana, loki, whatever)
   <img src="docs/images/example_04.png" alt="json with default caption">
   it printed as
   ```json
   {
     "Level": "info",
     "Message": {
       "request": { "method": "GET", "url": "/user/123456/balance" },
       "response": { "balance": 1000, "user": 123456 }
     },
     "TimeStamp": "2023-03-05T18:16:47.165173+09:00"
   }
   ```

#### Available Options

| key                      | ConsoleDefault | JsonDefault | description                                                                   |
| ------------------------ | -------------- | ----------- | ----------------------------------------------------------------------------- |
| TransportToJson          | false          | true        | if enabled, the whole log (include default caption) will be marshaled to json |
| UseDefaultCaption        | true           | true        | if enabled, timestamp and loglevel will be added                              |
| PrettyPrint              | true           | false       | if enabled, the message(object) is printed as pretty json                     |
| CustomCaption            | nil            | nil         | if has value, then the caption is attatched before message(object)            |
| CustomCaptionPrettyPrint | false          | false       | if enabled, the passed custom caption is printed as pretty json               |
| Colorize                 | true           | false       | if enabled, add predefined term color code by loglevel                        |

#### Source

- initalization

```go
// initialize log
var (
  logger = NewLogger()
)

func NewLogger() *logrus.Logger {
  log.SetFlags(log.LstdFlags | log.Lshortfile)
  log.SetOutput(os.Stdout)
  logger := logrus.New()
  logger.Level = logrus.TraceLevel
  // option 1. helper function Console makes default formatter for console
  logger.SetFormatter(caption_json_formatter.Console())

  // if you want, use default constructor and modify some option.
  // console := caption_json_formatter.Console()
  // console.CustomCaptionPrettyPrint = true
  // logger.SetFormatter(console)

  // use JsonFormatter
  // json := caption_json_formatter.Json()
  // logger.SetFormatter(json)

  // option 2. if you no need helper constructing, you can do yourself. (not recommended)
  // logger.SetFormatter(&caption_json_formatter.Formatter{ PrettyPrint: true })
  return logger
}

func Log() *caption_json_formatter.Entry {
	return &caption_json_formatter.Entry{ Entry: logrus.NewEntry(logger) }
}
```

- printLog

```go
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

  message := Message { Request{ "/user/123456/balance", "GET" }, Response{123456, 1000} }

  /* in current logrus implementation, there isn't way for set Custom Entry
  * hook or formatter, doesn't have opportunity for marshaling message.
  * so before pull request was made, you need to use custom function to use extended entry.
  */

  Log().Debug(message)
}
```

### When To Use

it is useful when shows 'message' that include json to human.
show as pretty format, eventhought there is nested json. **(not strigified)**

for example, for message

```go
// accessData is json type that has request / reponse data
log.Debug(accessData)
```

```
[2020-01-20T16:46:08.7452971+09:00] [Debug]
{
    "request": {
        "headers": { "content-type": "application/json" },
        "method": "GET",
        "route": "/user/{userId}/balance",
        "url": "/user/123456/balance"
    },
    "response": {
        "body": {
            "userId": 123456,
            "balance": 1000
        }
    }
}
```

not like other formatters doing. (like below)

```
[2020-01-20T16:46:08.7452971+09:00] [Debug] {/"request/": { /"headers/": { /"content-type/": /"application/json/" }, /"method/": /"GET/", /"route/": /"/user/{userId}/balance/", /"url/": /"/user/123456/balance/" }, /"response/": { /"body/": { /"userId/": 123456, /"balance/": 1000 } } }
```

or

```
{ "time": "2020-01-20T16:46:08.7452971+09:00", "msg": {/"request/": { /"headers/": { /"content-type/": /"application/json/" }, /"method/": /"GET/", /"route/": /"/user/{userId}/balance/", /"url/": /"/user/123456/balance/" }, /"response/": { /"body/": { /"userId/": 123456, /"balance/": 1000 } }
```

or

```
"2020-01-20T16:46:08.7452971+09:00" {{/user/123456/balance GET} {123456 1000}}"
```
