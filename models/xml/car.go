package xml

import (
	"errors"
	"gitee.com/ha666/golibs"
	"strconv"
	"time"
	db "yuniot/framework/mysql"
)

type Car struct {
	Id	int
	Carno	string	//车牌号
	Name	string	//名称
	Productionlineid	string	//产线
	Driver	string	//司机
	Time	time.Time	//时间
}

func ExistCar(id int) (bool, error) {
	rows, err := db.Xml.Query("select count(0) Count from car where id=?", id)
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

func InsertCar(car Car) (int64, error) {
	result, err := db.Xml.Exec("insert into car(carno,name,productionlineid,driver,time) values(?,?,?,?,?)", car.Carno,car.Name,car.Productionlineid,car.Driver,car.Time)
	if err != nil {
		return -1, err
	}
	return result.LastInsertId()
}

func UpdateCar(car Car) (bool, error) {
	result, err := db.Xml.Exec("update car set carno=?, name=?, productionlineid=?, driver=?, time=? where id=?", car.Carno, car.Name, car.Productionlineid, car.Driver, car.Time, car.Id)
	if err != nil {
		return false, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	return affected > 0, nil
}

func GetCar(id int) (car Car, err error) {
	rows, err := db.Xml.Query("select id, carno, name, productionlineid, driver, time from car where id=?", id)
	if err != nil {
		return car, err
	}
	if len(rows) <= 0 {
		return car, nil
	}
	cars, err := _CarRowsToArray(rows)
	if err != nil {
		return car, err
	}
	return cars[0], nil
}

func GetCarRowCount() (count int, err error) {
	rows, err := db.Xml.Query("select count(0) Count from car")
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

func _CarRowsToArray(maps []map[string][]byte) ([]Car, error) {
	models := make([]Car, len(maps))
	var err error
	for index, obj := range maps {
		model := Car{}
		model.Id, err = strconv.Atoi(string(obj["id"]))
		if err != nil {
			return nil, errors.New("parse Id error: " + err.Error())
		}
		model.Carno = string(obj["carno"])
		model.Name = string(obj["name"])
		model.Productionlineid = string(obj["productionlineid"])
		model.Driver = string(obj["driver"])
		model.Time, err = time.ParseInLocation(golibs.Time_TIMEMYSQL, string(obj["time"]), time.Local)
		if err != nil {
			return nil, errors.New("parse Time error: " + err.Error())
		}
		models[index] = model
	}
	return models, err
}
