package main

import (
	"yuniot/framework/rabbitmq"
	"yuniot/models/xml"
	"yuniot/core"
	"time"
	"yuniot/repository"
	"yuniot/framework/redis"
	db "yuniot/framework/mysql"
	"strconv"
	"fmt"
	"github.com/go-linq"
	"github.com/influxdata/influxdb/client/v2"
	"yuniot/core/extentions"
	"yuniot/framework/influx"
	"yuniot/models/others"
	"sort"
	"strings"
)
var (
	myConfig  = new(core.Config)
    MParms []xml.Motorparams

	)
func init() {
	myConfig.InitConfig("./config/config.txt")
	//var env = myConfig.Read("global", "env")
	redis.Redis.Init(5)
}
func main() {
	defer db.SqlDB.Close()
	//1.数据解析- 单独协程
	//频繁使用的集合数据
	var err error
	MParms,err=xml.GetMotorparamList()
	if err!=nil{
		core.Logger.Println("不存在任何MotorParams,err：%s",err)
	}
		mqs,err:=xml.GetMessagequeues()
		if err!=nil || len(mqs)==0{
			core.Logger.Panicln("获取队列列表出错:%s",err)
		}
	for   _,q :=range  mqs {
	//go func(mq *xml.Messagequeue) {
	//var remains =0
	//if q.Id>2{
	//	continue
	//}
	 go rabbitmq.Read2(DataParse, q.Username, q.Pwd, q.Host+":"+strconv.Itoa(q.Port), q.Routekey, 0,q.Id)
	 time.Sleep(3*time.Second)
	//}(&q)
	}
	//阻塞
	select {}
}

