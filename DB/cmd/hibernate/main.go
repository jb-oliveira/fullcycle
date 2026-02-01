package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jb-oliveira/fullcycle/DB/internal/db"
)

type Category struct {
	ID          string
	Name        string
	Description string
	Courses     []Course
}

type Course struct {
	ID          string
	Category    Category
	Name        string
	Description string
	Price       float64
}

type CategoryRepository struct {
	queries *db.Queries
}

func NewCategoryRepository(queries *db.Queries) *CategoryRepository {
	return &CategoryRepository{queries: queries}
}

func (r *CategoryRepository) ListCategoriesWithCourses(ctx context.Context) ([]Category, error) {
	rows, err := r.queries.ListCoursesWithCategory(ctx)
	if err != nil {
		return nil, err
	}
	var categoryMap = make(map[string]*Category)
	for _, row := range rows {
		category, ok := categoryMap[row.CategoryID]
		if !ok {
			category = &Category{
				ID:          row.CategoryID,
				Name:        row.CategoryName,
				Description: row.CategoryDescription.String,
			}
			categoryMap[row.CategoryID] = category
		}
		course := Course{
			ID:          row.ID,
			Category:    *category,
			Name:        row.Name,
			Description: row.Description.String,
			Price:       float64(row.Price.Int.Int64()) / 100,
		}
		category.Courses = append(category.Courses, course)
	}

	// Convert map to slice after all courses are added
	var categories []Category
	for _, category := range categoryMap {
		categories = append(categories, *category)
	}

	return categories, nil
}

func main() {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, "postgres://postgres:password@localhost:5432/myapp")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(ctx)

	queries := db.New(conn)

	categoryRepo := NewCategoryRepository(queries)

	categories, err := categoryRepo.ListCategoriesWithCourses(ctx)
	if err != nil {
		log.Fatal(err)
	}

	for _, category := range categories {
		log.Println(category)
	}
}
