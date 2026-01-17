import (
	"context"
	"io"
	"log"
	"time"

	"google.golang.org/protobuf/types/known/emptypb"
)

type CategoryClient struct {
	Client CategoryServiceClient
}

func NewCategoryClient(client CategoryServiceClient) *CategoryClient {
	return &CategoryClient{Client: client}
}

func (s *CategoryClient) Create(name string, desc string) (*CategoryResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	return s.Client.CreateCategory(ctx, &CreateCategoryRequest{Name: name, Description: desc})
}

func (s *CategoryClient) List() (*CategoryListResponse, error) {
	return s.Client.ListCategory(context.Background(), &emptypb.Empty{})
}

func (s *CategoryClient) CreateBulk(categoryDB []database.CategoryDB) (*CategoryListResponse, error) {
	stream, err := s.Client.CreateCategoryStream(context.Background())
	if err != nil {
		return nil, err
	}

	for _, category := range categoryDB {
		if err := stream.Send(&CreateCategoryRequest{Name: category.Name, Description: category.Description}); err != nil {
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
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				close(waitc)
				return
			}
			if err != nil {
				return err
			}
			log.Printf("Server processed: %s", res.Name)
		}
	}()

	// Send messages
	for _, category := range categoryDB {
		if err := stream.Send(&CreateCategoryRequest{Name: category.Name, Description: category.Description}); err != nil {
			return err
		}
	}
	stream.CloseSend()
	<-waitc
	return nil
}