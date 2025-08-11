package repositories

import (
	"github.com/google/uuid"
	"github.com/k0haku1/order-service/internal/models"
	"gorm.io/gorm"
)

type CustomerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) *CustomerRepository {
	return &CustomerRepository{db: db}
}

func (r *CustomerRepository) FindByID(id uuid.UUID) (*models.Customer, error) {
	var customer models.Customer
	if err := r.db.First(&customer, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &customer, nil
}
