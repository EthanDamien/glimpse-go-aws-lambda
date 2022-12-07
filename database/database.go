package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

// credentials
var (
	host     = os.Getenv("host")
	password = os.Getenv("password")
	port     = "3306"
	user     = "root"
	dbName   = "GLIMPSE"
)

// get the connection to MySQL
func GetConnection() (*sql.DB, error) {
	return sql.Open("mysql", dsn())

}

// get the DSN string
func dsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, password, host, port, dbName)
}
