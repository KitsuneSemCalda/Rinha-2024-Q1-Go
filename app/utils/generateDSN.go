package utils

var dsn string 

func init() {
	dsn = "host=db user=postgres password=postgres dbname=rinha sslmode=disable"
}

func GetDSN() string {
	return dsn
}
