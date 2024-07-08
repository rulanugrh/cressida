package main

import (
	"context"
	"log"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rulanugrh/cressida/config"
	"github.com/rulanugrh/cressida/internal/helper"
	handler "github.com/rulanugrh/cressida/internal/http"
	"github.com/rulanugrh/cressida/internal/repository"
	"github.com/rulanugrh/cressida/internal/router"
	"github.com/rulanugrh/cressida/internal/service"
)

func main() {
	// parsing connection opentelemetry
	opentelemetry, err := helper.InitTracer()
	if err != nil {
		log.Fatalf("Error while trace provider: %v", err)
	}

	// defer function for checking opentelemetry while running and trace function
	defer func(){
		if err := opentelemetry.Shutdown(context.Background()); err != nil {
			log.Fatalf("Error while trace provider: %v", err)
		}
	}()

	// parsing connection for golang
	mySQL := config.InitPostgreSQL()

	// connection for register
	register := prometheus.NewRegistry()

	// register for observability
	observability := helper.NewPrometheus(register, nil)

	// create new variabel for repository
	userRepository := repository.NewUserRepository(mySQL)
	orderRepository := repository.NewOrderRepository(mySQL)
	vehicleRepository := repository.NewVehicleRepository(mySQL)
	notificationRepository := repository.NewNotificationRepository(mySQL)

	// create new variabel for service
	userService := service.NewUserService(userRepository)
	vehicleService := service.NewVehicleService(vehicleRepository)
	notificationService := service.NewNotificationService(notificationRepository)
	orderService := service.NewOrderService(orderRepository, vehicleRepository, notificationRepository)

	// create new variabel for handler
	userHandler := handler.NewUserHandler(userService, observability)
	vehicleHandler := handler.NewVehicleHandler(vehicleService, observability)
	orderHandler := handler.NewOrderHandler(orderService, observability)
	notificationHandler := handler.NewNotificiationHandler(notificationService)

	// parsing endpoint user
	router.RouteEndpoint(userHandler, orderHandler, vehicleHandler, notificationHandler, register, observability)
}
