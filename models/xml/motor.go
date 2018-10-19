package xml

import (
	"errors"
	"strconv"
	db "yuniot/framework/mysql"
	"strings"
	"gitee.com/ha666/golibs"
)

type Motor struct {
	Id	int64
	EmbeddedDeviceId	int
	FeedSize	float64	//进料尺寸
	FinalSize	float64	//出料尺寸
	MotorId	string
	MotorPower	float64	//电机功率
	MotorTypeId	string
	Name	string
	ProductSpecification	string
	ProductionLineId	string
	SerialNumber	string
	StandValue	float64	//额定值
	Time	int64
	IsBeltWeight	bool
	IsMainBeltWeight	bool
	OffSet	float64
	Slope	float64
	UseCalc	bool
}

func ExistMotor(Id int64) (bool, error) {
	rows, err := db.Xml.Query("select count(0) Count from motor where Id=?", Id)
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

func InsertMotor(motor Motor) (int64, error) {
	result, err := db.Xml.Exec("insert into motor(EmbeddedDeviceId,FeedSize,FinalSize,MotorId,MotorPower,MotorTypeId,Name,ProductSpecification,ProductionLineId,SerialNumber,StandValue,Time,IsBeltWeight,IsMainBeltWeight,OffSet,Slope,UseCalc) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)", motor.EmbeddedDeviceId,motor.FeedSize,motor.FinalSize,motor.MotorId,motor.MotorPower,motor.MotorTypeId,motor.Name,motor.ProductSpecification,motor.ProductionLineId,motor.SerialNumber,motor.StandValue,motor.Time,motor.IsBeltWeight,motor.IsMainBeltWeight,motor.OffSet,motor.Slope,motor.UseCalc)
	if err != nil {
		return -1, err
	}
	return result.LastInsertId()
}

func UpdateMotor(motor Motor) (bool, error) {
	result, err := db.Xml.Exec("update motor set EmbeddedDeviceId=?, FeedSize=?, FinalSize=?, MotorId=?, MotorPower=?, MotorTypeId=?, Name=?, ProductSpecification=?, ProductionLineId=?, SerialNumber=?, StandValue=?, Time=?, IsBeltWeight=?, IsMainBeltWeight=?, OffSet=?, Slope=?, UseCalc=? where Id=?", motor.EmbeddedDeviceId, motor.FeedSize, motor.FinalSize, motor.MotorId, motor.MotorPower, motor.MotorTypeId, motor.Name, motor.ProductSpecification, motor.ProductionLineId, motor.SerialNumber, motor.StandValue, motor.Time, motor.IsBeltWeight, motor.IsMainBeltWeight, motor.OffSet, motor.Slope, motor.UseCalc, motor.Id)
	if err != nil {
		return false, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	return affected > 0, nil
}

func GetMotor(Id int64) (motor Motor, err error) {
	rows, err := db.Xml.Query("select Id, EmbeddedDeviceId, FeedSize, FinalSize, MotorId, MotorPower, MotorTypeId, Name, ProductSpecification, ProductionLineId, SerialNumber, StandValue, Time, IsBeltWeight, IsMainBeltWeight, OffSet, Slope, UseCalc from motor where Id=?", Id)
	if err != nil {
		return motor, err
	}
	if len(rows) <= 0 {
		return motor, nil
	}
	motors, err := _MotorRowsToArray(rows)
	if err != nil {
		return motor, err
	}
	return motors[0], nil
}
func GetMotorByMotorId(motorId string) (motor Motor, err error) {
	rows, err := db.Xml.Query("select Id, EmbeddedDeviceId, FeedSize, FinalSize, MotorId, MotorPower, MotorTypeId, Name, ProductSpecification, ProductionLineId, SerialNumber, StandValue, Time, IsBeltWeight, IsMainBeltWeight, OffSet, Slope, UseCalc from motor where motorId=?", motorId)
	if err != nil {
		return motor, err
	}
	if len(rows) <= 0 {
		return motor, nil
	}
	motors, err := _MotorRowsToArray(rows)
	if err != nil {
		return motor, err
	}
	return motors[0], nil
}
// IX_Motor_MotorId_MotorTypeId_EmbeddedDeviceId
func GetMotorByMotorIdMotorTypeIdEmbeddedDeviceId(MotorId string, MotorTypeId string, EmbeddedDeviceId int) (motors []Motor, err error) {
	rows, err := db.Xml.Query("select Id, EmbeddedDeviceId, FeedSize, FinalSize, MotorId, MotorPower, MotorTypeId, Name, ProductSpecification, ProductionLineId, SerialNumber, StandValue, Time, IsBeltWeight, IsMainBeltWeight, OffSet, Slope, UseCalc from motor where motorid=? and motortypeid=? and embeddeddeviceid=?",MotorId, MotorTypeId, EmbeddedDeviceId)
	if err != nil {
		return motors, err
	}
	if len(rows) <= 0 {
		return motors, nil
	}
	return _MotorRowsToArray(rows)
}

func GetMotorRowCount() (count int, err error) {
	rows, err := db.Xml.Query("select count(0) Count from motor")
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

// IX_Motor_MotorId_MotorTypeId_EmbeddedDeviceId
func GetMotorRowListByMotorIdMotorTypeIdEmbeddedDeviceId(MotorIdIsDesc bool, MotorTypeIdIsDesc bool, EmbeddedDeviceIdIsDesc bool, PageIndex, PageSize int) (motors []Motor, err error) {
	sqlstr := golibs.NewStringBuilder()
	sqlstr.Append("select Id, EmbeddedDeviceId, FeedSize, FinalSize, MotorId, MotorPower, MotorTypeId, Name, ProductSpecification, ProductionLineId, SerialNumber, StandValue, Time, IsBeltWeight, IsMainBeltWeight, OffSet, Slope, UseCalc from motor order by")
	sqlstr.Append(" MotorId")
	if MotorIdIsDesc {
		sqlstr.Append(" Desc")
	}
	sqlstr.Append(",")
	sqlstr.Append(" MotorTypeId")
	if MotorTypeIdIsDesc {
		sqlstr.Append(" Desc")
	}
	sqlstr.Append(",")
	sqlstr.Append(" EmbeddedDeviceId")
	if EmbeddedDeviceIdIsDesc {
		sqlstr.Append(" Desc")
	}
	sqlstr.Append(" limit ?,?")
	rows, err := db.Xml.Query(sqlstr.ToString(), (PageIndex-1)*PageSize, PageSize)
	sqlstr.Append(",")
	if err != nil {
		return motors, err
	}
	if len(rows) <= 0 {
		return motors, nil
	}
	return _MotorRowsToArray(rows)
}
func GetMotors() (motors []Motor, err error) {
	rows, err := db.Xml.Query("select * from motor ")
	if err != nil {
		return motors, err
	}
	if len(rows) <= 0 {
		return motors, nil
	}
	return _MotorRowsToArray(rows)
}
func _MotorRowsToArray(maps []map[string][]byte) ([]Motor, error) {
	models := make([]Motor, len(maps))
	var err error
	for index, obj := range maps {
		model := Motor{}
		model.Id, err = strconv.ParseInt(string(obj["Id"]), 10, 64)
		if err != nil {
			return nil, errors.New("parse Id error: " + err.Error())
		}
		model.EmbeddedDeviceId, err = strconv.Atoi(string(obj["EmbeddedDeviceId"]))
		if err != nil {
			return nil, errors.New("parse EmbeddedDeviceId error: " + err.Error())
		}
		model.FeedSize, err = strconv.ParseFloat(string(obj["FeedSize"]), 64)
		if err != nil {
			return nil, errors.New("parse FeedSize error: " + err.Error())
		}
		model.FinalSize, err = strconv.ParseFloat(string(obj["FinalSize"]), 64)
		if err != nil {
			return nil, errors.New("parse FinalSize error: " + err.Error())
		}
		model.MotorId = string(obj["MotorId"])
		model.MotorPower, err = strconv.ParseFloat(string(obj["MotorPower"]), 64)
		if err != nil {
			return nil, errors.New("parse MotorPower error: " + err.Error())
		}
		model.MotorTypeId = string(obj["MotorTypeId"])
		model.Name = string(obj["Name"])
		model.ProductSpecification = string(obj["ProductSpecification"])
		model.ProductionLineId = string(obj["ProductionLineId"])
		model.SerialNumber = string(obj["SerialNumber"])
		model.StandValue, err = strconv.ParseFloat(string(obj["StandValue"]), 64)
		if err != nil {
			return nil, errors.New("parse StandValue error: " + err.Error())
		}
		model.Time, err = strconv.ParseInt(string(obj["Time"]), 10, 64)
		if err != nil {
			return nil, errors.New("parse Time error: " + err.Error())
		}
		if strings.ToLower(string(obj["IsBeltWeight"])) == "true" || obj["IsBeltWeight"][0] == 1 {
			model.IsBeltWeight = true
		} else {
			model.IsBeltWeight = false
		}
		if strings.ToLower(string(obj["IsMainBeltWeight"])) == "true" || obj["IsMainBeltWeight"][0] == 1 {
			model.IsMainBeltWeight = true
		} else {
			model.IsMainBeltWeight = false
		}
		model.OffSet, err = strconv.ParseFloat(string(obj["OffSet"]), 64)
		if err != nil {
			return nil, errors.New("parse OffSet error: " + err.Error())
		}
		model.Slope, err = strconv.ParseFloat(string(obj["Slope"]), 64)
		if err != nil {
			return nil, errors.New("parse Slope error: " + err.Error())
		}
		if strings.ToLower(string(obj["UseCalc"])) == "true" || obj["UseCalc"][0] == 1 {
			model.UseCalc = true
		} else {
			model.UseCalc = false
		}
		models[index] = model
	}
	return models, err
}
