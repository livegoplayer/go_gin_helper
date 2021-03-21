package helper

import "strings"

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
