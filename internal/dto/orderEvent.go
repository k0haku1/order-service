package dto

import "github.com/google/uuid"

type OrderEvent struct {
	EventID    uuid.UUID           `json:"event_id"`
	OrderID    uuid.UUID           `json:"order_id"`
	CustomerID uuid.UUID           `json:"customer_id"`
	Products   []OrderEventProduct `json:"products"`
}
