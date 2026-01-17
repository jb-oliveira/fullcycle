package service

import (
	"context"
	"io"

	"github.com/jb-oliveira/fullcycle/gRPC/internal/database"
	"github.com/jb-oliveira/fullcycle/gRPC/internal/pb"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type CategoryService struct {
	pb.UnimplementedCategoryServiceServer
	categoryDB database.CategoryDB
}

func NewCategoryService(categoryDB database.CategoryDB) *CategoryService {
	return &CategoryService{categoryDB: categoryDB}
}

func (c *CategoryService) GetCategory(ctx context.Context, in *pb.GetCategoryRequest) (*pb.CategoryResponse, error) {
	category, err := c.categoryDB.FindById(in.Id)
	if err != nil {
		return nil, err
	}

	return &pb.CategoryResponse{
		Category: &pb.Category{
			Id:          category.Id,
			Name:        category.Name,
			Description: category.Description,
		},
	}, nil
}

func (c *CategoryService) CreateCategory(ctx context.Context, in *pb.CreateCategoryRequest) (*pb.CategoryResponse, error) {
	c.categoryDB.Name = in.Name
	c.categoryDB.Description = in.Description
	err := c.categoryDB.Create()
	if err != nil {
		return nil, err
	}

	return &pb.CategoryResponse{
		Category: &pb.Category{
			Id:          c.categoryDB.Id,
			Name:        c.categoryDB.Name,
			Description: c.categoryDB.Description,
		},
	}, nil
}

func (c *CategoryService) CreateCategoryStream(stream grpc.ClientStreamingServer[pb.CreateCategoryRequest, pb.CategoryListResponse]) error {
	categories := &pb.CategoryListResponse{}

	for {
		category, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(categories)
		}
		if err != nil {
			return err
		}
		c.categoryDB.Name = category.Name
		c.categoryDB.Description = category.Description
		err = c.categoryDB.Create()
		if err != nil {
			return err
		}
		categories.Categories = append(categories.Categories, &pb.Category{
			Id:          c.categoryDB.Id,
			Name:        c.categoryDB.Name,
			Description: c.categoryDB.Description,
		})
	}

}

func (c *CategoryService) CreateCategoryStreamBidirectional(stream grpc.BidiStreamingServer[pb.CreateCategoryRequest, pb.CategoryResponse]) error {
	for {
		category, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		c.categoryDB.Name = category.Name
		c.categoryDB.Description = category.Description
		err = c.categoryDB.Create()
		if err != nil {
			return err
		}
		err = stream.Send(&pb.CategoryResponse{
			Category: &pb.Category{
				Id:          c.categoryDB.Id,
				Name:        c.categoryDB.Name,
				Description: c.categoryDB.Description,
			},
		})
		if err != nil {
			return err
		}
	}
}

func (c *CategoryService) ListCategory(ctx context.Context, in *emptypb.Empty) (*pb.CategoryListResponse, error) {
	categories, err := c.categoryDB.FindAll()
	if err != nil {
		return nil, err
	}

	var response pb.CategoryListResponse
	for _, category := range categories {
		response.Categories = append(response.Categories, &pb.Category{
			Id:          category.Id,
			Name:        category.Name,
			Description: category.Description,
		})
	}

	return &response, nil
}
