package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	"github.com/jb-oliveira/fullcycle/UnitOfWork/internal/db"
	"github.com/jb-oliveira/fullcycle/UnitOfWork/internal/entity"
)

type CourseRepository interface {
	Create(ctx context.Context, course *entity.Course) error
	FindByID(ctx context.Context, id string) (*entity.Course, error)
}

type CourseRepositoryPGImpl struct {
	queries *db.Queries
}

func NewCourseRepositoryPGImpl(conn *sql.DB) *CourseRepositoryPGImpl {
	return &CourseRepositoryPGImpl{queries: db.New(conn)}
}

func (r *CourseRepositoryPGImpl) Create(ctx context.Context, course *entity.Course) error {
	price := fmt.Sprintf("%0.2f", course.Price)
	return r.queries.CreateCourse(ctx, db.CreateCourseParams{
		ID:          course.ID,
		CategoryID:  course.Category.ID,
		Name:        course.Name,
		Description: sql.NullString{String: course.Description, Valid: true},
		Price:       price,
	})
}

func (r *CourseRepositoryPGImpl) FindByID(ctx context.Context, id string) (*entity.Course, error) {
	course, err := r.queries.FindCourseByID(ctx, id)
	if err != nil {
		return nil, err
	}
	price, err := strconv.ParseFloat(course.Price, 64)
	if err != nil {
		return nil, err
	}
	return &entity.Course{
		ID:          course.ID,
		Category:    entity.Category{ID: course.CategoryID},
		Name:        course.Name,
		Description: course.Description.String,
		Price:       price,
	}, nil
}
