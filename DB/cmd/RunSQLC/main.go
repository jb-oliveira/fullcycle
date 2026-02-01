package main

import (
	"context"
	"fmt"
	"log"

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

	queries := db.New(conn)

	var params []db.CreateCategoriesBatchParams
	count, err := queries.CountCategories(ctx)
	if err != nil {
		log.Fatal(err)
	}
	for range 10000 {
		count++
		id := db.NewIDString()
		params = append(params, db.CreateCategoriesBatchParams{
			ID:          id,
			Name:        fmt.Sprintf("Category %d", count),
			Description: pgtype.Text{String: fmt.Sprintf("Category %d description", count), Valid: true},
		})
	}

	// queries.CreateCategoriesBatch(ctx, params)

	categories, err := queries.ListCategories(ctx, db.ListCategoriesParams{
		LimitVal:   5,
		OffsetVal:  0,
		SortOrder:  "ASC",
		SortColumn: "name",
	})
	if err != nil {
		log.Fatal(err)
	}

	for _, category := range categories {
		log.Println(category)
	}
}
