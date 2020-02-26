package model

import (
	. "github.com/EthereumHD/Scan/src/const"
	. "github.com/EthereumHD/Scan/src/util"
	"errors"
	_ "fmt"
	"github.com/jinzhu/gorm"
	"qoobing.com/utillib.golang/log"
	"time"
)

type Transaction struct {
	F_id          uint64 `gorm:"column:F_id"` //ID
	F_tx_hash     string `gorm:"column:F_tx_hash"`
	F_block       int64  `gorm:"column:F_block"`
	F_timestamp   int64  `gorm:"column:F_timestamp"`
	F_from        string `gorm:"column:F_from"`
	F_to          string `gorm:"column:F_to"`
	F_value       string `gorm:"column:F_value"`
	F_tx_fee      string `gorm:"column:F_tx_fee"`
	F_status      int    `gorm:"column:F_status"` //0 非法 ，1正常，2分叉
	F_tx_type     int64  `gorm:"column:F_tx_type"`
	F_tx_type_ext string `gorm:"column:F_tx_type_ext"`
	F_create_time string `gorm:"column:F_create_time"` //创建时间
	F_modify_time string `gorm:"column:F_modify_time"` //修改时间

}

const (
	TX_TYPE_UNDEFINED   = 0
	TX_TYPE_ALL         = 0
	TX_TYPE_FROM_ME     = 1
	TX_TYPE_TO_ME       = 2
	TX_TYPE_ME_MORTGAGE = 3
	TX_TYPE_ME_REDEEM   = 4
	TX_TYPE_QUERY_3OR4  = 5
)

func (t *Transaction) TableName() string {
	return "t_transaction"
}

func (t *Transaction) BeforeCreate(scope *gorm.Scope) error {
	currentTime := time.Now().Local()
	newFormat := currentTime.Format("2006-01-02 15:04:05.000")

	scope.SetColumn("F_create_time", newFormat)
	scope.SetColumn("F_modify_time", newFormat)

	scope.SetColumn("F_status", NORMAL)
	return nil
}

func (t *Transaction) BeforeUpdate(scope *gorm.Scope) error {
	currentTime := time.Now().Local()
	newFormat := currentTime.Format("2006-01-02 15:04:05.000")
	scope.SetColumn("F_modify_time", newFormat)
	return nil
}

func (t *Transaction) CreateTransaction(db *gorm.DB) (err error) {

	//ASSERT(t.F_block != 0, "CreateTransaction, F_block can't be nul")
	ASSERT(t.F_tx_hash != "", "CreateTransaction, F_tx_hash can't be nul")

	rdb := db.Create(&t)

	return rdb.Error
}

func (t *Transaction) FindTrasactionByHash(db *gorm.DB, hash string) (transcation Transaction, err error) {

	rdb := db.Where("F_tx_hash = ?", hash).First(&transcation)
	if rdb.RecordNotFound() {
		err = errors.New(DATA_NOT_EXIST)
	} else if rdb.Error != nil {
		panic("FindTrasactionByHash error:" + rdb.Error.Error())
	} else {
		err = nil
	}

	return transcation, err
}

func (t *Transaction) FindTrasactionByHeight(db *gorm.DB, height int64) (transcations []Transaction, err error) {

	rdb := db.Where("F_block = ? and F_status = ?", height, NORMAL).Find(&transcations)
	if rdb.RecordNotFound() {
		err = errors.New(DATA_NOT_EXIST)
	} else if rdb.Error != nil {
		panic("FindTrasactionByHeight error:" + rdb.Error.Error())
	} else {
		err = nil
	}

	return transcations, err
}

func (t *Transaction) UpdateTransactionStatus(db *gorm.DB) (err error) {
	updateinfo := map[string]interface{}{"F_status": t.F_status}
	return t.updateTransactionColumn(db, updateinfo)
}

func (t *Transaction) updateTransactionColumn(db *gorm.DB, updateinfo map[string]interface{}) (err error) {

	log.Debugf("updateTransactionColumn F_id:%d,block:%d,%+v", t.F_id, t.F_block, updateinfo)

	tx := db.Begin()

	rdb := tx.Where("F_id = ?", t.F_id).Model(&t).Update(updateinfo)

	if rdb.Error != nil {
		tx.Rollback()
		return rdb.Error
	}

	tx.Commit()
	return rdb.Error
}

func GetTransactionsCount(db *gorm.DB) (count int64, err error) {
	var rdb *gorm.DB
	num := Count_number{}
	rdb = db.Table("t_transaction").Where("F_status = ?", NORMAL).Select(" count(*) as count ").Find(&num)

	if rdb.Error != nil {
		//panic("find information error:" + rdb.Error.Error())
		err = errors.New("GetTransactionsCount error:" + rdb.Error.Error())
	} else {
		err = nil
	}
	return num.Count, err
}