//数据解析
func DataParse(bytes []byte,mqid int) (bool){
	var (
		_model     *xml.DataGram
		_values    []int
		_time    time.Time
		_unixTimeStr    string//当日日期
		_normalPoints []*client.Point
		//_alarmPoints []*client.Point
		//_otherDiPoints []*client.Point
	)
	t1 := time.Now() // get current time
	mq,err :=xml.GetMessagequeue(mqid)
	if err!=nil{
		core.Logger.Panicln("该队列不存在！data :%s , err：%s",mqid,err)
		return false
	}
	//repository.Init()
	_, newBytes := core.BytesSplit(bytes, 6)
	Collectdeviceindex := core.BytesConvertHexArr(newBytes)
	if strings.ToUpper(Collectdeviceindex)!=mq.Collectdeviceindex{
		return false
	}
	var datagram=repository.BytesParse(bytes)//数据接收器
	if strings.ToUpper(Collectdeviceindex)!=mq.Collectdeviceindex{
		return false
	}
	if len(datagram.PValues)<=0{
		core.Logger.Panicln("数据接收器无值数据！data :%s , err：%s",datagram.Collectdeviceindex,err)
		return false
	}
	var globalTime int64
	_values =make([]int,0)
	for key,v := range datagram.PValues {
		globalTime = key
		_values=v
	}
	timeLayout := "2006-01-02 15:04:05"
	loc, _ := time.LoadLocation("Local")   //重要：获取时区
	dataTimeStr := time.Unix(globalTime, 0).Format(timeLayout) //设置时间戳 使用模板格式化为日期字符串
	theTime, err := time.ParseInLocation(timeLayout, dataTimeStr, loc) //使用模板在对应时区转化为time.time类型
	if err!=nil{
		core.Logger.Panicln("时间转换错误！data : %s ,err：%s",dataTimeStr,err)
		return false
	}
	_model=datagram
	_time=theTime
	_unixTimeStr=strconv.FormatInt(time.Date(theTime.Year(), theTime.Month(), theTime.Day(),0, 0, 0, 0, theTime.Location()).Unix(),10)
	//todo
	collectDevice,err:=xml.GetCollectdeviceByIndex(_model.Collectdeviceindex)
	if err!=nil{
		core.Logger.Println("根据Index %s 获取相关采集设备err：%s",_model.Collectdeviceindex,err)
		return false
	}
	//存储至influxdb
	_normalPoints=make([]*client.Point,0)
	influxd,err:=xml.GetInfluxByProductionlineId(collectDevice.ProductionlineId)
	if err!=nil{
		core.Logger.Println("获取产线 %s 的influxdb ，err：%s",collectDevice.ProductionlineId,err)
		return false
	}
	conn:=influx.ConnInfluxParam(influxd.Addr,influxd.User,influxd.Pwd)
	defer conn.Close()
	//更新产线最新时间
	line,err:=xml.GetProductionlineById(collectDevice.ProductionlineId)
	if err!=nil{
		core.Logger.Println("该采集设备%s不存在任何产线err：%s",collectDevice.Index,err)
		//return false
	}
	line.Time=globalTime
	_,err=xml.UpdateProductionline(line)
	if err!=nil{
		core.Logger.Println("更新%s 失败 err：%s",line.Name,err)
		//return false
	}

	motors,err:=xml.GetMotorsByCollectDeviceId(collectDevice.Id)
	if err!=nil{
		core.Logger.Println("不存在任何电机err：%s",err)
		return false
	}
	if len(motors) > 0 {
		for _, motor := range motors {
			var id=motor.MotorId
			var code=motor.MotorTypeId
			tags := map[string]string{"motorid":id}
			dataformModels,err:=xml.GetDataformmodels(id)
			if err!=nil{
				core.Logger.Println("不存在数据映射表单！data : %s ,err：%s",id,err)
				return false
			}
			var sccParms []xml.Motorparams
			if MParms != nil && len(MParms) > 0 {
				linq.From(MParms).Where(func(i interface{}) bool {
					return i.(xml.Motorparams).MotorTypeId==code&& i.(xml.Motorparams).ProductionLineId==motor.ProductionLineId
				}).ToSlice(&sccParms)
			}
			bussiness,err:=xml.GetBussinessByMotorid(id)
			if err!=nil{
				core.Logger.Println("获取设备业务参数出错！data : %s ,err：%s",id,err)
				return false
			}
			fields := map[string]interface{}{}
			var normalDataModels=make([]*xml.Dataformmodel,0)
			var diDataModels=make([]*xml.Dataformmodel,0)
			var alarmDataModels=make([]*xml.Dataformmodel,0)

			if len(dataformModels)>0{
				for _,dataformmodel :=range  dataformModels  {
					var normalModel =dataformmodel
					var isDi=true
					if len(bussiness)>0{
						for _,mb :=range  bussiness  {
							if dataformmodel.FieldParamEn==mb.Param{
								normalDataModels=append(normalDataModels,&dataformmodel)
								isDi=false
							}
						}
					}
					if len(sccParms) > 0 {
						for _, sccparm := range sccParms {
							if normalModel.FieldParamEn==sccparm.Param{
								normalDataModels=append(normalDataModels,&normalModel)
								isDi=false
							}
						}
					}
					if  normalModel.DataPhysicalId==21{
						alarmDataModels=append(alarmDataModels,&normalModel)
						isDi=false
					}
					if isDi==true{
						diDataModels=append(diDataModels,&normalModel)
					}
				}}
			if len(normalDataModels)>0{
				var ais=make([]others.AI,0)
				for _,dataformmodel :=range  normalDataModels  {
					var param = core.StrFirstToUpper(dataformmodel.FieldParamEn)
					var paramValue=ConvertToNormal(dataformmodel,_values)
					fields[param]=paramValue
					ais=append(ais,others.AI{Param:param,Value:paramValue})
				}
				var tormorrowTime=_time.Add(24*time.Hour)//明天
				expireTime:= time.Date(tormorrowTime.Year(), tormorrowTime.Month(), tormorrowTime.Day(),
					0, 0, 0, 0, tormorrowTime.Location())//明天凌晨
				fmt.Println(expireTime.Format("2006-01-02 00:00:00"))
				//ai分析使用
				redis.Redis.Select(5)
				redis.Redis.HSet(id+"|"+_unixTimeStr,_time.Format("2006-01-02 15:04:05"),ais)
				//var lifeTime=core.GetSecs(_time,endTime)
				redis.Redis.ExpireAt(id+"|"+_unixTimeStr,expireTime.Unix()) //1 days
				//入influx库持久化
				pt, err := client.NewPoint(
					code,
					tags,
					fields,
					_time,
				)
				if err != nil {
					core.Logger.Fatal(err)
				}
				_normalPoints=append(_normalPoints,pt)
			}
			if len(alarmDataModels)>0{
				var alarmContents =""
				for _,dataformmodel :=range  alarmDataModels  {
					if _values[dataformmodel.Index]==1{
						alarmContents=dataformmodel.Remark+"|"
					}
				}
				var alarm=others.Alarm{Content:alarmContents,Time:_time}
				redis.Redis.Select(3)
				redis.Redis.Lpushs(id+"|"+_unixTimeStr,alarm)
				redis.Redis.Expire(id+"|"+_unixTimeStr,604800)//7 days
				//fields["Content"]=alarmContents
				//fields["MotorName"]=motor.Name
				//pt, err := client.NewPoint(
				//	code,
				//	tags,
				//	fields,
				//	_time,
				//)
				//if err != nil {
				//	core.Logger.Fatal(err)
				//}
				//_alarmPoints=append(_alarmPoints,pt)
			}
			if len(diDataModels)>0{
				//倒序排序,按照param排序
				sort.Slice(diDataModels, func(i, j int) bool {
					a:=diDataModels[i].Remark
					b:=diDataModels[j].Remark
					return  a>b
				})
				var diParams=""
				var diValues=make([]int,0)
				for _,dataformmodel :=range  diDataModels  {
					var param = core.StrFirstToUpper(dataformmodel.Remark)
					var paramValue=_values[dataformmodel.Index]
					diParams=diParams+param+"|"
					diValues=append(diValues,paramValue)
					//if paramValue>0{
					//}
					//fields[param]=paramValue
					//pt, err := client.NewPoint(
					//	code,
					//	tags,
					//	fields,
					//	_time,
					//)
					//if err != nil {
					//	core.Logger.Fatal(err)
					//	return false
					//}
					//_otherDiPoints=append(_otherDiPoints,pt)
				}
				var diData=others.DI{ others.DITitle{Params:diParams},others.DICache{Values:diValues,Time:_time}}
				redis.Redis.Select(15)
				exist,err:=redis.Redis.Exists(id)
				if err!=nil{
					core.Logger.Println("获取di参数出错！data : %s ,err：%s",id,err)
					return false
				}
				if !exist{
					redis.Redis.Set(id,diParams)//保存di参数
				}
				redis.Redis.Select(1)
				redis.Redis.HSet(id+"|"+_unixTimeStr,_time.Format("2006-01-02 15:04:05"),diData.DICaches.Values)
				redis.Redis.Expire(id+"|"+_unixTimeStr,7776000) //90 days
			}
		}
	}
	//存储至influxdb
	collect,err:=xml.GetCollectdeviceByIndex(mq.Collectdeviceindex)
	if err!=nil{
		core.Logger.Println("不存在任何采集设备err：%s",err)
		return false
	}
	//倒序排序,按照tag排序
	sort.Slice(_normalPoints, func(i, j int) bool {
		a:=_normalPoints[i].Tags()["motorid"]
		b:=_normalPoints[j].Tags()["motorid"]
		return  a>b
	})
	exist,err:=influx.ExistDB(conn,collect.ProductionlineId)
	if err==nil&&!exist{
		_,err=influx.CreateDB(conn,collect.ProductionlineId)
		if err!=nil{
			core.Logger.Println("不创建数据库%s err：%s",collect.ProductionlineId,err)
			return false
		}
	}
	influx.WritesPoints(conn,collect.ProductionlineId,_normalPoints)
	//influx.WritesPoints(conn,"alarm",_alarmPoints)
	//influx.WritesPoints(conn,"otherdi",_otherDiPoints)

	elapsed := time.Since(t1)
	fmt.Println("耗时:%d ", elapsed)
	//fmt.Println("write tps: %d ms", 1e14/elapsed)

	//存储原始数据
	xml.InsertOriginalbytes(xml.Originalbytes{Productionlineid:collect.ProductionlineId,Collectdeviceindex:mq.Collectdeviceindex,
		Time:theTime,Bytes:core.BytesConvertHexArr(bytes)})
	return true
}



