package xml

import (
	"errors"
	"gitee.com/ha666/golibs"
	"strconv"
	"time"
	db "yuniot/framework/mysql"
	"strings"
)

type Sim struct {
	Id	int
	Simid	string
	Name	string
	Productionlineid	string
	IsOn	bool	//是否在用
	Time	time.Time
}

func ExistSim(id int) (bool, error) {
	rows, err := db.Xml.Query("select count(0) Count from sim where id=?", id)
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

func InsertSim(sim Sim) (int64, error) {
	result, err := db.Xml.Exec("insert into sim(simid,name,productionlineid,IsOn,time) values(?,?,?,?,?)", sim.Simid,sim.Name,sim.Productionlineid,sim.IsOn,sim.Time)
	if err != nil {
		return -1, err
	}
	return result.LastInsertId()
}

func UpdateSim(sim Sim) (bool, error) {
	result, err := db.Xml.Exec("update sim set simid=?, name=?, productionlineid=?, IsOn=?, time=? where id=?", sim.Simid, sim.Name, sim.Productionlineid, sim.IsOn, sim.Time, sim.Id)
	if err != nil {
		return false, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	return affected > 0, nil
}

func GetSim(id int) (sim Sim, err error) {
	rows, err := db.Xml.Query("select id, simid, name, productionlineid, IsOn, time from sim where id=?", id)
	if err != nil {
		return sim, err
	}
	if len(rows) <= 0 {
		return sim, nil
	}
	sims, err := _SimRowsToArray(rows)
	if err != nil {
		return sim, err
	}
	return sims[0], nil
}

func GetSimRowCount() (count int, err error) {
	rows, err := db.Xml.Query("select count(0) Count from sim")
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

func _SimRowsToArray(maps []map[string][]byte) ([]Sim, error) {
	models := make([]Sim, len(maps))
	var err error
	for index, obj := range maps {
		model := Sim{}
		model.Id, err = strconv.Atoi(string(obj["id"]))
		if err != nil {
			return nil, errors.New("parse Id error: " + err.Error())
		}
		model.Simid = string(obj["simid"])
		model.Name = string(obj["name"])
		model.Productionlineid = string(obj["productionlineid"])
		if strings.ToLower(string(obj["IsOn"])) == "true" || obj["IsOn"][0] == 1 {
			model.IsOn = true
		} else {
			model.IsOn = false
		}
		model.Time, err = time.ParseInLocation(golibs.Time_TIMEMYSQL, string(obj["time"]), time.Local)
		if err != nil {
			return nil, errors.New("parse Time error: " + err.Error())
		}
		models[index] = model
	}
	return models, err
}
