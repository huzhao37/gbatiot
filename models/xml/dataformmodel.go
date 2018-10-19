package xml

import (
	"errors"
	"gitee.com/ha666/golibs"
	"strconv"
	"time"
	db "yuniot/framework/mysql"
)

//本地数据库
type Dataformmodel struct {
	ID	int	//编号
	MachineName	string
	Index	int
	FieldParam	string
	FieldParamEn	string
	MotorTypeName	string
	Unit	string
	DataType	string
	DataPhysicalFeature	string
	DataPhysicalAccuracy	string
	MachineId	string
	DeviceId	string
	Time	time.Time
	Value	float32
	Bit	int
	BitDesc	string
	LineId	string
	CollectdeviceIndex	string
	MotorId	string
	DataPhysicalId	int
	FormatId	int
	Remark	string
}

func ExistDataformmodel(ID int) (bool, error) {
	rows, err := db.Xml.Query("select count(0) Count from dataformmodel where ID=?", ID)
	if err != nil {
		return false, err
	}
	if len(rows) <= 0 {
		return false, nil
	}
	for _, obj := range rows {
		count, err := strconv.Atoi(string(obj["Count"]))
		if err != nil {
			return false, errors.New("parse Count error: " + err.Error())
		}
		return count > 0, nil
	}
	return false, nil
}

func InsertDataformmodel(dataformmodel Dataformmodel) (int64, error) {
	result, err := db.Xml.Exec("insert into dataformmodel(MachineName,Index,FieldParam,FieldParamEn,MotorTypeName,Unit,DataType,DataPhysicalFeature,DataPhysicalAccuracy,MachineId,DeviceId,Time,Value,Bit,BitDesc,LineId,collectdevice_index,MotorId,DataPhysicalId,FormatId,Remark) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)", dataformmodel.MachineName,dataformmodel.Index,dataformmodel.FieldParam,dataformmodel.FieldParamEn,dataformmodel.MotorTypeName,dataformmodel.Unit,dataformmodel.DataType,dataformmodel.DataPhysicalFeature,dataformmodel.DataPhysicalAccuracy,dataformmodel.MachineId,dataformmodel.DeviceId,dataformmodel.Time,dataformmodel.Value,dataformmodel.Bit,dataformmodel.BitDesc,dataformmodel.LineId,dataformmodel.CollectdeviceIndex,dataformmodel.MotorId,dataformmodel.DataPhysicalId,dataformmodel.FormatId,dataformmodel.Remark)
	if err != nil {
		return -1, err
	}
	return result.LastInsertId()
}

