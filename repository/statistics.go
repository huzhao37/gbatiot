package repository

import (
	"yuniot/models/xml"
	"yuniot/core"
	"github.com/influxdata/influxdb/client/v2"
	"fmt"
	"time"
	"yuniot/core/extentions"
)
//获取时间段内的统计信息
func GetStatistics(conn client.Client,motor xml.Motor,start string,end string)(map[string]float32,error){
	//motor,err:=xml.GetMotorByMotorId(motorid)
	//if err!=nil{
	//	core.Logger.Printf("获取设备列表出错：%s",err)
	//	return nil,err
	//}
	var err error
	var bussiness xml.Bussinesskind
	var accDatas=make([]string,0)
	var runningParam  =""
	var datamap=map[string]float32{}
	switch motor.MotorTypeId {
		case	"CY":
		//开机时间
		bussiness,err=xml.GetBussinesskindByKindAndTypeAndLineId("boottimes",motor.MotorTypeId,motor.ProductionLineId)
		if err!=nil{
			core.Logger.Printf("获取业务信息出错：%s",err)
			return nil,err
		}
		var  boottimes =0
		if bussiness.Defaultparam!=""{
			runningParam=bussiness.Defaultparam
			boots,err:=CalcBootTimes(conn,motor.ProductionLineId,motor.MotorTypeId,bussiness.Defaultparam,motor.MotorId,
				start,end,"")
			if err!=nil{
				core.Logger.Printf("获取开机时间出错：%s",err)
				return nil,err
			}
			boottimes=boots
			var runningTimes=float32(extensions.Round(float64(boots)/60,2)) //h
			datamap["boottimes"]=runningTimes}
		//datamaps=append(datamaps,datamap)
		//累加：产量,电量
		bussiness,err=xml.GetBussinesskindByKindAndTypeAndLineId("output",motor.MotorTypeId,motor.ProductionLineId)
		if err!=nil{
			core.Logger.Printf("获取业务信息出错：%s",err)
			return nil,err
		}
		var accumulativeweight float32=0
		if bussiness.Defaultparam!=""{
			accDatas=append(accDatas,bussiness.Defaultparam)
			accumulativeweightMap,err:=CalcAccuDiffNonNeg(conn,motor.ProductionLineId,motor.MotorTypeId,bussiness.Defaultparam,motor.MotorId,
				start,end)
			if err!=nil{
				core.Logger.Printf("计算产量出错：%s",err)
				return nil,err
			}
			accumulativeweight=accumulativeweightMap[bussiness.Defaultparam]
			datamap[bussiness.Defaultparam]=accumulativeweight
		}
		bussiness,err=xml.GetBussinesskindByKindAndTypeAndLineId("totalpower",motor.MotorTypeId,motor.ProductionLineId)
		if err!=nil{
				core.Logger.Printf("获取业务信息出错：%s",err)
				return nil,err
			}
		if bussiness.Defaultparam!=""{
			accDatas=append(accDatas,bussiness.Defaultparam)
			totalPowerMap,err:=CalcAccuDiffNonNeg(conn,motor.ProductionLineId,motor.MotorTypeId,bussiness.Defaultparam,motor.MotorId,
				start,end)
			if err!=nil{
				core.Logger.Printf("计算电量出错：%s",err)
				return nil,err
			}
			totalPower:=totalPowerMap[bussiness.Defaultparam]
			datamap[bussiness.Defaultparam]=totalPower
		}
		//平均：除了累加的监测点
		 avgMaps,err:=CalcAvgs(conn,motor.ProductionLineId,motor.MotorTypeId,motor.MotorId,runningParam,
			start,end," ")
		if err!=nil{
			core.Logger.Printf("获取平均值出错：%s",err)
			return nil,err
		}
		for _,row:=range avgMaps {
			for j:=0;j<len(accDatas) ;j++  {
				var is bool
				var param string
				for key:=range row {
					param=key
					if key==accDatas[j]{
						is=false
						continue
					}
				}
				if is{
					datamap["avg_"+param]=row[param]
					//datamaps=append(datamaps,avg)
				}
			}
		}
		//负荷
		bussiness,err=xml.GetBussinesskindByKindAndTypeAndLineId("loadstall",motor.MotorTypeId,motor.ProductionLineId)
		if err!=nil{
			core.Logger.Printf("获取业务信息出错：%s",err)
			return nil,err
		}
		load,err:=CalcLoadStall(datamap[bussiness.Defaultparam],motor.StandValue,boottimes)
		datamap["loadstall"]=load
		return datamap,nil
		default:
		//开机时间
		bussiness,err=xml.GetBussinesskindByKindAndTypeAndLineId("boottimes",motor.MotorTypeId,motor.ProductionLineId)
		if err!=nil{
			core.Logger.Printf("获取业务信息出错：%s",err)
			return nil,err
		}

		var boottimes=0
		if bussiness.Defaultparam!=""{
				runningParam=bussiness.Defaultparam
				boottimes,err=CalcBootTimes(conn,motor.ProductionLineId,motor.MotorTypeId,bussiness.Defaultparam,motor.MotorId,
					start,end,"")
				if err!=nil{
					core.Logger.Printf("获取开机时间出错：%s",err)
					return nil,err
				}
				datamap["boottimes"]=float32(boottimes)
			}
		//累加：电量 //可能存在多个电量累加的情况
		bussinesses,err:=xml.GetBussinesskindsByKindAndTypeAndLineId("totalpower",motor.MotorTypeId,motor.ProductionLineId)
		if err!=nil{
			core.Logger.Printf("获取业务信息出错：%s",err)
			return nil,err
		}
		if len(bussinesses)>0{
			var totalPower float32=0
			for i:=0;i<len(bussinesses) ;i++  {
				var buss=bussinesses[i]
				accDatas=append(accDatas,buss.Defaultparam)
				totalPowerMap,err:=CalcAccuDiffNonNeg(conn,motor.ProductionLineId,motor.MotorTypeId,buss.Defaultparam,motor.MotorId,
					start,end)
				if err!=nil{
					core.Logger.Printf("计算电量出错：%s",err)
					return nil,err
				}
				totalPower=totalPower+totalPowerMap[buss.Defaultparam]
			}
			datamap[bussiness.Defaultparam]=totalPower
			}
		//平均：除了累加的监测点
		avgMaps,err:=CalcAvgs(conn,motor.ProductionLineId,motor.MotorTypeId,motor.MotorId,runningParam,
			start,end," ")
		if err!=nil{
			core.Logger.Printf("获取平均值出错：%s",err)
			return nil,err
		}
		//负荷
		bussiness,err=xml.GetBussinesskindByKindAndTypeAndLineId("loadstall",motor.MotorTypeId,motor.ProductionLineId)
		if err!=nil{
			core.Logger.Printf("获取业务信息出错：%s",err)
			return nil,err
		}
		//var value float32=0
		for _,row:=range avgMaps {
			for j:=0;j<len(accDatas) ;j++  {
				var is bool
				var param string
				for key:=range row {
					param=key
					if key==accDatas[j]{
						is=false
						continue
					}
				}
				if is{
					//if bussiness.Defaultparam!=""{
					//	if bussiness.Defaultparam==param{
					//		//value=row[param]
					//	}
					//}
				//var avg=map[string]float32{"avg_"+param:row[param]}
				//datamaps=append(datamaps,avg)
				datamap["avg_"+param]=row[param]
				}
			}
		}
		var val float32=0
	    if bussiness.Defaultparam!="" {
			val=datamap["avg_"+bussiness.Defaultparam]
		}
		load,err:=CalcLoadStall(val,motor.StandValue,boottimes)
		if err!=nil{
				core.Logger.Printf("获取负荷信息：%s",err)
				return nil,err
			}
		datamap["loadstall"]=load
		return datamap,nil
		//region
		//case	"IC":
		////开机时间
		//bussiness,err=xml.GetBussinesskindByKindAndTypeAndLineId("boottimes",motor.MotorTypeId,motor.ProductionLineId)
		//if err!=nil{
		//	core.Logger.Printf("获取业务信息出错：%s",err)
		//	return nil,err
		//}
		//runningParam=bussiness.Defaultparam
		//var boottimes=0
		//if bussiness.Defaultparam!=""{
		//		boottimes,err=CalcBootTimes(conn,motor.ProductionLineId,motor.MotorTypeId,bussiness.Defaultparam,motor.MotorId,
		//			start,end,"")
		//		if err!=nil{
		//			core.Logger.Printf("获取开机时间出错：%s",err)
		//			return nil,err
		//		}
		//		var datamap=map[string]float32{}
		//		datamap["boottimes"]=float32(boottimes)
		//	}
		////累加：电量
		//bussiness,err=xml.GetBussinesskindByKindAndTypeAndLineId("totalpower",motor.MotorTypeId,motor.ProductionLineId)
		//if err!=nil{
		//	core.Logger.Printf("获取业务信息出错：%s",err)
		//	return nil,err
		//}
		//if bussiness.Defaultparam!=""{
		//		accDatas=append(accDatas,bussiness.Defaultparam)
		//		totalPowerMap,err:=CalcAccu(conn,motor.ProductionLineId,motor.MotorTypeId,bussiness.Defaultparam,motor.MotorId,
		//			start,end," ")
		//		if err!=nil{
		//			core.Logger.Printf("计算电量出错：%s",err)
		//			return nil,err
		//		}
		//		totalPower:=totalPowerMap[bussiness.Defaultparam]
		//		datamap[bussiness.Defaultparam]=totalPower
		//	}
		//
		////平均：除了累加的监测点
		//avgMaps,err:=CalcAvgs(conn,motor.ProductionLineId,motor.MotorTypeId,motor.MotorId,runningParam,
		//	start,end," ")
		//if err!=nil{
		//	core.Logger.Printf("获取平均值出错：%s",err)
		//	return nil,err
		//}
		////负荷
		//bussiness,err=xml.GetBussinesskindByKindAndTypeAndLineId("loadstall",motor.MotorTypeId,motor.ProductionLineId)
		//if err!=nil{
		//	core.Logger.Printf("获取业务信息出错：%s",err)
		//	return nil,err
		//}
		//
		//for _,row:=range avgMaps {
		//	for j:=0;j<len(accDatas) ;j++  {
		//		var is bool
		//		var param string
		//		for key:=range row {
		//			param=key
		//			if key==accDatas[j]{
		//				is=false
		//				continue
		//			}
		//		}
		//		if is{
		//				datamap["avg_"+param]=row[param]
		//			}
		//
		//	}
		//}
		//var val float32=0
		//if bussiness.Defaultparam!="" {
		//	val=datamap["avg_"+bussiness.Defaultparam]
		//}
		//load,err:=CalcLoadStall(val,motor.StandValue,boottimes)
		//datamap["loadstall"]=load
		//return datamap,nil
		//region
	}
	return nil,nil
}

