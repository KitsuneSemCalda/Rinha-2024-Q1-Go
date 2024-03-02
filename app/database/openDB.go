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
	dsn := utils.GenerateDSN()

	log.Printf("Unsecure show me the DSN: %s", dsn)

	for i := 3; i < 11; i++ {
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

		if err == nil {
			break
		}

		log.Printf("Error: Unknown error opening gorm: %s", err.Error())
		log.Printf("Retrying in %d seconds...", i)
		time.Sleep(time.Duration(i) * time.Second)
	}

	if err != nil {
		log.Fatalf("Error: Failed to connect to database after multiple attempts: %s", err.Error())
	}

	sqlDB, err := DB.DB()

	if err != nil {
		log.Fatalf("Error: Unknown error getting sql.DB: %s", err.Error())
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)
}
