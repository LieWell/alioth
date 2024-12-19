package models

import (
	"time"

	"liewell.fun/alioth/core"
)

var (
	userTableName = "user"
	EmptyUser     = &User{}

	// 用户状态
	UserStatusEnable  = 0
	UserStatusDisable = 1
	UserStatusDelete  = 2
)

type User struct {
	ID        int       `gorm:"primaryKey;autoIncrement"`
	Username  string    `gorm:"column:user_name;type:varchar(32);not null;index"`
	NickName  string    `gorm:"type:varchar(64);not null"`
	Email     string    `gorm:"type:varchar(255);not null"`
	Password  string    `gorm:"type:varchar(128);not null"`
	Salt      string    `gorm:"type:varchar(32);not null"`
	Status    int       `gorm:"type:tinyint(1);default:0;not null"` // 0:启用, 1:禁用, 2:删除
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (u *User) TableName() string {
	return userTableName
}

func FindUserByUsername(username string) (*User, error) {
	var u User
	if err := core.MYSQL.Where("user_name = ?", username).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func SaveUser(u *User) (int, error) {
	err := core.MYSQL.Save(u).Error
	return u.ID, err
}
