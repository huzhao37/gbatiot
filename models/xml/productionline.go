package xml

import (
	"errors"
	"strconv"
	db "yuniot/framework/mysql"
)

type Productionline struct {
	Id	int64
	Name	string
	ProductionLineId	string
	Time	int64
}

func ExistProductionline(Id int64) (bool, error) {
	rows, err := db.Xml.Query("select count(0) Count from productionline where Id=?", Id)
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

func InsertProductionline(productionline Productionline) (int64, error) {
	result, err := db.Xml.Exec("insert into productionline(Name,ProductionLineId,Time) values(?,?,?)", productionline.Name,productionline.ProductionLineId,productionline.Time)
	if err != nil {
		return -1, err
	}
	return result.LastInsertId()
}

func UpdateProductionline(productionline Productionline) (bool, error) {
	result, err := db.Xml.Exec("update productionline set Name=?, ProductionLineId=?, Time=? where Id=?", productionline.Name, productionline.ProductionLineId, productionline.Time, productionline.Id)
	if err != nil {
		return false, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	return affected > 0, nil
}

func GetProductionline(Id int64) (productionline Productionline, err error) {
	rows, err := db.Xml.Query("select Id, Name, ProductionLineId, Time from productionline where Id=?", Id)
	if err != nil {
		return productionline, err
	}
	if len(rows) <= 0 {
		return productionline, nil
	}
	productionlines, err := _ProductionlineRowsToArray(rows)
	if err != nil {
		return productionline, err
	}
	return productionlines[0], nil
}
func GetProductionlineById(productionLineId string) (productionline Productionline, err error) {
	rows, err := db.Xml.Query("select Id, Name, ProductionLineId, Time from productionline where ProductionLineId=?", productionLineId)
	if err != nil {
		return productionline, err
	}
	if len(rows) <= 0 {
		return productionline, nil
	}
	productionlines, err := _ProductionlineRowsToArray(rows)
	if err != nil {
		return productionline, err
	}
	return productionlines[0], nil
}

func GetProductionlineRowCount() (count int, err error) {
	rows, err := db.Xml.Query("select count(0) Count from productionline")
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

func _ProductionlineRowsToArray(maps []map[string][]byte) ([]Productionline, error) {
	models := make([]Productionline, len(maps))
	var err error
	for index, obj := range maps {
		model := Productionline{}
		model.Id, err = strconv.ParseInt(string(obj["Id"]), 10, 64)
		if err != nil {
			return nil, errors.New("parse Id error: " + err.Error())
		}
		model.Name = string(obj["Name"])
		model.ProductionLineId = string(obj["ProductionLineId"])
		model.Time, err = strconv.ParseInt(string(obj["Time"]), 10, 64)
		if err != nil {
			return nil, errors.New("parse Time error: " + err.Error())
		}
		models[index] = model
	}
	return models, err
}
