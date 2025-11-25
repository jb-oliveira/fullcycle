package database

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/jb-oliveira/fullcycle/tree/main/APIS/internal/entity"
	pkgEntity "github.com/jb-oliveira/fullcycle/tree/main/APIS/pkg/entity"
)

type Product struct {
	db *gorm.DB
}

func NewProductDB(db *gorm.DB) *Product {
	return &Product{db: db}
}

func (p *Product) Create(product *entity.Product) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
	defer cancel()
	return gorm.G[entity.Product](p.db).Create(ctx, product)
}

func (p *Product) FindAll(page, limit int, sort string) ([]entity.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
	defer cancel()
	offset := (page - 1) * limit
	return gorm.G[entity.Product](p.db).
		Order(sort).
		Limit(limit).
		Offset(offset).
		Find(ctx)
}

func (p *Product) FindByID(id string) (*entity.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
	defer cancel()

	productID, err := pkgEntity.ParseID(id)
	if err != nil {
		return nil, err
	}

	product, err := gorm.G[entity.Product](p.db).Where("prd_id = ?", productID).First(ctx)
	return &product, err
}

func (p *Product) Update(product *entity.Product) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
	defer cancel()
	_, err := gorm.G[entity.Product](p.db).Updates(ctx, *product)
	return err
}

func (p *Product) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
	defer cancel()

	productID, err := pkgEntity.ParseID(id)
	if err != nil {
		return err
	}

	_, err = gorm.G[entity.Product](p.db).Where("prd_id = ?", productID).Delete(ctx)
	return err
}

func (p *Product) Count() (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
	defer cancel()

	return gorm.G[entity.Product](p.db).Count(ctx, "prd_id")
}
