package apis

import (
	"net/http"
	"gopkg.in/gin-gonic/gin.v1"
	"yuniot/models/xml"
)

//获取设备信息
func GetMotorsApi(c *gin.Context){
	isPass := c.GetBool("isPass")
	if !isPass {
		return
	}
	motors,err:=xml.GetMotorsByProductionlineId(c.Request.FormValue("productionlineid"))
	if err!=nil {
		c.JSON(http.StatusOK,gin.H{
			"status":0,
			"msg":"获取设备信息"+err.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status":1,
		"msg": motors,
	})
}
