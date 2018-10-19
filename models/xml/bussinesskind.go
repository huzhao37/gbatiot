package xml

import (
	"errors"
	"gitee.com/ha666/golibs"
	"strconv"
	"time"
	db "yuniot/framework/mysql"
)

type Bussinesskind struct {
	Id	int
	Bussinessid	string	//业务ID
	Name	string	//业务名称
	Calcfomula	string	//计算公式
	Defaultparam	string	//默认参数英文名称
	Remark	string	//备注
	Time	time.Time	//时间
}

func ExistBussinesskind(id int) (bool, error) {
	rows, err := db.Xml.Query("select count(0) Count from bussinesskind where id=?", id)
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

func InsertBussinesskind(bussinesskind Bussinesskind) (int64, error) {
	result, err := db.Xml.Exec("insert into bussinesskind(bussinessid,name,calcfomula,defaultparam,remark,time) values(?,?,?,?,?,?)", bussinesskind.Bussinessid,bussinesskind.Name,bussinesskind.Calcfomula,bussinesskind.Defaultparam,bussinesskind.Remark,bussinesskind.Time)
	if err != nil {
		return -1, err
	}
	return result.LastInsertId()
}

func UpdateBussinesskind(bussinesskind Bussinesskind) (bool, error) {
	result, err := db.Xml.Exec("update bussinesskind set bussinessid=?, name=?, calcfomula=?, defaultparam=?, remark=?, time=? where id=?", bussinesskind.Bussinessid, bussinesskind.Name, bussinesskind.Calcfomula, bussinesskind.Defaultparam, bussinesskind.Remark, bussinesskind.Time, bussinesskind.Id)
	if err != nil {
		return false, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	return affected > 0, nil
}

func GetBussinesskind(id int) (bussinesskind Bussinesskind, err error) {
	rows, err := db.Xml.Query("select id, bussinessid, name, calcfomula, defaultparam, remark, time from bussinesskind where id=?", id)
	if err != nil {
		return bussinesskind, err
	}
	if len(rows) <= 0 {
		return bussinesskind, nil
	}
	bussinesskinds, err := _BussinesskindRowsToArray(rows)
	if err != nil {
		return bussinesskind, err
	}
	return bussinesskinds[0], nil
}

func GetBussinesskindRowCount() (count int, err error) {
	rows, err := db.Xml.Query("select count(0) Count from bussinesskind")
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

func _BussinesskindRowsToArray(maps []map[string][]byte) ([]Bussinesskind, error) {
	models := make([]Bussinesskind, len(maps))
	var err error
	for index, obj := range maps {
		model := Bussinesskind{}
		model.Id, err = strconv.Atoi(string(obj["id"]))
		if err != nil {
			return nil, errors.New("parse Id error: " + err.Error())
		}
		model.Bussinessid = string(obj["bussinessid"])
		model.Name = string(obj["name"])
		model.Calcfomula = string(obj["calcfomula"])
		model.Defaultparam = string(obj["defaultparam"])
		model.Remark = string(obj["remark"])
		model.Time, err = time.ParseInLocation(golibs.Time_TIMEMYSQL, string(obj["time"]), time.Local)
		if err != nil {
			return nil, errors.New("parse Time error: " + err.Error())
		}
		models[index] = model
	}
	return models, err
}
