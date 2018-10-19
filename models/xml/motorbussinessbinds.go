package xml

import (
	"errors"
	"gitee.com/ha666/golibs"
	"strconv"
	"time"
	db "yuniot/framework/mysql"
)

type Motorbussinessbinds struct {
	Id	int
	Motorid	string	//电机ID
	Bussinessid	string	//业务ID
	Param	string	//业务关键字段英文名称
	Remark	string	//备注
	Time	time.Time	//时间
}

func ExistMotorbussinessbinds(id int) (bool, error) {
	rows, err := db.Xml.Query("select count(0) Count from motorbussinessbinds where id=?", id)
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

func InsertMotorbussinessbinds(motorbussinessbinds Motorbussinessbinds) (int64, error) {
	result, err := db.Xml.Exec("insert into motorbussinessbinds(motorid,bussinessid,param,remark,time) values(?,?,?,?,?)", motorbussinessbinds.Motorid,motorbussinessbinds.Bussinessid,motorbussinessbinds.Param,motorbussinessbinds.Remark,motorbussinessbinds.Time)
	if err != nil {
		return -1, err
	}
	return result.LastInsertId()
}

func UpdateMotorbussinessbinds(motorbussinessbinds Motorbussinessbinds) (bool, error) {
	result, err := db.Xml.Exec("update motorbussinessbinds set motorid=?, bussinessid=?, param=?, remark=?, time=? where id=?", motorbussinessbinds.Motorid, motorbussinessbinds.Bussinessid, motorbussinessbinds.Param, motorbussinessbinds.Remark, motorbussinessbinds.Time, motorbussinessbinds.Id)
	if err != nil {
		return false, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	return affected > 0, nil
}

func GetMotorbussinessbinds(id int) (motorbussinessbinds Motorbussinessbinds, err error) {
	rows, err := db.Xml.Query("select id, motorid, bussinessid, param, remark, time from motorbussinessbinds where id=?", id)
	if err != nil {
		return motorbussinessbinds, err
	}
	if len(rows) <= 0 {
		return motorbussinessbinds, nil
	}
	motorbussinessbindss, err := _MotorbussinessbindsRowsToArray(rows)
	if err != nil {
		return motorbussinessbinds, err
	}
	return motorbussinessbindss[0], nil
}

func GetMotorbussinessbindsRowCount() (count int, err error) {
	rows, err := db.Xml.Query("select count(0) Count from motorbussinessbinds")
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
func GetBussinessByMotorid(motorid string) (motorbussinessbinds []Motorbussinessbinds, err error) {
	rows, err := db.Xml.Query("select * from motorbussinessbinds where motorid=?", motorid)
	if err != nil {
		return motorbussinessbinds, err
	}
	if len(rows) <= 0 {
		return motorbussinessbinds, nil
	}
	motorbussinessbinds, err = _MotorbussinessbindsRowsToArray(rows)
	if err != nil {
		return motorbussinessbinds, err
	}
	return motorbussinessbinds, nil
}
func _MotorbussinessbindsRowsToArray(maps []map[string][]byte) ([]Motorbussinessbinds, error) {
	models := make([]Motorbussinessbinds, len(maps))
	var err error
	for index, obj := range maps {
		model := Motorbussinessbinds{}
		model.Id, err = strconv.Atoi(string(obj["id"]))
		if err != nil {
			return nil, errors.New("parse Id error: " + err.Error())
		}
		model.Motorid = string(obj["motorid"])
		model.Bussinessid = string(obj["bussinessid"])
		model.Param = string(obj["param"])
		model.Remark = string(obj["remark"])
		model.Time, err = time.ParseInLocation(golibs.Time_TIMEStandard, string(obj["time"]), time.Local)
		if err != nil {
			return nil, errors.New("parse Time error: " + err.Error())
		}
		models[index] = model
	}
	return models, err
}
