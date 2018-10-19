package xml

import (
	"errors"
	"strconv"
	"time"
	db "yuniot/framework/mysql"
)

type Originalbytes struct {
	Id	int64
	Bytes	string
	Collectdeviceindex	string
	Productionlineid	string
	Time	time.Time
}

func ExistOriginalbytes(id int64) (bool, error) {
	rows, err := db.Xml.Query("select count(0) Count from originalbytes where id=?", id)
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

func InsertOriginalbytes(originalbytes Originalbytes) (int64, error) {
	result, err := db.Xml.Exec("insert into originalbytes(bytes,Collectdeviceindex,productionlineid,time) values(?,?,?,?)", originalbytes.Bytes,originalbytes.Collectdeviceindex,originalbytes.Productionlineid,originalbytes.Time)
	if err != nil {
		return -1, err
	}
	return result.LastInsertId()
}

func UpdateOriginalbytes(originalbytes Originalbytes) (bool, error) {
	result, err := db.Xml.Exec("update originalbytes set bytes=?, Collectdeviceindex=?, productionlineid=?, time=? where id=?", originalbytes.Bytes, originalbytes.Collectdeviceindex, originalbytes.Productionlineid, originalbytes.Time, originalbytes.Id)
	if err != nil {
		return false, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	return affected > 0, nil
}

func GetOriginalbytes(id int64) (originalbytes Originalbytes, err error) {
	rows, err := db.Xml.Query("select id, bytes, Collectdeviceindex, productionlineid, time from originalbytes where id=?", id)
	if err != nil {
		return originalbytes, err
	}
	if len(rows) <= 0 {
		return originalbytes, nil
	}
	originalbytess, err := _OriginalbytesRowsToArray(rows)
	if err != nil {
		return originalbytes, err
	}
	return originalbytess[0], nil
}

func GetOriginalbytess() (originalbytes []Originalbytes, err error) {
	rows, err := db.Xml.Query("select * from Originalbytes ")
	if err != nil {
		return originalbytes, err
	}
	if len(rows) <= 0 {
		return originalbytes, nil
	}
	return _OriginalbytesRowsToArray(rows)
}
func GetOriginalbytesRowCount() (count int, err error) {
	rows, err := db.Xml.Query("select count(0) Count from originalbytes")
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

func _OriginalbytesRowsToArray(maps []map[string][]byte) ([]Originalbytes, error) {
	models := make([]Originalbytes, len(maps))
	var err error
	for index, obj := range maps {
		model := Originalbytes{}
		model.Id, err = strconv.ParseInt(string(obj["id"]), 10, 64)
		if err != nil {
			return nil, errors.New("parse Id error: " + err.Error())
		}
		model.Bytes = string(obj["bytes"])
		model.Collectdeviceindex = string(obj["Collectdeviceindex"])
		//if err != nil {
		//	return nil, errors.New("parse Collectdeviceindex error: " + err.Error())
		//}
		model.Productionlineid = string(obj["productionlineid"])
		timeLayout := "2006-01-02 15:04:05"
		model.Time, err = time.ParseInLocation(timeLayout, string(obj["time"]), time.Local)
		if err != nil {
			return nil, errors.New("parse Time error: " + err.Error())
		}
		models[index] = model
	}
	return models, err
}
