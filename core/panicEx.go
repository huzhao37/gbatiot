package core

//实现 try catch 例子
func Try(fun func(), handler func(interface{})) {
	defer func() {
		if err := recover(); err != nil {
			handler(err)
		}
	}()
	fun()
}

//demo
func main() {
	Try(func() {
		panic("foo") //正常执行的函数
	}, func(e interface{}) {
		print(e) //异常处理委托
	})
}
