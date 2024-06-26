package service

import (
	"github.com/rulanugrh/cressida/internal/entity/web"
	"github.com/rulanugrh/cressida/internal/middleware"
	"github.com/rulanugrh/cressida/internal/repository"
)

type OrderService interface {
	// interface for create new order
	CreateOrder(request web.OrderRequest) (*web.OrderResponse, error)
	// interface for get order
	GetOrder(uuid string) (*web.OrderResponse, error)
	// interface for get history
	GetHistory(userID uint) (*[]web.OrderResponse, error)
}

type order struct {
	repository repository.OrderRepository
	vehicle repository.VehicleRepository
	validate middleware.IValidation
}

func NewOrderService(repository repository.OrderRepository, vehicle repository.VehicleRepository) OrderService {
	return &order{
		repository: repository,
		vehicle: vehicle,
		validate: middleware.NewValidation(),
	}
}

func(o *order) CreateOrder(request web.OrderRequest) (*web.OrderResponse, error) {
	// validate request struct
	err := o.validate.Validate(request)
	if err != nil {
		return nil, o.validate.ValidationMessage(err)
	}

	// create new variable for check data
	check, err := o.vehicle.FindByIDTransporter(request.TransporterID)
	if err != nil {
		return nil, web.NotFound("sorry transporter with this id not found")
	}

	// check the weight if it exceeds then return an error
	if request.Weight > check.MaxWeight {
		return nil, web.BadRequest("weight exceeds the maximum weight")
	}

	// check the distance if it exceeds then return error
	if request.Distance > check.MaxDistance {
		return nil, web.BadRequest("distance exceeds the specified limits")
	}

	// input data into db
	data, err := o.repository.CreateOrder(request)
	if err != nil {
		return nil, web.BadRequest("cannot create request order")
	}


	// parsing response into new struct
	response := web.OrderResponse{
		UserName: data.User.FName + " " + data.User.LName,
		Weight: data.Weight,
		Distance: data.Distance,
		PickupLocation: data.PickupLocation,
		DeliveryLocation: data.DeliveryLocation,
		Description: data.Description,
	}

	return &response, nil
}

func(o *order) GetOrder(uuid string) (*web.OrderResponse, error) {
	// get order by uuid
	data, err := o.repository.GetOrder(uuid)
	if err != nil {
		return nil, web.NotFound("sorry data with this uuid not found")
	}

	// parsing response into new struct
	response := web.OrderResponse{
		UserName: data.User.FName + " " + data.User.LName,
		Weight: data.Weight,
		Distance: data.Distance,
		PickupLocation: data.PickupLocation,
		DeliveryLocation: data.DeliveryLocation,
		Description: data.Description,
	}
	
	return &response, nil
}

func(o *order) GetHistory(userID uint) (*[]web.OrderResponse, error) {
	// get history by user id
	data, err := o.repository.GetHistory(userID)
	if err != nil {
		return nil, web.BadRequest("sorry history with this id not found")
	}

	// looping data for passing into new variable response
	var response []web.OrderResponse
	for _, v := range *data {
		response = append(response, web.OrderResponse{
			Distance: v.Distance,
			Weight: v.Weight,
			DeliveryLocation: v.DeliveryLocation,
			PickupLocation: v.PickupLocation,
			Description: v.Description,
			UserName: v.User.FName + " " + v.User.LName,
		})
	}

	return &response, nil
}
