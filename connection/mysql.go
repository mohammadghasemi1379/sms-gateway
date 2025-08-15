package connection

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/mohammadghasemi1379/sms-gateway/config"
	"github.com/mohammadghasemi1379/sms-gateway/pkg/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

func MysqlConnection(ctx context.Context, logger *logger.Logger, cfg *config.Config) (*gorm.DB, *sql.DB) {
	mc := cfg.Database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		mc.User, mc.Password, mc.Host, mc.Port, mc.DBName,
	)

	logLevel := gormLogger.Info
	if cfg.App.IsProduction() {
		logLevel = gormLogger.Silent
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger:      gormLogger.Default.LogMode(logLevel),
		PrepareStmt: true,
	})
	if err != nil {
		logger.Panic(ctx, "failed to connect to core Mysql", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.Panic(ctx, "failed to get sqlDB", err)
	}

	sqlDB.SetMaxIdleConns(40)
	sqlDB.SetMaxOpenConns(60)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	logger.Info(ctx, "Connected to Mysql")

	return db, sqlDB
}
