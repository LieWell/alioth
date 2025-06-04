package models

import (
	"gorm.io/datatypes"
	"liewell.fun/alioth/core"
)

var (
	rplaceTableName = "rplace"
	EmptyRplace     = &Rplace{}
)

type Rplace struct {
	ID   int            `gorm:"primaryKey;autoIncrement"`
	Date string         `gorm:"type:date"` // 创建日期, 格式必须为: YYYY-MM-DD
	Data datatypes.JSON `gorm:"type:json"` // 数据
}

func (u *Rplace) TableName() string {
	return rplaceTableName
}

func FindRplaceByDate(dateNow string) (*Rplace, error) {
	var r Rplace
	if err := core.MYSQL.Where("date = ?", dateNow).First(&r).Error; err != nil {
		return nil, err
	}
	return &r, nil
}

func SaveRplace(r *Rplace) (int, error) {
	err := core.MYSQL.Save(r).Error
	return r.ID, err
}
