package influx

import (
	"github.com/influxdata/influxdb/client/v2"
	"log"
	corex "yuniot/core"
)
var (
	myConfig = new(corex.Config)
)
const (
	MyDB     = "test1"
	MyMeasurement = "cpu_usage"
)

func ConnInflux() client.Client {
	// 从配置文件获取redis的ip以及db
	myConfig.InitConfig("./config/config.txt")
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
