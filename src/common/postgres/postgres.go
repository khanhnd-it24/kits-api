package apppostgres

import (
	"context"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"kits/api/src/common/configs"
)

type DBProvider struct {
	db *gorm.DB
}

func NewPostgresProvider(config *configs.Config) (*DBProvider, error) {
	cf := config.Postgresql
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s password=%s", cf.Host,
		cf.Port, cf.User, cf.DbName, cf.SslMode, cf.Password)
	logMode := logger.Info
	if config.Mode.IsProd() {
		logMode = logger.Silent
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logMode),
	})

	if err != nil {
		return nil, fmt.Errorf("[Postgres] failed to connect to DB %w", err)
	}

	return &DBProvider{
		db: db,
	}, nil
}

func (provider *DBProvider) Stop(ctx context.Context) error {
	db, err := provider.db.DB()
	if err != nil {
		return fmt.Errorf("[Postgres] failed to disconnect %w", err)
	}
	return db.Close()
}

func (provider *DBProvider) DB() *gorm.DB {
	return provider.db
}
