package database

import (
	"RinhaBackend/app/utils"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	time.Sleep(5 * time.Second)

	var err error

	dsn := utils.GetDSN()

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Error: Unknown error opening gorm: %s", err.Error())
	}

	sqlDB, err := DB.DB()

	if err != nil {
		log.Fatalf("Error: Unknown error getting sql.DB: %s", err.Error())
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(10 * time.Minute)
}
