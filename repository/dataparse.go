package repository

import (
	"fmt"
	"github.com/go-linq"
	"time"
	"yuniot/core"
	model "yuniot/models/xml"
	"strings"
)

var (
	DataTypes   []model.Datatype
	//DataConfigs []model.Dataconfig

)

func Init() {
	//engine.Form.Close() //关闭数据库连接
}

//原始字节解析协议
func BytesParse(bytes []byte) (dataGram *model.DataGram) {

	fmt.Printf("当前解析数据：%s",core.BytesConvertHexArr(bytes))

	var remainsByte,newBytes []byte
	dataGram = new(model.DataGram)
	dataGram.PValues = make(map[int64][]int,0)
	remainsByte, newBytes = core.BytesSplit(bytes, 6)
	dataGram.Collectdeviceindex = core.BytesConvertHexArr(newBytes)
	remainsByte, newBytes = core.BytesSplit(remainsByte, 2)
	dataGram.DataConfigId = core.BytesToInt(newBytes)
	remainsByte, newBytes = core.BytesSplit(remainsByte, 1)
	dataGram.Count = core.BytesToInt(newBytes)

	fmt.Printf("[dataparse]主控设备ID %s \n", dataGram.Collectdeviceindex)
	fmt.Printf("[dataparse]表单数据配置ID %d \n", dataGram.DataConfigId)
	var configModels =make([]model.Dataconfig,0)

	DataTypes,err := model.GetDatatypes()
	//DataConfigs,err2 := model.GetDataconfigs()
	if err != nil  {
		core.Logger.Println("[dataparse]查询form数据库出错: %s \n", err.Error())
	}

	if len(DataTypes) == 0  {
		core.Logger.Println("[dataparse]Error: Cannot find config Info in SqlDb \n")
	}
	configModels,err=model.GetDataconfigByCollectdeviceindex( strings.ToUpper(dataGram.Collectdeviceindex))
	if err!=nil{
		core.Logger.Println("[dataparse]Error: %s; \n",err)
	}
	if len(configModels) == 0 {
		core.Logger.Println("[dataparse]Error: 没有对应的数据表单; \n")
	}

	for i := 0; i < dataGram.Count; i++ {
		var parmValues []int=make([]int,0)
		var timeValue = time.Now().Unix()
		for i := 0; i < len(configModels); i++ {
			config := configModels[i]
			//在此处与协议规定的以大数优先
			var typed = config.Datatypeid
			var typeModels []model.Datatype
			core.Try(func() {
				linq.From(DataTypes).Where(func(i interface{}) bool {
					return i.(model.Datatype).Id == typed
				}).ToSlice(&typeModels)
			}, func(e interface{}) {
				core.Logger.Println("[dataparse]出错:%s \n", e)
			})
			typeModel := typeModels[0]
			if typeModel.Id == 0 {
				return dataGram
			}
			var temp =make([]byte,0)
			switch typed {
			case 4:
				{
					//4个字节(记录ID);
					for j := 0; j < int(config.Count); j++ {
						remainsByte, temp = core.BytesSplit(remainsByte, int(typeModel.Inbyte))
						var value = core.BytesToUInt(temp)
						var resultValue int
						if value == 4294967295 {
							resultValue = -1
						}
						resultValue = int(value)
						parmValues = append(parmValues, resultValue)
						fmt.Printf("[dataparse] %d \n", resultValue)
					}
					continue
				}
			case 5:
				{
					//4个字节(时间);
					remainsByte, temp = core.BytesSplit(remainsByte, int(typeModel.Inbyte))
					var value = core.BytesToTime(temp)
					timeValue = value
					parmValues = append(parmValues, int(timeValue))
					fmt.Printf("[dataparse] %s \n", time.Unix(value, 0).Format("2006-01-02 15:04:05"))
					continue
				}
			case 7:
				{
					//4个字节(整形模拟量)        FFFFFFFF - 4294967295 (uint)
					for j := 0; j < int(config.Count); j++ {
						remainsByte, temp = core.BytesSplit(remainsByte, int(typeModel.Inbyte))
						var value = core.BytesToUInt(temp)
						var resultValue int
						if value == 4294967295 {
							resultValue = -1
						}
						resultValue = int(value)
						parmValues = append(parmValues, resultValue)
						fmt.Printf("[dataparse] %d \n", resultValue)
					}
					continue
				}
			case 9:
				{
					//4个字节
					for j := 0; j < int(config.Count); j++ {
						remainsByte, temp = core.BytesSplit(remainsByte, int(typeModel.Inbyte))
						var value = core.BytesToUInt(temp)
						var resultValue int
						if value == 4294967295 {
							resultValue = -1
						}else{
							resultValue = int(value)
						}
						parmValues = append(parmValues, resultValue)
						fmt.Printf("[dataparse] %d \n", resultValue)
					}
					continue
				}
			case 11:
				{
					//2个字节(整形模拟量)   FFFF - 65535
					for j := 0; j < int(config.Count); j++ {
						remainsByte, temp = core.BytesSplit(remainsByte, int(typeModel.Inbyte))
						var value = core.BytesToInt(temp)
						var resultValue int
						if value == 65535 {
							resultValue = -1
						} else {
							resultValue = value
						}
						parmValues = append(parmValues, resultValue)
						fmt.Printf("[dataparse] %d \n", resultValue)
					}
					continue
				}
			case 12:
				{
					//12位，每次读取3个字节，2个数据;  FFF - 4095
					var total int
					if config.Count%2 == 0 {
						total = int(config.Count) / 2
					} else {
						total = int(config.Count)/2 + 1
					}
					var tempCount = 0
					for j := 0; j < total; j++ {
						remainsByte, temp = core.BytesSplit(remainsByte, int(typeModel.Inbyte))
						var values = core.ByteToInts(temp)
						var tempValue int
						for _, value := range values {
							if value == 4095 {
								tempValue = -1
							} else {
								tempValue = value
							}
							tempCount++
							if tempCount> config.Count {
								continue
							}
							parmValues = append(parmValues, tempValue)
							fmt.Printf("[dataparse] %d \n", tempValue)
						}
					}
					continue
				}
			case 13:
				{
					//1个字节(整形模拟量)    FF - 255
					for j := 0; j < int(config.Count); j++ {
						remainsByte, temp = core.BytesSplit(remainsByte, int(typeModel.Inbyte))
						var value = core.BytesToInt(temp)
						var resultValue int
						if value == 255 {
							resultValue = -1
						} else {
							resultValue = value
						}
						parmValues = append(parmValues, resultValue)
						fmt.Printf("[dataparse] %d \n", resultValue)
					}
					continue
				}
			case 14:
				{
					/*一个8个字节，7*8+1=57 57个参数(最后一个字节只有一位有效位)*/
					//1位，8个一位组成一个字节,每次读取一个字节，返回8个数据;
					var total int
					if config.Count%8 == 0 {
						total = int(config.Count) / 8
					} else {
						total = int(config.Count)/8 + 1
					}
					var tempCount = 0
					for j := 0; j < total; j++ {
						remainsByte, temp = core.BytesSplit(remainsByte, int(typeModel.Inbyte))
						newBytes := temp[0]
						var values = core.ByteToInts1(newBytes)
						for _, value := range values {
							tempCount++
							if tempCount > int(config.Count) {
								continue
							}
							parmValues = append(parmValues, value)
							fmt.Printf("[dataparse] %d \n", value)
						}
					}
					continue
				}
			default:
				continue
			}
		}
		dataGram.PValues[timeValue] = parmValues
	}
	return dataGram
}
