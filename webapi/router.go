package main

import (
	"gopkg.in/gin-gonic/gin.v1"
	. "yuniot/apis"
	"yuniot/middleware"
	"net/http"
)

func initRouter() *gin.Engine {

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()//gin.Default()
	//router.Use(gin.Logger())
	router.Use(gin.Recovery())
	//cors
	router.Use(Cors())
	router.POST("/auth", Auth)

	//router.OPTIONS("/login", LoginApi)

	//jwt token
	taR := router.Group("/cqdy")
	taR.Use(middleware.JWTAuth())
	{
		taR.GET("/dataByTime",GetDataByTime)
		taR.GET("/login", GetUsersApi)
		taR.GET("/motors", GetMotorsApi)
		//taR.OPTIONS("/motors", GetMotorsApi)
		taR.GET("/currentmaincymonth", GetMainCyCurrentMonthApi)
	//	taR.OPTIONS("/currentmaincymonth", GetMainCyCurrentMonthApi)
		taR.GET("/currentmaincyday", GetMainCyCurrentDayApi)
	//	taR.OPTIONS("/currentmaincyday", GetMainCyCurrentDayApi)
		taR.GET("/currentmaincy", GetMainCyCurrentApi)
	//	taR.OPTIONS("/currentmaincy", GetMainCyCurrentApi)
		taR.GET("/currentbeltcys", GetBeltCysCurrentDayApi)
		//taR.OPTIONS("/currentbeltcys", GetBeltCysCurrentDayApi)
		taR.GET("/currentdevice", GetDeviceCurrentApi)
		//taR.OPTIONS("/currentdevice", GetDeviceCurrentApi)
		taR.GET("/status", GetStatusApi)
		//taR.OPTIONS("/status", GetStatusApi)
		//taR.POST("/login", LoginApi)
	}
	return router
}
// 处理跨域请求,支持options访问
func Cors() gin.HandlerFunc {

	return func(c *gin.Context) {

		method := c.Request.Method

		c.Header("Access-Control-Allow-Origin", "*")

		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")

		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")

		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")

		c.Header("Access-Control-Allow-Credentials", "true")

		//放行所有OPTIONS方法

		if method == "OPTIONS" {

			c.AbortWithStatus(http.StatusOK)

		}

		// 处理请求

		c.Next()
		}
}
