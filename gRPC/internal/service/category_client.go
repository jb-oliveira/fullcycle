package service

import (
	"context"
	"io"
	"log"
	"time"

	"github.com/jb-oliveira/fullcycle/gRPC/internal/database"
	"github.com/jb-oliveira/fullcycle/gRPC/internal/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

type CategoryClient struct {
	Client pb.CategoryServiceClient
}

func NewCategoryClient(client pb.CategoryServiceClient) *CategoryClient {
	return &CategoryClient{Client: client}
}

func (s *CategoryClient) Create(name string, desc string) (*pb.CategoryResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	return s.Client.CreateCategory(ctx, &pb.CreateCategoryRequest{Name: name, Description: desc})
}

func (s *CategoryClient) List() (*pb.CategoryListResponse, error) {
	return s.Client.ListCategory(context.Background(), &emptypb.Empty{})
}

func (s *CategoryClient) CreateBulk(categoryDB []database.CategoryDB) (*pb.CategoryListResponse, error) {
	stream, err := s.Client.CreateCategoryStream(context.Background())
	if err != nil {
		return nil, err
	}

	for _, category := range categoryDB {
		if err := stream.Send(&pb.CreateCategoryRequest{Name: category.Name, Description: category.Description}); err != nil {
			return nil, err
		}
	}

	// CloseAndRecv tells the server we are done sending and waits for the final result
	return stream.CloseAndRecv()
}

func (s *CategoryClient) CreateChatty(categoryDB []database.CategoryDB) error {
	stream, err := s.Client.CreateCategoryStreamBidirectional(context.Background())
	if err != nil {
		return err
	}

	waitc := make(chan struct{})

	// Receive goroutine
	go s.receive(stream, waitc)

	// Send messages
	for _, category := range categoryDB {
		if err := stream.Send(&pb.CreateCategoryRequest{Name: category.Name, Description: category.Description}); err != nil {
			return err
		}
	}
	stream.CloseSend()
	<-waitc
	return nil
}

func (s *CategoryClient) receive(stream pb.CategoryService_CreateCategoryStreamBidirectionalClient, waitc chan struct{}) {
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			close(waitc)
			return
		}
		if err != nil {
			log.Printf("Error receiving from stream: %v", err)
			close(waitc)
			return
		}
		log.Printf("Server processed: %s", res.Category.Name)
	}
}
