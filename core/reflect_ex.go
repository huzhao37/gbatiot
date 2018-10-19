package core

import "reflect"

//获取结构体中字段的名称
func GetFieldName(structName interface{}) []string  {
	t := reflect.TypeOf(structName)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		Logger.Println("Check type error not Struct")
		return nil
	}
	fieldNum := t.NumField()
	result := make([]string,0,fieldNum)
	for i:= 0;i<fieldNum;i++ {
		result = append(result,t.Field(i).Name)
	}
	return result
}
