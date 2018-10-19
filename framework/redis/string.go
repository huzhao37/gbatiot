package redis

import (
	"encoding/json"
	"github.com/garyburd/redigo/redis"
)

func (r *RedisClient) do(commandName string, args ...interface{}) (reply interface{}, err error) {
	c := r.Pool.Get()
	defer c.Close()
	return c.Do(commandName, args...)
}

//INFO 以一种易于解释（parse）且易于阅读的格式，返回关于 Redis 服务器的各种信息和统计数值。
func (c *RedisClient) INFO(key string) (string, error) {
	info, err := redis.String(c.do("INFO", key))
	return info, err
}

//Exists 检查给定 key 是否存在。
func (c *RedisClient) Exists(key string) (bool, error) {
	count, err := redis.Int(c.do("EXISTS", key))
	var exist = count == 1
	return exist, err
}

//Set 设置字符串值
func (c *RedisClient) Set(key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		//插入基本类型
		_, err = c.do("SET", key, value)
	} else {
		//插入对象或数组
		_, err = c.do("SET", key, data)
	}
	return err
}

//Get 获取字符串值
func (c *RedisClient) Get(key string, value interface{}) error {
	buffer, err := redis.Bytes(c.do("GET", key))
	if err != nil {
		return err
	}
	return json.Unmarshal(buffer, value)
}
//DEL删除给定 key
func (c *RedisClient) Del(key string) (bool, error) {
	count, err := redis.Int(c.do("DEL", key))
	var exist = count == 1
	return exist, err
}


//选择数据库
func (c *RedisClient) Select(dbindex int) (bool, error) {
	count, err := redis.Int(c.do("SELECT", dbindex))
	var exist = count == 1
	return exist, err
}
