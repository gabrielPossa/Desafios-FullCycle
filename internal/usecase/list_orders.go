package usecase

import (
	"github.com/gabrielPossa/Desafios-FullCycle/cleanArch/internal/entity"
)

type ListOrderOutputDTO struct {
	Orders []entity.Order
}

type ListOrdersUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
}

func NewListOrdersUseCase(
	OrderRepository entity.OrderRepositoryInterface,
) *ListOrdersUseCase {
	return &ListOrdersUseCase{
		OrderRepository: OrderRepository,
	}
}

func (c *ListOrdersUseCase) Execute() (ListOrderOutputDTO, error) {
	orders, err := c.OrderRepository.GetAllOrders()
	if err != nil {
		return ListOrderOutputDTO{}, err
	}

	dto := ListOrderOutputDTO{
		Orders: orders,
	}

	return dto, nil
}
