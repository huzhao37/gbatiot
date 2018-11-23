package repository

import (
	"fmt"
	"time"
	"github.com/influxdata/influxdb/client/v2"
	"yuniot/framework/influx"
	"strconv"
	"encoding/json"
	"strings"
	"yuniot/core/extentions"
	"yuniot/core"
)

//计算累加值---所有监测点最后一条数据减去第一条数据
func CalcAccus(conn client.Client,database string,measure string,motorid string,
	start string,end string,where string)([]map[string]float32,error){
	var err error
	qs,err:= fmt.Sprintf("SELECT First(/./),Last(/./) FROM %s WHERE time >= '%s' AND time <= '%s'  AND motorid = '%s' ",
		measure,start,end,motorid) +where,nil
	res,err:=influx.QueryDB(conn,database,qs)
	if len(res)==0||len(res[0].Series)==0{
		return nil,nil
	}
	var colums=make([]string,0)
	var length=0
	for _,c:=range  res[0].Series[0].Columns{
		if c!=""&&c!="time"{
			var cols=strings.Split(c,"first_")
			if len(cols)==2{
				colums=append(colums, cols[1])
			}else{
				continue
			}
			length=len(colums)
		}
	}
	if err==nil{
		var accuMaps=make([]map[string]float32,0)
		for i:=1;i< length+1;i++  {
			frow:=res[0].Series[0].Values[0][i]
			lrow:=res[0].Series[0].Values[0][i+length]
			//if i==0{
			//	_, err := time.Parse(time.RFC3339, row.(string))
			//	if err != nil {
			//		return nil,err
			//	}
			//	continue
			//}
			first ,err:=strconv.ParseFloat(fmt.Sprintf("%s", frow.(json.Number)),32)
			if err!=nil{
				return nil,err
			}
			last ,err:=strconv.ParseFloat( fmt.Sprintf("%s", lrow.(json.Number)),32)
			if err!=nil{
				return nil,err
			}
			var accuMap=map[string]float32{colums[i-1]:float32(last-first)}
			accuMaps=append(accuMaps,accuMap)
		}
		return accuMaps,nil
	}
	return nil,err
}

//计算累加值---单个监测点最后一条数据减去第一条数据
func CalcAccu(conn client.Client,database string,measure string,feildName string,motorid string,
	start string,end string,where string)(map[string]float32,error){
	var err error
	qs,err:= fmt.Sprintf("SELECT First(%s),Last(%s) FROM %s WHERE time >= '%s' AND time <= '%s'  AND motorid = '%s'  AND %s>0 ",
		feildName,feildName,measure,start,end,motorid,feildName) +where,nil
	res,err:=influx.QueryDB(conn,database,qs)
	if len(res)==0||len(res[0].Series)==0{
		return nil,nil
	}
	if err==nil{
		for _,row:=range  res[0].Series[0].Values {
			_, err := time.Parse(time.RFC3339, row[0].(string))
			if err != nil {
				return nil,err
			}
			first ,err:=strconv.ParseFloat(fmt.Sprintf("%s", row[1].(json.Number)),32)
			if err!=nil{
				return nil,err
			}
			last ,err:=strconv.ParseFloat( fmt.Sprintf("%s", row[2].(json.Number)),32)
			if err!=nil{
				return nil,err
			}
			return map[string]float32{feildName:float32(last-first)},nil
		}
	}
	return nil,err
}

//计算累加值---单个监测点差异值的和
func CalcAccuDiff(conn client.Client,database string,measure string,feildName string,motorid string,
	start string,end string)(map[string]float32,error){
	var err error
	qs,err:= fmt.Sprintf("SELECT Sum(%s) FROM (SELECT DIFFERENCE(%s) FROM %s WHERE time >= '%s' AND time <= '%s'  AND motorid = '%s'  AND %s>-1  %s) where difference > 0",
		"difference",feildName,measure,start,end,motorid) ,nil
	res,err:=influx.QueryDB(conn,database,qs)
	if len(res)==0||len(res[0].Series)==0{
		return nil,nil
	}
	if err==nil{
		for _,row:=range  res[0].Series[0].Values {
			_, err := time.Parse(time.RFC3339, row[0].(string))
			if err != nil {
				return nil,err
			}
			val ,err:=strconv.ParseFloat(fmt.Sprintf("%s", row[1].(json.Number)),32)
			if err!=nil{
				return nil,err
			}

			return map[string]float32{feildName:float32(val)},nil
		}
	}
	return nil,err
}

//计算累加值---单个监测点非负差异值的和
func CalcAccuDiffNonNeg(conn client.Client,database string,measure string,feildName string,motorid string,
	start string,end string)(map[string]float32,error){
	var err error
	qs,err:= fmt.Sprintf("SELECT Sum(%s) FROM (SELECT NON_NEGATIVE_DIFFERENCE(%s) FROM %s WHERE time >= '%s' AND time <= '%s'  AND motorid = '%s'  ) ",
		"non_negative_difference",feildName,measure,start,end,motorid) ,nil
	res,err:=influx.QueryDB(conn,database,qs)
	if len(res)==0||len(res[0].Series)==0{
		return nil,nil
	}
	if err==nil{
		for _,row:=range  res[0].Series[0].Values {
			_, err := time.Parse(time.RFC3339, row[0].(string))
			if err != nil {
				return nil,err
			}
			val ,err:=strconv.ParseFloat(fmt.Sprintf("%s", row[1].(json.Number)),32)
			if err!=nil{
				return nil,err
			}

			return map[string]float32{feildName:float32(val)},nil
		}
	}
	return nil,err
}

