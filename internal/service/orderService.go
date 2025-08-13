package service

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/k0haku1/order-service/internal/models"
	"github.com/k0haku1/order-service/internal/repositories"
	"gorm.io/gorm"
	"log"
)

type OrderService struct {
	orderRepository    *repositories.OrderRepository
	customerRepository *repositories.CustomerRepository
	productRepository  *repositories.ProductRepository
}

func NewOrderService(
	orderRepository *repositories.OrderRepository,
	customerRepository *repositories.CustomerRepository,
	productRepository *repositories.ProductRepository,
) *OrderService {
	return &OrderService{
		orderRepository:    orderRepository,
		customerRepository: customerRepository,
		productRepository:  productRepository,
	}
}

func (s *OrderService) CreateOrder(customerID uuid.UUID, products []models.OrderProduct) (*models.Order, error) {
	_, err := s.customerRepository.FindByID(customerID)
	if err != nil {
		return nil, err
	}

	productIDs := make([]uuid.UUID, len(products))
	for i, p := range products {
		productIDs[i] = p.ProductID
	}

	existingProducts, err := s.productRepository.FindByID(productIDs)
	if err != nil {
		return nil, err
	}

	productMap := make(map[uuid.UUID]*models.Product)
	for i := range existingProducts {
		productMap[existingProducts[i].ID] = &existingProducts[i]
	}

	for _, p := range products {
		prod := productMap[p.ProductID]
		if prod == nil || prod.Amount < p.Quantity {
			return nil, fmt.Errorf("product %s not available", p.ProductID)
		}
	}

	order := &models.Order{
		CustomerID:    customerID,
		Status:        "CREATED",
		OrderProducts: products,
	}

	err = s.orderRepository.WithTransaction(func(tx *gorm.DB) error {
		if err := s.orderRepository.CreateWithTx(tx, order); err != nil {
			return err
		}

		for _, p := range products {
			prod := productMap[p.ProductID]
			prod.Amount -= p.Quantity
			if err := s.productRepository.UpdateWithTx(tx, prod); err != nil {
				return err
			}
		}

		return nil
	})

	fullOrder, err := s.orderRepository.FindByID(order.ID)
	if err != nil {
		return nil, err
	}

	return fullOrder, nil

}

func (s *OrderService) UpdateOrder(customerID, orderID uuid.UUID, products []models.OrderProduct) (*models.Order, error) {
	_, err := s.customerRepository.FindByID(customerID)
	if err != nil {
		return nil, err
	}
	order, err := s.orderRepository.FindByID(orderID)
	if err != nil {
		return nil, err
	}
	if order.CustomerID != customerID {
		return nil, fmt.Errorf("order customer %s is not owned by %s", order.CustomerID, customerID)
	}

	productIDs := make([]uuid.UUID, len(products))
	for i, p := range products {
		productIDs[i] = p.ProductID
	}

	existingProducts, err := s.productRepository.FindByID(productIDs)
	if err != nil {
		return nil, err
	}

	productMap := make(map[uuid.UUID]*models.Product)
	for i := range existingProducts {
		productMap[existingProducts[i].ID] = &existingProducts[i]
	}

	for _, p := range products {
		prod := productMap[p.ProductID]
		if prod == nil || prod.Amount < p.Quantity {
			return nil, fmt.Errorf("product %s not available", p.ProductID)
		}
	}

	err = s.orderRepository.WithTransaction(func(tx *gorm.DB) error {
		log.Printf("=== Начало обновления заказа %s для клиента %s ===", orderID, customerID)

		orderProductMap := make(map[uuid.UUID]*models.OrderProduct)
		for i := range order.OrderProducts {
			op := &order.OrderProducts[i]
			orderProductMap[op.ProductID] = op
			log.Printf("Текущий продукт в заказе: %s (кол-во: %d)", op.ProductID, op.Quantity)
		}

		for _, p := range products {
			log.Printf("Обрабатываем продукт %s (новое кол-во: %d)", p.ProductID, p.Quantity)

			prod := productMap[p.ProductID]
			if prod == nil {
				log.Printf("❌ Продукт %s не найден в базе", p.ProductID)
				return fmt.Errorf("product %s not found", p.ProductID)
			}

			log.Printf("Текущий остаток на складе продукта %s: %d", prod.ID, prod.Amount)

			if existingOP, ok := orderProductMap[p.ProductID]; ok {
				totalQuantity := existingOP.Quantity + p.Quantity

				if prod.Amount < p.Quantity {
					return fmt.Errorf("not enough stock for product %s", p.ProductID)
				}

				existingOP.Quantity = totalQuantity

				prod.Amount -= p.Quantity
			} else {
				if prod.Amount < p.Quantity {
					return fmt.Errorf("not enough stock for product %s", p.ProductID)
				}
				order.OrderProducts = append(order.OrderProducts, p)
				prod.Amount -= p.Quantity
			}

			log.Printf("Сохраняем продукт %s (остаток: %d)", prod.ID, prod.Amount)
			if err := s.productRepository.UpdateWithTx(tx, prod); err != nil {
				log.Printf("❌ Ошибка при сохранении продукта %s: %v", prod.ID, err)
				return err
			}
		}

		log.Printf("=== Завершено обновление заказа %s ===", orderID)

		return s.orderRepository.UpdateWithTx(tx, order)
	})

	if err != nil {
		return nil, err
	}

	return s.orderRepository.FindByID(orderID)
}

func (s *OrderService) GetOrder(id uuid.UUID) (*models.Order, error) {
	return s.orderRepository.FindByID(id)
}
