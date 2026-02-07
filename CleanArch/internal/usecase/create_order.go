package usecase

import (
	"github.com/jb-oliveira/fullcycle/CleanArch/internal/entity"
	"github.com/jb-oliveira/fullcycle/CleanArch/pkg/events"
)

type CreateOrderInput struct {
	ID    string  `json:"id"`
	Price float64 `json:"price"`
	Tax   float64 `json:"tax"`
}

type CreateOrderOutput struct {
	ID         string  `json:"id"`
	Price      float64 `json:"price"`
	Tax        float64 `json:"tax"`
	FinalPrice float64 `json:"final_price"`
}

type CreateOrderUseCase struct {
	OrderRepository entity.OrderRepository
	OrderCreated    events.EventInterface
	EventDispatcher events.EventDispatcherInterface
}

func NewCreateOrderUseCase(orderRepository entity.OrderRepository, orderCreated events.EventInterface, eventDispatcher events.EventDispatcherInterface) *CreateOrderUseCase {
	return &CreateOrderUseCase{
		OrderRepository: orderRepository,
		OrderCreated:    orderCreated,
		EventDispatcher: eventDispatcher,
	}
}

func (c *CreateOrderUseCase) Execute(input CreateOrderInput) (CreateOrderOutput, error) {
	order, err := entity.NewOrder(input.ID, input.Price, input.Tax)
	if err != nil {
		return CreateOrderOutput{}, err
	}
	finalPrice := order.CalculateFinalPrice()
	err = c.OrderRepository.Save(order)
	if err != nil {
		return CreateOrderOutput{}, err
	}
	dto := CreateOrderOutput{
		ID:         order.ID,
		Price:      order.Price,
		Tax:        order.Tax,
		FinalPrice: finalPrice,
	}
	c.OrderCreated.SetPayload(dto)
	c.EventDispatcher.Dispatch(c.OrderCreated)
	return dto, nil
}
