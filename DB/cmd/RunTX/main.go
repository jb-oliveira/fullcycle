package main

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jb-oliveira/fullcycle/DB/internal/db"
)

func main() {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, "postgres://postgres:password@localhost:5432/myapp")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(ctx)

	err = db.ExecuteInTransaction(ctx, conn, func(queries *db.Queries) error {
		category, err := queries.CreateCategory(ctx, db.CreateCategoryParams{
			ID:          db.NewIDString(),
			Name:        "Category 2",
			Description: pgtype.Text{String: "Category 2 description", Valid: true},
		})
		if err != nil {
			return err
		}

		for i := 100; i <= 110; i++ {
			_, err := queries.CreateCourse(ctx, db.CreateCourseParams{
				ID:          db.NewIDString(),
				CategoryID:  category.ID,
				Name:        fmt.Sprintf("Course %d", i),
				Description: pgtype.Text{String: fmt.Sprintf("Course %d description", i), Valid: true},
				Price:       pgtype.Numeric{Int: big.NewInt(int64(i * 89)), Valid: true},
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}
