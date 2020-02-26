package gotest

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"testing"
)

var (
	//block_now int64
	rediscoon *RedisConn
)

type RedisConn struct {
	redis.Conn
}

func get_redis() *RedisConn {

	if rediscoon == nil {
		rconn, err := redis.Dial("tcp", "127.0.0.1:6379")
		rediscoon = &RedisConn{Conn: rconn}
		fmt.Println("abc",err)
		if err != nil {
			rediscoon = nil
			fmt.Println("connect redis[%s] failed [%s]")
			panic("connect redis failed")
		}
	}

	return rediscoon
}

func Test_Redis(t *testing.T) {

	type Input struct {
		Boxid string
	}

	//err := model.SetRedisData(get_redis(), "FDSAF", Input{"fdsa"}, 111)
	//if err != nil {
	//	t.Error("err:%s", err.Error())
	//}
}