func UpdateDataformmodel(dataformmodel Dataformmodel) (bool, error) {
	result, err := db.Xml.Exec("update dataformmodel set MachineName=?, Index=?, FieldParam=?, FieldParamEn=?, MotorTypeName=?, Unit=?, DataType=?, DataPhysicalFeature=?, DataPhysicalAccuracy=?, MachineId=?, DeviceId=?, Time=?, Value=?, Bit=?, BitDesc=?, LineId=?, collectdevice_index=?, MotorId=?, DataPhysicalId=?, FormatId=?, Remark=? where ID=?", dataformmodel.MachineName, dataformmodel.Index, dataformmodel.FieldParam, dataformmodel.FieldParamEn, dataformmodel.MotorTypeName, dataformmodel.Unit, dataformmodel.DataType, dataformmodel.DataPhysicalFeature, dataformmodel.DataPhysicalAccuracy, dataformmodel.MachineId, dataformmodel.DeviceId, dataformmodel.Time, dataformmodel.Value, dataformmodel.Bit, dataformmodel.BitDesc, dataformmodel.LineId, dataformmodel.CollectdeviceIndex, dataformmodel.MotorId, dataformmodel.DataPhysicalId, dataformmodel.FormatId, dataformmodel.Remark, dataformmodel.ID)
	if err != nil {
		return false, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	return affected > 0, nil
}

func GetDataformmodel(ID int) (dataformmodel Dataformmodel, err error) {
	rows, err := db.Xml.Query("select ID, MachineName, Index, FieldParam, FieldParamEn, MotorTypeName, Unit, DataType, DataPhysicalFeature, DataPhysicalAccuracy, MachineId, DeviceId, Time, Value, Bit, BitDesc, LineId, collectdevice_index, MotorId, DataPhysicalId, FormatId, Remark from dataformmodel where ID=?", ID)
	if err != nil {
		return dataformmodel, err
	}
	if len(rows) <= 0 {
		return dataformmodel, nil
	}
	dataformmodels, err := _DataformmodelRowsToArray(rows)
	if err != nil {
		return dataformmodel, err
	}
	return dataformmodels[0], nil
}

func GetDataformmodels(motorId string) (dataformmodels []Dataformmodel, err error) {
	rows, err := db.Xml.Query("select * from dataformmodel where MotorId=? ",motorId)
	if err != nil {
		return dataformmodels, err
	}
	if len(rows) <= 0 {
		return dataformmodels, nil
	}
	return _DataformmodelRowsToArray(rows)
}

func GetDataformmodelByMotorIdAndParamEn(motorId string,paramen string) (dataformmodel Dataformmodel, err error) {
	rows, err := db.Xml.Query("select * from dataformmodel where MotorId=?  and FieldParamEn=?",motorId,paramen)
	if err != nil {
		return dataformmodel, err
	}
	if len(rows) <= 0 {
		return dataformmodel, nil
	}
	dataformmodels, err := _DataformmodelRowsToArray(rows)
	if err != nil {
		return dataformmodel, err
	}
	return dataformmodels[0], nil
}
func GetDataformmodelRowCount() (count int, err error) {
	rows, err := db.Xml.Query("select count(0) Count from dataformmodel")
	if err != nil {
		return -1, err
	}
	if len(rows) <= 0 {
		return -1, nil
	}
	for _, obj := range rows {
		count, err := strconv.Atoi(string(obj["Count"]))
		if err != nil {
			return -1, errors.New("parse Count error: " + err.Error())
		}
		return count, nil
	}
	return -1, nil
}

func _DataformmodelRowsToArray(maps []map[string][]byte) ([]Dataformmodel, error) {
	models := make([]Dataformmodel, len(maps))
	var err error
	for index, obj := range maps {
		model := Dataformmodel{}
		model.ID, err = strconv.Atoi(string(obj["ID"]))
		if err != nil {
			return nil, errors.New("parse ID error: " + err.Error())
		}
		model.MachineName = string(obj["MachineName"])
		model.Index, err = strconv.Atoi(string(obj["Index"]))
		if err != nil {
			return nil, errors.New("parse Index error: " + err.Error())
		}
		model.FieldParam = string(obj["FieldParam"])
		model.FieldParamEn = string(obj["FieldParamEn"])
		model.MotorTypeName = string(obj["MotorTypeName"])
		model.Unit = string(obj["Unit"])
		model.DataType = string(obj["DataType"])
		model.DataPhysicalFeature = string(obj["DataPhysicalFeature"])
		model.DataPhysicalAccuracy = string(obj["DataPhysicalAccuracy"])
		model.MachineId = string(obj["MachineId"])
		model.DeviceId = string(obj["DeviceId"])
		model.Time, err = time.ParseInLocation(golibs.Time_TIMEStandard, string(obj["Time"]), time.Local)
		if err != nil {
			return nil, errors.New("parse Time error: " + err.Error())
		}
		value, err := strconv.ParseFloat(string(obj["Value"]), 32)
		model.Value=float32(value)
		if err != nil {
			return nil, errors.New("parse Value error: " + err.Error())
		}
		model.Bit, err = strconv.Atoi(string(obj["Bit"]))
		if err != nil {
			return nil, errors.New("parse Bit error: " + err.Error())
		}
		model.BitDesc = string(obj["BitDesc"])
		model.LineId = string(obj["LineId"])
		model.CollectdeviceIndex = string(obj["collectdevice_index"])
		model.MotorId = string(obj["MotorId"])
		model.DataPhysicalId, err = strconv.Atoi(string(obj["DataPhysicalId"]))
		if err != nil {
			return nil, errors.New("parse DataPhysicalId error: " + err.Error())
		}
		model.FormatId, err = strconv.Atoi(string(obj["FormatId"]))
		if err != nil {
			return nil, errors.New("parse FormatId error: " + err.Error())
		}
		model.Remark = string(obj["Remark"])
		models[index] = model
	}
	return models, err
}
