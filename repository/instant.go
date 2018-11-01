package repository

import (
	"fmt"
	"yuniot/framework/influx"
	"time"
	"strconv"
	"github.com/influxdata/influxdb/client/v2"
	"encoding/json"
	"strings"
	"yuniot/models/xml"
)

//获取原始数据，默认按时间正序--所有监测点的数据
func GetInstants(conn client.Client,motor xml.Motor,
	start string,end string,where string)([]string,[]map[time.Time][]float32,error){
	var err error
	qs,err:= fmt.Sprintf("SELECT * FROM %s WHERE time >= '%s' AND time <= '%s'  AND motorid = '%s' ",
		motor.MotorTypeId,start,end,motor.MotorId) +where,nil
	res,err:=influx.QueryDB(conn,motor.ProductionLineId,qs)
	if len(res)==0||len(res[0].Series)==0{
		return nil,nil,nil
	}
	var colums=make([]string,0)
	for _,c:=range  res[0].Series[0].Columns{
		if c!=""{
			colums=append(colums, c)
		}
	}
	if err==nil{
		var valuesMaps=make([]map[time.Time][]float32,0)
		for i:=0;i<len(res[0].Series[0].Values);i++ {
			var values=make([]float32,0)
			var row=res[0].Series[0].Values[i]
			var timeStamp time.Time
			for j:=0;j<len(row);j++ {
				//time
				if j==0{
					timeStamp, err= time.Parse(time.RFC3339, row[j].(string))
					if err != nil {
						return nil,nil,err
					}
					continue
				}
				//tag
				if j==len(row)-1{
					continue
				}
				//values
				value ,err:=strconv.ParseFloat(fmt.Sprintf("%s", row[j].(json.Number)),32)
				if err!=nil{
					return colums,nil,err
				}
				values=append(values,float32(value))
			}
			var valuesMap=map[time.Time][]float32{timeStamp:values}
			valuesMaps=append(valuesMaps,valuesMap)
		}
		return colums,valuesMaps,err
	}
	return nil,nil,err
}
//获取最新一条记录
func GetLasted(conn client.Client,motor xml.Motor,
	where string)([]string,map[time.Time][]float32,error){
	var err error
	qs,err:= fmt.Sprintf("SELECT Last(*) FROM %s WHERE  motorid = '%s'  ",
		motor.MotorTypeId,motor.MotorId) +where,nil
	res,err:=influx.QueryDB(conn,motor.ProductionLineId,qs)
	if len(res)==0||len(res[0].Series)==0{
		return nil,nil,nil
	}
	var colums=make([]string,0)
	for _,c:=range  res[0].Series[0].Columns{
		if c!=""&&c!="time"{
			var cols=strings.Split(c,"last_")
			if len(cols)==2{
				colums=append(colums, cols[1])
			}
		}
	}
	if err==nil{
		//var valuesMaps=make([]map[time.Time][]float32,0)
		for i:=0;i<len(res[0].Series[0].Values);i++ {
			var values=make([]float32,0)
			var row=res[0].Series[0].Values[i]
			var timeStamp time.Time
			for j:=0;j<len(row);j++ {
				//time
				if j==0{
					timeStamp, err= time.Parse(time.RFC3339, row[j].(string))
					if err != nil {
						return nil,nil,err
					}
					continue
				}
				//tag
				//if j==len(row)-1{
				//	continue
				//}
				//values
				value ,err:=strconv.ParseFloat(fmt.Sprintf("%s", row[j].(json.Number)),32)
				if err!=nil{
					return colums,nil,err
				}
				values=append(values,float32(value))
			}
			var valuesMap=map[time.Time][]float32{timeStamp:values}
			//valuesMaps=append(valuesMaps,valuesMap)
			return colums,valuesMap,err
		}
		return colums,nil,err
	}
	return nil,nil,err
}
