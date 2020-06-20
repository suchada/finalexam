package database

import (
	"database/sql"
	"log"
	"os"
	// "fmt"

	_ "github.com/lib/pq"
)

var db *sql.DB

func init() {

	url := os.Getenv("DATABASE_URL")
	if len(url) == 0 {
		url = "postgres://hvbetmpa:k2OW2QByCH97dBP2c9fqcBZp3YiDpTeY@john.db.elephantsql.com:5432/hvbetmpa"
	}

	var error error
	db, error = sql.Open("postgres", url)
	if error != nil {
		log.Fatal(error)
	}
	createTb := `CREATE TABLE IF NOT EXISTS customer (
		id SERIAL PRIMARY KEY,
		name TEXT,
		email TEXT,
		status TEXT
	   );`
	  
	   _, error = db.Exec(createTb)
	   if error != nil {
		log.Fatal("Cannot create table to database:", error)
	   }
}

//Conn...
func Conn() *sql.DB {
	return db
}
