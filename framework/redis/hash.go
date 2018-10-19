package redis

import (
	"encoding/json"
	"github.com/garyburd/redigo/redis"
)

//HSet 将哈希表 key 中的域 field 的值设为 value 。
func (c *RedisClient) HSet(key string, field string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	_, err = c.do("HSET", key, field, data)
	return err
}

//HGet 返回哈希表 key 中给定域 field 的值。
func (c *RedisClient) HGet(key string, field string, value interface{}) error {
	buffer, err := redis.Bytes(c.do("HGET", key, field))
	if err != nil {
		return err
	}
	return json.Unmarshal(buffer, value)
}

//HGetAll 返回哈希表 key 中，所有的域和值。
func (c *RedisClient) HGetAll(key string) (error,[]interface{}) {
	buffer, err := redis.ByteSlices(c.do("hvals", key))
	if err != nil {
		return err,nil
	}
	if len(buffer) > 0 {
		value:=make([]interface{},len(buffer))
		for i,b:=range buffer {
			err = json.Unmarshal(b, &value[i])
			if err != nil {
				return err,nil
			}
		}
		return  err,value
	}
	return err,nil
	//return json.Unmarshal(buffer,value)
}
//HGetAll 返回哈希表 key 中，所有的域和值。
func (c *RedisClient) HVals(key string) (error,[][]byte) {
	buffer, err := redis.ByteSlices(c.do("hvals", key))
	if err != nil {
		return err,nil
	}
	return err,buffer
}
//HDel 删除哈希表 key 中的一个或多个指定域，不存在的域将被忽略。
func (c *RedisClient) HDel(key string, field string) (int, error) {
	count, err := redis.Int(c.do("HDEL", key, field))
	if err != nil {
		return 0, err
	}
	return count, err
}



