package model

import (
	"github.com/EthereumHD/Scan/src/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"qoobing.com/utillib.golang/log"
)

var Schema = config.Config().DB.Schema

var Table = map[string]string{
	"t_transaction": "CREATE TABLE IF NOT EXISTS " + Schema + ".t_transaction (" +
		"`F_id` bigint(20) unsigned NOT NULL AUTO_INCREMENT," +
		"`F_tx_hash` varchar(128) NOT NULL DEFAULT ''," +
		"`F_block` int(64)  NOT NULL DEFAULT -1," +
		"`F_timestamp` int(64)   NOT NULL DEFAULT -1," +
		"`F_from` varchar(128) NOT NULL DEFAULT ''," +
		"`F_to` varchar(128) NOT NULL DEFAULT ''," +
		"`F_value` varchar(128) NOT NULL DEFAULT ''," +
		"`F_tx_fee` varchar(128) NOT NULL DEFAULT ''," +
		"`F_status` int(4)  NOT NULL DEFAULT 0," +
		"`F_tx_type` bigint(20)  NOT NULL DEFAULT 0," +
		"`F_tx_type_ext` varchar(128) NOT NULL DEFAULT ''," +
		"`F_create_time` datetime NOT NULL," +
		"`F_modify_time` datetime NOT NULL," +

		"PRIMARY KEY (`F_id`)," +
		"UNIQUE KEY (`F_tx_hash`)," +
		"INDEX (`F_from`)," +
		"INDEX (`F_to`)," +
		"INDEX (`F_block`)" +
		") ENGINE=InnoDB  DEFAULT CHARSET=utf8 ;",

	"t_pending": "CREATE TABLE IF NOT EXISTS " + Schema + ".t_pending (" +
		"`F_id` bigint(20) unsigned NOT NULL AUTO_INCREMENT," +
		"`F_tx_hash` varchar(128) NOT NULL DEFAULT ''," +
		"`F_from` varchar(128) NOT NULL DEFAULT ''," +
		"`F_to` varchar(128) NOT NULL DEFAULT ''," +
		"`F_value` varchar(128) NOT NULL DEFAULT ''," +
		"`F_tx_fee` varchar(128) NOT NULL DEFAULT ''," +
		"`F_status` int(4)  NOT NULL DEFAULT 0," +
		"`F_create_time` datetime NOT NULL," +
		"`F_modify_time` datetime NOT NULL," +

		"PRIMARY KEY (`F_id`)," +
		"UNIQUE KEY (`F_tx_hash`)," +
		"INDEX (`F_from`)," +
		"INDEX (`F_to`)" +
		") ENGINE=InnoDB  DEFAULT CHARSET=utf8 ;",

	"t_block": "CREATE TABLE IF NOT EXISTS " + Schema + ".t_block (" +
		"`F_id` bigint(20) unsigned NOT NULL AUTO_INCREMENT," +
		"`F_block` int(64)  NOT NULL DEFAULT -1," +
		"`F_timestamp` int(64)   NOT NULL DEFAULT -1," +
		"`F_txn` int(64)  NOT NULL DEFAULT -1," +
		"`F_miner` varchar(128) NOT NULL DEFAULT ''," +
		"`F_gas_used` varchar(128) NOT NULL DEFAULT ''," +
		"`F_gas_limit` varchar(128) NOT NULL DEFAULT ''," +
		"`F_hash` varchar(128) NOT NULL DEFAULT ''," +
		"`F_parent_hash` varchar(128) NOT NULL DEFAULT ''," +
		"`F_reward` varchar(128) NOT NULL DEFAULT ''," +
		"`F_fees` varchar(128) NOT NULL DEFAULT ''," +
		"`F_status` int(4)  NOT NULL DEFAULT 0," +
		"`F_create_time` datetime NOT NULL," +
		"`F_modify_time` datetime NOT NULL," +

		"PRIMARY KEY (`F_id`)," +
		"UNIQUE KEY (`F_hash`)," +
		"INDEX (`F_miner`)," +
		"INDEX (`F_block`)" +
		") ENGINE=InnoDB  DEFAULT CHARSET=utf8 ;",

	"t_miner_reward": "CREATE TABLE IF NOT EXISTS " + Schema + ".t_miner_reward (" +
		"`F_id` bigint(20) unsigned NOT NULL AUTO_INCREMENT," +
		"`F_miner` varchar(128) NOT NULL DEFAULT ''," +
		"`F_total_reward` varchar(128) NOT NULL DEFAULT ''," +
		"`F_total_fees` varchar(128) NOT NULL DEFAULT ''," +
		"`F_create_time` datetime NOT NULL," +
		"`F_modify_time` datetime NOT NULL," +

		"PRIMARY KEY (`F_id`)," +
		"UNIQUE KEY (`F_miner`)" +
		") ENGINE=InnoDB  DEFAULT CHARSET=utf8 ;",

	"t_rate": "CREATE TABLE IF NOT EXISTS " + Schema + ".t_rate (" +
		"`F_id` bigint(20) unsigned NOT NULL AUTO_INCREMENT," +
		"`F_eth` double  NOT NULL DEFAULT 0," +
		"`F_btc` double  NOT NULL DEFAULT 0," +
		"`F_usd` double  NOT NULL DEFAULT 0," +
		"`F_kwr` double  NOT NULL DEFAULT 0," +
		"`F_timestamp` int(64)  NOT NULL DEFAULT 0," +
		"`F_create_time` datetime NOT NULL," +
		"`F_modify_time` datetime NOT NULL," +

		"PRIMARY KEY (`F_id`)," +
		"UNIQUE KEY (`F_timestamp`)" +
		") ENGINE=InnoDB  DEFAULT CHARSET=utf8 ;",
}

func InitDatabase() {
	db, err := gorm.Open("mysql", config.Config().DB.Database)
	defer db.Close()
	if err != nil {
		log.Fatalf("connect mysql[%s] failed [%s]", config.Config().DB.Database, err)
		panic("connect mysql failed")
	}

	for _, value := range Table {
		db.Exec(value)
	}
}
