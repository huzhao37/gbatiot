package main

import (
	xml "yuniot/framework/mysql"
)


func init()  {
}
//共享变量有一个读通道和一个写通道组成
//var total int
//var mutex sync.RWMutex
//func Write(){
//	mutex.Lock()
//	total=total+1
//	mutex.Unlock()
//}
func main() {

	//var start=time.Now()
	//for i:=1;i<1001 ;i++  {
	//	go Write()
	//}
	//time.Sleep(3 * time.Second)
	//
	//fmt.Printf("total：%d",total)
	//var cost=time.Since(start)
	//fmt.Printf("cost：%d ms",cost/1e6)

	defer xml.SqlDB.Close()
	router := initRouter()
	router.Run(":5200")
	//outer.RunTLS(":8000", "./testdata/server.pem", "./testdata/server.key")
}

