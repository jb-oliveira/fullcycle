package main

import (
	"context"
	"fmt"

	"io"

	"github.com/jb-oliveira/fullcycle/gRPC/internal/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("Failed to connect:", err)
		return
	}
	defer conn.Close()

	client := pb.NewCategoryServiceClient(conn)

	// CreateCategory
	category, err := client.CreateCategory(context.Background(), &pb.CreateCategoryRequest{
		Name:        "Test Category",
		Description: "Test Category Description",
	})
	if err != nil {
		fmt.Println("Failed to create category:", err)
		return
	}
	fmt.Println("Category created:", category.Category)

	// CreateCategoryStream
	stream, err := client.CreateCategoryStream(context.Background())
	if err != nil {
		fmt.Println("Failed to create stream:", err)
		return
	}
	for i := 0; i < 5; i++ {
		if err := stream.Send(&pb.CreateCategoryRequest{
			Name:        fmt.Sprintf("Test Category %d", i),
			Description: fmt.Sprintf("Test Category Description %d", i),
		}); err != nil {
			fmt.Println("Failed to send:", err)
			return
		}
	}
	categoryList, err := stream.CloseAndRecv()
	if err != nil {
		fmt.Println("Failed to receive:", err)
		return
	}
	fmt.Println("Category stream created:", categoryList.Categories)

	// CreateCategoryStreamBidirectional
	streamBidi, err := client.CreateCategoryStreamBidirectional(context.Background())
	if err != nil {
		fmt.Println("Failed to create stream:", err)
		return
	}

	waitc := make(chan struct{})
	go func() {
		for {
			res, err := streamBidi.Recv()
			if err == io.EOF {
				close(waitc)
				return
			}
			if err != nil {
				fmt.Println("Failed to receive:", err)
				return
			}
			fmt.Println("Category stream bidirectional created:", res.Category)
		}
	}()

	for i := 0; i < 5; i++ {
		if err := streamBidi.Send(&pb.CreateCategoryRequest{
			Name:        fmt.Sprintf("Test Category %d", i),
			Description: fmt.Sprintf("Test Category Description %d", i),
		}); err != nil {
			fmt.Println("Failed to send:", err)
			return
		}
	}
	streamBidi.CloseSend()
	<-waitc

	// ListCategory
	categoryList, err = client.ListCategory(context.Background(), &emptypb.Empty{})
	if err != nil {
		fmt.Println("Failed to list categories:", err)
		return
	}
	fmt.Println("Categories:", categoryList.Categories)

	// GetCategory
	category, err = client.GetCategory(context.Background(), &pb.GetCategoryRequest{
		Id: "1",
	})
	if err != nil {
		fmt.Println("Failed to get category:", err)
		return
	}
	fmt.Println("Category:", category.Category)
}
