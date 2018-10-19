package db

import (
	"errors"
	"gitee.com/ha666/golibs"
	"strconv"
	"time"
	db "yuniot/framework/mysql"
)

type User struct {
	Id	int
	Userid	int
	Account	string
	Pwd	string
	Mail	string
	Mobileno	string
	Remark	string
	Time	time.Time
	Nick	string
	Userroleid	int
}

func ExistUser(id int) (bool, error) {
	rows, err := db.Auth.Query("select count(0) Count from user where id=?", id)
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

func InsertUser(user User) (int64, error) {
	result, err := db.Auth.Exec("insert into user(userid,account,pwd,mail,mobileno,remark,time,nick,userroleid) values(?,?,?,?,?,?,?,?,?)", user.Userid,user.Account,user.Pwd,user.Mail,user.Mobileno,user.Remark,user.Time,user.Nick,user.Userroleid)
	if err != nil {
		return -1, err
	}
	return result.LastInsertId()
}

func UpdateUser(user User) (bool, error) {
	result, err := db.Auth.Exec("update user set userid=?, account=?, pwd=?, mail=?, mobileno=?, remark=?, time=?, nick=?, userroleid=? where id=?", user.Userid, user.Account, user.Pwd, user.Mail, user.Mobileno, user.Remark, user.Time, user.Nick, user.Userroleid, user.Id)
	if err != nil {
		return false, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	return affected > 0, nil
}

func GetUser(id int) (user User, err error) {
	rows, err := db.Auth.Query("select id, userid, account, pwd, mail, mobileno, remark, time, nick, userroleid from user where id=?", id)
	if err != nil {
		return user, err
	}
	if len(rows) <= 0 {
		return user, nil
	}
	users, err := _UserRowsToArray(rows)
	if err != nil {
		return user, err
	}
	return users[0], nil
}

func GetUserRowCount() (count int, err error) {
	rows, err := db.Auth.Query("select count(0) Count from user")
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
func UserLogin(name string,pwd string) (user User, err error ) {
	rows, err := db.Auth.Query("select * from user where account=? and pwd=?", name,pwd)
	if err != nil {
		return user, err
	}
	if len(rows) <= 0 {
		return user, nil
	}
	users, err := _UserRowsToArray(rows)
	if err != nil {
		return user, err
	}
	return users[0], nil
}
func GetUsers() (users []User, err error) {
	rows, err := db.Auth.Query("select * from user ")
	if err != nil {
		return users, err
	}
	if len(rows) <= 0 {
		return users, nil
	}
	return _UserRowsToArray(rows)
}
func _UserRowsToArray(maps []map[string][]byte) ([]User, error) {
	models := make([]User, len(maps))
	var err error
	for index, obj := range maps {
		model := User{}
		model.Id, err = strconv.Atoi(string(obj["id"]))
		if err != nil {
			return nil, errors.New("parse Id error: " + err.Error())
		}
		model.Userid, err = strconv.Atoi(string(obj["userid"]))
		if err != nil {
			return nil, errors.New("parse userid error: " + err.Error())
		}
		//model.Userid = string(obj["userid"])
		model.Account = string(obj["account"])
		model.Pwd = string(obj["pwd"])
		model.Mail = string(obj["mail"])
		model.Mobileno = string(obj["mobileno"])
		model.Remark = string(obj["remark"])
		model.Time, err = time.ParseInLocation(golibs.Time_TIMEStandard, string(obj["time"]), time.Local)
		if err != nil {
			return nil, errors.New("parse Time error: " + err.Error())
		}
		model.Nick = string(obj["nick"])
		model.Userroleid, err = strconv.Atoi(string(obj["userroleid"]))
		if err != nil {
			return nil, errors.New("parse Userroleid error: " + err.Error())
		}
		models[index] = model
	}
	return models, err
}
