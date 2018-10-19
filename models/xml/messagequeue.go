package xml

import (
	"errors"
	"strconv"
	"time"
	db "yuniot/framework/mysql"
)

type Messagequeue struct {
	Id	int
	Host	string	//主机地址
	Port	int	//端口号
	Routekey	string	//关键字
	Name	string	//名称
	Timer	int	//采集频率
	Collectdeviceindex	string	//采集设备ID
	Writeread	int	//读写功能
	Comtype	string	//通讯类型
	Time	time.Time	//时间
	Remark	string	//备注
	Username	string	//用户名
	Pwd	string	//密码
}

func ExistMessagequeue(id int) (bool, error) {
	rows, err := db.Xml.Query("select count(0) Count from messagequeue where id=?", id)
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

func InsertMessagequeue(messagequeue Messagequeue) (int64, error) {
	result, err := db.Xml.Exec("insert into messagequeue(host,port,routekey,name,timer,Collectdeviceindex,writeread,comtype,time,remark,username,pwd) values(?,?,?,?,?,?,?,?,?,?,?,?)", messagequeue.Host,messagequeue.Port,messagequeue.Routekey,messagequeue.Name,messagequeue.Timer,messagequeue.Collectdeviceindex,messagequeue.Writeread,messagequeue.Comtype,messagequeue.Time,messagequeue.Remark,messagequeue.Username,messagequeue.Pwd)
	if err != nil {
		return -1, err
	}
	return result.LastInsertId()
}

func UpdateMessagequeue(messagequeue Messagequeue) (bool, error) {
	result, err := db.Xml.Exec("update messagequeue set host=?, port=?, routekey=?, name=?, timer=?, Collectdeviceindex=?, writeread=?, comtype=?, time=?, remark=?, username=?, pwd=? where id=?", messagequeue.Host, messagequeue.Port, messagequeue.Routekey, messagequeue.Name, messagequeue.Timer, messagequeue.Collectdeviceindex, messagequeue.Writeread, messagequeue.Comtype, messagequeue.Time, messagequeue.Remark, messagequeue.Username, messagequeue.Pwd, messagequeue.Id)
	if err != nil {
		return false, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	return affected > 0, nil
}

func GetMessagequeue(id int) (messagequeue Messagequeue, err error) {
	rows, err := db.Xml.Query("select id, host, port, routekey, name, timer, Collectdeviceindex, writeread, comtype, time, remark, username, pwd from messagequeue where id=?", id)
	if err != nil {
		return messagequeue, err
	}
	if len(rows) <= 0 {
		return messagequeue, nil
	}
	messagequeues, err := _MessagequeueRowsToArray(rows)
	if err != nil {
		return messagequeue, err
	}
	return messagequeues[0], nil
}

func GetMessagequeueByMasterstrid(Collectdeviceindex string) (messagequeue Messagequeue, err error) {
	rows, err := db.Xml.Query("select id, host, port, routekey, name, timer, Collectdeviceindex, writeread, comtype, time, remark, username, pwd from messagequeue where Collectdeviceindex=?", Collectdeviceindex)
	if err != nil {
		return messagequeue, err
	}
	if len(rows) <= 0 {
		return messagequeue, nil
	}
	messagequeues, err := _MessagequeueRowsToArray(rows)
	if err != nil {
		return messagequeue, err
	}
	return messagequeues[0], nil
}
func GetMessagequeueRowCount() (count int, err error) {
	rows, err := db.Xml.Query("select count(0) Count from messagequeue")
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
func GetMessagequeues() (messagequeues []Messagequeue, err error) {
	rows, err := db.Xml.Query("select * from messagequeue ")
	if err != nil {
		return messagequeues, err
	}
	if len(rows) <= 0 {
		return messagequeues, nil
	}
	return _MessagequeueRowsToArray(rows)
}
func _MessagequeueRowsToArray(maps []map[string][]byte) ([]Messagequeue, error) {
	models := make([]Messagequeue, len(maps))
	var err error
	for index, obj := range maps {
		model := Messagequeue{}
		model.Id, err = strconv.Atoi(string(obj["id"]))
		if err != nil {
			return nil, errors.New("parse Id error: " + err.Error())
		}
		model.Host = string(obj["host"])
		model.Port, err = strconv.Atoi(string(obj["port"]))
		if err != nil {
			return nil, errors.New("parse Port error: " + err.Error())
		}
		model.Routekey = string(obj["routekey"])
		model.Name = string(obj["name"])
		model.Timer, err = strconv.Atoi(string(obj["timer"]))
		if err != nil {
			return nil, errors.New("parse Timer error: " + err.Error())
		}
		model.Collectdeviceindex = string(obj["Collectdeviceindex"])
		model.Writeread, err = strconv.Atoi(string(obj["writeread"]))
		if err != nil {
			return nil, errors.New("parse Writeread error: " + err.Error())
		}
		model.Comtype = string(obj["comtype"])
		timeLayout := "2006-01-02 15:04:05"
		model.Time, err = time.ParseInLocation(timeLayout, string(obj["time"]), time.Local)
		if err != nil {
			return nil, errors.New("parse Time error: " + err.Error())
		}
		model.Remark = string(obj["remark"])
		model.Username = string(obj["username"])
		model.Pwd = string(obj["pwd"])
		models[index] = model
	}
	return models, err
}
