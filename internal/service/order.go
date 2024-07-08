package service

import (
	"context"
	"fmt"

	"github.com/rulanugrh/cressida/internal/entity/domain"
	"github.com/rulanugrh/cressida/internal/entity/web"
	"github.com/rulanugrh/cressida/internal/helper"
	"github.com/rulanugrh/cressida/internal/middleware"
	"github.com/rulanugrh/cressida/internal/repository"
)

type OrderService interface {
	// interface for create new order
	CreateOrder(request web.OrderRequest) (*web.OrderResponse, error)
	// interface for get order
	GetOrder(uuid string, userID uint) (*web.OrderResponse, error)
	// interface for get history
	GetHistory(userID uint) (*[]web.OrderResponse, error)
	// interfce for get order with status process
	GetOrderProcess(perPage int, page int) (*[]web.OrderResponse, error)
	// interface for update status
	OrderSuccess(uuid string) (*web.OrderResponse, error)
}

type order struct {
	repository   repository.OrderRepository
	vehicle      repository.VehicleRepository
	validate     middleware.IValidation
	log          helper.ILog
	trace        helper.IOpenTelemetry
	notification repository.NotificationRepository
}

func NewOrderService(repository repository.OrderRepository, vehicle repository.VehicleRepository, notification repository.NotificationRepository) OrderService {
	return &order{
		repository:   repository,
		vehicle:      vehicle,
		validate:     middleware.NewValidation(),
		log:          helper.NewLogger(),
		trace:        helper.NewOpenTelemetry(),
		notification: notification,
	}
}

func (o *order) CreateOrder(request web.OrderRequest) (*web.OrderResponse, error) {
	// span for tracing request this endpoint
	span := o.trace.StartTracer(context.Background(), "CreateOrder")
	defer span.End()

	// validate request struct
	err := o.validate.Validate(request)
	if err != nil {
		o.log.Error(fmt.Sprintf("[SERVICE] - [CreateOrder] Error while validate request: %s", err.Error()))
		return nil, o.validate.ValidationMessage(err)
	}

	// create new variable for check data
	check, err := o.vehicle.FindByIDTransporter(request.TransporterID)
	if err != nil {
		o.log.Debug(fmt.Sprintf("[SERVICE] - [CreateOrder] id: %d, transporter with this id not found", request.TransporterID))
		return nil, web.NotFound("sorry transporter with this id not found")
	}

	// check the weight if it exceeds then return an error
	if request.Weight > check.MaxWeight {
		o.log.Debug(fmt.Sprintf("[SERVICE] - [CreateOrder] %d kg, weight has been reach max weight", request.Weight))
		return nil, web.BadRequest("weight exceeds the maximum weight")
	}

	// check the distance if it exceeds then return error
	if request.Distance > check.MaxDistance {
		o.log.Debug(fmt.Sprintf("[SERVICE] - [CreateOrder] %d km, distance has been reach max limit", request.Distance))
		return nil, web.BadRequest("distance exceeds the specified limits")
	}

	// input data into db
	data, errCreate := o.repository.CreateOrder(request)
	if errCreate != nil {
		o.log.Error(fmt.Sprintf("[SERVICE] - [CreateOrder] Error while input into db: %s", errCreate.Error()))
		return nil, web.BadRequest("cannot create request order")
	}

	// parsing response into new struct
	response := web.OrderResponse{
		UserName:         data.User.FName + " " + data.User.LName,
		Weight:           data.Weight,
		Distance:         data.Distance,
		PickupLat:        data.PickupLat,
		PickupLang:       data.PickupLang,
		PickupCoordinate: data.PickupCoordinate,
		PickupAddress:    data.PickupAddress,
		DropLat:          data.DropLat,
		DropLang:         data.DropLang,
		DropCoordinate:   data.DropCoordinate,
		DropAddress:      data.DropAddress,
		Description:      data.Description,
	}

	// create go routine for notification after create order
	go o.notificationAfterCreateOrder(data)

	o.log.Info(fmt.Sprintf("[SERVICE] - [CreateOrder] new order id %s append into db", data.ID.String()))
	return &response, nil
}

