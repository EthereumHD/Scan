package sync

import (
	"github.com/EthereumHD/Scan/src/config"
	. "github.com/EthereumHD/Scan/src/model"
	"github.com/jinzhu/gorm"
	"math/big"
	"qoobing.com/utillib.golang/log"
	"time"
)

func CheckReward() {
	for {
		checkReward()
		time.Sleep(3 * time.Hour)
	}
}

func checkReward() {
	mysql, err := gorm.Open("mysql", config.Config().DB.Database)
	if err != nil {
		//fatal_list("connect mysql[%s] failed [%s]", config.Config().DB.Database, err)
		panic("connect mysql failed,err: " + err.Error())
	}

	miners, err := (&MinerReward{}).GetAddrList(mysql)
	if err != nil {
		panic("GetAddrList,err: " + err.Error())
	}

	for _, miner := range miners {
		blocks, err := (&Block{}).FindMinedBlockByAddrAndTime(mysql, miner.F_miner, 0, 9540364631)
		if err != nil {
			panic("FindMinedBlockByAddrAndTime,err: " + err.Error())
		}

		rewards := big.NewInt(0)
		fees := big.NewInt(0)

		for _, block := range blocks {
			reward, _ := big.NewInt(0).SetString(block.F_reward, 10)
			fee, _ := big.NewInt(0).SetString(block.F_fees, 10)

			rewards.Add(rewards, reward)
			fees.Add(fees, fee)
		}

		if miner.F_total_reward != rewards.String() || miner.F_total_fees != fees.String() {
			log.Debugf("addr :%s,Old total_reward:%s != all block rewards:%s", miner.F_miner, miner.F_total_reward, rewards.String())

			miner.F_total_reward = rewards.String()
			miner.F_total_fees = fees.String()

			err := miner.UpdateMinerReward(mysql)
			if err != nil {
				panic("UpdateMinerReward,err: " + err.Error())
			}
		} else {
			log.Debugf("addr :%s check is right", miner.F_miner)
		}

	}
}