//辅助方法


/// <summary>
/// 根据数据精度和数据参数将数值转化为实际值
/// </summary>
/// <param name="form">数据表单集合</param>
/// <param name="values">数据值集合</param>
func  ConvertToNormal(form *xml.Dataformmodel , values []int )( float32 ){
	if form.Index >= len(values) {
	core.Logger.Println("[Normalize]excite values index")
	return 0
}
	 var accur float64=1
	accur,err :=strconv.ParseFloat(form.DataPhysicalAccuracy,32)
	if err!=nil{
		core.Logger.Println("数据物理精度错误！data : %s ,err：%s",form.DataPhysicalAccuracy,err)
		return 0
	}
var oldValue = values[form.Index]
switch (form.DataPhysicalFeature){
case "温度":
var des = extensions.TempTranster(oldValue)
if oldValue != -1{
	return float32(extensions.Round(float64(des) * accur, 2))
}
case "电流":
	if oldValue != -1{
		return float32(extensions.Round(float64(oldValue) * accur, 2))
	}
case "配置":
if form.FieldParamEn=="Unit"{
	var tempInt=0
	if oldValue != -1{
		tempInt =int(oldValue)
	}
var  value = tempInt
if (tempInt != -1){
	value = tempInt & 7
}
return float32(value)
}
case "称重":
	unitForm,err:=xml.GetDataformmodelByMotorIdAndParamEn(form.MotorId,"Unit")
	if err!=nil||unitForm.ID == 0{
		core.Logger.Println("该电机不存在称重单位！data : %s ,err：%s",form.MotorId,err)
		return 0
	}
unitForm.Value =ConvertToNormal(&unitForm,values)
	var originalValue float32=0
if oldValue != -1{
	 originalValue =float32( extensions.Round(float64(oldValue)  *accur, 2))
	}

return ConveyorWeightConvert(int(unitForm.Value), form.FieldParam, originalValue)
case "瞬时称重":
	unitForm2,err:=xml.GetDataformmodelByMotorIdAndParamEn(form.MotorId,"Unit")
	if err!=nil||unitForm2.ID == 0{
		core.Logger.Println("该电机不存在称重单位！data : %s ,err：%s",form.MotorId,err)
		return 0
	}
unitForm2.Value = ConvertToNormal(&unitForm2, values)
	var originalValue2 float32=0
	if oldValue != -1{
		 originalValue2 =float32(extensions.Round(float64(oldValue) * accur, 2))
	}
	return ConveyorWeightConvert(int(unitForm2.Value), form.FieldParam, originalValue2)
}
if oldValue != -1{
	return float32(extensions.Round(float64(oldValue) * accur, 2))
	}
   return 0
}

