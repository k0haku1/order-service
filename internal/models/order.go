package models

import (
	"github.com/google/uuid"
)

type Order struct {
	ID            uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	CustomerID    uuid.UUID `gorm:"type:uuid;not null" json:"customer_id"`
	Customer      Customer  `gorm:"foreignKey:CustomerID;references:ID" json:"customer"`
	Status        string    `gorm:"type:text;not null" json:"status"`
	OrderProducts []OrderProduct
}
