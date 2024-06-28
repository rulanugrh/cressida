package service

import (
	"github.com/rulanugrh/cressida/internal/entity/web"
	"github.com/rulanugrh/cressida/internal/helper"
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
	vehicle    repository.VehicleRepository
	validate   middleware.IValidation
	log        helper.ILog
}

func NewOrderService(repository repository.OrderRepository, vehicle repository.VehicleRepository) OrderService {
	return &order{
		repository: repository,
		vehicle:    vehicle,
		validate:   middleware.NewValidation(),
		log:        helper.NewLogger(),
	}
}

func (o *order) CreateOrder(request web.OrderRequest) (*web.OrderResponse, error) {
	// validate request struct
	err := o.validate.Validate(request)
	if err != nil {
		o.log.Error(err)
		return nil, o.validate.ValidationMessage(err)
	}

	// create new variable for check data
	check, err := o.vehicle.FindByIDTransporter(request.TransporterID)
	if err != nil {
		o.log.Debug("transporter with this id not found")
		return nil, web.NotFound("sorry transporter with this id not found")
	}

	// check the weight if it exceeds then return an error
	if request.Weight > check.MaxWeight {
		o.log.Debug("weight has been reach max weight")
		return nil, web.BadRequest("weight exceeds the maximum weight")
	}

	// check the distance if it exceeds then return error
	if request.Distance > check.MaxDistance {
		o.log.Debug("distance has been reach max limit")
		return nil, web.BadRequest("distance exceeds the specified limits")
	}

	// input data into db
	data, errCreate := o.repository.CreateOrder(request)
	if errCreate != nil {
		o.log.Error(errCreate)
		return nil, web.BadRequest("cannot create request order")
	}

	// parsing response into new struct
	response := web.OrderResponse{
		UserName:         data.User.FName + " " + data.User.LName,
		Weight:           data.Weight,
		Distance:         data.Distance,
		PickupLocation:   data.PickupLocation,
		DeliveryLocation: data.DeliveryLocation,
		Description:      data.Description,
	}

	o.log.Info("Have been Append to DB", data.ID.String())
	return &response, nil
}

func (o *order) GetOrder(uuid string) (*web.OrderResponse, error) {
	// get order by uuid
	data, err := o.repository.GetOrder(uuid)
	if err != nil {
		o.log.Error(err)
		return nil, web.NotFound("sorry data with this uuid not found")
	}

	// parsing response into new struct
	response := web.OrderResponse{
		UserName:         data.User.FName + " " + data.User.LName,
		Weight:           data.Weight,
		Distance:         data.Distance,
		PickupLocation:   data.PickupLocation,
		DeliveryLocation: data.DeliveryLocation,
		Description:      data.Description,
	}

	o.log.Info("Found", data.ID.String())
	return &response, nil
}

func (o *order) GetHistory(userID uint) (*[]web.OrderResponse, error) {
	// get history by user id
	data, err := o.repository.GetHistory(userID)
	if err != nil {
		o.log.Error(err)
		return nil, web.BadRequest("sorry history with this id not found")
	}

	// looping data for passing into new variable response
	var response []web.OrderResponse
	for _, v := range *data {
		response = append(response, web.OrderResponse{
			Distance:         v.Distance,
			Weight:           v.Weight,
			DeliveryLocation: v.DeliveryLocation,
			PickupLocation:   v.PickupLocation,
			Description:      v.Description,
			UserName:         v.User.FName + " " + v.User.LName,
		})
	}

	o.log.Info("success get his history", userID)
	return &response, nil
}

func (o *order) UpdateStatus(uuid string, status string) (*web.OrderResponse, error) {
	// get order by uuid and update status
	data, err := o.repository.UpdateStatus(uuid, status)
	if err != nil {
		o.log.Error(err)
		return nil, web.BadRequest("sorry you cant update status with this order uuid")
	}

	// parsing response into new struct
	response := web.OrderResponse{
		UserName:         data.User.FName + " " + data.User.LName,
		Weight:           data.Weight,
		Distance:         data.Distance,
		PickupLocation:   data.PickupLocation,
		DeliveryLocation: data.DeliveryLocation,
		Description:      data.Description,
	}

	o.log.Info("Have been Update to DB", uuid)
	return &response, nil
}
