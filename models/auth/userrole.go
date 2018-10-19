package db

import (
	"errors"
	"gitee.com/ha666/golibs"
	"strconv"
	"time"
	db "yuniot/framework/mysql"
)

type Userrole struct {
	Id	int
	Desc	string
	Remark	string
	Time	time.Time
}

func ExistUserrole(id int) (bool, error) {
	rows, err := db.Auth.Query("select count(0) Count from userrole where id=?", id)
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

func InsertUserrole(userrole Userrole) (int64, error) {
	result, err := db.Auth.Exec("insert into userrole(desc,remark,time) values(?,?,?)", userrole.Desc,userrole.Remark,userrole.Time)
	if err != nil {
		return -1, err
	}
	return result.LastInsertId()
}

func UpdateUserrole(userrole Userrole) (bool, error) {
	result, err := db.Auth.Exec("update userrole set desc=?, remark=?, time=? where id=?", userrole.Desc, userrole.Remark, userrole.Time, userrole.Id)
	if err != nil {
		return false, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	return affected > 0, nil
}

func GetUserrole(id int) (userrole Userrole, err error) {
	rows, err := db.Auth.Query("select id, desc, remark, time from userrole where id=?", id)
	if err != nil {
		return userrole, err
	}
	if len(rows) <= 0 {
		return userrole, nil
	}
	userroles, err := _UserroleRowsToArray(rows)
	if err != nil {
		return userrole, err
	}
	return userroles[0], nil
}

func GetUserroleRowCount() (count int, err error) {
	rows, err := db.Auth.Query("select count(0) Count from userrole")
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

func _UserroleRowsToArray(maps []map[string][]byte) ([]Userrole, error) {
	models := make([]Userrole, len(maps))
	var err error
	for index, obj := range maps {
		model := Userrole{}
		model.Id, err = strconv.Atoi(string(obj["id"]))
		if err != nil {
			return nil, errors.New("parse Id error: " + err.Error())
		}
		model.Desc = string(obj["desc"])
		model.Remark = string(obj["remark"])
		model.Time, err = time.ParseInLocation(golibs.Time_TIMEStandard, string(obj["time"]), time.Local)
		if err != nil {
			return nil, errors.New("parse Time error: " + err.Error())
		}
		models[index] = model
	}
	return models, err
}
