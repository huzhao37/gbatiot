package redis

import (
	"encoding/json"
)

//
func (c *RedisClient) Lpush(key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	_, err = c.do("LPUSH", key, data)
	return err
}

func (c *RedisClient) Lpushs(key string, value ...interface{}) error {
	_, err := c.do("LPUSH", key, value)
	return err
}

func (c *RedisClient) Lranges(key string) (interface{},error ){
	return c.do("LRANGE", key)
}

func (c *RedisClient) Latest(key string) (interface{},error ){
	return c.do("LINDEX", key,0)
}
func (c *RedisClient) Expire(key string,sec int) (interface{},error ){
	return c.do("EXPIRE", key,sec)
}
func (c *RedisClient) ExpireAt(key string,unixtime int64) (interface{},error ){
	return c.do("EXPIREAT", key,unixtime)
}
