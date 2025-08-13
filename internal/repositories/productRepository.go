package repositories

import (
	"github.com/google/uuid"
	"github.com/k0haku1/order-service/internal/models"
	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) FindByID(id []uuid.UUID) ([]models.Product, error) {
	var products []models.Product
	if err := r.db.Where("id IN ?", id).Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (r *ProductRepository) UpdateWithTx(tx *gorm.DB, product *models.Product) error {
	return tx.Save(product).Error
}