//获取瞬时统计信息
func GetInstantStatistics(conn client.Client,motor xml.Motor)(map[string]interface{},error){
	var err error
	var bussiness xml.Bussinesskind
	loc, _ := time.LoadLocation("Local")   //重要：获取时区
	timeLayout := "2006-01-02 15:04:05"
	start:=time.Date(time.Now().Year(),time.Now().Month(),time.Now().Day(),0,0,0,0,loc)
	startStr := start.Format(timeLayout)
	cols,res, err :=GetLasted(conn,motor,fmt.Sprintf(" AND time > '%s'  ",startStr))
	if err!=nil{
		core.Logger.Printf("获取最新一条数据出错：%s",err)
		return nil,err
	}
	var m=map[string]interface{}{}
	if len(cols)>0{
		for i:=0;i<len(cols) ;i++  {
			if len(res)>0{
				for key,val:= range  res{
					m["time"]=key
					if len(val)>0{
						m[cols[i]]=val[i]
					}
				}
			}
		}
	}
	//瞬时负荷
	if len(m)>0{
		//瞬时负荷
		bussiness,err=xml.GetBussinesskindByKindAndTypeAndLineId("instantloadstall",motor.MotorTypeId,motor.ProductionLineId)
		if err!=nil{
			core.Logger.Printf("获取业务信息出错：%s",err)
			return nil,err
		}
		var val float32=0
		if bussiness.Defaultparam!=""{
			val=m[bussiness.Defaultparam].(float32)
		}
		load,err:=CalcInstantLoadStall(val,motor.StandValue)
		if err != nil {
			core.Logger.Printf("计算%s瞬时负荷：%s",motor.MotorId,err)
			return m,err
		}
		m["loadstall"]=load
	}
	return m,nil
}

//获取产线状态
func GetLineStatus(conn client.Client,productionlineId string)(bool,error){
	line,err:=xml.GetProductionlineById(productionlineId)
	if err!=nil{
		return false,err
	}
	return (time.Now().Unix()-line.Time)<=600,nil
}

//获取该产线下所有设备状态
func GetDevicesStatus(conn client.Client,productionlineId string)([]map[string]bool,error){
	motors,err:=xml.GetMotorsByProductionlineId(productionlineId)
	if err!=nil||len(motors)==0{
		return nil,err
	}
	loc, _ := time.LoadLocation("Local")   //重要：获取时区
	timeLayout := "2006-01-02 15:04:05"
	start:=time.Date(time.Now().Year(),time.Now().Month(),time.Now().Day(),time.Now().Hour(),0,0,0,loc)
	startStr := start.Format(timeLayout)
	var devsStatus=make([]map[string]bool,0)
	for _,m:= range motors {
		_,last,err:=GetLasted(conn,m,startStr)
		var status=false
		if err==nil{
			for key,_ := range last {
				status=time.Since(key).Minutes()<=10
			}
		}
		var devStatus=map[string]bool{m.MotorId:status}
		devsStatus=append(devsStatus,devStatus)
	}
	return devsStatus,nil
}