func (o *order) GetOrder(uuid string, userID uint) (*web.OrderResponse, error) {
	// span for tracing request this endpoint
	span := o.trace.StartTracer(context.Background(), "GetOrder")
	defer span.End()

	// get order by uuid
	data, err := o.repository.GetOrder(uuid, userID)
	if err != nil {
		o.log.Error(fmt.Sprintf("[SERVICE] - [GetOrder] Error while get data in db: %s", err.Error()))
		return nil, web.NotFound("sorry data with this uuid not found")
	}

	// parsing response into new struct
	response := web.OrderResponse{
		UserName:         data.User.FName + " " + data.User.LName,
		Weight:           data.Weight,
		Distance:         data.Distance,
		PickupLat:        data.PickupLat,
		PickupLang:       data.PickupLang,
		PickupCoordinate: data.PickupCoordinate,
		PickupAddress:    data.PickupAddress,
		DropLat:          data.DropLat,
		DropLang:         data.DropLang,
		DropCoordinate:   data.DropCoordinate,
		DropAddress:      data.DropAddress,
		Description:      data.Description,
	}

	o.log.Info(fmt.Sprintf("[SERVICE] - [GetOrder] order id %s success found", uuid))
	return &response, nil
}

func (o *order) GetHistory(userID uint) (*[]web.OrderResponse, error) {
	// span for tracing request this endpoint
	span := o.trace.StartTracer(context.Background(), "GetHistory")
	defer span.End()

	// get history by user id
	data, err := o.repository.GetHistory(userID)
	if err != nil {
		o.log.Error(fmt.Sprintf("[SERVICE] - [GetHistory] Error while get data in db: %s", err.Error()))
		return nil, web.BadRequest("sorry history with this id not found")
	}

	// looping data for passing into new variable response
	var response []web.OrderResponse
	for _, v := range *data {
		response = append(response, web.OrderResponse{
			Distance:         v.Distance,
			Weight:           v.Weight,
			PickupLat:        v.PickupLat,
			PickupLang:       v.PickupLang,
			PickupCoordinate: v.PickupCoordinate,
			PickupAddress:    v.PickupAddress,
			DropLat:          v.DropLat,
			DropLang:         v.DropLang,
			DropCoordinate:   v.DropCoordinate,
			DropAddress:      v.DropAddress,
			Description:      v.Description,
			UserName:         v.User.FName + " " + v.User.LName,
		})
	}

	o.log.Info(fmt.Sprintf("[SERVICE] - [GetHistory] userID: %d success get history", userID))
	return &response, nil
}

func (o *order) OrderSuccess(uuid string) (*web.OrderResponse, error) {
	// span for tracing request this endpoint
	span := o.trace.StartTracer(context.Background(), "UpdateStatus")
	defer span.End()

	// get order by uuid and update status
	data, err := o.repository.OrderSuccess(uuid)
	if err != nil {
		o.log.Error(fmt.Sprintf("[SERVICE] - [UpdateStatus] Error while get data in db: %s", err.Error()))
		return nil, web.BadRequest("sorry you cant update status with this order uuid")
	}

	// parsing response into new struct
	response := web.OrderResponse{
		UserName:         data.User.FName + " " + data.User.LName,
		Weight:           data.Weight,
		Distance:         data.Distance,
		PickupLat:        data.PickupLat,
		PickupLang:       data.PickupLang,
		PickupCoordinate: data.PickupCoordinate,
		PickupAddress:    data.PickupAddress,
		DropLat:          data.DropLat,
		DropLang:         data.DropLang,
		DropCoordinate:   data.DropCoordinate,
		DropAddress:      data.DropAddress,
		Description:      data.Description,
	}

	tax := data.Transporter.Price * (0.11 * data.Transporter.Price)
	_, _ = o.repository.SaveTransaction(domain.Transaction{
		OrderID:     data.ID,
		Weight:      data.Weight,
		Distance:    data.Distance,
		UserEmail:   data.User.Email,
		Subtotal:    int64(data.Transporter.Price) + int64(tax),
		TypePayment: data.TypePayment,
	})

	// start for notification after update order
	go o.notificationUpdateStatusOrder(data)

	o.log.Info(fmt.Sprintf("[SERVICE] - [UpdateStatus] order id %s success update status: %s", uuid, data.Status))
	return &response, nil
}

