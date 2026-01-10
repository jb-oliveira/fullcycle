package database

import (
	"database/sql"

	"github.com/google/uuid"
)

type CategoryDB struct {
	db          *sql.DB
	Id          string
	Name        string
	Description string
}

func NewCategory(db *sql.DB) *CategoryDB {
	return &CategoryDB{db: db}
}

func (c *CategoryDB) Create() error {
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

func (c *CategoryDB) FindAll() ([]CategoryDB, error) {
	rows, err := c.db.Query("SELECT id, name, description FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []CategoryDB
	for rows.Next() {
		var category CategoryDB
		err := rows.Scan(&category.Id, &category.Name, &category.Description)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func (c *CategoryDB) FindByCourseID(courseID string) (CategoryDB, error) {
	var category CategoryDB
	row := c.db.QueryRow("SELECT c.id, c.name, c.description FROM categories c JOIN courses cc ON c.id = cc.category_id WHERE cc.id = $1", courseID)
	err := row.Scan(&category.Id, &category.Name, &category.Description)
	if err != nil {
		return CategoryDB{}, err
	}
	return category, nil
}

func (c *CategoryDB) FindById(id string) (CategoryDB, error) {
	var category CategoryDB
	row := c.db.QueryRow("SELECT id, name, description FROM categories WHERE id = $1", id)
	err := row.Scan(&category.Id, &category.Name, &category.Description)
	if err != nil {
		return CategoryDB{}, err
	}
	return category, nil
}
