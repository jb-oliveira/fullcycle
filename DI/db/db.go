package db

import (
	"database/sql"
	"log"
)

func NewDB() *sql.DB {
	db, err := sql.Open("postgres", "postgres://postgres:password@localhost:5432/myapp?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS PRODUCTS(
			ID VARCHAR(36) PRIMARY KEY,
			NAME VARCHAR(255),
			PRICE DECIMAL(10,2)
		);
		`)
	if err != nil {
		log.Fatal(err)
	}
	return db
}
