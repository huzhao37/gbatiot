package core

import (
	"strings"
	"strconv"
)

func StrFirstToUpper(str string) string {
	var values []string=make([]string,0)
	temp := []byte(str)
	values=append(values,strings.ToUpper(string(temp[0])))
	for i := 1; i < len(temp); i++ {
		values=append(values,strings.ToLower(string(temp[i])))
		}
	return strings.Join(values,"")
}
//float32 转 String工具类，保留6位小数
func FloatToString(input_num float32) string {
	// to convert a float number to a string
	return strconv.FormatFloat(float64(input_num), 'f', 6, 64)
}
