package usecase

import (
	"github.com/jb-oliveira/fullcycle/CleanArch/internal/entity"
)

type ListOrdersUseCase struct {
	OrderRepository entity.OrderRepository
}

func NewListOrdersUseCase(orderRepository entity.OrderRepository) *ListOrdersUseCase {
	return &ListOrdersUseCase{
		OrderRepository: orderRepository,
	}
}

func (l *ListOrdersUseCase) Execute(page, limit int, sort, sortDir string) ([]entity.Order, error) {
	orders, err := l.OrderRepository.FindAll(page, limit, sort, sortDir)
	if err != nil {
		return nil, err
	}
	return orders, nil
}