//计算平均值--所有开机状态下的监测点的平均值
func CalcAvgs(conn client.Client,database string,measure string,motorid string,bootParam string,
	start string,end string,where string)([]map[string]float32,error){
	var err error
	qs,err:= fmt.Sprintf("SELECT Mean(/./) FROM %s WHERE time >= '%s' AND time <= '%s'  AND motorid = '%s' AND %s >0 ",
		measure,start,end,motorid,bootParam) +where,nil
	res,err:=influx.QueryDB(conn,database,qs)
	if len(res)==0||len(res[0].Series)==0{
		return nil,nil
	}
	var colums=make([]string,0)
	for _,c:=range  res[0].Series[0].Columns{
		if c!=""&&c!="time"{
			var cols=strings.Split(c,"mean_")
			if len(cols)==2{
				colums=append(colums, cols[1])
			}
		}
	}
	if err==nil{
		var avgMaps=make([]map[string]float32,0)
		for i:=0;i<len(res[0].Series[0].Values[0]);i++ {
			var row=res[0].Series[0].Values[0][i]
			if i==0{
				_, err := time.Parse(time.RFC3339, row.(string))
				if err != nil {
					return nil,err
				}
				continue
			}
			avg ,err:=strconv.ParseFloat(fmt.Sprintf("%s", row.(json.Number)),32)
			if err!=nil{
				return nil,err
			}
			var avgMap=map[string]float32{colums[i-1]:float32(avg)}
			avgMaps=append(avgMaps,avgMap)
		}
		return avgMaps,err
	}
	return nil,err
}

//计算平均值--开机状态下的单个监测点的平均值
func CalcAvg(conn client.Client,database string,measure string,feildName string,motorid string,bootParam string,
	start string,end string,where string)(map[string]float32,error){
	var err error
	qs,err:= fmt.Sprintf("SELECT Mean(%s) FROM %s WHERE time >= '%s' AND time <= '%s'  AND motorid = '%s' AND %s >0  ",
		feildName,measure,start,end,motorid,bootParam) +where,nil
	res,err:=influx.QueryDB(conn,database,qs)
	if len(res)==0||len(res[0].Series)==0{
		return nil,nil
	}
	if err==nil{
		for _,row:=range res[0].Series[0].Values{
			_, err := time.Parse(time.RFC3339, row[0].(string))
			if err != nil {
				return nil,err
			}
			avg ,err:=strconv.ParseFloat(fmt.Sprintf("%s", row[1].(json.Number)),32)
			if err!=nil{
				return nil,err
			}
			return map[string]float32{feildName:float32(avg)},nil
		}
	}
	return nil,err
}

//计算开机时间--输入开机标识字段,单位：分钟
func CalcBootTimes(conn client.Client,database string,measure string,feildName string,motorid string,
	start string,end string,where string)(int,error){
	var err error
	qs,err:= fmt.Sprintf("SELECT Count(%s) FROM %s WHERE time >= '%s' AND time <= '%s'  AND motorid = '%s' AND %s >0 ",
		feildName,measure,start,end,motorid,feildName) +where,nil
	res,err:=influx.QueryDB(conn,database,qs)
	if len(res)==0||len(res[0].Series)==0{
		return 0,nil
	}
	if err==nil{
		for _,row:=range res[0].Series[0].Values{
			_, err := time.Parse(time.RFC3339, row[0].(string))
			if err != nil {
				return 0,err
			}
			boots ,err:=strconv.Atoi(fmt.Sprintf("%s", row[1].(json.Number)))
			if err!=nil{
				return 0,err
			}
			return boots,nil
		}
	}
	return 0,err
}

//计算平均负荷--校准值，开机时间，额定值，单位：顿/小时
func CalcLoadStall(value float32,capacity float32,bootTimes int)( loadStall float32,err error){
	loadStall=0
	if capacity*float32(bootTimes)*value==0{
		return loadStall,err
	}
	core.Try(func() {
		loadStall= float32(extensions.Round(float64((value*60/float32(bootTimes))/capacity),3))
	}, func(e interface{}) {
		core.Logger.Printf("获取产量负荷出错：%s",e)
	})
	return loadStall*100,err //100%
}

//计算瞬时负荷--校准值，开机时间，额定值，单位：顿/小时
func CalcInstantLoadStall(value float32,capacity float32)( loadStall float32,err error){
	loadStall=0
	if capacity*value==0{
		return loadStall,err
	}
	core.Try(func() {
		loadStall= float32(extensions.Round(float64((value)/capacity),3))//value*60
	}, func(e interface{}) {
		core.Logger.Printf("获取产量负荷出错：%s",e)
	})
	return loadStall*100,err //100%
}


