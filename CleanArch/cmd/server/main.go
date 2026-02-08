package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	graphql_handler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	config "github.com/jb-oliveira/fullcycle/CleanArch/configs"
	"github.com/jb-oliveira/fullcycle/CleanArch/internal/event/handler"
	"github.com/jb-oliveira/fullcycle/CleanArch/internal/infra/database"
	"github.com/jb-oliveira/fullcycle/CleanArch/internal/infra/graph"
	"github.com/jb-oliveira/fullcycle/CleanArch/internal/infra/grpc/pb"
	"github.com/jb-oliveira/fullcycle/CleanArch/internal/infra/grpc/service"
	"github.com/jb-oliveira/fullcycle/CleanArch/internal/infra/web/webserver"
	"github.com/jb-oliveira/fullcycle/CleanArch/pkg/events"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/streadway/amqp"
)

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	db, err := database.NewDatabase(cfg)
	if err != nil {
		log.Fatal("cannot connect to database:", err)
	}
	defer db.Close()

	channel, err := getRabbitMQConnection()
	if err != nil {
		log.Fatal("cannot connect to rabbitmq:", err)
	}
	defer channel.Close()

	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("OrderCreated", handler.NewOrderCreatedHandlerRabbitMQ(channel))

	createOrderUseCase := NewCreateOrderUseCase(db, eventDispatcher)
	listOrdersUseCase := NewListOrdersUseCase(db)

	handler := NewWebOrderHandler(db, eventDispatcher)
	webserver := webserver.NewWebServer(cfg.WebServerPort)
	webserver.RegisterHandler(http.MethodPost, "/api/v1/orders", handler.CreateOrder)
	webserver.RegisterHandler(http.MethodGet, "/api/v1/orders", handler.ListOrders)
	fmt.Println("Starting web server on port", cfg.WebServerPort)
	go func() {
		if err := webserver.Start(); err != nil {
			panic(err)
		}
	}()

	grpcServer := grpc.NewServer()
	orderService := service.NewOrderService(*createOrderUseCase, *listOrdersUseCase)
	pb.RegisterOrderServiceServer(grpcServer, orderService)
	// This line is only necessary for evans
	reflection.Register(grpcServer)

	fmt.Println("Starting gRPC server on port", cfg.GrpcServerPort)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GrpcServerPort))
	if err != nil {
		panic(err)
	}
	go grpcServer.Serve(lis)

	srv := graphql_handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		CreateOrderUseCase: *createOrderUseCase,
		ListOrdersUseCase:  *listOrdersUseCase,
	}}))
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	fmt.Println("Starting GraphQL server on port", cfg.GrapQLServerPort)
	http.ListenAndServe(":"+cfg.GrapQLServerPort, nil)
}

func getRabbitMQConnection() (*amqp.Channel, error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672")
	if err != nil {
		return nil, err
	}
	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	return channel, nil
}
