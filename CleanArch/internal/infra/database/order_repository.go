package database

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jb-oliveira/fullcycle/CleanArch/internal/entity"
)

type OrderRepositoryPG struct {
	db *sql.DB
}

func NewOrderRepositoryPG(db *sql.DB) entity.OrderRepository {
	return &OrderRepositoryPG{db: db}
}

func (r *OrderRepositoryPG) Save(order *entity.Order) error {
	stmt, err := r.db.Prepare("INSERT INTO orders (id, price, tax) VALUES ($1, $2, $3)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(order.ID, order.Price, order.Tax)
	return err
}

// func (r *OrderRepositoryPG) FindAll(page, limit int, sort string) ([]entity.Order, error) {
// 	stmt, err := r.db.Prepare("SELECT id, price, tax FROM orders order by $1 asc limit $2 offset $3")
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer stmt.Close()
// 	rows, err := stmt.Query(sort, limit, page)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()
// 	var orders []entity.Order
// 	for rows.Next() {
// 		var order entity.Order
// 		err = rows.Scan(&order.ID, &order.Price, &order.Tax)
// 		if err != nil {
// 			return nil, err
// 		}
// 		orders = append(orders, order)
// 	}
// 	return orders, nil
// }

func (r *OrderRepositoryPG) FindAll(page, limit int, sort, sortDir string) ([]entity.Order, error) {
	// 1. Validar o sort para evitar SQL Injection (Importante!)
	allowedSorts := map[string]bool{"id": true, "price": true, "tax": true}
	if !allowedSorts[sort] {
		return nil, errors.New("invalid sort")
	}

	allowedSortDirs := map[string]bool{"ASC": true, "DESC": true}
	direction := strings.ToUpper(sortDir)
	if !allowedSortDirs[direction] {
		return nil, errors.New("invalid sort direction")
	}

	// 2. Calcular o offset corretamente
	offset := (page - 1) * limit

	// 3. Montar a query com o sort injetado (placeholders não funcionam aqui)
	query := fmt.Sprintf("SELECT id, price, tax FROM orders ORDER BY %s %s LIMIT $1 OFFSET $2", sort, direction)

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// 4. Passar apenas limit e offset como parâmetros
	rows, err := stmt.Query(limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []entity.Order
	for rows.Next() {
		var order entity.Order
		if err := rows.Scan(&order.ID, &order.Price, &order.Tax); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}
