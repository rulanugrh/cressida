package repository

import (
	"github.com/google/uuid"
	"github.com/rulanugrh/cressida/config"
	"github.com/rulanugrh/cressida/internal/entity/domain"
	"github.com/rulanugrh/cressida/internal/entity/web"
)

type OrderRepository interface {
	CreateOrder(request web.OrderRequest) (*domain.Order, error)
	GetOrder(uuid string) (*domain.Order, error)
	SaveTransaction(request domain.Transaction) (*domain.Transaction, error)
	GetHistory(userID uint) (*[]domain.Order, error)
}

type order struct {
	conn *config.SDatabase
}

func NewOrderRepository(conn *config.SDatabase) OrderRepository {
	return &order{conn: conn}
}

func(o *order) CreateOrder(request web.OrderRequest) (*domain.Order, error) {
	var response domain.Order
	err := o.conn.DB.Exec("INSERT INTO orders(id, transporter_id, user_id, weight, distance, pickup_location, delivery_location, description, status) VALUES (?,?,?,?,?,?,?,?,?)",
		uuid.New().String(),
		request.TransporterID,
		request.UserID,
		request.Weight,
		request.Distance,
		request.PickupLocation,
		request.DeliveryLocation,
		request.Description,
		"pending",
	).Scan(&response).Error

	if err != nil {
		return nil, err
	}

	return &response, nil
}

func(o *order) GetOrder(uuid string) (*domain.Order, error) {
	var response domain.Order
	err := o.conn.DB.Exec("SELECT * FROM orders WHERE id = ?", uuid).Preload("Transporter").Find(&response).Error
	if err != nil {
		return nil, err
	}

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
		return nil, err
	}

	return &response, nil
}

func(o *order) GetHistory(userID uint) (*[]domain.Order, error) {
	var response []domain.Order
	err := o.conn.DB.Exec("SELECT * FROM transactions WHERE user_id = ?", userID).Preload("User").Preload("Transporter").Find(&response).Error
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (o *order) UpdateStatus(uuid string, status string) (*domain.Order, error) {
	var response domain.Order
	err := o.conn.DB.Exec("UPDATE orders SET status = ? WHERE id = ?", status, uuid).Preload("Transporter").Preload("User").Find(&response).Error
	if err != nil {
		return nil, err
	}

	return &response, nil
}