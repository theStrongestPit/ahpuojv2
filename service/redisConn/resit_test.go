package redisConn

import (
	"testing"

	_ "github.com/gomodule/redigo/redis"
)

func TestRedis(t *testing.T) {
	var err error
	pool := newPool()
	conn := pool.Get()
	t.Log(conn)
	if err != nil {
		// handle error
		t.Errorf(err.Error())
	} else {
		tt, err := conn.Do("get", "1")
		t.Log(tt, err)
		t.Logf("%s", tt)
	}
	defer conn.Close()

}
