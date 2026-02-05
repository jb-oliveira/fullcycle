package usecase

import (
	"context"
	"database/sql"

	"github.com/jb-oliveira/fullcycle/UnitOfWork/internal/entity"
	"github.com/jb-oliveira/fullcycle/UnitOfWork/internal/repository"
	"github.com/jb-oliveira/fullcycle/UnitOfWork/pkg/uow"
)

type AddCourseUseUowCaseImpl struct {
	uow uow.UnitOfWork
}

func NewAddCourseUseUowCaseImpl(uow uow.UnitOfWork) AddCourseUseCase {
	return &AddCourseUseUowCaseImpl{uow: uow}
}

func (u *AddCourseUseUowCaseImpl) Execute(ctx context.Context, input AddCourseInputDTO) (AddCourseOutputDTO, error) {
	var output AddCourseOutputDTO
	err := u.uow.Do(ctx, func(uow uow.UnitOfWork) error {
		categoryRepository := u.getCategoryRepository(ctx, uow)
		category, err := categoryRepository.FindByID(ctx, input.CategoryID)
		if err != nil && err != sql.ErrNoRows {
			return err
		}
		if err == sql.ErrNoRows {
			category = &entity.Category{
				ID:          input.CategoryID,
				Name:        input.CategoryName,
				Description: input.CategoryDescription,
			}
			categoryRepository.Create(ctx, category)
		}
		course := entity.Course{
			ID:          input.ID,
			Category:    *category,
			Name:        input.Name,
			Description: input.Description,
			Price:       input.Price,
		}
		courseRepository := u.getCourseRepository(ctx, uow)
		err = courseRepository.Create(ctx, &course)
		if err != nil {
			return err
		}
		output = AddCourseOutputDTO{
			ID:          course.ID,
			CategoryID:  course.Category.ID,
			Name:        course.Name,
			Description: course.Description,
			Price:       course.Price,
		}
		return nil
	})
	if err != nil {
		return AddCourseOutputDTO{}, err
	}
	return output, nil
}

func (u *AddCourseUseUowCaseImpl) getCategoryRepository(ctx context.Context, uow uow.UnitOfWork) repository.CategoryRepository {
	repo, err := uow.Get(ctx, "category_repository")
	if err != nil {
		panic(err)
	}
	return repo.(repository.CategoryRepository)
}

func (u *AddCourseUseUowCaseImpl) getCourseRepository(ctx context.Context, uow uow.UnitOfWork) repository.CourseRepository {
	repo, err := uow.Get(ctx, "course_repository")
	if err != nil {
		panic(err)
	}
	return repo.(repository.CourseRepository)
}
