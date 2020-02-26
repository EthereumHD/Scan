/***********************************************************************
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.
//******
// Filename:
// Description:
// Author:
// CreateTime:
/***********************************************************************/
package apicontext

import (
	//"../config"
	//. "../const"
	"github.com/EthereumHD/Scan/src/config"
	. "github.com/EthereumHD/Scan/src/const"
	"encoding/json"
	"encoding/xml"
	"errors"
	"github.com/fatih/structs"
	"github.com/gomodule/redigo/redis"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"go-web3"
	"go-web3/providers"
	"gopkg.in/go-playground/validator.v9"
	"qoobing.com/utillib.golang/gls"
	"qoobing.com/utillib.golang/log"
	"reflect"
	"runtime/debug"
	"time"
)

type ApiContext interface {
	echo.Context
	Mysql() *gorm.DB
	//MiscMysql() *gorm.DB
	Redis() *RedisConn
	Web3() *web3.Web3

	BindInput(i interface{}) error
	PANIC_RECOVER()
	RESULT(output interface{}) error
	XMLRESULT(output interface{}) error
	RESULT_ERROR(eno int, err string) error
	RESULT_PARAMETER_ERROR(err string) error
	RecordTime()
}

type RedisConn struct {
	redis.Conn
}

type apiContext struct {
	echo.Context
	mysql *gorm.DB
	redis *RedisConn
	web3  *web3.Web3
	start time.Time
}

type BaseOutput struct {
	ErrNo  int    `json:"err_no"`
	ErrMsg string `json:"err_msg"`
}

var (
	//db *gorm.DB
	pool = newPool()
	//red *RedisConn
)

func (c *apiContext) Mysql() *gorm.DB {
	if c.mysql == nil {
		var err error
		c.mysql, err = gorm.Open("mysql", config.Config().DB.Database)
		if err != nil {
			c.mysql = nil
			log.Fatalf("connect mysql[%s] failed [%s]", config.Config().DB.Database, err)
			panic("connect mysql failed,err: " + err.Error())
		}
		c.mysql.DB().SetMaxOpenConns(config.Config().DB.MaxOpenCoons)
		c.mysql.DB().SetMaxIdleConns(config.Config().DB.MaxIdleCoons)

		gls.SetGlsValue("mysql", c.mysql)
	}

	return c.mysql
}

func (c *apiContext) Redis() *RedisConn {
	if c.redis == nil {
		c.redis = &RedisConn{Conn: pool.Get()}
		if c.redis == nil {
			panic(" c.redis is nil")
		}

		gls.SetGlsValue("redis", c.redis)
	}
	return c.redis
}

func (c *apiContext) Web3() *web3.Web3 {
	if c.web3 == nil {
		c.web3 = web3.NewWeb3(providers.NewHTTPProvider(config.Config().Gate, config.Config().TimeOut.RPCTimeOut, false))
	}
	return c.web3
}

func (c *apiContext) RecordTime() {
	c.start = time.Now()
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
func (c *apiContext) BindInput(i interface{}) error {
	err := c.Bind(i)
	b, _ := json.Marshal(i)
	req := c.Request()
	str := "[" + req.RequestURI + " " + req.Method + " " + req.Host + " " + c.RealIP() + " " + "]"
	log.Debugf("from:%s  input:%s ", str, string(b))
	if err == nil {
		err = c.Validate(i)
	}
	return err
}

func (c *apiContext) Validate(i interface{}) error {
	return paramvalidator.Struct(i)
}

func (c *apiContext) PANIC_RECOVER() {
	//Step 1. clean goroutine local storage
	gls.CleanGlsValues()

	//Step 2. recover panic
	if err := recover(); err != nil {
		t := reflect.TypeOf(err)

		log.Debugf("PANIC_RECOVER:%s", string(debug.Stack()))
		switch t.Kind() {
		case reflect.String:
			log.Fatalf("panic err:%s", err.(string))
			c.RESULT_ERROR(ERR_INNER_ERROR, err.(string))
		default:
			log.Fatalf("panic err type: %s:%s", t.Name(), t.String())
			c.RESULT_ERROR(ERR_INNER_ERROR, "unkown error")
			panic(err)
		}
	}
}

func (c *apiContext) RESULT(output interface{}) error {
	var t = reflect.TypeOf(output)
	defer func(o *interface{}) {
		b, err := json.Marshal(o)
		if err != nil {
			return
		}
		if c.mysql != nil {
			c.mysql.Close()
			c.mysql = nil
		}
		log.Debugf("output:" + string(b) + "\n")
	}(&output)

	if _, err := t.FieldByName("ErrNo"); err != true {
		errstr := "Result MUST have 'ErrNo' field"
		output = BaseOutput{ERR_INNER_ERROR, errstr}
		c.JSON(HTTPOK, output)
		return errors.New(errstr)
	}

	if _, err := t.FieldByName("ErrMsg"); err != true {
		errstr := "Result MUST have 'ErrMsg' field"
		output = BaseOutput{ERR_INNER_ERROR, errstr}
		c.JSON(HTTPOK, output)
		return errors.New(errstr)
	}

	log.Noticef("{\"Api\":\"%s\", \"Cost\":%d,\"ErrNo\":%d,\"TimeStamp\":%d,\"ProcessorName\":\"%s\"}",
		c.Request().RequestURI, time.Now().Sub(c.start).Nanoseconds(), structs.Map(output)["ErrNo"], c.start.Unix(), "scan")

	return c.JSON(HTTPOK, output)
}

func (c *apiContext) XMLRESULT(output interface{}) error {
	b, err := xml.MarshalIndent(output, "", "")
	if err != nil {
		return errors.New("xml.MarshalIndent")
	}

	log.Noticef("output:" + string(b) + "\n")
	return c.XMLPretty(HTTPOK, output, "  ")
}

func (c *apiContext) RESULT_ERROR(eno int, err string) error {
	result := BaseOutput{eno, err}
	return c.RESULT(result)
}

func (c *apiContext) RESULT_PARAMETER_ERROR(err string) error {
	return c.RESULT_ERROR(ERR_PARAMETER_INVALID, err)
}

func New(c echo.Context) ApiContext {
	return &apiContext{Context: c}
}

var (
	paramvalidator *validator.Validate
)

func init() {
	paramvalidator = validator.New()
}
