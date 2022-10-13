package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var (
	host     = os.Getenv("host")
	password = os.Getenv("password")
	port     = "3306"
	user     = "root"
	dbName   = "GLIMPSE"
)

func GetConnection() (*sql.DB, error) {
	return sql.Open("mysql", dsn())

}

func dsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, password, host, port, dbName)
}
