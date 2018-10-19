package xml

import (
	"errors"
	"gitee.com/ha666/golibs"
	"strconv"
	"time"
	db "yuniot/framework/mysql"
)

type Motorparams struct {
	Id	int
	Param	string
	Description	string
	MotorTypeId	string
	PhysicId	int
	Time	time.Time
}

func ExistMotorparams(Id int) (bool, error) {
	rows, err := db.Xml.Query("select count(0) Count from motorparams where Id=?", Id)
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

func InsertMotorparams(motorparams Motorparams) (int64, error) {
	result, err := db.Xml.Exec("insert into motorparams(Param,Description,MotorTypeId,PhysicId,Time) values(?,?,?,?,?)", motorparams.Param,motorparams.Description,motorparams.MotorTypeId,motorparams.PhysicId,motorparams.Time)
	if err != nil {
		return -1, err
	}
	return result.LastInsertId()
}

func UpdateMotorparams(motorparams Motorparams) (bool, error) {
	result, err := db.Xml.Exec("update motorparams set Param=?, Description=?, MotorTypeId=?, PhysicId=?, Time=? where Id=?", motorparams.Param, motorparams.Description, motorparams.MotorTypeId, motorparams.PhysicId, motorparams.Time, motorparams.Id)
	if err != nil {
		return false, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	return affected > 0, nil
}

func GetMotorparams(Id int) (motorparams Motorparams, err error) {
	rows, err := db.Xml.Query("select Id, Param, Description, MotorTypeId, PhysicId, Time from motorparams where Id=?", Id)
	if err != nil {
		return motorparams, err
	}
	if len(rows) <= 0 {
		return motorparams, nil
	}
	motorparamss, err := _MotorparamsRowsToArray(rows)
	if err != nil {
		return motorparams, err
	}
	return motorparamss[0], nil
}

func GetMotorparamsRowCount() (count int, err error) {
	rows, err := db.Xml.Query("select count(0) Count from motorparams")
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
func GetMotorparamList() (motorparams []Motorparams, err error) {
	rows, err := db.Xml.Query("select * from motorparams ")
	if err != nil {
		return motorparams, err
	}
	if len(rows) <= 0 {
		return motorparams, nil
	}
	return _MotorparamsRowsToArray(rows)
}
func _MotorparamsRowsToArray(maps []map[string][]byte) ([]Motorparams, error) {
	models := make([]Motorparams, len(maps))
	var err error
	for index, obj := range maps {
		model := Motorparams{}
		model.Id, err = strconv.Atoi(string(obj["Id"]))
		if err != nil {
			return nil, errors.New("parse Id error: " + err.Error())
		}
		model.Param = string(obj["Param"])
		model.Description = string(obj["Description"])
		model.MotorTypeId = string(obj["MotorTypeId"])
		model.PhysicId, err = strconv.Atoi(string(obj["PhysicId"]))
		if err != nil {
			return nil, errors.New("parse PhysicId error: " + err.Error())
		}
		model.Time, err = time.ParseInLocation(golibs.Time_TIMEStandard, string(obj["Time"]), time.Local)
		if err != nil {
			return nil, errors.New("parse Time error: " + err.Error())
		}
		models[index] = model
	}
	return models, err
}
