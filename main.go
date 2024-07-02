package main

import (
	"github.com/rulanugrh/cressida/config"
	handler "github.com/rulanugrh/cressida/internal/http"
	"github.com/rulanugrh/cressida/internal/repository"
	"github.com/rulanugrh/cressida/internal/router"
	"github.com/rulanugrh/cressida/internal/service"
)

func main() {
	// parsing connection for golang
	mySQL := config.InitPostgreSQL()

	// create new variabel for repository
	userRepository := repository.NewUserRepository(mySQL)
	orderRepository := repository.NewOrderRepository(mySQL)
	vehicleRepository := repository.NewVehicleRepository(mySQL)

	// create new variabel for service
	userService := service.NewUserService(userRepository)
	vehicleService := service.NewVehicleService(vehicleRepository)
	orderService := service.NewOrderService(orderRepository, vehicleRepository)

	// create new variabel for handler
	userHandler := handler.NewUserHandler(userService)
	vehicleHandler := handler.NewVehicleHandler(vehicleService)
	orderHandler := handler.NewOrderHandler(orderService)

	// parsing endpoint user
	router.RouteEndpoint(userHandler, orderHandler, vehicleHandler)
}
