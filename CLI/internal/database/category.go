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

func (c *Category) Create(name string, description string) error {
	stmt, err := c.db.Prepare("INSERT INTO categories (id, name, description) VALUES ($1, $2, $3)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	uuid := uuid.Must(uuid.NewV7()).String()
	_, err = stmt.Exec(uuid, name, description)
	if err != nil {
		return err
	}
	c.Id = uuid
	c.Name = name
	c.Description = description
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

func (c *Category) FindByCourseID(courseID string) (Category, error) {
	var category Category
	row := c.db.QueryRow("SELECT c.id, c.name, c.description FROM categories c JOIN courses cc ON c.id = cc.category_id WHERE cc.id = $1", courseID)
	err := row.Scan(&category.Id, &category.Name, &category.Description)
	if err != nil {
		return Category{}, err
	}
	return category, nil
}
