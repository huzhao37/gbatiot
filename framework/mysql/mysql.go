package mysql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"github.com/go-xorm/xorm"
	"time"
	"yuniot/core"
)

var (
	SqlDB *sql.DB
	Xml *xorm.Engine
	Auth *xorm.Engine
	myConfig  = new(core.Config)
)
func init() {
	myConfig.InitConfig("./config/config.txt")
	init_Xml()
	init_Auth()
}

func init_Xml() {
	var env = myConfig.Read("global", "env")
	connectionstring := myConfig.Read(env, "xmlconn")
	maxidleconns := 10
	maxopenconns :=100
	db_conn, err := xorm.NewEngine("mysql", connectionstring)
	if err != nil {
		log.Fatalf("【xml.NewEngine】ex:%s\n", err.Error())
		return
	}
	err = db_conn.Ping()
	if err != nil {
		log.Fatalf("【xml.Ping】ex:%s\n", err.Error())
		return
	}
	db_conn.TZLocation = time.Local
	db_conn.SetMaxIdleConns(maxidleconns)
	db_conn.SetMaxOpenConns(maxopenconns)
	Xml = db_conn
}

func init_Auth() {
	var env = myConfig.Read("global", "env")
	connectionstring := myConfig.Read(env, "authconn")
	maxidleconns := 10
	maxopenconns :=100
	db_conn, err := xorm.NewEngine("mysql", connectionstring)
	if err != nil {
		log.Fatalf("【xml.NewEngine】ex:%s\n", err.Error())
		return
	}
	err = db_conn.Ping()
	if err != nil {
		log.Fatalf("【xml.Ping】ex:%s\n", err.Error())
		return
	}
	db_conn.TZLocation = time.Local
	db_conn.SetMaxIdleConns(maxidleconns)
	db_conn.SetMaxOpenConns(maxopenconns)
	Auth = db_conn
}
