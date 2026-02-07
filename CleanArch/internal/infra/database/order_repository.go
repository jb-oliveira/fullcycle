package database

import (
	"database/sql"

	"github.com/jb-oliveira/fullcycle/CleanArch/internal/entity"
)

type orderRepositoryPG struct {
	db *sql.DB
}

func NewOrderRepositoryPG(db *sql.DB) entity.OrderRepository {
	return &orderRepositoryPG{db: db}
}

func (r *orderRepositoryPG) Save(order *entity.Order) error {
	stmt, err := r.db.Prepare("INSERT INTO orders (id, price, tax) VALUES ($1, $2, $3)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(order.ID, order.Price, order.Tax)
	return err
}
