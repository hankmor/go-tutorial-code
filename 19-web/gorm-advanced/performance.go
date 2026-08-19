package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// CreateAuthorsInBatches 使用分批写入减少单条 INSERT 的往返次数。
func CreateAuthorsInBatches(db *gorm.DB, authors []Author, batchSize int) error {
	if batchSize <= 0 {
		return fmt.Errorf("batch size must be positive")
	}
	return db.CreateInBatches(authors, batchSize).Error
}

// FindAuthorsWithPosts 使用预加载批量读取作者及其文章。
func FindAuthorsWithPosts(db *gorm.DB) ([]Author, error) {
	var authors []Author
	err := db.Preload("Posts").Order("id ASC").Find(&authors).Error
	return authors, err
}

// ConfigureConnectionPool 设置 database/sql 连接池参数。
func ConfigureConnectionPool(db *gorm.DB, maxOpen, maxIdle int, lifetime, idleTime time.Duration) (*sql.DB, error) {
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(maxOpen)
	sqlDB.SetMaxIdleConns(maxIdle)
	sqlDB.SetConnMaxLifetime(lifetime)
	sqlDB.SetConnMaxIdleTime(idleTime)
	return sqlDB, nil
}

// NewSlowQueryLogger 创建按耗时记录 SQL 的 GORM Logger。
func NewSlowQueryLogger(slowThreshold time.Duration) logger.Interface {
	return logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             slowThreshold,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)
}
