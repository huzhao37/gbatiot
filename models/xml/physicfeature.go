package xml

import (
	"errors"
	"strconv"
	"time"
	db "yuniot/framework/mysql"
)

type Physicfeature struct {
	Id	int
	Physictype	string	//名称
	Unit	string	//单位
	Format	int	//格式
	Accur	float64	//精度
	Time	time.Time
}

func ExistPhysicfeature(id int) (bool, error) {
	rows, err := db.Xml.Query("select count(0) Count from physicfeature where id=?", id)
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

func InsertPhysicfeature(physicfeature Physicfeature) (int64, error) {
	result, err := db.Xml.Exec("insert into physicfeature(physictype,unit,format,accur,time) values(?,?,?,?,?)", physicfeature.Physictype,physicfeature.Unit,physicfeature.Format,physicfeature.Accur,physicfeature.Time)
	if err != nil {
		return -1, err
	}
	return result.LastInsertId()
}

func UpdatePhysicfeature(physicfeature Physicfeature) (bool, error) {
	result, err := db.Xml.Exec("update physicfeature set physictype=?, unit=?, format=?, accur=?, time=? where id=?", physicfeature.Physictype, physicfeature.Unit, physicfeature.Format, physicfeature.Accur, physicfeature.Time, physicfeature.Id)
	if err != nil {
		return false, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	return affected > 0, nil
}

func GetPhysicfeature(id int) (physicfeature Physicfeature, err error) {
	rows, err := db.Xml.Query("select id, physictype, unit, format, accur, time from physicfeature where id=?", id)
	if err != nil {
		return physicfeature, err
	}
	if len(rows) <= 0 {
		return physicfeature, nil
	}
	physicfeatures, err := _PhysicfeatureRowsToArray(rows)
	if err != nil {
		return physicfeature, err
	}
	return physicfeatures[0], nil
}

func GetPhysicfeatureRowCount() (count int, err error) {
	rows, err := db.Xml.Query("select count(0) Count from physicfeature")
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

func GetPhysicfeatures() (physicfeatures []Physicfeature, err error) {
	rows, err := db.Xml.Query("select * from Physicfeature ")
	if err != nil {
		return physicfeatures, err
	}
	if len(rows) <= 0 {
		return physicfeatures, nil
	}
	return _PhysicfeatureRowsToArray(rows)
}

func _PhysicfeatureRowsToArray(maps []map[string][]byte) ([]Physicfeature, error) {
	models := make([]Physicfeature, len(maps))
	var err error
	for index, obj := range maps {
		model := Physicfeature{}
		model.Id, err = strconv.Atoi(string(obj["id"]))
		if err != nil {
			return nil, errors.New("parse Id error: " + err.Error())
		}
		model.Physictype = string(obj["physictype"])
		model.Unit = string(obj["unit"])
		model.Format, err = strconv.Atoi(string(obj["format"]))
		if err != nil {
			return nil, errors.New("parse Format error: " + err.Error())
		}
		model.Accur, err = strconv.ParseFloat(string(obj["accur"]), 64)
		if err != nil {
			return nil, errors.New("parse Accur error: " + err.Error())
		}
		timeLayout := "2006-01-02 15:04:05"
		model.Time, err = time.ParseInLocation(timeLayout, string(obj["time"]), time.Local)
		if err != nil {
			return nil, errors.New("parse Time error: " + err.Error())
		}
		models[index] = model
	}
	return models, err
}
