package utils

import (
	"fmt"
	"os"
)

var dsn string

func GenerateDSN() string {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")

	if !(host == "") && !(user == "") && !(password == "") && !(dbname == "") {
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
			host, user, password, dbname)
	} else {
		dsn = "host=db user=postgres password=postgres dbname=rinha sslmode=disable"
	}
}
