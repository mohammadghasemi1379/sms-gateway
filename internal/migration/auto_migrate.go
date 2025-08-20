package migration

import (
	"context"

	"github.com/mohammadghasemi1379/sms-gateway/internal/entity"
	"github.com/mohammadghasemi1379/sms-gateway/pkg/logger"
	"gorm.io/gorm"
)

// AutoMigrateEntities runs GORM's AutoMigrate for all entities
// This is an alternative to SQL file-based migrations
func AutoMigrateEntities(ctx context.Context, db *gorm.DB, logger *logger.Logger) error {
	logger.Info(ctx, "Starting GORM AutoMigrate for entities")

	// List all your entities here
	entities := []interface{}{
		&entity.User{},
		&entity.SMS{},
		&entity.Transaction{},
	}

	if err := db.AutoMigrate(entities...); err != nil {
		logger.Error(ctx, "Failed to auto-migrate entities", "error", err.Error())
		return err
	}

	logger.Info(ctx, "GORM AutoMigrate completed successfully")
	return nil
}
