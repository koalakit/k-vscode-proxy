package main

import (
	"encoding/json"
	"time"

	"github.com/gomodule/redigo/redis"
)

// 全局redis实例池
var redisdb redis.Pool

// Redis数据库初始化
func RedisInit(url string) (err error) {
	redisdb.MaxIdle = 5
	redisdb.MaxActive = 10
	redisdb.IdleTimeout = 1 * time.Hour

	redisdb.Dial = func() (c redis.Conn, err error) {
		c, err = redis.DialURL(url)
		if err != nil {
			Log("[Redis] DialURL: ", err)
			return
		}

		return
	}

	go func() {
		ticker := time.NewTicker(2 * time.Minute)

		for {
			<-ticker.C

			if err := RedisPing(); err != nil {
				Log("[Redis] ping: ", err)
				continue
			}
		}
	}()

	err = RedisPing()

	return
}

// Redis ping测试
func RedisPing() (err error) {
	db := redisdb.Get()
	defer db.Close()

	_, err = db.Do("PING")
	return
}

func RedisDel(key string) (err error) {
	db := redisdb.Get()
	defer db.Close()

	_, err = db.Do("DEL", key)
	return
}

func RedisSet(key string, value any, timeout time.Duration) (err error) {
	db := redisdb.Get()
	defer db.Close()

	if timeout > 0 {
		_, err = db.Do("SETEX", key, int64(timeout/time.Second), value)
		return
	}

	_, err = db.Do("SET", key, value)

	return
}

func RedisSetJSON(key string, value any, timeout time.Duration) (err error) {
	db := redisdb.Get()
	defer db.Close()

	data, err := json.Marshal(value)
	if err != nil {
		return
	}

	if timeout > 0 {
		_, err = db.Do("SETEX", key, int64(timeout/time.Second), data)
		return
	}

	_, err = db.Do("SET", key, data)

	return
}

func RedisGet(key string) (value string, err error) {
	db := redisdb.Get()
	defer db.Close()

	value, err = redis.String(db.Do("GET", key))
	return
}

func RedisGetBytes(key string) (value []byte, err error) {
	db := redisdb.Get()
	defer db.Close()

	value, err = redis.Bytes(db.Do("GET", key))
	return
}

func RedisGetInt(key string) (value int, err error) {
	db := redisdb.Get()
	defer db.Close()

	value, err = redis.Int(db.Do("GET", key))
	return
}

func RedisGetInt64(key string) (value int64, err error) {
	db := redisdb.Get()
	defer db.Close()

	value, err = redis.Int64(db.Do("GET", key))
	return
}

func RedisGetUint(key string) (value uint, err error) {
	db := redisdb.Get()
	defer db.Close()

	x, err := redis.Uint64(db.Do("GET", key))
	value = uint(x)
	return
}

func RedisGetUint64(key string) (value uint64, err error) {
	db := redisdb.Get()
	defer db.Close()

	value, err = redis.Uint64(db.Do("GET", key))
	return
}

func RedisGetFloat64(key string) (value float64, err error) {
	db := redisdb.Get()
	defer db.Close()

	value, err = redis.Float64(db.Do("GET", key))
	return
}

func RedisGetJSON[ValueT any](key string) (value ValueT, err error) {
	db := redisdb.Get()
	defer db.Close()

	x, err := redis.Bytes(db.Do("GET", key))
	if err != nil {
		return
	}

	err = json.Unmarshal(x, &value)

	return
}

func RedisGetJSONEx(key string, value any) (err error) {
	db := redisdb.Get()
	defer db.Close()

	x, err := redis.Bytes(db.Do("GET", key))
	if err != nil {
		return
	}

	err = json.Unmarshal(x, value)

	return
}