/// <summary>
/// 根据单位计算皮带机瞬时称重、累计称重
/// </summary>
/// <param name="unit">单位</param>
/// <param name="param">称重参数</param>
/// <param name="oldValue">称重原始值</param>
func ConveyorWeightConvert(unit int ,  param string ,oldValue float32 )(float32) {
switch unit{
case 0:
if param=="瞬时称重"{
	if oldValue != -1{
		return float32(extensions.Round(float64(oldValue / 3600), 2))
	}
	return oldValue
	}
case 1:
if param=="累计称重" {
	if oldValue != -1{
		oldValue=float32(extensions.Round(float64(oldValue / 1000), 2))
	}
if oldValue < -1{
//4294967295
oldValue =float32(extensions.Round(float64((4294967295 + oldValue * 1000) / 1000), 2))
}
return oldValue
}
if param=="瞬时称重"{
	if oldValue != -1{
		oldValue=float32(extensions.Round(float64(oldValue / 3.6), 2))
	}
	return oldValue
}
case 2:
if param=="瞬时称重"{
	if oldValue != -1{
		oldValue=float32(extensions.Round(float64(oldValue / 1000), 2))
	}
	return oldValue
}
case 3:
if param=="累计称重"{
	if oldValue != -1{
		oldValue=float32(extensions.Round(float64(oldValue / 1000), 2))
	}
	if oldValue < -1{
		//4294967295
		oldValue =float32(extensions.Round(float64((4294967295 + oldValue * 1000) / 1000), 2))
	}
	return oldValue
}
case 4:
if param=="瞬时称重"{
	if oldValue != -1{
		oldValue=float32(extensions.Round(float64(oldValue* 60), 2))
	}
	return oldValue
}
case 5:
if param=="累计称重"{
	if oldValue != -1{
		oldValue=float32(extensions.Round(float64(oldValue/ 1000), 2))
	}
	if oldValue < -1{
		//4294967295
		oldValue = float32(extensions.Round(float64((4294967295 + oldValue * 1000) / 1000), 2))
	}
	return oldValue
}
if param=="瞬时称重"{
	if oldValue != -1{
		oldValue=float32(extensions.Round(float64(oldValue * 0.06), 2))
	}
	return oldValue
}
case 6:
break
case -1:
if param=="累计称重"||param=="瞬时称重"{
	return -1
}
break
default:
return oldValue
}
return  oldValue
}


