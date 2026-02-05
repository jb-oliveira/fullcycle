package main

import (
	"database/sql"
	"log"

	"github.com/jb-oliveira/fullcycle/DI/product"
	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", "postgres://postgres:password@localhost:5432/myapp?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.Exec(`
	CREATE TABLE IF NOT EXISTS PRODUCTS(
		ID VARCHAR(36) PRIMARY KEY,
		NAME VARCHAR(255),
		PRICE DECIMAL(10,2)
	);
	`)

	product := &product.ProductInputDto{
		ID:    "3",
		Name:  "Product 3",
		Price: 30.0,
	}

	useCase := NewInsertProductUseCase(db)
	err = useCase.Execute(product)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Product inserted successfully")
}
