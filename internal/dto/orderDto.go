package dto

import (
	"github.com/google/uuid"
	"github.com/k0haku1/order-service/internal/models"
)

type CreateOrderRequest struct {
	CustomerID uuid.UUID             `json:"customer_id"`
	Products   []models.OrderProduct `json:"products"`
}

type UpdateOrderRequest struct {
	CustomerID uuid.UUID             `json:"customer_id"`
	OrderID    uuid.UUID             `json:"order_id"`
	Products   []models.OrderProduct `json:"products"`
}

type CreateOrderResponse struct {
	ID         uuid.UUID             `json:"id"`
	CustomerID uuid.UUID             `json:"customer_id"`
	Status     string                `json:"status"`
	Products   []models.OrderProduct `json:"products"`
}
