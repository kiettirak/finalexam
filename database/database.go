package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

func Conn() *sql.DB {

	if db != nil {
		return db
	}

	var err error
	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("can't connect database : ", err)
	}

	return db
}
func InsertCustomer(name, email, status string) *sql.Row {
	return Conn().QueryRow("INSERT INTO customers(name, email, status) VALUES ($1, $2, $3) RETURNING id", name, email, status)
}
