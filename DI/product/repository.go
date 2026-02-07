package product

import "database/sql"

type ProductRepository interface {
	Insert(product *Product) error
}

type productRepositoryPostgres struct {
	db *sql.DB
}

func NewProductRepositoryPostgres(db *sql.DB) ProductRepository {
	return &productRepositoryPostgres{db: db}
}

func (r *productRepositoryPostgres) Insert(product *Product) error {
	_, err := r.db.Exec("INSERT INTO products (id, name, price) VALUES ($1, $2, $3)", product.ID, product.Name, product.Price)
	if err != nil {
		return err
	}
	return nil
}
