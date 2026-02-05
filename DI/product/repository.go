package product

import "database/sql"

type ProductRepository interface {
	Insert(product *Product) error
}

type ProductRepositoryPostgres struct {
	DB *sql.DB
}

func NewProductRepositoryPostgres(db *sql.DB) ProductRepository {
	return &ProductRepositoryPostgres{DB: db}
}

func (r *ProductRepositoryPostgres) Insert(product *Product) error {
	_, err := r.DB.Exec("INSERT INTO products (id, name, price) VALUES ($1, $2, $3)", product.ID, product.Name, product.Price)
	if err != nil {
		return err
	}
	return nil
}
