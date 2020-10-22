package gredis

import (
	"QuickPass/pkg/setting"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
	"time"
)

func Setup() {
	idle := setting.RedisSetting.MaxIdle
	sec := setting.RedisSetting.IdleTimeoutSec
	url := setting.RedisSetting.Url
	db := setting.RedisSetting.Db
	password := setting.RedisSetting.Password
	clientPool = newRedisPool(idle, sec, db, url, password)
}

// redis客户端连接池
var clientPool *redis.Pool

// NewRedisPool 返回redis连接池
func newRedisPool(maxIdle, idleTimeout, db int, url, password string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     maxIdle,
		IdleTimeout: time.Duration(idleTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.DialURL(fmt.Sprintf("redis://%s", url))
			if err != nil {
				log.Fatal("redis connection error: ", err)
				return nil, err
			}
			//验证redis密码
			if password != "" {
				if _, authErr := c.Do("AUTH", password); authErr != nil {
					log.Println("redis auth password error: ", authErr)
					return nil, authErr
				}
			}
			if _, err := c.Do("SELECT", db); err != nil {
				c.Close()
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			if err != nil {
				log.Fatal("ping redis error: ", err)
				return err
			}
			return nil
		},
	}
}

func Set(k, v string) error {
	c := clientPool.Get()
	defer c.Close()
	_, err := c.Do("SET", k, v)
	return err
}

func GetStringValue(k string) (string, error) {
	c := clientPool.Get()
	defer c.Close()
	username, err := redis.String(c.Do("GET", k))
	if err != nil {
		return "", err
	}
	return username, nil
}

func SetKeyExpire(k string, second int64) error {
	c := clientPool.Get()
	defer c.Close()
	_, err := c.Do("EXPIRE", k, second)
	return err
}

func CheckKey(k string) (bool, error) {
	c := clientPool.Get()
	defer c.Close()
	exist, err := redis.Bool(c.Do("EXISTS", k))
	if err != nil {
		return false, err
	}

	return exist, nil
}

func DelKey(k string) error {
	c := clientPool.Get()
	defer c.Close()
	_, err := c.Do("DEL", k)
	return err
}

func SetJson(k string, data interface{}) error {
	c := clientPool.Get()
	defer c.Close()
	value, _ := json.Marshal(data)
	n, err := c.Do("SETNX", k, value)
	if err != nil {
		return err
	}
	if n != int64(1) {
		return errors.New("set failed")
	}
	return nil
}

func GetJsonByte(k string) ([]byte, error) {
	c := clientPool.Get()
	defer c.Close()
	return redis.Bytes(c.Do("GET", k))
}
