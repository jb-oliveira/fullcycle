package usecase

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/testcontainers/testcontainers-go"
	pgmodule "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupTestContainer(t *testing.T) (*sql.DB, func()) {
	ctx := context.Background()

	// Start Postgres Container
	container, err := pgmodule.Run(ctx, "postgres:18.1-alpine",
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

	// Get Connection String
	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	sqlDB, err := sql.Open("pgx", connStr)
	if err != nil {
		t.Fatal(err)
	}

	// Run Migrations
	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		t.Fatal(err)
	}

	migrationsPath, _ := filepath.Abs("../../db/migrations")
	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
		"postgres", driver)
	if err != nil {
		t.Fatal(err)
	}
	m.Up()

	// Return cleanup function
	return sqlDB, func() {
		sqlDB.Close()
		container.Terminate(ctx)
	}
}
