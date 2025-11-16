package main

import (
	"context"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/plugin/optimisticlock"
)

type Category struct {
	ID       int `gorm:"primaryKey"`
	Name     string
	Products []Product `gorm:"many2many:products_categories;"`
}

type Product struct {
	ID           int `gorm:"primaryKey"`
	Name         string
	Price        float64
	Categories   []Category `gorm:"many2many:products_categories;"`
	SerialNumber SerialNumber
}

type SerialNumber struct {
	ID        int `gorm:"primaryKey"`
	Number    string
	ProductID int
}

type User struct {
	ID      int
	Name    string
	Age     uint
	Version optimisticlock.Version
}

func main() {
	_, err := executaTudo()
	if err != nil {
		panic(err)
	}
}

func executaTudo() (*Product, error) {
	connStr := "user=postgres password=password dbname=myapp host=localhost port=5432 sslmode=disable"

	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&Category{}, &Product{}, &SerialNumber{})

	ctx := context.Background()

	// result, err := criaValores(ctx, err, db)
	// if err != nil {
	// 	return result, err
	// }

	// result, err := exibeValor(ctx, db)
	// if err != nil {
	// 	return result, err
	// }

	_, err = lockPessimista(db, ctx)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func lockPessimista(db *gorm.DB, ctx context.Context) (any, error) {

	tx := db.Begin()

	var c Category
	err := tx.Debug().Clauses(clause.Locking{Strength: "UPDATE"}).First(&c, nil).Error
	if err != nil {
		return nil, err
	}
	c.Name = "Teste Transaction"
	tx.Debug().Save(&c)
	tx.Commit()

	return nil, nil
}

func exibeValor(ctx context.Context, db *gorm.DB) (*Product, error) {
	cats, err := gorm.G[Category](db).
		Preload("Products.SerialNumber", nil).
		Find(ctx)
	if err != nil {
		return nil, err
	}

	for _, cat := range cats {
		fmt.Printf("Categoria %s:\n", cat.Name)
		for _, prod := range cat.Products {
			fmt.Printf(" - Produto: %s, Serial: %s  \n", prod.Name, prod.SerialNumber.Number)
		}
	}
	return nil, nil
}

func criaValores(ctx context.Context, err error, db *gorm.DB) (*Product, error) {
	categoria := Category{Name: "Cozinha"}
	err = gorm.G[Category](db).Create(ctx, &categoria)
	if err != nil {
		return nil, err
	}
	categoria2 := Category{Name: "Eletronico"}
	err = gorm.G[Category](db).Create(ctx, &categoria2)
	if err != nil {
		return nil, err
	}

	produto := Product{
		Name:       "Panela",
		Price:      99.4,
		Categories: []Category{categoria, categoria2},
	}
	err = gorm.G[Product](db).Create(ctx, &produto)
	if err != nil {
		return nil, err
	}

	serial := SerialNumber{Number: "1234", ProductID: produto.ID}
	err = gorm.G[SerialNumber](db).Create(ctx, &serial)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
