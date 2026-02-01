migrate-up:
	migrate -path=db/migrations -database "postgres://postgres:password@localhost:5432/myapp?sslmode=disable" up

migrate-down:
	migrate -path=db/migrations -database "postgres://postgres:password@localhost:5432/myapp?sslmode=disable" down


.PHONY: migrate-up migrate-down