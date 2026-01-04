package database

import (
	"database/sql"

	"github.com/google/uuid"
)

type Category struct {
	db          *sql.DB
	Id          string
	Name        string
	Description string
}

func NewCategory(db *sql.DB) *Category {
	return &Category{db: db}
}

func (c *Category) Create() error {
	stmt, err := c.db.Prepare("INSERT INTO categories (id, name, description) VALUES ($1, $2, $3)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	uuid := uuid.Must(uuid.NewV7()).String()
	c.Id = uuid
	_, err = stmt.Exec(c.Id, c.Name, c.Description)
	if err != nil {
		return err
	}

	return nil
}

func (c *Category) FindAll() ([]Category, error) {
	rows, err := c.db.Query("SELECT id, name, description FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var category Category
		err := rows.Scan(&category.Id, &category.Name, &category.Description)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}
