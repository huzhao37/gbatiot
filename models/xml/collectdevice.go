package xml

import (
	"errors"
	"gitee.com/ha666/golibs"
	"strconv"
	"time"
	db "yuniot/framework/mysql"
)

type Collectdevice struct {
	Id	int
	Index	string
	ProductionlineId	string
	Time	time.Time
	Remark	string
}

func ExistCollectdevice(id int) (bool, error) {
	rows, err := db.Xml.Query("select count(0) Count from collectdevice where id=?", id)
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

func InsertCollectdevice(collectdevice Collectdevice) (int64, error) {
	result, err := db.Xml.Exec("insert into collectdevice(`index`,productionline_id,time,remark) values(?,?,?,?)", collectdevice.Index,collectdevice.ProductionlineId,collectdevice.Time,collectdevice.Remark)
	if err != nil {
		return -1, err
	}
	return result.LastInsertId()
}

func UpdateCollectdevice(collectdevice Collectdevice) (bool, error) {
	result, err := db.Xml.Exec("update collectdevice set `index`=?, productionline_id=?, time=?, remark=? where id=?", collectdevice.Index, collectdevice.ProductionlineId, collectdevice.Time, collectdevice.Remark, collectdevice.Id)
	if err != nil {
		return false, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	return affected > 0, nil
}

func GetCollectdevice(id int) (collectdevice Collectdevice, err error) {
	rows, err := db.Xml.Query("select id, `index`, productionline_id, time, remark from collectdevice where id=?", id)
	if err != nil {
		return collectdevice, err
	}
	if len(rows) <= 0 {
		return collectdevice, nil
	}
	collectdevices, err := _CollectdeviceRowsToArray(rows)
	if err != nil {
		return collectdevice, err
	}
	return collectdevices[0], nil
}
func GetCollectdeviceByProductionlineid(productionlineid string) (collectdevices []Collectdevice, err error) {
	rows, err := db.Xml.Query("select id, `index`, productionline_id, time, remark from collectdevice where productionlineid=?", productionlineid)
	if err != nil {
		return collectdevices, err
	}
	if len(rows) <= 0 {
		return collectdevices, nil
	}
	collectdevices, err = _CollectdeviceRowsToArray(rows)
	if err != nil {
		return collectdevices, err
	}
	return collectdevices, nil
}
func GetCollectdeviceByIndex(index string) (collectdevice Collectdevice, err error) {
	rows, err := db.Xml.Query("select id, `index`, productionline_id, time, remark from collectdevice where `index`=?", index)
	if err != nil {
		return collectdevice, err
	}
	if len(rows) <= 0 {
		return collectdevice, nil
	}
	collectdevices, err := _CollectdeviceRowsToArray(rows)
	if err != nil {
		return collectdevice, err
	}
	return collectdevices[0], nil
}
func GetCollectdeviceRowCount() (count int, err error) {
	rows, err := db.Xml.Query("select count(0) Count from collectdevice")
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

func _CollectdeviceRowsToArray(maps []map[string][]byte) ([]Collectdevice, error) {
	models := make([]Collectdevice, len(maps))
	var err error
	for index, obj := range maps {
		model := Collectdevice{}
		model.Id, err = strconv.Atoi(string(obj["id"]))
		if err != nil {
			return nil, errors.New("parse Id error: " + err.Error())
		}
		model.Index = string(obj["index"])
		model.ProductionlineId = string(obj["productionline_id"])
		model.Time, err = time.ParseInLocation(golibs.Time_TIMEStandard, string(obj["time"]), time.Local)
		if err != nil {
			return nil, errors.New("parse Time error: " + err.Error())
		}
		model.Remark = string(obj["remark"])
		models[index] = model
	}
	return models, err
}
