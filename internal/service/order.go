package service

import "github.com/rulanugrh/cressida/internal/entity/web"

type OrderService interface {
	// Interface for create new order
	CreateOrder(request web.OrderRequest) (*web.OrderResponse, error)
}