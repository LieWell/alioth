package core

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var MYSQL *gorm.DB

func InitMysql() {

	cfg := GlobalConfig.Mysql
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database)

	// mysql 控制项
	mysqlConfig := mysql.Config{
		DSN:                       dsn,
		DefaultStringSize:         256, // 默认字符串长度
		SkipInitializeWithVersion: false,
	}

	// gorm 控制项
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn), // 日志级别
	}

	db, err := gorm.Open(mysql.New(mysqlConfig), gormConfig)
	if err != nil {
		Logger.Fatalf("[mysql] connected error: %v\n", err)
	}

	// 连接池
	// TODO 这样设置有效吗?
	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConnections)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConnections)

	MYSQL = db
	Logger.Infof("[mysql] connected success: %v\n", dsn)
}
