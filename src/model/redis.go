package model

import (
	"github.com/gomodule/redigo/redis"
	//"config"
	. "github.com/EthereumHD/Scan/src/const"
	"qoobing.com/utillib.golang/log"
	"github.com/EthereumHD/Scan/src/config"
	"encoding/json"
)


func SetRate(rds redis.Conn, rate UbbeyRate) error {

	data, err := json.Marshal(rate)
	if err != nil {
		return err
	}

	value := string(data[:])
	_, err = rds.Do("SETEX", UBBEY_RATE, config.Config().RateInRedis,value)
	if err != nil {
		log.Fatalf("SETEX Rate-%s success, value:[%+v],time:%d",UBBEY_RATE, rate, config.Config().RateInRedis)
		return err
	}

	log.Debugf("SETEX Rate-%s success, value:[%+v],time:%d",UBBEY_RATE, rate, config.Config().RateInRedis)
	return nil
}
