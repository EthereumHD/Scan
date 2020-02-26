package model

import (
	_ "fmt"
	"github.com/jinzhu/gorm"
	"time"
)

type Rate struct {
	F_id          uint64  `gorm:"column:F_id"` //ID
	F_eth         float64 `gorm:"column:F_eth"`
	F_btc         float64 `gorm:"column:F_btc"`
	F_usd         float64 `gorm:"column:F_usd"`
	F_kwr         float64 `gorm:"column:F_kwr"`
	F_timestamp   int64   `gorm:"column:F_timestamp"`
	F_create_time string  `gorm:"column:F_create_time"` //创建时间
	F_modify_time string  `gorm:"column:F_modify_time"` //修改时间
}


func (r *Rate) TableName() string {
	return "t_rate"
}

func (r *Rate) BeforeCreate(scope *gorm.Scope) error {
	currentTime := time.Now().Local()
	newFormat := currentTime.Format("2006-01-02 15:04:05.000")

	scope.SetColumn("F_timestamp", time.Now().Unix())
	scope.SetColumn("F_create_time", newFormat)
	scope.SetColumn("F_modify_time", newFormat)

	return nil
}

func (r *Rate) BeforeUpdate(scope *gorm.Scope) error {
	currentTime := time.Now().Local()
	newFormat := currentTime.Format("2006-01-02 15:04:05.000")
	scope.SetColumn("F_modify_time", newFormat)
	return nil
}

func (r *Rate) CreateRate(db *gorm.DB) (err error) {
	rdb := db.Create(&r)

	return rdb.Error
}
