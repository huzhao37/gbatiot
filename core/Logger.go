package core

import (
	"os"
	"time"
	"fmt"
	"log"
)

var (
	Logger *log.Logger
)
func init() {
	//创建输出日志文件
	logFile, errLog := os.Create("./" + time.Now().Format("20060102") + ".txt")
	if errLog != nil {
		fmt.Println(errLog)
	}

	//创建一个Logger
	//参数1：日志写入目的地
	//参数2：每条日志的前缀
	//参数3：日志属性
	Logger = log.New(logFile, "test_", log.Ldate|log.Ltime|log.Lshortfile)
}

