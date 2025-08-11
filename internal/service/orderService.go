package service

import (
	"github.com/google/uuid"
	"github.com/k0haku1/order-service/internal/models"
	"github.com/k0haku1/order-service/internal/repositories"
)

type OrderService struct {
	orderRepository *repositories.OrderRepository
}

func NewOrderService(orderRepository *repositories.OrderRepository) *OrderService {
	return &OrderService{orderRepository: orderRepository}
}

func (s *OrderService) CreateOrder(order *models.Order) error {
	return s.orderRepository.CreateOrder(order)
}
func (s *OrderService) GetOrder(id uuid.UUID) (*models.Order, error) {
	return s.orderRepository.FindByID(id)
}
