package main

import (
	db "yuniot/framework/mysql"
	"yuniot/core"
	"fmt"
)


func init()  {
}
func main() {

	core.Try(func() {
		defer db.Xml.DB().Close()
		defer  db.Auth.DB().Close()
		router := initRouter()
		router.Run(":5200")
	}, func(i interface{}) {
		fmt.Printf("%s",i)
	})

	//outer.RunTLS(":8000", "./testdata/server.pem", "./testdata/server.key")
}

