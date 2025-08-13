package repositories

import (
	"github.com/google/uuid"
	"github.com/k0haku1/order-service/internal/models"
	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
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
func (r *OrderRepository) CreateWithTx(tx *gorm.DB, order *models.Order) error {
	return tx.Create(order).Error
}

func (r *OrderRepository) UpdateWithTx(tx *gorm.DB, order *models.Order) error {
	if err := tx.Save(order).Error; err != nil {
		return err
	}

	for i := range order.OrderProducts {
		if err := tx.Save(&order.OrderProducts[i]).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *OrderRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Order{}, "id = ?", id).Error
}

func (r *OrderRepository) WithTransaction(fn func(tx *gorm.DB) error) error {
	return r.db.Transaction(fn)
}
