package models

import "github.com/google/uuid"

type Customer struct {
	ID     uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name   string    `gorm:"type:varchar(255);not null"`
	Email  string    `gorm:"type:varchar(255);not null"`
	Orders []Order   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"orders,omitempty"`
}
