package redis

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"strconv"
	"time"
	corex "yuniot/core"
)

var (
	myConfig  = new(corex.Config)
	Redis = new(RedisClient)
)

type RedisClient struct {
	// 定义常量
	Pool       *redis.Pool
	REDIS_HOST string
	REDIS_DB   int
	REDIS_AUTH string
	DataType   int //0:json,1:proto
}

func (r *RedisClient) Init(dbIndex int) {
	// 从配置文件获取redis的ip以及db
	myConfig.InitConfig("./config/config.txt")
	var env = myConfig.Read("global", "env")
	r.REDIS_HOST = myConfig.Read(env, "redis.host")
	r.REDIS_DB = dbIndex //strconv.Atoi(myConfig.Read("test", "redis.db"))
	r.REDIS_AUTH= myConfig.Read(env, "redis.auth")
	maxidle, _ := strconv.Atoi(myConfig.Read(env, "redis.maxidle"))
	maxactive, _ := strconv.Atoi(myConfig.Read(env, "redis.maxactive"))
	// 建立连接池
	r.Pool = &redis.Pool{
		// 从配置文件获取maxidle以及maxactive，取不到则用后面的默认值
		MaxIdle:     maxidle,
		MaxActive:   maxactive,
		IdleTimeout: 180 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", r.REDIS_HOST)
			if err != nil {
				return nil, err
			}
			if r.REDIS_AUTH != "" {
				if _, err := c.Do("AUTH", r.REDIS_AUTH); err != nil {
					c.Close()
					return nil, err
				}
			}
			// 选择db
			c.Do("SELECT", r.REDIS_DB)
			return c, nil
		},
	}
}

//demo
func Test() {
	// 从池里获取连接
	client := new(RedisClient)
	rc := client.Pool.Get()
	//rc.Do("set","testkey","testvalue")
	repl, err := redis.String(rc.Do("get", "1013"))
	if err != nil {
		fmt.Printf(err.Error())
	} else {
		fmt.Printf(repl)
	}
	// 用完后将连接放回连接池
	defer rc.Close()
}
