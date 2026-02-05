//go:build wireinject
// +build wireinject

package main

import (
	"database/sql"

	"github.com/google/wire"
	"github.com/jb-oliveira/fullcycle/DI/product"
)

var setRepositoryDependency = wire.NewSet(
	product.NewProductRepositoryPostgres,
	product.NewProductUseCaseImpl,
)

func NewInsertProductUseCase(db *sql.DB) product.InsertProductUseCase {
	wire.Build(setRepositoryDependency)
	return &product.ProductUseCaseImpl{}
}
