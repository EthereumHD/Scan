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

type Block struct {
	F_id          uint64 `gorm:"column:F_id"` //ID
	F_block       int64  `gorm:"column:F_block"`
	F_timestamp   int64  `gorm:"column:F_timestamp"`
	F_txn         int64  `gorm:"column:F_txn"` //区块交易个数
	F_miner       string `gorm:"column:F_miner"`
	F_gas_used    string `gorm:"column:F_gas_used"`
	F_gas_limit   string `gorm:"column:F_gas_limit"`
	F_hash        string `gorm:"column:F_hash"`
	F_parent_hash string `gorm:"column:F_parent_hash"`
	F_reward      string `gorm:"column:F_reward"`      //区块奖励
	F_fees        string `gorm:"column:F_fees"`        //区块手续费总和
	F_status      int    `gorm:"column:F_status"`      //0 非法 ，1正常，2分叉
	F_create_time string `gorm:"column:F_create_time"` //创建时间
	F_modify_time string `gorm:"column:F_modify_time"` //修改时间

}

type Count_number struct {
	Count int64 `gorm:"column:count"`
}

type MinedBlocksGroupByDate struct {
	Date   string `gorm:"column:date"`
	Num    string `gorm:"column:num"`
	Reward string `gorm:"column:reward"`
	Fees   string `gorm:"column:fees"`
}

func (b *Block) TableName() string {
	return "t_block"
}

//FindBlock

//func FindBox(db *gorm.DB, inputboxid string) (box Block, err error) {
//	var rdb *gorm.DB
//
//	rdb = db.Where("F_boxid = ?", inputboxid).First(&box)
//	if rdb.RecordNotFound() {
//		err = errors.New(DATA_NOT_EXIST)
//	} else if rdb.Error != nil {
//		panic("find box error:" + rdb.Error.Error())
//	} else {
//		err = nil
//	}
//
//	return box, err
//}

func (block *Block) BeforeCreate(scope *gorm.Scope) error {
	currentTime := time.Now().Local()
	newFormat := currentTime.Format("2006-01-02 15:04:05.000")

	scope.SetColumn("F_create_time", newFormat)
	scope.SetColumn("F_modify_time", newFormat)

	scope.SetColumn("F_status", NORMAL)
	return nil
}

func (block *Block) BeforeUpdate(scope *gorm.Scope) error {
	currentTime := time.Now().Local()
	newFormat := currentTime.Format("2006-01-02 15:04:05.000")
	scope.SetColumn("F_modify_time", newFormat)
	return nil
}

func (block *Block) CreateBlock(db *gorm.DB) (err error) {

	//ASSERT(block.F_block != 0, "create block, F_block can't be nul")
	util.ASSERT(block.F_miner != "", "create block, F_miner can't be nul")
	util.ASSERT(block.F_hash != "", "create block, F_hash can't be nul")
	util.ASSERT(block.F_parent_hash != "", "create block, F_parent_hash can't be nul")

	rdb := db.Create(&block)

	return rdb.Error
}

func (b *Block) FindBlockByHash(db *gorm.DB, hash string) (block Block, err error) {

	rdb := db.Where("F_hash = ?", hash).First(&block)
	if rdb.RecordNotFound() {
		err = errors.New(DATA_NOT_EXIST)
	} else if rdb.Error != nil {
		panic("FindBlockByHash error:" + rdb.Error.Error())
	} else {
		err = nil
	}

	return block, err
}

func (b *Block) FindBlockByHeight(db *gorm.DB, height int64) (block Block, err error) {

	rdb := db.Where("F_block = ? and F_status = ?", height, NORMAL).First(&block)
	if rdb.RecordNotFound() {
		err = errors.New(DATA_NOT_EXIST)
	} else if rdb.Error != nil {
		panic("FindBlockByHeight error:" + rdb.Error.Error())
	} else {
		err = nil
	}

	return block, err
}

func (b *Block) UpdateBlockStatus(db *gorm.DB) (err error) {
	updateinfo := map[string]interface{}{"F_status": b.F_status}
	return b.updateBlockColumn(db, updateinfo)
}

func (b *Block) updateBlockColumn(db *gorm.DB, updateinfo map[string]interface{}) (err error) {

	log.Debugf("updateBlockColumn F_id:%d,hash:%s,%+v", b.F_id, b.F_hash, updateinfo)

	tx := db.Begin()

	rdb := tx.Where("F_id = ?", b.F_id).Model(&b).Update(updateinfo)

	if rdb.Error != nil {
		tx.Rollback()
		return rdb.Error
	}

	tx.Commit()
	return rdb.Error
}

func GetRecentBlocks(db *gorm.DB, offset int, size int) (blocks []Block, err error) {
	rdb := db.Where("F_status = ?", NORMAL).Order("F_block desc").Offset(offset).Limit(size).Find(&blocks)
	if rdb.Error != nil {
		err = errors.New("GetRecentBlocks error:" + rdb.Error.Error())
	} else {
		err = nil
	}

	return blocks, err
}

