package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/k0haku1/order-service/database"
	"github.com/k0haku1/order-service/internal/handlers"
	"github.com/k0haku1/order-service/internal/models"
	"github.com/k0haku1/order-service/internal/repositories"
	"github.com/k0haku1/order-service/internal/service"
	"log"
)

func main() {

	db, err := database.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	if err := db.AutoMigrate(&models.Customer{}, &models.Order{}, &models.Product{}, &models.OrderProduct{}); err != nil {
		log.Fatal(err)
	}

	orderRepository := repositories.NewOrderRepository(db)
	orderService := service.NewOrderService(orderRepository)
	orderHandler := handlers.NewOrderHandler(orderService)

	app := fiber.New()
	app.Post("orders/create", orderHandler.CreateOrder)
	app.Get("orders/:id", orderHandler.GetOrder)

	err = app.Listen(":8080")
	if err != nil {
		return
	}

}
