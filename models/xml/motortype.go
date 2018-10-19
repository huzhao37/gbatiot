package xml

import (
	"errors"
	"gitee.com/ha666/golibs"
	"strconv"
	"time"
	db "yuniot/framework/mysql"
)

type Motortype struct {
	Id	int
	MotorTypeName	string
	MotorTypeId	string
	Time	time.Time
}

func ExistMotortype(Id int) (bool, error) {
	rows, err := db.Xml.Query("select count(0) Count from motortype where Id=?", Id)
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

func InsertMotortype(motortype Motortype) (int64, error) {
	result, err := db.Xml.Exec("insert into motortype(MotorTypeName,MotorTypeId,Time) values(?,?,?)", motortype.MotorTypeName,motortype.MotorTypeId,motortype.Time)
	if err != nil {
		return -1, err
	}
	return result.LastInsertId()
}

func UpdateMotortype(motortype Motortype) (bool, error) {
	result, err := db.Xml.Exec("update motortype set MotorTypeName=?, MotorTypeId=?, Time=? where Id=?", motortype.MotorTypeName, motortype.MotorTypeId, motortype.Time, motortype.Id)
	if err != nil {
		return false, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	return affected > 0, nil
}

func GetMotortype(Id int) (motortype Motortype, err error) {
	rows, err := db.Xml.Query("select Id, MotorTypeName, MotorTypeId, Time from motortype where Id=?", Id)
	if err != nil {
		return motortype, err
	}
	if len(rows) <= 0 {
		return motortype, nil
	}
	motortypes, err := _MotortypeRowsToArray(rows)
	if err != nil {
		return motortype, err
	}
	return motortypes[0], nil
}

func GetMotortypeRowCount() (count int, err error) {
	rows, err := db.Xml.Query("select count(0) Count from motortype")
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

func _MotortypeRowsToArray(maps []map[string][]byte) ([]Motortype, error) {
	models := make([]Motortype, len(maps))
	var err error
	for index, obj := range maps {
		model := Motortype{}
		model.Id, err = strconv.Atoi(string(obj["Id"]))
		if err != nil {
			return nil, errors.New("parse Id error: " + err.Error())
		}
		model.MotorTypeName = string(obj["MotorTypeName"])
		model.MotorTypeId = string(obj["MotorTypeId"])
		model.Time, err = time.ParseInLocation(golibs.Time_TIMEMYSQL, string(obj["Time"]), time.Local)
		if err != nil {
			return nil, errors.New("parse Time error: " + err.Error())
		}
		models[index] = model
	}
	return models, err
}
