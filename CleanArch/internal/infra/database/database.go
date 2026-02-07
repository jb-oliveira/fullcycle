package database

import (
	"database/sql"

	config "github.com/jb-oliveira/fullcycle/CleanArch/configs"
	_ "github.com/lib/pq"
)

func NewDatabase(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open(cfg.DBDriver, cfg.DBSource)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS PRODUCTS(
			ID VARCHAR(36) PRIMARY KEY,
			NAME VARCHAR(255),
			PRICE DECIMAL(10,2)
		);
	`)
	if err != nil {
		return nil, err
	}

	return db, nil
}
