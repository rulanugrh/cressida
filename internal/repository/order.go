package repository

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/rulanugrh/cressida/config"
	"github.com/rulanugrh/cressida/internal/entity/domain"
	"github.com/rulanugrh/cressida/internal/entity/web"
	"github.com/rulanugrh/cressida/internal/helper"
)

type OrderRepository interface {
	CreateOrder(request web.OrderRequest) (*domain.Order, error)
	GetOrder(uuid string, userID uint) (*domain.Order, error)
	SaveTransaction(request domain.Transaction) (*domain.Transaction, error)
	GetHistory(userID uint) (*[]domain.Order, error)
	OrderSuccess(uuid string) (*domain.Order, error)
	CheckOrderProcess(perPage int, page int) (*[]domain.Order, error)
	TakeOrder(uuid string) (*domain.Order, error)
}

type order struct {
	conn *config.SDatabase
	log helper.ILog
}

func NewOrderRepository(conn *config.SDatabase) OrderRepository {
	return &order{conn: conn, log: helper.NewLogger()}
}

func(o *order) CreateOrder(request web.OrderRequest) (*domain.Order, error) {
	// parsing value into domain order
	response := domain.Order{
		ID: uuid.New(),
		PickupLat: request.PickupLat,
		PickupLang: request.PickupLang,
		PickupCoordinate: request.PickupCoordinate,
		PickupAddress: request.PickupAddress,
		DropLat: request.DropLat,
		DropLang: request.DropLang,
		DropCoordinate: request.DropCoordinate,
		DropAddress: request.DropAddress,
		Distance: request.Distance,
		Weight: request.Weight,
		TransporterID: request.TransporterID,
		UserID: request.UserID,
		Description: request.Description,
		Status: "Process",
		TypePayment: request.TypePayment,
	}
	err := o.conn.DB.Create(&response).Error

	if err != nil {
		o.log.Error(fmt.Sprintf("[REPOSITORY] - [CreateOrder] %s", err.Error()))
		return nil, err
	}

	o.log.Info("[REPOSITORY] - [CreateOrder] success add order to DB")
	return &response, nil
}

func(o *order) GetOrder(uuid string, userID uint) (*domain.Order, error) {
	var response domain.Order
	err := o.conn.DB.Exec("SELECT * FROM orders WHERE id = ? AND user_id = ?", uuid, userID).Preload("Transporter").Find(&response).Error
	if err != nil {
		o.log.Error(fmt.Sprintf("[REPOSITORY] - [GetOrder] %s", err.Error()))
		return nil, err
	}

	o.log.Info("[REPOSITORY] - [GetOrder] success get order with uuid")
	return &response, nil
}

func(o *order) SaveTransaction(request domain.Transaction) (*domain.Transaction, error) {
	var response domain.Transaction
	err := o.conn.DB.Exec("INSERT INTO transactions(order_id, type_payment, user_email, weight, distance, subtotal) VALUES(?,?,?,?,?,?)",
		request.OrderID,
		request.TypePayment,
		request.UserEmail,
		request.Weight,
		request.Distance,
		request.Subtotal,
	).Scan(&response).Error

	if err != nil {
		o.log.Error(fmt.Sprintf("[REPOSITORY] - [SaveTransaction] %s", err.Error()))
		return nil, err
	}

	o.log.Info("[REPOSITORY] - [SaveTransaction] success save transaction into db")
	return &response, nil
}

func(o *order) GetHistory(userID uint) (*[]domain.Order, error) {
	var response []domain.Order
	err := o.conn.DB.Exec("SELECT * FROM transactions WHERE user_id = ?", userID).Preload("User").Preload("Transporter").Find(&response).Error
	if err != nil {
		o.log.Error(fmt.Sprintf("[REPOSITORY] - [GetHistory] %s", err.Error()))
		return nil, err
	}

	o.log.Info(fmt.Sprintf("[REPOSITORY] - [GetHistory] userID: %d success get history", userID))
	return &response, nil
}

func (o *order) OrderSuccess(uuid string) (*domain.Order, error) {
	var response domain.Order
	err := o.conn.DB.Exec("UPDATE orders SET status = ? WHERE id = ?", "Success", uuid).Error
	if err != nil {
		o.log.Error(fmt.Sprintf("[REPOSITORY] - [UpdateStatus] %s", err.Error()))
		return nil, err
	}

	err = o.conn.DB.Exec("SELECT * FROM orders WHERE = ?", uuid).Preload("Transporter").Preload("User").Find(&response).Error
	if err != nil {
		o.log.Error(fmt.Sprintf("[REPOSITORY] - [UpdateStatus] %s", err.Error()))
		return nil, err
	}

	o.log.Info(fmt.Sprintf("[REPOSITORY] - [UpdateStatus] orderID: %s success update status with status: %s", uuid, "Success"))
	return &response, nil
}

func (o *order) CheckOrderProcess(perPage int, page int) (*[]domain.Order, error) {
	var response []domain.Order
	err := o.conn.DB.Exec("SELECT * FROM orders WHERE status = ?", "Process").Offset((page - 1) * page).Limit(perPage).Preload("Transporter").Preload("User").Find(&response).Error
	if err != nil {
		o.log.Error(fmt.Sprintf("[REPOSITORY] - [CheckOrderProcess] %s", err.Error()))
		return nil, err
	}

	o.log.Info("[REPOSITORY] - [CheckOrderProcess] success get order")
	return &response, nil
}

func (o *order) TakeOrder(uuid string) (*domain.Order, error) {
	var response domain.Order
	err := o.conn.DB.Exec("UPDATE orders SET status = ? WHERE id = ?", "Confirmed", uuid).Error
	if err != nil {
		o.log.Error(fmt.Sprintf("[REPOSITORY] - [TakeOrder] %s", err.Error()))
		return nil, err
	}

	err = o.conn.DB.Exec("SELECT * FROM orders WHERE id = ?", uuid).Preload("Transporter").Preload("Transporter.Driver").Find(&response).Error
	if err != nil {
		return nil, err
	}

	o.log.Info(fmt.Sprintf("[REPOSITORY] - [UpdateStatus] orderID: %s success take order", uuid))
	return &response, nil

}