func GetRecentOneDayReward(db *gorm.DB, offset int, size int) (blocks []Block, err error) {
	rdb := db.Where("F_status = ?", NORMAL).Order("F_block desc").Offset(offset).Limit(size).Find(&blocks)
	if rdb.Error != nil {
		err = errors.New("GetRecentBlocks error:" + rdb.Error.Error())
	} else {
		err = nil
	}

	return blocks, err
}

func GetRecentOneDayBlockNumber(db *gorm.DB) (n int64, err error) {
	end := time.Now().Unix()
	start := time.Now().Add(-24 * time.Hour).Unix()
	num := Count_number{}
	rdb := db.Table("t_block").
		Where("F_timestamp >=? and F_timestamp <= ? and F_status = ? ", start, end, NORMAL).
		Select(" count(*) as count ").Find(&num)
	if rdb.RecordNotFound() {
		err = errors.New(DATA_NOT_EXIST)
	} else if rdb.Error != nil {
		panic("FindMinedBlockByAddrAndTime error:" + rdb.Error.Error())
	} else {
		err = nil
	}
	return num.Count, err
}

func GetBlocksByMinerAddr(db *gorm.DB, addr string, offset int, size int) (blocks []Block, err error) {
	rdb := db.Where("F_miner = ? and F_status = ?", addr, NORMAL).Order("F_block desc").Offset(offset).Limit(size).Find(&blocks)
	if rdb.Error != nil {
		err = errors.New("GetRecentBlocks error:" + rdb.Error.Error())
	} else {
		err = nil
	}

	return blocks, err
}

func GetBlockByHash(db *gorm.DB, hash string) (blocks []Block, err error) {
	rdb := db.Where("F_hash = ? and F_status = ?", hash, NORMAL).Find(&blocks)
	if rdb.Error != nil {
		err = errors.New("GetBlockByHash error:" + rdb.Error.Error())
	} else {
		err = nil
	}

	return blocks, err
}

func GetActiveBlockNum(db *gorm.DB) (count int64, err error) {
	var rdb *gorm.DB
	num := Count_number{}
	rdb = db.Table("t_block").Where("F_status = ?", NORMAL).Select(" count(*) as count ").Find(&num)

	if rdb.Error != nil {
		//panic("find information error:" + rdb.Error.Error())
		err = errors.New("GetActiveBlockNum error:" + rdb.Error.Error())
	} else {
		err = nil
	}
	return num.Count, err
}

func GetActiveBlockNumByAddr(db *gorm.DB, addr string) (count int64, err error) {
	var rdb *gorm.DB
	num := Count_number{}
	rdb = db.Table("t_block").Where("F_miner = ? and F_status = ?", addr, NORMAL).Select(" count(*) as count ").Find(&num)

	if rdb.Error != nil {
		//panic("find information error:" + rdb.Error.Error())
		err = errors.New("GetActiveBlockNum error:" + rdb.Error.Error())
	} else {
		err = nil
	}
	return num.Count, err
}

func (b *Block) GetMaxBlocNumber(db *gorm.DB) (number int64, err error) {

	type Number struct {
		Number int64 `gorm:"column:max_block"`
	}

	n := Number{}

	rdb := db.Table(b.TableName()).Where("F_status = ?", NORMAL).Select("MAX(F_block) as max_block").Find(&n)
	if rdb.RecordNotFound() {
		err = errors.New(DATA_NOT_EXIST)
		log.Fatalf("DATA_NOT_EXIST")
	} else if rdb.Error != nil {
		panic("GetMaxBlocNumber error:" + rdb.Error.Error())
	} else {
		err = nil
	}

	return n.Number, err
}

func (b *Block) FindMinedBlockByAddrAndTime(db *gorm.DB, addr string, start, end int64) (blocks []Block, err error) {

	rdb := db.Where("F_miner = ? and F_timestamp >=? and F_timestamp <= ? and F_status = ? ", addr, start, end, NORMAL).Find(&blocks)
	if rdb.RecordNotFound() {
		err = errors.New(DATA_NOT_EXIST)
	} else if rdb.Error != nil {
		panic("FindMinedBlockByAddrAndTime error:" + rdb.Error.Error())
	} else {
		err = nil
	}

	return blocks, err
}

func (b *Block) FindMinedBlockByAddrAndGroupByDate(db *gorm.DB, addr string, start, end int64) (blocks []MinedBlocksGroupByDate, err error) {
	rdb := db.Table("t_block").
		Where("F_miner = ? and F_timestamp >=? and F_timestamp <= ? and F_status = ? ", addr, start, end, NORMAL).
		Select("FROM_UNIXTIME(F_timestamp, '%Y-%m-%d') as date, count(F_hash) as num, sum(F_reward) as reward, sum(F_fees) as fees").
		Group("date").
		Scan(&blocks)
	if rdb.RecordNotFound() {
		err = errors.New(DATA_NOT_EXIST)
	} else if rdb.Error != nil {
		panic("FindMinedBlockByAddrAndTime error:" + rdb.Error.Error())
	} else {
		err = nil
	}

	return blocks, err
}
