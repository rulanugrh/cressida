package handler

import (
	"net/http"

	"github.com/rulanugrh/cressida/internal/middleware"
	"github.com/rulanugrh/cressida/internal/service"
)

type OrderHandler interface {
	// interface for create order
	CreateOrder(w http.ResponseWriter, r *http.Request)
	// interface for get history by user
	GetHistory(w http.ResponseWriter, r *http.Request)
	// interface for update status
	UpdateStatus(w http.ResponseWriter, r *http.Request)
	// interface for get order with status process
	GetOrderProcess(w http.ResponseWriter, r *http.Request)
	// interface for get order with uuid and userid
	GetOrder(w http.ResponseWriter, r *http.Request)
}

type order struct {
	service service.OrderService
	middleware middleware.InterfaceJWT
}

func NewOrderHandler(service service.OrderService) OrderHandler {
	return &order{
		service: service,
	}
}
