package repository

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jb-oliveira/fullcycle/UnitOfWork/internal/entity"

	"github.com/testcontainers/testcontainers-go"
	pgmodule "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func setupTestContainer(t *testing.T) (*sql.DB, func()) {
	ctx := context.Background()

	// 1. Start Postgres Container
	container, err := pgmodule.Run(ctx, "postgres:15-alpine",
		pgmodule.WithDatabase("testdb"),
		pgmodule.WithUsername("user"),
		pgmodule.WithPassword("pass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(30*time.Second)),
	)
	if err != nil {
		t.Fatal(err)
	}

	// 2. Get Connection String
	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	sqlDB, err := sql.Open("pgx", connStr)
	if err != nil {
		t.Fatal(err)
	}

	// 3. Run Migrations
	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		t.Fatal(err)
	}

	// Adjust the path to your actual migrations folder
	migrationsPath, _ := filepath.Abs("../../db/migrations")
	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
		"postgres", driver)
	if err != nil {
		t.Fatal(err)
	}
	m.Up()

	// 4. Return cleanup function
	return sqlDB, func() {
		sqlDB.Close()
		container.Terminate(ctx)
	}
}

func TestCategoryRepository_Create(t *testing.T) {
	sqlDB, cleanup := setupTestContainer(t)
	defer cleanup()

	// Initialize sqlc queries and your repository
	repo := NewCategoryRepositoryPGImpl(sqlDB) // passing sqlDB and internally calling db.New

	ctx := context.Background()

	t.Run("Successfully create a category", func(t *testing.T) {
		category := &entity.Category{
			ID:          "cat-123",
			Name:        "Electronics",
			Description: "Tech gadgets",
		}

		err := repo.Create(ctx, category)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		// Verify using sqlc directly to check DB state
		saved, err := repo.FindByID(ctx, "cat-123")
		if err != nil {
			t.Errorf("could not find category in DB: %v", err)
		}
		if saved.Name != "Electronics" {
			t.Errorf("expected Electronics, got %s", saved.Name)
		}
	})
}
