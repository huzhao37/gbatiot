package xml

import (
	"errors"
	"strconv"
	db "yuniot/framework/mysql"
	"gitee.com/ha666/golibs"
)

type Dataconfig struct {
	Id	int
	Datatypeid	int
	Count	int
	Collectdeviceindex	string	//数据表单的Index
}

func ExistDataconfig(id int) (bool, error) {
	rows, err := db.Xml.Query("select count(0) Count from dataconfig where id=?", id)
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

func InsertDataconfig(dataconfig Dataconfig) (int64, error) {
	result, err := db.Xml.Exec("insert into dataconfig(datatypeid,count,collectdeviceindex) values(?,?,?)", dataconfig.Datatypeid,dataconfig.Count,dataconfig.Collectdeviceindex)
	if err != nil {
		return -1, err
	}
	return result.LastInsertId()
}

func UpdateDataconfig(dataconfig Dataconfig) (bool, error) {
	result, err := db.Xml.Exec("update dataconfig set datatypeid=?, count=?, collectdeviceindex=? where id=?", dataconfig.Datatypeid, dataconfig.Count, dataconfig.Collectdeviceindex, dataconfig.Id)
	if err != nil {
		return false, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	return affected > 0, nil
}

func GetDataconfig(id int) (dataconfig Dataconfig, err error) {
	rows, err := db.Xml.Query("select id, datatypeid, count, collectdeviceindex from dataconfig where id=?", id)
	if err != nil {
		return dataconfig, err
	}
	if len(rows) <= 0 {
		return dataconfig, nil
	}
	dataconfigs, err := _DataconfigRowsToArray(rows)
	if err != nil {
		return dataconfig, err
	}
	return dataconfigs[0], nil
}

// IX_DataConfig_DataTypeId
func GetDataconfigBydatatypeid(datatypeid int) (dataconfigs []Dataconfig, err error) {
	rows, err := db.Xml.Query("select id, datatypeid, count, collectdeviceindex from dataconfig where datatypeid=?",datatypeid)
	if err != nil {
		return dataconfigs, err
	}
	if len(rows) <= 0 {
		return dataconfigs, nil
	}
	return _DataconfigRowsToArray(rows)
}

func GetDataconfigRowCount() (count int, err error) {
	rows, err := db.Xml.Query("select count(0) Count from dataconfig")
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

// IX_DataConfig_DataTypeId
func GetDataconfigRowListBydatatypeid(datatypeidIsDesc bool, PageIndex, PageSize int) (dataconfigs []Dataconfig, err error) {
	sqlstr := golibs.NewStringBuilder()
	sqlstr.Append("select id, datatypeid, count, collectdeviceindex from dataconfig order by")
	sqlstr.Append(" datatypeid")
	if datatypeidIsDesc {
		sqlstr.Append(" Desc")
	}
	sqlstr.Append(" limit ?,?")
	rows, err := db.Xml.Query(sqlstr.ToString(), (PageIndex-1)*PageSize, PageSize)
	sqlstr.Append(",")
	if err != nil {
		return dataconfigs, err
	}
	if len(rows) <= 0 {
		return dataconfigs, nil
	}
	return _DataconfigRowsToArray(rows)
}
func GetDataconfigByCollectdeviceindex(collectdeviceindex string) (dataconfigs []Dataconfig, err error) {
	rows, err := db.Xml.Query("select id, datatypeid, count, collectdeviceindex from dataconfig where Collectdeviceindex=?",collectdeviceindex)
	if err != nil {
		return dataconfigs, err
	}
	if len(rows) <= 0 {
		return dataconfigs, nil
	}
	return _DataconfigRowsToArray(rows)
}
func _DataconfigRowsToArray(maps []map[string][]byte) ([]Dataconfig, error) {
	models := make([]Dataconfig, len(maps))
	var err error
	for index, obj := range maps {
		model := Dataconfig{}
		model.Id, err = strconv.Atoi(string(obj["id"]))
		if err != nil {
			return nil, errors.New("parse Id error: " + err.Error())
		}
		model.Datatypeid, err = strconv.Atoi(string(obj["datatypeid"]))
		if err != nil {
			return nil, errors.New("parse Datatypeid error: " + err.Error())
		}
		model.Count, err = strconv.Atoi(string(obj["count"]))
		if err != nil {
			return nil, errors.New("parse Count error: " + err.Error())
		}
		model.Collectdeviceindex = string(obj["collectdeviceindex"])
		models[index] = model
	}
	return models, err
}
