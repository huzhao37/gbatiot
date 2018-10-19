package xml

import (
	"errors"
	"gitee.com/ha666/golibs"
	"strconv"
	"time"
	db "yuniot/framework/mysql"
	"strings"
)

type Simcarbinds struct {
	Id	int
	Simid	string	//sim卡
	Carno	string	//车牌号
	Isbinding	bool	//是否处于绑定状态
	Time	time.Time	//时间
}

func ExistSimcarbinds(id int) (bool, error) {
	rows, err := db.Xml.Query("select count(0) Count from simcarbinds where id=?", id)
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

func InsertSimcarbinds(simcarbinds Simcarbinds) (int64, error) {
	result, err := db.Xml.Exec("insert into simcarbinds(simid,carno,isbinding,time) values(?,?,?,?)", simcarbinds.Simid,simcarbinds.Carno,simcarbinds.Isbinding,simcarbinds.Time)
	if err != nil {
		return -1, err
	}
	return result.LastInsertId()
}

func UpdateSimcarbinds(simcarbinds Simcarbinds) (bool, error) {
	result, err := db.Xml.Exec("update simcarbinds set simid=?, carno=?, isbinding=?, time=? where id=?", simcarbinds.Simid, simcarbinds.Carno, simcarbinds.Isbinding, simcarbinds.Time, simcarbinds.Id)
	if err != nil {
		return false, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	return affected > 0, nil
}

func GetSimcarbinds(id int) (simcarbinds Simcarbinds, err error) {
	rows, err := db.Xml.Query("select id, simid, carno, isbinding, time from simcarbinds where id=?", id)
	if err != nil {
		return simcarbinds, err
	}
	if len(rows) <= 0 {
		return simcarbinds, nil
	}
	simcarbindss, err := _SimcarbindsRowsToArray(rows)
	if err != nil {
		return simcarbinds, err
	}
	return simcarbindss[0], nil
}

func GetSimcarbindsRowCount() (count int, err error) {
	rows, err := db.Xml.Query("select count(0) Count from simcarbinds")
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

func _SimcarbindsRowsToArray(maps []map[string][]byte) ([]Simcarbinds, error) {
	models := make([]Simcarbinds, len(maps))
	var err error
	for index, obj := range maps {
		model := Simcarbinds{}
		model.Id, err = strconv.Atoi(string(obj["id"]))
		if err != nil {
			return nil, errors.New("parse Id error: " + err.Error())
		}
		model.Simid = string(obj["simid"])
		model.Carno = string(obj["carno"])
		if strings.ToLower(string(obj["Isbinding"])) == "true" || obj["Isbinding"][0] == 1 {
			model.Isbinding = true
		} else {
			model.Isbinding = false
		}
		model.Time, err = time.ParseInLocation(golibs.Time_TIMEMYSQL, string(obj["time"]), time.Local)
		if err != nil {
			return nil, errors.New("parse Time error: " + err.Error())
		}
		models[index] = model
	}
	return models, err
}
