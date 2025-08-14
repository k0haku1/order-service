package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/k0haku1/order-service/internal/dto"
	"github.com/k0haku1/order-service/internal/service"
)

type OrderHandler struct {
	orderService *service.OrderService
}

func NewOrderHandler(orderService *service.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

func (h *OrderHandler) CreateOrder(c *fiber.Ctx) error {
	var req dto.CreateOrderRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	order, err := h.orderService.CreateOrder(req.CustomerID, req.Products)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(dto.CreateOrderResponse{
		ID:         order.ID,
		CustomerID: order.CustomerID,
		Status:     order.Status,
		Products:   order.OrderProducts,
	})
}

func (h *OrderHandler) UpdateOrder(c *fiber.Ctx) error {
	var req dto.UpdateOrderRequest
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	order := h.orderService.UpdateOrder(req.CustomerID, id, req.Products)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(dto.CreateOrderResponse{
		ID:         order.ID,
		CustomerID: order.CustomerID,
		Status:     order.Status,
		Products:   order.OrderProducts,
	})
}
func (h *OrderHandler) GetOrder(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id")
	}

	order, err := h.orderService.GetOrder(id)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(dto.CreateOrderResponse{
		ID:         order.ID,
		CustomerID: order.CustomerID,
		Status:     order.Status,
		Products:   order.OrderProducts,
	})
}
