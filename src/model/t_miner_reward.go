package model

import (
	. "github.com/EthereumHD/Scan/src/const"
	"github.com/EthereumHD/Scan/src/util"
	"errors"
	_ "fmt"
	"github.com/jinzhu/gorm"
	"qoobing.com/utillib.golang/log"
	"time"
)

//统计地址的挖矿总收益
type MinerReward struct {
	F_id           uint64 `gorm:"column:F_id"` //ID
	F_miner        string `gorm:"column:F_miner"`
	F_total_reward string `gorm:"column:F_total_reward"`
	F_total_fees   string `gorm:"column:F_total_fees"`
	F_create_time  string `gorm:"column:F_create_time"` //创建时间
	F_modify_time  string `gorm:"column:F_modify_time"` //修改时间

}

func (r *MinerReward) TableName() string {
	return "t_miner_reward"
}

func (r *MinerReward) BeforeCreate(scope *gorm.Scope) error {
	currentTime := time.Now().Local()
	newFormat := currentTime.Format("2006-01-02 15:04:05.000")

	scope.SetColumn("F_create_time", newFormat)
	scope.SetColumn("F_modify_time", newFormat)
	return nil
}

func (r *MinerReward) BeforeUpdate(scope *gorm.Scope) error {
	currentTime := time.Now().Local()
	newFormat := currentTime.Format("2006-01-02 15:04:05.000")
	scope.SetColumn("F_modify_time", newFormat)
	return nil
}

func (r *MinerReward) CreateMinerReward(db *gorm.DB) (err error) {

	log.Debugf("CreateMinerReward,miner:%s,total_reward:%s", r.F_miner, r.F_total_reward)
	//ASSERT(block.F_block != 0, "create block, F_block can't be nul")
	util.ASSERT(r.F_miner != "", "CreateMinerReward, F_miner can't be nul")
	util.ASSERT(r.F_total_reward != "", "CreateMinerReward, F_total_reward can't be nul")

	rdb := db.Create(&r)

	return rdb.Error
}

func (r *MinerReward) FindRewardByMiner(db *gorm.DB, addr string) (reward MinerReward, err error) {

	rdb := db.Where("F_miner = ?", addr).First(&reward)
	if rdb.RecordNotFound() {
		err = errors.New(DATA_NOT_EXIST)
	} else if rdb.Error != nil {
		panic("FindRewardByMiner error:" + rdb.Error.Error())
	} else {
		err = nil
	}

	return reward, err
}

func (r *MinerReward) UpdateMinerReward(db *gorm.DB) (err error) {
	updateinfo := map[string]interface{}{"F_total_reward": r.F_total_reward}
	return r.updateMinerRewardColumn(db, updateinfo)
}

func (r *MinerReward) updateMinerRewardColumn(db *gorm.DB, updateinfo map[string]interface{}) (err error) {

	log.Debugf("updateMinerRewardColumn F_id:%d,miner:%s,%+v", r.F_id, r.F_miner, updateinfo)

	tx := db.Begin()

	rdb := tx.Where("F_miner = ?", r.F_miner).Model(&r).Update(updateinfo)

	if rdb.Error != nil {
		tx.Rollback()
		return rdb.Error
	}

	tx.Commit()
	return rdb.Error
}

func (r *MinerReward) GetAddrList(db *gorm.DB) (addrlist []MinerReward, err error) {

	rdb := db.Find(&addrlist)
	if rdb.RecordNotFound() {
		err = errors.New(DATA_NOT_EXIST)
	} else if rdb.Error != nil {
		panic("FindRewardByMiner error:" + rdb.Error.Error())
	} else {
		err = nil
	}

	return addrlist, err
}
