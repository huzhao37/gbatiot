package apis

import (
	"net/http"
	"gopkg.in/gin-gonic/gin.v1"
	"yuniot/core"
	"time"
	"yuniot/repository"
	"yuniot/framework/influx"
	"yuniot/models/xml"
)

type Status struct {
	productionlinestatus bool
	devicesStatus []map[string]bool
}

//获取当月主皮带累计情况概览
func GetMainCyCurrentMonthApi(c *gin.Context){
	isPass := c.GetBool("isPass")
	if !isPass {
		return
	}
	motor,err:=xml.GetMainCyByProductionlineId(c.Request.FormValue("productionlineid"))
	if err!=nil {
		c.JSON(http.StatusOK,gin.H{
			"status":0,
			"msg":"设备不存在"+err.Error(),
		})
	}
	influxd,err:=xml.GetInfluxByProductionlineId(motor.ProductionLineId)
	if err!=nil{
		core.Logger.Println("获取产线 %s 的influxdb ，err：%s",motor.ProductionLineId,err)
		c.JSON(http.StatusOK,gin.H{
			"status":0,
			"msg":"获取产线没有数据库"+err.Error(),
		})
	}
	conn:=influx.ConnInfluxParam(influxd.Addr,influxd.User,influxd.Pwd)
	defer conn.Close()
	timeLayout := "2006-01-02 15:04:05"
	loc, _ := time.LoadLocation("Local")   //重要：获取时区
	start:=time.Date(time.Now().Year(),time.Now().Month(),1,0,0,0,0,loc)
	end:=time.Date(time.Now().Year(),time.Now().Month(),time.Now().Day(),time.Now().Hour(),time.Now().Minute(),0,0,loc)
	startStr := start.Format(timeLayout) //设置时间戳 使用模板格式化为日期字符串
	endStr := end.Format(timeLayout) //设置时间戳 使用模板格式化为日期字符串
	res, err :=repository.GetStatistics(conn,motor,startStr,endStr)
	if err != nil {
		c.JSON(http.StatusOK,gin.H{
			"status":0,
			"msg":"获取数据失败"+err.Error(),
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"status":1,
		"msg": res,
	})
}

//获取当日主皮带累计情况概览
func GetMainCyCurrentDayApi(c *gin.Context){
	isPass := c.GetBool("isPass")
	if !isPass {
		return
	}
	motor,err:=xml.GetMainCyByProductionlineId(c.Request.FormValue("productionlineid"))
	if err!=nil {
		c.JSON(http.StatusOK,gin.H{
			"status":0,
			"msg":"设备不存在"+err.Error(),
		})
	}
	influxd,err:=xml.GetInfluxByProductionlineId(motor.ProductionLineId)
	if err!=nil{
		core.Logger.Println("获取产线 %s 的influxdb ，err：%s",motor.ProductionLineId,err)
		c.JSON(http.StatusOK,gin.H{
			"status":0,
			"msg":"获取产线没有数据库"+err.Error(),
		})
	}
	conn:=influx.ConnInfluxParam(influxd.Addr,influxd.User,influxd.Pwd)
	defer conn.Close()
	timeLayout := "2006-01-02 15:04:05"
	loc, _ := time.LoadLocation("Local")   //重要：获取时区
	start:=time.Date(time.Now().Year(),time.Now().Month(),time.Now().Day(),0,0,0,0,loc)
	end:=time.Date(time.Now().Year(),time.Now().Month(),time.Now().Day(),time.Now().Hour(),time.Now().Minute(),0,0,loc)
	startStr := start.Format(timeLayout) //设置时间戳 使用模板格式化为日期字符串
	endStr := end.Format(timeLayout) //设置时间戳 使用模板格式化为日期字符串
	res, err :=repository.GetStatistics(conn,motor,startStr,endStr)
	if err != nil {
		c.JSON(http.StatusOK,gin.H{
			"status":0,
			"msg":"获取数据失败"+err.Error(),
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"status":1,
		"msg": res,
	})
}

//获取当日主皮带瞬时情况概览
func GetMainCyCurrentApi(c *gin.Context){
	isPass := c.GetBool("isPass")
	if !isPass {
		return
	}
	motor,err:=xml.GetMainCyByProductionlineId(c.Request.FormValue("productionlineid"))
	if err!=nil {
		c.JSON(http.StatusOK,gin.H{
			"status":0,
			"msg":"设备不存在"+err.Error(),
		})
	}
	influxd,err:=xml.GetInfluxByProductionlineId(motor.ProductionLineId)
	if err!=nil{
		core.Logger.Println("获取产线 %s 的influxdb ，err：%s",motor.ProductionLineId,err)
		c.JSON(http.StatusOK,gin.H{
			"status":0,
			"msg":"获取产线没有数据库"+err.Error(),
		})
	}
	conn:=influx.ConnInfluxParam(influxd.Addr,influxd.User,influxd.Pwd)
	defer conn.Close()
	res, err :=repository.GetInstantStatistics(conn,motor)
	c.JSON(http.StatusOK, gin.H{
		"status":1,
		"msg": res,
	})
}

//获取当日成品皮带产量饼图
func GetBeltCysCurrentDayApi(c *gin.Context){
	isPass := c.GetBool("isPass")
	if !isPass {
		return
	}
	var productionlineid=c.Request.FormValue("productionlineid")
	influxd,err:=xml.GetInfluxByProductionlineId(productionlineid)
	if err!=nil{
		core.Logger.Println("获取产线 %s 的influxdb ，err：%s",productionlineid,err)
		c.JSON(http.StatusOK,gin.H{
			"status":0,
			"msg":"获取产线没有数据库"+err.Error(),
		})
	}
	conn:=influx.ConnInfluxParam(influxd.Addr,influxd.User,influxd.Pwd)
	defer conn.Close()
	motors,err:=xml.GetBeltCyAndNotMainByProductionlineId(productionlineid)
	if err!=nil {
		c.JSON(http.StatusOK,gin.H{
			"status":0,
			"msg":"获取成品皮带"+err.Error(),
		})
	}
	var data=map[string]float32{}
	if len(motors)>0{
		for i:=0;i<len(motors) ;i++  {
			var motor=motors[i]
			timeLayout := "2006-01-02 15:04:05"
			loc, _ := time.LoadLocation("Local")   //重要：获取时区
			start:=time.Date(time.Now().Year(),time.Now().Month(),time.Now().Day(),0,0,0,0,loc)
			end:=time.Date(time.Now().Year(),time.Now().Month(),time.Now().Day(),time.Now().Hour(),time.Now().Minute(),0,0,loc)
			startStr := start.Format(timeLayout) //设置时间戳 使用模板格式化为日期字符串
			endStr := end.Format(timeLayout) //设置时间戳 使用模板格式化为日期字符串
			res, err :=repository.GetStatistics(conn,motor,startStr,endStr)
			if err != nil {
				c.JSON(http.StatusOK,gin.H{
					"status":0,
					"msg":"获取产量数据失败"+err.Error(),
				})
			}
			bussinesskind,err:=xml.GetBussinesskindByKindAndTypeAndLineId("output",motor.MotorTypeId,motor.ProductionLineId)
			if err != nil {
				c.JSON(http.StatusOK,gin.H{
					"status":0,
					"msg":"获取产量绑定业务参数失败"+err.Error(),
				})
			}
			data[motor.Name]=res[bussinesskind.Defaultparam]
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"status":1,
		"msg": data,
	})
}

//其他设备通用
//获取当日主皮带瞬时情况概览
func GetDeviceCurrentApi(c *gin.Context){
	isPass := c.GetBool("isPass")
	if !isPass {
		return
	}
	motor,err:=xml.GetMotorByMotorId(c.Request.FormValue("motorid"))
	if err!=nil {
		c.JSON(http.StatusOK,gin.H{
			"status":0,
			"msg":"获取设备信息"+err.Error(),
		})
	}
	influxd,err:=xml.GetInfluxByProductionlineId(motor.ProductionLineId)
	if err!=nil{
		core.Logger.Println("获取产线 %s 的influxdb ，err：%s",motor.ProductionLineId,err)
		c.JSON(http.StatusOK,gin.H{
			"status":0,
			"msg":"获取产线没有数据库"+err.Error(),
		})
	}
	conn:=influx.ConnInfluxParam(influxd.Addr,influxd.User,influxd.Pwd)
	defer conn.Close()
	res, err :=repository.GetInstantStatistics(conn,motor)
	c.JSON(http.StatusOK, gin.H{
		"status":1,
		"msg": res,
	})
}

//获取产线和设备状态
func GetStatusApi(c *gin.Context){
	isPass := c.GetBool("isPass")
	if !isPass {
		return
	}
	var lineid=c.Request.FormValue("productionlineid")
	influxd,err:=xml.GetInfluxByProductionlineId(lineid)
	if err!=nil{
		core.Logger.Println("获取产线 %s 的influxdb ，err：%s",lineid,err)
		c.JSON(http.StatusOK,gin.H{
			"status":0,
			"msg":"该产线没有数据库"+err.Error(),
		})
	}
	conn:=influx.ConnInfluxParam(influxd.Addr,influxd.User,influxd.Pwd)
	defer conn.Close()
	lineStatus,err:=repository.GetLineStatus(conn,lineid)
	if err!=nil {
		c.JSON(http.StatusOK,gin.H{
			"status":0,
			"msg":"获取产线状态"+err.Error(),
		})
	}
	devsStatus,err:=repository.GetDevicesStatus(conn,lineid)
	if err!=nil {
		c.JSON(http.StatusOK,gin.H{
			"status":0,
			"msg":"获取设备状态"+err.Error(),
		})
	}
	var linestatus=map[string]bool{lineid:lineStatus}
	devsStatus=append(devsStatus,linestatus)
	//var res=Status{lineStatus,devsStatus}
	c.JSON(http.StatusOK,gin.H{
		"status":1,
		"msg":devsStatus,
	})
}
