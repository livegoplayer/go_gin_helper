package helper

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"strconv"
	"strings"
)

//获取source的子串,如果start小于0或者end大于source长度则返回""
//start:开始index，从0开始，包括0
//end:结束index，以end结束，但不包括end
func Substring(source string, start int, end int) string {
	var r = []rune(source)
	length := len(r)

	if start < 0 || end > length || start > end {
		return ""
	}

	if start == 0 && end == length {
		return source
	}

	return string(r[start:end])
}

//获取两个字符串中间的字符串
func GetSubStringBetween(source string, startString string, endString string) string {
	//先拿到第一个字符串到最后的子字符串
	start := strings.Index(source, startString)
	startIndex := 0
	if start != -1 {
		startIndex = start + strings.Count(startString, "") - 1
	}
	source = Substring(source, startIndex, strings.Count(source, "")-1)
	if endString == "" {
		return source
	}
	return Substring(source, 0, strings.Index(source, endString))
}

func AsString(v interface{}) string {
	switch v.(type) {
	case uint32:
		return strconv.FormatInt(int64(v.(uint32)), 10)
	case uint64:
		return strconv.FormatInt(int64(v.(uint64)), 10)
	case int:
		return strconv.Itoa(v.(int))
	case int32:
		return strconv.Itoa(int(v.(int32)))
	case int64:
		return strconv.FormatInt(v.(int64), 10)
	case float64:
		return strconv.FormatFloat(v.(float64), 'E', -1, 64)
	case float32:
		return strconv.FormatFloat(float64(v.(float32)), 'E', -1, 64)
	case string:
		return v.(string)
	case bool:
		if v.(bool) {
			return "true"
		} else {
			return "false"
		}
	case map[string]interface{}:
		return JsonEncode(v)
	default:
		return ""
	}
}

func JsonEncode(data interface{}) string {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	jsonByte, err := json.Marshal(&data)
	if err != nil {
		fmt.Printf("json加密出错:" + err.Error())
	}
	return string(jsonByte[:])
}
