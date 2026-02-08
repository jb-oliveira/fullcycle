package service

import (
	"context"

	"github.com/jb-oliveira/fullcycle/CleanArch/internal/infra/grpc/pb"
	"github.com/jb-oliveira/fullcycle/CleanArch/internal/usecase"
)

type OrderService struct {
	pb.UnimplementedOrderServiceServer
	CreateOrderUseCase usecase.CreateOrderUseCase
	ListOrdersUseCase  usecase.ListOrdersUseCase
}

func NewOrderService(createOrderUseCase usecase.CreateOrderUseCase, listOrdersUseCase usecase.ListOrdersUseCase) *OrderService {
	return &OrderService{
		CreateOrderUseCase: createOrderUseCase,
		ListOrdersUseCase:  listOrdersUseCase,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.Order, error) {
	dto := usecase.CreateOrderInput{
		ID:    req.Id,
		Price: req.Price,
		Tax:   req.Tax,
	}
	output, err := s.CreateOrderUseCase.Execute(dto)
	if err != nil {
		return nil, err
	}
	return &pb.Order{
		Id:         output.ID,
		Price:      output.Price,
		Tax:        output.Tax,
		FinalPrice: output.FinalPrice,
	}, nil
}

func (s *OrderService) ListOrders(ctx context.Context, req *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	output, err := s.ListOrdersUseCase.Execute(int(req.Page), int(req.Limit), req.Sort, req.SortDir)
	if err != nil {
		return nil, err
	}
	var orders []*pb.Order
	for _, order := range output {
		orders = append(orders, &pb.Order{
			Id:         order.ID,
			Price:      order.Price,
			Tax:        order.Tax,
			FinalPrice: order.CalculateFinalPrice(),
		})
	}
	return &pb.ListOrdersResponse{
		Orders: orders,
	}, nil
}
