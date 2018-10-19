package xml

import (
	"errors"
	"strconv"
	db "yuniot/framework/mysql"
)

type Datatype struct {
	Id	int
	Description	string
	Bit	int
	Inbyte	int
	Outintarray	int
}

func ExistDatatype(id int) (bool, error) {
	rows, err := db.Xml.Query("select count(0) Count from datatype where id=?", id)
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

func InsertDatatype(datatype Datatype) (bool, error) {
	result, err := db.Xml.Exec("insert into datatype(id,description,bit,inbyte,outintarray) values(?,?,?,?,?)", datatype.Id,datatype.Description,datatype.Bit,datatype.Inbyte,datatype.Outintarray)
	if err != nil {
		return false, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	return affected > 0, nil
}

func UpdateDatatype(datatype Datatype) (bool, error) {
	result, err := db.Xml.Exec("update datatype set description=?, bit=?, inbyte=?, outintarray=? where id=?", datatype.Description, datatype.Bit, datatype.Inbyte, datatype.Outintarray, datatype.Id)
	if err != nil {
		return false, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	return affected > 0, nil
}

func InsertUpdateDatatype(datatype Datatype) (bool, error) {
	result, err := db.Xml.Exec("insert into datatype(id,description,bit,inbyte,outintarray) values(?,?,?,?,?) ON DUPLICATE KEY UPDATE description=?,bit=?,inbyte=?,outintarray=?", datatype.Id,datatype.Description,datatype.Bit,datatype.Inbyte,datatype.Outintarray,datatype.Description,datatype.Bit,datatype.Inbyte,datatype.Outintarray)
	if err != nil {
		return false, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	return affected > 0, nil
}

func GetDatatype(id int) (datatype Datatype, err error) {
	rows, err := db.Xml.Query("select id, description, bit, inbyte, outintarray from datatype where id=?", id)
	if err != nil {
		return datatype, err
	}
	if len(rows) <= 0 {
		return datatype, nil
	}
	datatypes, err := _DatatypeRowsToArray(rows)
	if err != nil {
		return datatype, err
	}
	return datatypes[0], nil
}

func GetDatatypeRowCount() (count int, err error) {
	rows, err := db.Xml.Query("select count(0) Count from datatype")
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

func GetDatatypes() (datatypes []Datatype, err error) {
	rows, err := db.Xml.Query("select * from datatype ")
	if err != nil {
		return datatypes, err
	}
	if len(rows) <= 0 {
		return datatypes, nil
	}
	return _DatatypeRowsToArray(rows)
}

func _DatatypeRowsToArray(maps []map[string][]byte) ([]Datatype, error) {
	models := make([]Datatype, len(maps))
	var err error
	for index, obj := range maps {
		model := Datatype{}
		model.Id, err = strconv.Atoi(string(obj["id"]))
		if err != nil {
			return nil, errors.New("parse Id error: " + err.Error())
		}
		model.Description = string(obj["description"])
		model.Bit, err = strconv.Atoi(string(obj["bit"]))
		if err != nil {
			return nil, errors.New("parse Bit error: " + err.Error())
		}
		model.Inbyte, err = strconv.Atoi(string(obj["inbyte"]))
		if err != nil {
			return nil, errors.New("parse Inbyte error: " + err.Error())
		}
		model.Outintarray, err = strconv.Atoi(string(obj["outintarray"]))
		if err != nil {
			return nil, errors.New("parse Outintarray error: " + err.Error())
		}
		models[index] = model
	}
	return models, err
}
