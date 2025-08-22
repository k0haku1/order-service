package service

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/k0haku1/order-service/internal/kafka"
	"github.com/k0haku1/order-service/internal/models"
	"github.com/k0haku1/order-service/internal/repositories"
	"gorm.io/gorm"
	"log"
)

type OrderService struct {
	orderRepository    *repositories.OrderRepository
	customerRepository *repositories.CustomerRepository
	productRepository  *repositories.ProductRepository
	dispatcher         *kafka.Dispatcher
}

func NewOrderService(
	orderRepository *repositories.OrderRepository,
	customerRepository *repositories.CustomerRepository,
	productRepository *repositories.ProductRepository,
	dispatcher *kafka.Dispatcher,
) *OrderService {
	return &OrderService{
		orderRepository:    orderRepository,
		customerRepository: customerRepository,
		productRepository:  productRepository,
		dispatcher:         dispatcher,
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

	payload := struct {
		OrderID    uuid.UUID `json:"order_id"`
		CustomerID uuid.UUID `json:"customer_id"`
		Products   []struct {
			ID       uuid.UUID `json:"id"`
			Name     string    `json:"name"`
			Quantity int32     `json:"quantity"`
		} `json:"products"`
	}{
		OrderID:    fullOrder.ID,
		CustomerID: fullOrder.CustomerID,
		Products: []struct {
			ID       uuid.UUID `json:"id"`
			Name     string    `json:"name"`
			Quantity int32     `json:"quantity"`
		}{},
	}

	for _, op := range fullOrder.OrderProducts {
		payload.Products = append(payload.Products, struct {
			ID       uuid.UUID `json:"id"`
			Name     string    `json:"name"`
			Quantity int32     `json:"quantity"`
		}{
			ID:       op.ProductID,
			Name:     op.Product.Name,
			Quantity: int32(op.Quantity),
		})
	}

	b, _ := json.Marshal(payload)

	s.dispatcher.Publish("order.created", b)

	return fullOrder, nil

}

func (s *OrderService) UpdateOrder(customerID, orderID uuid.UUID, products []models.OrderProduct) *models.Order {
	_, err := s.customerRepository.FindByID(customerID)
	if err != nil {
		return nil
	}
	order, err := s.orderRepository.FindByID(orderID)
	if err != nil {
		return nil
	}
	if order.CustomerID != customerID {
		return nil
	}

	productIDs := make([]uuid.UUID, len(products))
	for i, p := range products {
		productIDs[i] = p.ProductID
	}

	existingProducts, err := s.productRepository.FindByID(productIDs)
	if err != nil {
		return nil
	}

	productMap := make(map[uuid.UUID]*models.Product)
	for i := range existingProducts {
		productMap[existingProducts[i].ID] = &existingProducts[i]
	}

	for _, p := range products {
		prod := productMap[p.ProductID]
		if prod == nil || prod.Amount < p.Quantity {
			return nil
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
				log.Printf("Ошибка при сохранении продукта %s: %v", prod.ID, err)
				return err
			}
		}

		log.Printf("=== Завершено обновление заказа %s ===", orderID)

		return s.orderRepository.UpdateWithTx(tx, order)
	})

	if err != nil {
		return nil
	}

	fullOrder, err := s.orderRepository.FindByID(order.ID)
	if err != nil {
		return nil
	}

	orderJSON, _ := json.Marshal(fullOrder)
	s.dispatcher.Publish(
		"order.created",
		orderJSON,
	)

	return fullOrder
}

func (s *OrderService) GetOrder(id uuid.UUID) (*models.Order, error) {
	return s.orderRepository.FindByID(id)
}
