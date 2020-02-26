/***********************************************************************
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.
//******
// Filename:
// Description:
// Author:
// CreateTime:
/***********************************************************************/
package config

import (
	"github.com/pelletier/go-toml"
	"io/ioutil"
	"qoobing.com/utillib.golang/log"
	"sync"
)

type appConfig struct {
	Server string
	IP     string
	Port   string
	Gate   string

	DB    database `toml:"database"`
	Redis string

	TimeOut timeout

	RateSyncInterval int64
	RateInRedis      int64

	Stats stats
}

type database struct {
	Database     string
	MaxOpenCoons int
	MaxIdleCoons int
	Schema       string
}

type timeout struct {
	BlockchainTimeout int64
	RPCTimeOut        int32
}

type stats struct {
	StatAddr string
	ServerId string
}

//

var (
	cfg  appConfig
	once sync.Once
)

func Config() *appConfig {
	once.Do(func() {
		doc, err := ioutil.ReadFile("./conf/scan.conf")
		if err != nil {
			panic("initial config, read config file error:" + err.Error())
		}
		if err := toml.Unmarshal(doc, &cfg); err != nil {
			panic("initial config, unmarshal config file error:" + err.Error())
		}

		if cfg.Stats.StatAddr == "" {
			cfg.Stats.StatAddr = ":3000"
		}

		if cfg.Stats.ServerId == "" {
			cfg.Stats.ServerId = "Scan&Stats"
		}

		log.Debugf("config:%+v\n", cfg)
	})
	return &cfg
}
