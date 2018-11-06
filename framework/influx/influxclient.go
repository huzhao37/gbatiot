package influx

import (
	"github.com/influxdata/influxdb/client/v2"
	"log"
	corex "yuniot/core"
	"fmt"
)
var (
	myConfig = new(corex.Config)
)
const (
	MyDB     = "test1"
	MyMeasurement = "cpu_usage"
)
func init(){
	// 从配置文件获取redis的ip以及db
	myConfig.InitConfig("./config/config.txt")
}
func ConnInflux() client.Client {
	// 从配置文件获取redis的ip以及db
	//myConfig.InitConfig("./config/config.txt")
	var env = myConfig.Read("global", "env")
	var addr = myConfig.Read(env, "influx.addr")
	var user = myConfig.Read(env, "influx.user")
	var password = myConfig.Read(env, "influx.pwd")
	cli, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     addr,
		Username: user,
		Password: password,
	})
	if err != nil {
		log.Fatal(err)
	}
	return cli
}

func ConnInfluxParam(addr string,user string,password string) client.Client {
	cli, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     addr,
		Username: user,
		Password: password,
	})
	if err != nil {
		log.Fatal(err)
	}
	return cli
}

func ConnInfluxUdp(addr string)client.Client{
	// Make client
	config := client.UDPConfig{Addr: addr}
	c, err := client.NewUDPClient(config)
	if err != nil {
		fmt.Println("Error: ", err.Error())
	}
	return c
}
//batch Insert
func WritesPoints(cli client.Client,database string ,points []*client.Point) {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  database,
		Precision: "m",//时间精度分钟//["n", "u", "ms", "s", "m", "h"]
		//RetentionPolicy: "yearpolicy",
	})
	if err != nil {
		log.Fatal(err)
	}
	//批量写入
	bp.AddPoints(points)
	if err := cli.Write(bp); err != nil {
		log.Fatal(err)
	}
}
//query
func QueryDB(cli client.Client,database string, cmd string) (res []client.Result, err error) {
	q := client.Query{
		Command:  cmd,
		Database: database,
	}
	if response, err := cli.Query(q); err == nil {
		if response.Error() != nil {
			return res, response.Error()
		}
		res = response.Results
	} else {
		return res, err
	}
	return res, nil
}

//create database
func CreateDB(cli client.Client,database string) (res []client.Result, err error) {
	q := client.Query{
		Command:  "CREATE DATABASE  "+database,
	}
	if response, err := cli.Query(q); err == nil {
		if response.Error() != nil {
			return res, response.Error()
		}
		res = response.Results
	} else {
		return res, err
	}
	return res, nil
}

//exist database
func ExistDB(cli client.Client,database string) (bool, error) {
	q := client.Query{
		Command:  "Show DATABASES ",
	}
	if response, err := cli.Query(q); err == nil {
		if response.Error() != nil {
			return false, response.Error()
		}
		res := response.Results
		for i:=0;i<len(res[0].Series[0].Values);i++ {
			var row=res[0].Series[0].Values[i][0]
			dbName :=fmt.Sprintf("%s", row.(string))
			if dbName==database{
				return true,err
			}
		}
	} else {
		return false, err
	}
	return false, nil
}
//过期策略
func Repl(cli client.Client,database string,days int) (res []client.Result, err error) {
	q := client.Query{
		Command:  fmt.Sprintf("ALTER RETENTION POLICY %s ON  %s DURATION %dd DEFAULT  ","autogen",database,days),
	}
	if response, err := cli.Query(q); err == nil {
		if response.Error() != nil {
			return res, response.Error()
		}
		res = response.Results
	} else {
		return res, err
	}
	return res, nil
}
