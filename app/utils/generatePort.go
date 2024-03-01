package utils

import (
	"log"
	"os"
)

var port string

func init() {
	os_port := os.Getenv("PORT")

	if os_port == "" {
		log.Fatalln("Error: missing enviroment variable PORT")
	}

	port = ":" + os_port
}

func GetPort() string {
	return port
}
