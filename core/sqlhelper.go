package core

import "bytes"

//生成sql 语句
func Query(tb string,names []string)(r string){
	var buf bytes.Buffer
	buf.WriteString("SELECT *  FROM ")
	buf.WriteString(tb)
	buf.WriteString(" where ")
	for _,n := range names {
		buf.WriteString(n+"=? And ")
	}
	return buf.String()
}
func Update(tb string,i interface{})(r string){
	var names= GetFieldName(i)
	var buf bytes.Buffer
	buf.WriteString("update ")
	buf.WriteString(tb)
	buf.WriteString(" set ")
	for _,name := range names {
		buf.WriteString(name+"=?, ")
	}
	buf.WriteString(" where id=?")
	return buf.String()
}