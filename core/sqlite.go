package core

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var SQLite *gorm.DB

func InitSqlite() {

	cfg := GlobalConfig.Sqlite

	// gorm 控制项
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn), // 日志级别
	}

	db, err := gorm.Open(sqlite.Open(cfg.Path), gormConfig)
	if err != nil {
		Logger.Fatalf("[sqlite] connected error: %v\n", err)
	}

	// SQLite 单文件数据库,串行化使用避免锁竞争
	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(1)

	SQLite = db
	Logger.Infof("[sqlite] connected success: %v\n", cfg.Path)
}