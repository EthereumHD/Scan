package sync

import (
	"github.com/EthereumHD/Scan/src/config"
	"github.com/gomodule/redigo/redis"
	"github.com/jinzhu/gorm"
	"go-web3"
	"go-web3/providers"
	"qoobing.com/utillib.golang/log"
	"reflect"
	"runtime/debug"
)

type Connect struct {
	mysql *gorm.DB
	web3  *web3.Web3
	redis *RedisConn
}

type RedisConn struct {
	redis.Conn
}

type BaseOutput struct {
	ErrNo  int    `json:"err_no"`
	ErrMsg string `json:"err_msg"`
}

var (
	block_now int64
	//db *gorm.DB
	pool = newPool()
	//red *RedisConn
)

func init() {
	block_now = -1
}

func (c *Connect) Web3() *web3.Web3 {
	if c.web3 == nil {
		c.web3 = web3.NewWeb3(providers.NewHTTPProvider(config.Config().Gate, config.Config().TimeOut.RPCTimeOut, false))
	}
	return c.web3
}

func (c *Connect) Mysql() *gorm.DB {
	if c.mysql == nil {
		var err error
		c.mysql, err = gorm.Open("mysql", config.Config().DB.Database)
		if err != nil {
			c.mysql = nil
			log.Fatalf("connect mysql[%s] failed [%s]", config.Config().DB.Database, err.Error())
			panic("connect mysql failed,err: " + err.Error())
		}
		c.mysql.DB().SetMaxOpenConns(config.Config().DB.MaxOpenCoons)
		c.mysql.DB().SetMaxIdleConns(config.Config().DB.MaxIdleCoons)

		//SetGlsValue("mysql", c.mysql)
	}

	return c.mysql
}

func (c *Connect) Redis() *RedisConn {
	if c.redis == nil {
		c.redis = &RedisConn{Conn: pool.Get()}
		if c.redis == nil {
			panic(" c.redis is nil")
		}

		//SetGlsValue("redis", c.redis)
	}
	return c.redis
}

func newPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:   80,
		MaxActive: 12000, // max number of connections
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", config.Config().Redis)
			if err != nil {
				log.Fatalf("connect redis[%s] failed [%s]", config.Config().Redis, err)
				panic(err.Error())
			}
			return c, err
		},
	}

}

func (c *Connect) Close() {

	if c.mysql != nil {
		err := c.mysql.Close()
		if err != nil {
			log.Debugf("c.mysql.Close,err:%s", err.Error())
		}

		c.mysql = nil
	}
	if c.redis != nil {
		c.redis.Close()
		c.redis = nil
	}

	c.PanicRecover()
}

func (c *Connect) PanicRecover() {
	//Step 1. clean goroutine local storage
	//CleanGlsValues()

	//Step 2. recover panic
	if err := recover(); err != nil {
		t := reflect.TypeOf(err)

		log.Debugf("PANIC_RECOVER:%s", string(debug.Stack()))
		switch t.Kind() {
		case reflect.String:
			log.Fatalf("panic err:%s", err.(string))
		default:
			log.Fatalf("panic err type: %s:%s", t.Name(), t.String())
			panic(err)
		}
	}
}

func (c *Connect) SetBlockNow(number int64) {
	block_now = number
}

func (c *Connect) AddBlockNow(num int64) {
	block_now = block_now + num
}

func (c *Connect) GetBlockNOw() int64 {
	return block_now
}
