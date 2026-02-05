package usecase

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jb-oliveira/fullcycle/UnitOfWork/internal/entity"
	"github.com/jb-oliveira/fullcycle/UnitOfWork/internal/repository"
	"github.com/jb-oliveira/fullcycle/UnitOfWork/pkg/uow"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func TestAddCourseUseUowCase_Execute(t *testing.T) {
	sqlDB, cleanup := setupTestContainer(t)
	defer cleanup()

	uow := uow.NewUnitOfWork(context.Background(), sqlDB)

	uow.Register("category_repository", func(tx *sql.Tx) any {
		return repository.NewCategoryRepositoryPGImpl(tx)
	})
	uow.Register("course_repository", func(tx *sql.Tx) any {
		return repository.NewCourseRepositoryPGImpl(tx)
	})

	// Initialize repositories for verification
	categoryRepo := repository.NewCategoryRepositoryPGImpl(sqlDB)
	courseRepo := repository.NewCourseRepositoryPGImpl(sqlDB)

	// Initialize use case
	useCase := NewAddCourseUseUowCaseImpl(uow)

	ctx := context.Background()

	t.Run("Successfully add course with new category", func(t *testing.T) {
		input := AddCourseInputDTO{
			ID:                  "course-123",
			CategoryID:          "cat-123",
			CategoryName:        "Programming",
			CategoryDescription: "Software development courses",
			Name:                "Go Programming",
			Description:         "Learn Go programming language",
			Price:               99.99,
		}

		output, err := useCase.Execute(ctx, input)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		// Verify output
		if output.ID != input.ID {
			t.Errorf("expected ID %s, got %s", input.ID, output.ID)
		}
		if output.CategoryID != input.CategoryID {
			t.Errorf("expected CategoryID %s, got %s", input.CategoryID, output.CategoryID)
		}
		if output.Name != input.Name {
			t.Errorf("expected Name %s, got %s", input.Name, output.Name)
		}
		if output.Description != input.Description {
			t.Errorf("expected Description %s, got %s", input.Description, output.Description)
		}
		if output.Price != input.Price {
			t.Errorf("expected Price %.2f, got %.2f", input.Price, output.Price)
		}

		// Verify category was created in database
		savedCategory, err := categoryRepo.FindByID(ctx, input.CategoryID)
		if err != nil {
			t.Errorf("could not find category in DB: %v", err)
		}
		if savedCategory.Name != input.CategoryName {
			t.Errorf("expected category name %s, got %s", input.CategoryName, savedCategory.Name)
		}
		if savedCategory.Description != input.CategoryDescription {
			t.Errorf("expected category description %s, got %s", input.CategoryDescription, savedCategory.Description)
		}

		// Verify course was created in database
		savedCourse, err := courseRepo.FindByID(ctx, input.ID)
		if err != nil {
			t.Errorf("could not find course in DB: %v", err)
		}
		if savedCourse.Name != input.Name {
			t.Errorf("expected course name %s, got %s", input.Name, savedCourse.Name)
		}
		if savedCourse.Description != input.Description {
			t.Errorf("expected course description %s, got %s", input.Description, savedCourse.Description)
		}
		if savedCourse.Price != input.Price {
			t.Errorf("expected course price %.2f, got %.2f", input.Price, savedCourse.Price)
		}
		if savedCourse.Category.ID != input.CategoryID {
			t.Errorf("expected course category ID %s, got %s", input.CategoryID, savedCourse.Category.ID)
		}
	})

	t.Run("Handle duplicate category creation", func(t *testing.T) {
		// First, create a category
		firstInput := AddCourseInputDTO{
			ID:                  "course-456",
			CategoryID:          "cat-456",
			CategoryName:        "Web Development",
			CategoryDescription: "Frontend and backend development",
			Name:                "React Basics",
			Description:         "Learn React fundamentals",
			Price:               79.99,
		}

		_, err := useCase.Execute(ctx, firstInput)
		if err != nil {
			t.Errorf("expected no error on first execution, got %v", err)
		}

		// Try to create another course with the same category
		secondInput := AddCourseInputDTO{
			ID:                  "course-789",
			CategoryID:          "cat-456", // Same category ID
			CategoryName:        "Web Development",
			CategoryDescription: "Frontend and backend development",
			Name:                "Advanced React",
			Description:         "Advanced React patterns",
			Price:               129.99,
		}

		// This should fail due to unique constraint on category name
		_, err = useCase.Execute(ctx, secondInput)
		if err == nil {
			t.Error("expected error due to duplicate category name, got nil")
		}
	})

	t.Run("Handle duplicate course creation", func(t *testing.T) {
		input := AddCourseInputDTO{
			ID:                  "course-duplicate",
			CategoryID:          "cat-duplicate",
			CategoryName:        "Database",
			CategoryDescription: "Database management courses",
			Name:                "SQL Fundamentals",
			Description:         "Learn SQL basics",
			Price:               59.99,
		}

		// First execution should succeed
		_, err := useCase.Execute(ctx, input)
		if err != nil {
			t.Errorf("expected no error on first execution, got %v", err)
		}

		// Second execution with same course ID should fail
		_, err = useCase.Execute(ctx, input)
		if err == nil {
			t.Error("expected error due to duplicate course ID, got nil")
		}
	})

	t.Run("Handle zero price course", func(t *testing.T) {
		input := AddCourseInputDTO{
			ID:                  "course-free",
			CategoryID:          "cat-free",
			CategoryName:        "Free Courses",
			CategoryDescription: "Free educational content",
			Name:                "Introduction to Programming",
			Description:         "Basic programming concepts",
			Price:               0.00,
		}

		output, err := useCase.Execute(ctx, input)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if output.Price != 0.00 {
			t.Errorf("expected price 0.00, got %.2f", output.Price)
		}

		// Verify in database
		savedCourse, err := courseRepo.FindByID(ctx, input.ID)
		if err != nil {
			t.Errorf("could not find course in DB: %v", err)
		}
		if savedCourse.Price != 0.00 {
			t.Errorf("expected saved course price 0.00, got %.2f", savedCourse.Price)
		}
	})

	t.Run("Handle high precision price", func(t *testing.T) {
		input := AddCourseInputDTO{
			ID:                  "course-precision",
			CategoryID:          "cat-precision",
			CategoryName:        "Premium Courses",
			CategoryDescription: "High-value educational content",
			Name:                "Advanced Architecture",
			Description:         "Software architecture patterns",
			Price:               199.999, // High precision price
		}

		output, err := useCase.Execute(ctx, input)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		// The price should be preserved as-is in the output
		if output.Price != input.Price {
			t.Errorf("expected price %.3f, got %.3f", input.Price, output.Price)
		}

		// Verify in database (may be rounded depending on DB precision)
		savedCourse, err := courseRepo.FindByID(ctx, input.ID)
		if err != nil {
			t.Errorf("could not find course in DB: %v", err)
		}
		// Allow for small floating point differences
		if savedCourse.Price < 199.99 || savedCourse.Price > 200.01 {
			t.Errorf("expected saved course price around 200.00, got %.3f", savedCourse.Price)
		}
	})

	t.Run("Handle empty strings in input", func(t *testing.T) {
		input := AddCourseInputDTO{
			ID:                  "course-empty",
			CategoryID:          "cat-empty",
			CategoryName:        "", // Empty category name
			CategoryDescription: "Valid description",
			Name:                "Valid Course Name",
			Description:         "Valid course description",
			Price:               49.99,
		}

		output, err := useCase.Execute(ctx, input)
		// Current implementation doesn't validate, so it should succeed
		// but create a category with empty name
		if err != nil {
			t.Errorf("expected no error with current implementation, got %v", err)
		}

		// Verify the category was created with empty name
		savedCategory, err := categoryRepo.FindByID(ctx, input.CategoryID)
		if err != nil {
			t.Errorf("could not find category in DB: %v", err)
		}
		if savedCategory.Name != "" {
			t.Errorf("expected empty category name, got %s", savedCategory.Name)
		}

		// Verify course was still created successfully
		if output.ID != input.ID {
			t.Errorf("expected course ID %s, got %s", input.ID, output.ID)
		}
	})

	t.Run("Validate category entity validation method", func(t *testing.T) {
		// Test the entity validation method directly
		category := &entity.Category{
			ID:          "test-id",
			Name:        "", // Empty name
			Description: "Valid description",
		}

		err := category.Validate()
		if err == nil {
			t.Error("expected validation error for empty name, got nil")
		}

		// Test with empty description
		category.Name = "Valid name"
		category.Description = ""
		err = category.Validate()
		if err == nil {
			t.Error("expected validation error for empty description, got nil")
		}

		// Test with valid data
		category.Description = "Valid description"
		err = category.Validate()
		if err != nil {
			t.Errorf("expected no validation error for valid data, got %v", err)
		}
	})

	t.Run("Handle negative price", func(t *testing.T) {
		input := AddCourseInputDTO{
			ID:                  "course-negative",
			CategoryID:          "cat-negative",
			CategoryName:        "Test Category",
			CategoryDescription: "Test description",
			Name:                "Negative Price Course",
			Description:         "Course with negative price",
			Price:               -50.00,
		}

		output, err := useCase.Execute(ctx, input)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		// Verify negative price is preserved
		if output.Price != -50.00 {
			t.Errorf("expected price -50.00, got %.2f", output.Price)
		}

		// Verify in database
		savedCourse, err := courseRepo.FindByID(ctx, input.ID)
		if err != nil {
			t.Errorf("could not find course in DB: %v", err)
		}
		if savedCourse.Price != -50.00 {
			t.Errorf("expected saved course price -50.00, got %.2f", savedCourse.Price)
		}
	})

	t.Run("Handle very long strings", func(t *testing.T) {
		longName := "This is a very long course name that might exceed normal expectations and could potentially cause issues with database storage or application logic if not handled properly"
		longDescription := "This is an extremely long description that goes on and on and on, providing extensive details about the course content, learning objectives, prerequisites, target audience, and expected outcomes. It continues with more information about the curriculum, assessment methods, and certification details."

		input := AddCourseInputDTO{
			ID:                  "course-long",
			CategoryID:          "cat-long",
			CategoryName:        "Long Content Category",
			CategoryDescription: "Category for testing long content",
			Name:                longName,
			Description:         longDescription,
			Price:               299.99,
		}

		output, err := useCase.Execute(ctx, input)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		// Verify long strings are preserved
		if output.Name != longName {
			t.Errorf("expected long name to be preserved")
		}
		if output.Description != longDescription {
			t.Errorf("expected long description to be preserved")
		}

		// Verify in database
		savedCourse, err := courseRepo.FindByID(ctx, input.ID)
		if err != nil {
			t.Errorf("could not find course in DB: %v", err)
		}
		if savedCourse.Name != longName {
			t.Errorf("expected saved course name to match input")
		}
		if savedCourse.Description != longDescription {
			t.Errorf("expected saved course description to match input")
		}
	})

	t.Run("Handle special characters in strings", func(t *testing.T) {
		input := AddCourseInputDTO{
			ID:                  "course-special",
			CategoryID:          "cat-special",
			CategoryName:        "Special Characters & Symbols",
			CategoryDescription: "Category with émojis 🚀 and spëcial chars: @#$%^&*()",
			Name:                "Advanced C++ & Go Programming",
			Description:         "Learn C++ & Go with special focus on memory management, pointers -> references, and more!",
			Price:               149.99,
		}

		output, err := useCase.Execute(ctx, input)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		// Verify special characters are preserved
		if output.Name != input.Name {
			t.Errorf("expected name with special characters to be preserved")
		}

		// Verify in database
		savedCategory, err := categoryRepo.FindByID(ctx, input.CategoryID)
		if err != nil {
			t.Errorf("could not find category in DB: %v", err)
		}
		if savedCategory.Name != input.CategoryName {
			t.Errorf("expected saved category name to match input with special characters")
		}

		savedCourse, err := courseRepo.FindByID(ctx, input.ID)
		if err != nil {
			t.Errorf("could not find course in DB: %v", err)
		}
		if savedCourse.Name != input.Name {
			t.Errorf("expected saved course name to match input with special characters")
		}
	})
}
