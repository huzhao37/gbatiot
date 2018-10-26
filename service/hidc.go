package main

import (
	"yuniot/framework/influx"
	"fmt"
	"yuniot/core"
	"yuniot/repository"
	"time"
)

//小时统计
func main(){
	conn:=influx.ConnInflux()
	//accu
	//var where= fmt.Sprintf("AND %s >0","Accumulativeweight")
	//res,err := repository.CalcAccus(conn,"JXD001","CY","JXD001-CY1083",
	//	"2018-10-25 13:00:00","2018-10-25 14:00:00",where)

	//avg
	//var where= ""//fmt.Sprintf("AND %s >0","Accumulativeweight")
	//res,err := repository.CalcAvgs(conn,"JXD001","CY","JXD001-CY1083",
	//	"2018-10-25 13:00:00","2018-10-25 14:00:00",where)

	//boots
	//res,err := repository.CalcBootTimes(conn,"JXD001","CY","Bootflagbit","JXD001-CY1083",
	//	"2018-10-25 13:00:00","2018-10-25 14:00:00","")

	//getinstant
	//cols,res,err:=repository.GetInstants(conn,"JXD001","CY","JXD001-CY1083",
	//	"2018-10-25 13:00:00","2018-10-25 14:00:00","")

	//statistics
	t:=time.Now()
	res,err:=repository.GetStatistics(conn,"JXD001-CY1083","2018-10-25 13:00:00","2018-10-25 14:00:00")
	if err != nil {
		core.Logger.Printf("转换sql出错%s",err)
	}

	fmt.Printf(" result : %d \n 耗时：%d ms", res,time.Since(t)/1e6)
}

