package repositories

import (
	"github.com/google/uuid"
	"github.com/k0haku1/order-service/internal/models"
	"gorm.io/gorm"
)

type OrderRepository struct {
	db                 *gorm.DB
	customerRepository CustomerRepository
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{
		db: db,
	}
}

func (r *OrderRepository) FindByID(id uuid.UUID) (*models.Order, error) {
	var order models.Order

	if err := r.db.Preload("Customer").Preload("OrderProducts.Product").
		First(&order, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &order, nil
}
func (r *OrderRepository) CreateOrder(customerID uuid.UUID, order *models.Order) error {
	_, err := r.customerRepository.FindByID(customerID)
	if err != nil {
		return err
	}
	order.CustomerID = customerID
	order.Status = "CREATED"

	if err := r.db.Create(order).Error; err != nil {
		return err
	}

	return r.db.Preload("OrderProducts.Product").First(order, "id = ?", order.ID).Error
}

func (r *OrderRepository) addProductToOrder(id uuid.UUID) error {

}
