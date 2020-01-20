package caption_json_formatter

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"reflect"
)
func Stringify(v interface{}) string {
	var ret string
	if IsObject(v) || isMap(v) {
		ret = marshal(v)
	} else {
		ret = fmt.Sprint(v)
	}
	return ret
}

func IsObject(v interface{}) bool {
	return reflect.ValueOf(v).Kind() == reflect.Struct
}

func marshal(o interface{}) string {
	str, ok := o.(string)
	if ok {
		return str
	}
	data, err := json.Marshal(o)
	if err != nil {
		return fmt.Sprint(o)
	}
	return string(data)
}


func marshalIndent(o interface{}) string {
	str, ok := o.(string)
	if ok {
		return str
	}
	m, err := json.MarshalIndent(o, "", " ")
	if err != nil {
		return fmt.Sprint(o)
	}
	return string(m)
}


func isMap(v interface{}) bool {
	_, isMap := v.(map[string]interface{})
	_, isLogFields := v.(logrus.Fields)

	return isMap || isLogFields
}
