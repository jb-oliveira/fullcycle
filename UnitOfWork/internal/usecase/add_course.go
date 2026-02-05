package usecase

import (
	"context"

	"github.com/jb-oliveira/fullcycle/UnitOfWork/internal/entity"
	"github.com/jb-oliveira/fullcycle/UnitOfWork/internal/repository"
)

type AddCourseInputDTO struct {
	ID                  string  `json:"id"`
	CategoryID          string  `json:"category_id"`
	CategoryName        string  `json:"category_name"`
	CategoryDescription string  `json:"category_description"`
	Name                string  `json:"name"`
	Description         string  `json:"description"`
	Price               float64 `json:"price"`
}

type AddCourseOutputDTO struct {
	ID          string  `json:"id"`
	CategoryID  string  `json:"category_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

type AddCourseUseCase interface {
	Execute(ctx context.Context, input AddCourseInputDTO) (AddCourseOutputDTO, error)
}

type AddCourseUseCaseImpl struct {
	CategoryRepository repository.CategoryRepository
	CourseRepository   repository.CourseRepository
}

func NewAddCourseUseCaseImpl(categoryRepository repository.CategoryRepository, courseRepository repository.CourseRepository) AddCourseUseCase {
	return &AddCourseUseCaseImpl{CategoryRepository: categoryRepository, CourseRepository: courseRepository}
}

func (u *AddCourseUseCaseImpl) Execute(ctx context.Context, input AddCourseInputDTO) (AddCourseOutputDTO, error) {
	// Isso aqui é o certo, mas vamos ver o UOW
	// category, err := u.CategoryRepository.FindByID(ctx, input.CategoryID)
	// if err != nil {
	// 	return AddCourseOutputDTO{}, err
	// }

	category := entity.Category{
		ID:          input.CategoryID,
		Name:        input.CategoryName,
		Description: input.CategoryDescription,
	}
	err := u.CategoryRepository.Create(ctx, &category)
	if err != nil {
		return AddCourseOutputDTO{}, err
	}

	course := entity.Course{
		ID:          input.ID,
		Category:    category,
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
	}
	err = u.CourseRepository.Create(ctx, &course)
	if err != nil {
		return AddCourseOutputDTO{}, err
	}
	return AddCourseOutputDTO{
		ID:          course.ID,
		CategoryID:  course.Category.ID,
		Name:        course.Name,
		Description: course.Description,
		Price:       course.Price,
	}, nil
}
