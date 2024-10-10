package postgres

import (
	"context"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"gitlab.com/sunflower4026/interview-flash-mobile/toolkit/db"
)

func NewPostgresDatabase(opt *db.Option) (*gorm.DB, error) {
	connURL := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		opt.Host, opt.Port, opt.Username, opt.Password, opt.DatabaseName)

	// Open the GORM database connection
	gormDB, err := gorm.Open(postgres.Open(connURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Configure GORM logger (optional)
	})
	if err != nil {
		return nil, err
	}

	// Configure connection pool settings
	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(opt.ConnectionOption.MaxIdle)
	sqlDB.SetConnMaxLifetime(opt.ConnectionOption.MaxLifetime)
	sqlDB.SetMaxOpenConns(opt.ConnectionOption.MaxOpen)

	// Check the connection
	ctx, cancel := context.WithTimeout(context.Background(), opt.ConnectionOption.ConnectTimeout)
	defer cancel()

	if err = sqlDB.PingContext(ctx); err != nil {
		return nil, err
	}

	log.Println("successfully connected to postgres", opt.Host)

	return gormDB, nil
}