func (o *order) GetOrderProcess(perPage int, page int) (*[]web.OrderResponse, error) {
	// span for tracing request this endpoint
	span := o.trace.StartTracer(context.Background(), "GetOrderProcess")
	defer span.End()

	// get history by user id
	data, err := o.repository.CheckOrderProcess(perPage, page)
	if err != nil {
		o.log.Error(fmt.Sprintf("[SERVICE] - [GetOrderProcess] Error while get data in db: %s", err.Error()))
		return nil, web.BadRequest("sorry order not found")
	}

	// looping data for passing into new variable response
	var response []web.OrderResponse
	for _, v := range *data {
		response = append(response, web.OrderResponse{
			Distance:         v.Distance,
			Weight:           v.Weight,
			PickupLat:        v.PickupLat,
			PickupLang:       v.PickupLang,
			PickupCoordinate: v.PickupCoordinate,
			PickupAddress:    v.PickupAddress,
			DropLat:          v.DropLat,
			DropLang:         v.DropLang,
			DropCoordinate:   v.DropCoordinate,
			DropAddress:      v.DropAddress,
			Description:      v.Description,
			UserName:         v.User.FName + " " + v.User.LName,
		})
	}

	o.log.Info("[SERVICE] - [GetOrderProcess] success get all order process")
	return &response, nil
}

func (o *order) TakeOrder(uuid string) error {
	// span for tracing request this endpoint
	span := o.trace.StartTracer(context.Background(), "TakeOrder")
	defer span.End()

	// create request for take order
	data, err := o.repository.TakeOrder(uuid)
	if err != nil {
		return web.BadRequest(err.Error())
	}

	// start for notification after take order
	go o.notificationWhileDriverTakeOrder(data)

	return nil
}

func (o *order) notificationAfterCreateOrder(req *domain.Order) {
	// span for tracing request this endpoint
	span := o.trace.StartTracer(context.Background(), "notificationAfterCreateOrder")
	defer span.End()

	var data domain.NotificationStreamAfterCreateOrder
	notification := domain.Notification{
		UserID:  req.UserID,
		Content: "Transaction Succesfull",
		Status:  req.Status,
		OrderID: req.ID.String(),
	}

	channel, oke := data.UserID[req.Transporter.DriverID]
	if oke {
		channel <- domain.Notification{
			UserID:  req.Transporter.DriverID,
			Content: "New order incoming",
			Status:  req.Status,
			OrderID: req.ID.String(),
		}
	}

	_ = o.notification.Insert(notification)
}

func (o *order) notificationUpdateStatusOrder(req *domain.Order) {
	// span for tracing request this endpoint
	span := o.trace.StartTracer(context.Background(), "notificationUpdateStatusOrder")
	defer span.End()

	var data domain.NotificationStreamAfterUpdateOrder
	channel, oke := data.UserID[req.UserID]
	if oke {
		channel <- domain.NotificationUpdateOrder{
			Content: "Order success, thank you for ordering me",
			Status:  req.Status,
			OrderID: req.ID.String(),
		}
	}
	_ = o.notification.Update(req.ID.String(), req.Status)
}

func (o *order) notificationWhileDriverTakeOrder(req *domain.Order) {
	// span for tracing request this endpoint
	span := o.trace.StartTracer(context.Background(), "notificationWhileDriverTakeOrder")
	defer span.End()

	var data domain.NotificationStreamAfterTakeOrder
	channel, oke := data.UserID[req.UserID]
	if oke {
		channel <- domain.NotificationTakeOrder{
			Content:    "The order was successfully taken by the driver",
			Status:     req.Status,
			OrderID:    req.ID.String(),
			DriverName: req.Transporter.Driver.FName + req.Transporter.Driver.LName,
		}
	}

	_ = o.notification.UpdateWhileTakeOrder(req.ID.String(), req.Status, req.Transporter.Driver.FName+" "+req.Transporter.Driver.LName)
}
