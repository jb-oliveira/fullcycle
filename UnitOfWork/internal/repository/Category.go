package repository

import (
	"context"
	"database/sql"

	"github.com/jb-oliveira/fullcycle/UnitOfWork/internal/db"
	"github.com/jb-oliveira/fullcycle/UnitOfWork/internal/entity"
)

type CategoryRepository interface {
	Create(ctx context.Context, category *entity.Category) error
	FindByID(ctx context.Context, id string) (*entity.Category, error)
}

type CategoryRepositoryPGImpl struct {
	queries *db.Queries
}

func NewCategoryRepositoryPGImpl(dbTX db.DBTX) *CategoryRepositoryPGImpl {
	return &CategoryRepositoryPGImpl{queries: db.New(dbTX)}
}

func (r *CategoryRepositoryPGImpl) Create(ctx context.Context, category *entity.Category) error {
	return r.queries.CreateCategory(ctx, db.CreateCategoryParams{
		ID:          category.ID,
		Name:        category.Name,
		Description: sql.NullString{String: category.Description, Valid: true},
	})
}

func (r *CategoryRepositoryPGImpl) FindByID(ctx context.Context, id string) (*entity.Category, error) {
	category, err := r.queries.FindCategoryByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &entity.Category{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description.String,
	}, nil
}
