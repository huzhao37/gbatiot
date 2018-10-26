package repository

import (
	"yuniot/models/xml"
	"yuniot/core"
	"github.com/influxdata/influxdb/client/v2"
)
//获取时间段内的统计信息
func GetStatistics(conn client.Client,motorid string,start string,end string)(interface{},error){
	motor,err:=xml.GetMotorByMotorId(motorid)
	if err!=nil{
		core.Logger.Printf("获取设备列表出错：%s",err)
		return nil,err
	}
	var bussiness xml.Bussinesskind
	var accDatas=make([]string,0)
	var runningParam  =""
	var datamaps=make([]map[string]float32,0)
	switch motor.MotorTypeId {
		case	"CY":
		//开机时间
		bussiness,err=xml.GetBussinesskindByKindAndType("boottimes",motor.MotorTypeId)
		if err!=nil{
			core.Logger.Printf("获取业务信息出错：%s",err)
			return nil,err
		}
		runningParam=bussiness.Defaultparam
		boottimes,err:=CalcBootTimes(conn,motor.ProductionLineId,motor.MotorTypeId,bussiness.Defaultparam,motor.MotorId,
			start,end,"")
		if err!=nil{
			core.Logger.Printf("获取开机时间出错：%s",err)
			return nil,err
		}
		var datamap=map[string]float32{"boottimes":float32(boottimes)}
		datamaps=append(datamaps,datamap)
		//累加：产量,电量
		bussiness,err=xml.GetBussinesskindByKindAndType("output",motor.MotorTypeId)
		if err!=nil{
			core.Logger.Printf("获取业务信息出错：%s",err)
			return nil,err
		}
		accDatas=append(accDatas,bussiness.Defaultparam)
		accumulativeweightMap,err:=CalcAccu(conn,motor.ProductionLineId,motor.MotorTypeId,bussiness.Defaultparam,motor.MotorId,
			start,end," ")
		if err!=nil{
			core.Logger.Printf("计算产量出错：%s",err)
			return nil,err
		}
		accumulativeweight:=accumulativeweightMap[bussiness.Defaultparam]
		datamap=map[string]float32{bussiness.Defaultparam:accumulativeweight}
		datamaps=append(datamaps,datamap)
		bussiness,err=xml.GetBussinesskindByKindAndType("totalpower",motor.MotorTypeId)
		if err!=nil{
			core.Logger.Printf("获取业务信息出错：%s",err)
			return nil,err
		}
		accDatas=append(accDatas,bussiness.Defaultparam)
		totalPowerMap,err:=CalcAccu(conn,motor.ProductionLineId,motor.MotorTypeId,bussiness.Defaultparam,motor.MotorId,
			start,end," ")
		if err!=nil{
			core.Logger.Printf("计算产量出错：%s",err)
			return nil,err
		}
		totalPower:=totalPowerMap[bussiness.Defaultparam]
		datamap=map[string]float32{bussiness.Defaultparam:totalPower}
		datamaps=append(datamaps,datamap)

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
					var avg=map[string]float32{"avg_"+param:row[param]}
					datamaps=append(datamaps,avg)
				}
			}
		}
			//负荷
			bussiness,err=xml.GetBussinesskindByKindAndType("loadstall",motor.MotorTypeId)
			if err!=nil{
				core.Logger.Printf("获取业务信息出错：%s",err)
				return nil,err
			}
			load,err:=CalcLoadStall(accumulativeweight,motor.StandValue,boottimes)
			datamap=map[string]float32{"loadstall":load}
			datamaps=append(datamaps,datamap)
		return datamaps,nil
		case	"JC":
		//开机时间
		bussiness,err=xml.GetBussinesskindByKindAndType("boottimes",motor.MotorTypeId)
		if err!=nil{
			core.Logger.Printf("获取业务信息出错：%s",err)
			return nil,err
		}
		runningParam=bussiness.Defaultparam
		boottimes,err:=CalcBootTimes(conn,motor.ProductionLineId,motor.MotorTypeId,bussiness.Defaultparam,motor.MotorId,
			start,end,"")
		if err!=nil{
			core.Logger.Printf("获取开机时间出错：%s",err)
			return nil,err
		}
		var datamap=map[string]float32{"boottimes":float32(boottimes)}
		datamaps=append(datamaps,datamap)
		//累加：电量
		bussiness,err=xml.GetBussinesskindByKindAndType("totalpower",motor.MotorTypeId)
		if err!=nil{
			core.Logger.Printf("获取业务信息出错：%s",err)
			return nil,err
		}
		accDatas=append(accDatas,bussiness.Defaultparam)
		totalPowerMap,err:=CalcAccu(conn,motor.ProductionLineId,motor.MotorTypeId,bussiness.Defaultparam,motor.MotorId,
			start,end," ")
		if err!=nil{
			core.Logger.Printf("计算电量出错：%s",err)
			return nil,err
		}
		totalPower:=totalPowerMap[bussiness.Defaultparam]
		datamap=map[string]float32{bussiness.Defaultparam:totalPower}
		datamaps=append(datamaps,datamap)

		//平均：除了累加的监测点
		avgMaps,err:=CalcAvgs(conn,motor.ProductionLineId,motor.MotorTypeId,motor.MotorId,runningParam,
			start,end," ")
		if err!=nil{
			core.Logger.Printf("获取平均值出错：%s",err)
			return nil,err
		}
		//负荷
		bussiness,err=xml.GetBussinesskindByKindAndType("loadstall",motor.MotorTypeId)
		if err!=nil{
			core.Logger.Printf("获取业务信息出错：%s",err)
			return nil,err
		}
		var value float32=0
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
					if bussiness.Defaultparam==param{
						value=row[param]
					}
					var avg=map[string]float32{"avg_"+param:row[param]}
					datamaps=append(datamaps,avg)
				}
			}
		}
		load,err:=CalcLoadStall(value,motor.StandValue,boottimes)
		datamap=map[string]float32{"loadstall":load}
		datamaps=append(datamaps,datamap)
		return datamaps,nil
	}
	return nil,nil
}
