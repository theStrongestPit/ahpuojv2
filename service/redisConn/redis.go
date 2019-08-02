package redisConn

import (
	"ahpuoj/utils"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"
	_ "github.com/gomodule/redigo/redis"
)

func NewPool() *redis.Pool {
	cfg := utils.GetCfg()
	// cfg := utils.GetTestCfg("../../config/config.ini")
	rediscfg, _ := cfg.GetSection("redis")
	path := strings.Join([]string{rediscfg["host"], ":", rediscfg["port"]}, "")

	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", path)
			if err != nil {
				return nil, err
			}
			// if _, err := c.Do("AUTH", password); err != nil {
			// 	c.Close()
			// 	return nil, err
			// }
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}
