package main

import (
	"gopkg.in/gin-gonic/gin.v1"
	. "yuniot/apis"
	"yuniot/middleware"
)

func initRouter() *gin.Engine {

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	//cors
	router.Use(MiddleWare())
	router.POST("/auth", Auth)

	//router.OPTIONS("/login", LoginApi)

	//jwt token
	taR := router.Group("/cqdy")
	taR.Use(middleware.JWTAuth())
	{
		taR.GET("/dataByTime",GetDataByTime)
		taR.GET("/login", GetUsersApi)
		taR.GET("/motors", GetMotorsApi)
		taR.OPTIONS("/motors", GetMotorsApi)
		taR.GET("/currentmaincymonth", GetMainCyCurrentMonthApi)
		taR.OPTIONS("/currentmaincymonth", GetMainCyCurrentMonthApi)
		taR.GET("/currentmaincyday", GetMainCyCurrentDayApi)
		taR.OPTIONS("/currentmaincyday", GetMainCyCurrentDayApi)
		taR.GET("/currentmaincy", GetMainCyCurrentApi)
		taR.OPTIONS("/currentmaincy", GetMainCyCurrentApi)
		taR.GET("/currentbeltcys", GetBeltCysCurrentDayApi)
		taR.OPTIONS("/currentbeltcys", GetBeltCysCurrentDayApi)
		taR.GET("/currentdevice", GetDeviceCurrentApi)
		taR.OPTIONS("/currentdevice", GetDeviceCurrentApi)
		taR.GET("/status", GetStatusApi)
		taR.OPTIONS("/status", GetStatusApi)
		//taR.POST("/login", LoginApi)
	}
	return router
}
//cors
func MiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
			c.Request.SetBasicAuth("x","x")
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")//允许访问所有域
			c.Writer.Header().Add("Access-Control-Allow-Headers","Content-Type")//header的类型
			//c.Writer.Header().Set("content-type","application/json") //返回数据格式是json
			c.Next()
		//}

	}
}
