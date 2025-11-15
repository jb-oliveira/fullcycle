package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq" // Import the driver, aliased with _ to run its init function
)

func LogError(err error) {
	log.SetPrefix("ERROR: ")                             // Add a prefix to log entries
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile) // Include date, time, and file info
	log.Printf("An error occurred: %v", err)
}

type Product struct {
	ID    string
	Name  string
	Price float64
}

func NewProduct(name string, price float64) *Product {
	return &Product{
		ID:    "",
		Name:  name,
		Price: price,
	}
}

func insertProduct(db *sql.DB, product *Product) (*Product, error) {
	stmt, err := db.Prepare("insert into products(name,price) values($1,$2) returning id;")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(product.Name, product.Price).Scan(&product.ID)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func deleteProduct(db *sql.DB, id string) error {
	stmt, err := db.Prepare("delete from products where id = $1")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(id)
	return err
}

func getAllProducts(db *sql.DB) ([]Product, error) {
	// não precisa de statement por que não tem parametros
	rows, err := db.Query("select id,name,price from products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		err = rows.Scan(&p.ID, &p.Name, &p.Price)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func getProductByID(db *sql.DB, id string) (*Product, error) {
	stmt, err := db.Prepare("select name,price from products where id =$1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	var p Product
	err = stmt.QueryRow(id).Scan(&p.Name, &p.Price)
	if err != nil {
		return nil, err
	}
	p.ID = id
	return &p, nil
}

func updateProduct(db *sql.DB, product *Product) error {
	stmt, err := db.Prepare("update products set name = $1, price = $2 where id = $3")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(product.Name, product.Price, product.ID)
	if err != nil {
		return err
	}
	return nil
}

func main() {

	connStr := "user=postgres password=password dbname=myapp host=localhost port=5432 sslmode=disable"

	// Alternatively, you can use the URL format:
	// connStr := "postgres://postgres:yourpassword@localhost:5432/testdb?sslmode=disable"

	// 2. Open the connection
	// sql.Open() establishes the connection pool but DOES NOT actually connect to the database.
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// check if connection is right
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	product, err := insertProduct(db, NewProduct("inserido", 210.70))
	if err != nil {
		LogError(err)
		return
	}
	// // fmt.Printf("Produto: %v\n", product)
	product.Name = "alterado"
	err = updateProduct(db, product)
	if err != nil {
		LogError(err)
		return
	}

	newProd, err := getProductByID(db, product.ID)
	if err != nil {
		LogError(err)
		return
	}
	// fmt.Printf("Produto: %v\n", newProd)

	products, err := getAllProducts(db)
	if err != nil {
		LogError(err)
		return
	}
	for _, p := range products {
		fmt.Printf("Produto: %v\n", p)
	}

	err = deleteProduct(db, newProd.ID)
	if err != nil {
		LogError(err)
		return
	}

}
