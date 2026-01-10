package main

import (
	"database/sql"
	"net"

	"github.com/jb-oliveira/fullcycle/gRPC/internal/database"
	"github.com/jb-oliveira/fullcycle/gRPC/internal/pb"
	"github.com/jb-oliveira/fullcycle/gRPC/internal/service"
	_ "github.com/mattn/go-sqlite3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	db, err := sql.Open("sqlite3", "./db.sqlite")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	categoryDB := database.NewCategory(db)
	categoryService := service.NewCategoryService(*categoryDB)

	grpcServer := grpc.NewServer()
	pb.RegisterCategoryServiceServer(grpcServer, categoryService)

	// essa linha é so parqa usar o client não é necessaria
	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}
	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}
