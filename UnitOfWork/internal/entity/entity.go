package entity

import "errors"

type Category struct {
	ID          string
	Name        string
	Description string
	Courses     []Course
}

func (c *Category) Validate() error {
	if c.Name == "" {
		return errors.New("Name is required")
	}
	if c.Description == "" {
		return errors.New("Description is required")
	}
	return nil
}

func (c *Category) AddCourse(course Course) error {
	c.Courses = append(c.Courses, course)
	return nil
}

type Course struct {
	ID          string
	Category    Category
	Name        string
	Description string
	Price       float64
}
