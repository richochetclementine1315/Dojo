package database

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB is the global database instance
var DB *gorm.DB

// Connect function establishes a connection to the PostgreSQL uisng GORM.
func Connect(dsn string, isDebug bool) (*gorm.DB, error) {
	var err error

	// Configgure GORM logger based on isDebug flag
	logLevel := logger.Silent
	if isDebug {
		logLevel = logger.Info
	}
	// Connect to PostgreSQL database
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	// Get underlying sql.DB to configure connection pool
	sqlDB, err := DB.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance:%w", err)
	}
	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("✅ Database connection established")
	return DB, nil
}

// AutoMigrate function performs GORM automatic migration for the provided models.
func AutoMigrate(models ...interface{}) error {
	err := DB.AutoMigrate(models...)
	if err != nil {
		return fmt.Errorf("auto migration failed: %w", err)
	}
	log.Println("✅ Database migration completed")
	return nil

}

// Close function closes the database connection.
func Close() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
