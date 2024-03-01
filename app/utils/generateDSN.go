package utils

import (
	"fmt"
	"log"
	"os"
)

var dsn string = "host=db user=postgres password=postgres dbname=rinha sslmode=disable"

func init() {
	if dsn == "" {
		host := os.Getenv("DB_HOST")

		if host == "" {
			log.Fatalln("Error: missing enviroment variable: DB_HOST")
		}

		user := os.Getenv("POSTGRES_USER")

		if user == "" {
			log.Fatalln("Error: missing enviroment variable: POSTGRES_USER")
		}

		password := os.Getenv("POSTGRES_PASSWORD")

		if password == "" {
			log.Fatalln("Error: missing enviroment variable: POSTGRES_PASSWORD")
		}

		dbname := os.Getenv("POSTGRES_DB")

		if dbname == "" {
			log.Fatalln("Error: missing enviroment variable: POSTGRES_DB")
		}

		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
			host, user, password, dbname)
	}
}

func GetDSN() string {
	return dsn
}
