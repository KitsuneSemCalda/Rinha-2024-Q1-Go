package database

import (
	"RinhaBackend/app/utils"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var err error

func StartDB() {
	time.Sleep(3)

	dsn := utils.GenerateDSN()

	log.Printf("Unsecure show me the DSN: %s", dsn)

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
	sqlDB.SetConnMaxLifetime(5 * time.Minute)
}
