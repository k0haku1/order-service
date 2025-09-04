package dto

import "github.com/google/uuid"

type OrderEventProduct struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Quantity int32     `json:"quantity"`
}