func GetTransactionsCountByAddr(db *gorm.DB, addr string) (count int64, err error) {
	var rdb *gorm.DB
	num := Count_number{}
	rdb = db.Table("t_transaction").Where("(F_from = ?  or F_to = ?) and F_status = ?", addr, addr, NORMAL).Select(" count(*) as count ").Find(&num)

	if rdb.Error != nil {
		//panic("find information error:" + rdb.Error.Error())
		err = errors.New("GetTransactionsCount error:" + rdb.Error.Error())
	} else {
		err = nil
	}
	return num.Count, err
}

func GetTransactionsCountByHeight(db *gorm.DB, height int64) (count int64, err error) {
	var rdb *gorm.DB
	num := Count_number{}
	rdb = db.Table("t_transaction").Where("F_block = ? and F_status = ?", height, NORMAL).Select(" count(*) as count ").Find(&num)

	if rdb.Error != nil {
		//panic("find information error:" + rdb.Error.Error())
		err = errors.New("GetTransactionsCount error:" + rdb.Error.Error())
	} else {
		err = nil
	}
	return num.Count, err
}

func GetTransactions(db *gorm.DB, offset int, size int) (transList []Transaction, err error) {
	rdb := db.Where("F_status = ?", NORMAL).Order("F_timestamp desc").Offset(offset).Limit(size).Find(&transList)
	if rdb.Error != nil {
		err = errors.New("GetRecentBlocks error:" + rdb.Error.Error())
	} else {
		err = nil
	}

	return transList, err
}

func GetTransactionsByHeight(db *gorm.DB, height int64, offset int, size int) (transList []Transaction, err error) {
	rdb := db.Where("F_status = ? and F_block = ?", NORMAL, height).Order("F_timestamp desc").Offset(offset).Limit(size).Find(&transList)
	if rdb.Error != nil {
		err = errors.New("GetRecentBlocks error:" + rdb.Error.Error())
	} else {
		err = nil
	}

	return transList, err
}

func GetTransactionsByAddr(db *gorm.DB, addr string, offset int, size int) (transList []Transaction, err error) {
	rdb := db.Where("F_status = ? and (F_from = ? or F_to = ?)", NORMAL, addr, addr).Order("F_timestamp desc").Offset(offset).Limit(size).Find(&transList)
	if rdb.Error != nil {
		err = errors.New("GetRecentBlocks error:" + rdb.Error.Error())
	} else {
		err = nil
	}

	return transList, err
}

func GetTransactionsByAddrAndType(db *gorm.DB, addr string, txtype int64, offset int, size int) (transList []Transaction, count int64, err error) {
	if size <= 0 {
		return transList, 0, nil
	}

	var rdb *gorm.DB = nil
	if txtype == TX_TYPE_ALL {
		rdb = db.Where("F_status = ? and (F_from = ? or F_to = ?)", NORMAL, addr, addr)
	} else if txtype == TX_TYPE_ME_MORTGAGE {
		rdb = db.Where("F_status = ? and (F_from = ? or F_to = ?) and F_tx_type=?", NORMAL, addr, addr, TX_TYPE_ME_MORTGAGE)
	} else if txtype == TX_TYPE_ME_REDEEM {
		rdb = db.Where("F_status = ? and (F_from = ? or F_to = ?) and F_tx_type=?", NORMAL, addr, addr, TX_TYPE_ME_REDEEM)
	} else if txtype == TX_TYPE_QUERY_3OR4 {
		rdb = db.Where("F_status = ? and (F_from = ? or F_to = ?) and (F_tx_type = ? or F_tx_type = ?)", NORMAL, addr, addr, TX_TYPE_ME_MORTGAGE, TX_TYPE_ME_REDEEM)
	} else if txtype == TX_TYPE_FROM_ME {
		rdb = db.Where("F_status = ? and (F_from = ?) and F_tx_type=?", NORMAL, addr, TX_TYPE_UNDEFINED)
	} else if txtype == TX_TYPE_TO_ME {
		rdb = db.Where("F_status = ? and (F_to = ?) and F_tx_type=?", NORMAL, addr, TX_TYPE_UNDEFINED)
	} else {
		err = errors.New("Unknow txtype")
		return
	}

	num := Count_number{}
	cdb := rdb.Table("t_transaction")
	cdb = cdb.Select(" count(*) as count ").Find(&num)
	if cdb.Error != nil {
		err = errors.New("GetTransactionsByAddrAndType error:" + cdb.Error.Error())
		return
	}

	rdb = rdb.Order("F_timestamp desc").Offset(offset).Limit(size).Find(&transList)
	if rdb.Error != nil {
		err = errors.New("GetTransactionsByAddrAndType error:" + rdb.Error.Error())
		return
	}

	return transList, num.Count, nil
}

//func (b *Block) GetMaxBlocNumber(db *gorm.DB) (number int64, err error) {
//
//	type Number struct {
//		Number int64 `gorm:"column:max_block"`
//	}
//
//	n := Number{}
//
//	rdb := db.Table(b.TableName()).Where("F_status = ?", NORMAL).Select("MAX(F_block) as max_block").Find(&n)
//	if rdb.RecordNotFound() {
//	} else if rdb.Error != nil {
//		panic("GetMaxBlocNumber error:" + rdb.Error.Error())
//	}
//
//	return n.Number, nil
//}
