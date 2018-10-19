package xml

import (
	"errors"
	"gitee.com/ha666/golibs"
	"strconv"
	"time"
	db "yuniot/framework/mysql"
)

type Influx struct {
	Id	int
	Productionlineid	string
	Addr	string
	User	string
	Pwd	string
	Remark	string
	Time	time.Time
}

func ExistInflux(id int) (bool, error) {
	rows, err := db.Xml.Query("select count(0) Count from influx where id=?", id)
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

func InsertInflux(influx Influx) (int64, error) {
	result, err := db.Xml.Exec("insert into influx(productionlineid,addr,user,pwd,remark,time) values(?,?,?,?,?,?)", influx.Productionlineid,influx.Addr,influx.User,influx.Pwd,influx.Remark,influx.Time)
	if err != nil {
		return -1, err
	}
	return result.LastInsertId()
}

func UpdateInflux(influx Influx) (bool, error) {
	result, err := db.Xml.Exec("update influx set productionlineid=?, addr=?, user=?, pwd=?, remark=?, time=? where id=?", influx.Productionlineid, influx.Addr, influx.User, influx.Pwd, influx.Remark, influx.Time, influx.Id)
	if err != nil {
		return false, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	return affected > 0, nil
}

func GetInflux(id int) (influx Influx, err error) {
	rows, err := db.Xml.Query("select id, productionlineid, addr, user, pwd, remark, time from influx where id=?", id)
	if err != nil {
		return influx, err
	}
	if len(rows) <= 0 {
		return influx, nil
	}
	influxs, err := _InfluxRowsToArray(rows)
	if err != nil {
		return influx, err
	}
	return influxs[0], nil
}

func GetInfluxRowCount() (count int, err error) {
	rows, err := db.Xml.Query("select count(0) Count from influx")
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

func _InfluxRowsToArray(maps []map[string][]byte) ([]Influx, error) {
	models := make([]Influx, len(maps))
	var err error
	for index, obj := range maps {
		model := Influx{}
		model.Id, err = strconv.Atoi(string(obj["id"]))
		if err != nil {
			return nil, errors.New("parse Id error: " + err.Error())
		}
		model.Productionlineid = string(obj["productionlineid"])
		model.Addr = string(obj["addr"])
		model.User = string(obj["user"])
		model.Pwd = string(obj["pwd"])
		model.Remark = string(obj["remark"])
		model.Time, err = time.ParseInLocation(golibs.Time_TIMEMYSQL, string(obj["time"]), time.Local)
		if err != nil {
			return nil, errors.New("parse Time error: " + err.Error())
		}
		models[index] = model
	}
	return models, err
}
