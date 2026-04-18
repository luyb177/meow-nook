package database

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type MySQLClient struct {
	DB *gorm.DB
}

// NewMySQLClient
// dsn user:password@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Asia%2FShanghai&timeout=10s&readTimeout=30s&writeTimeout=30s
func NewMySQLClient(dsn string) (*MySQLClient, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		PrepareStmt:            true, // 启用预编译语句
		SkipDefaultTransaction: true, // 提升性能
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold: time.Millisecond * 200,
				LogLevel:      logger.Warn,
				Colorful:      true,
			},
		), // 开启慢查询日志记录
	})
	if err != nil {
		return nil, err
	}

	// 获取底层 sql.DB
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// 连接池配置
	sqlDB.SetMaxOpenConns(50)                  // 最大连接数
	sqlDB.SetMaxIdleConns(10)                  // 最大空闲连接
	sqlDB.SetConnMaxLifetime(time.Hour)        // 连接最大存活时间
	sqlDB.SetConnMaxIdleTime(30 * time.Minute) // 空闲最大时间

	return &MySQLClient{DB: db}, nil
}
