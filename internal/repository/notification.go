package repository

import (
	"github.com/rulanugrh/cressida/config"
	"github.com/rulanugrh/cressida/internal/entity/domain"
)

type NotificationRepository interface {
	// GetNotificationByUserId get notification by user id
	GetNotificationByUserId(userId int) (*[]domain.Notification, error)
	// Insert notification into db
	Insert(notification domain.Notification) error
	// Update notification
	Update(orderID string, status string) error
	// Update Notification While Take Order
	UpdateWhileTakeOrder(orderID string, status string, driverName string) error
}

type notification struct {
	db *config.SDatabase
}

func NewNotificationRepository(db *config.SDatabase) NotificationRepository {
	return &notification{db: db}
}

func(n *notification) GetNotificationByUserId(userId int) (*[]domain.Notification, error) {
	var notifications []domain.Notification
	err := n.db.DB.Exec("SELECT * FROM notifications WHERE user_id = ?", userId).Find(&notifications).Error
	if err != nil {
		return nil, err
	}

	return &notifications, nil
}

func(n *notification) Insert(notification domain.Notification) error {
	err := n.db.DB.Exec("INSERT INTO notifications (user_id, content, status, order_id) VALUES (?,?,?,?)", notification.UserID, notification.Content, notification.Status, notification.OrderID).Error
	if err != nil {
		return err
	}

	return nil
}

func(n *notification) Update(orderID string, status string) error {
	err := n.db.DB.Exec("UPDATE notifications SET status = ? WHERE order_id = ?", status, orderID).Error
	if err != nil {
		return err
	}

	return nil
}

func(n *notification) UpdateWhileTakeOrder(orderID string, status string, driverName string) error {
	err := n.db.DB.Exec("UPDATE notifications SET status = ?, driver_name = ? WHERE order_id = ?", status, driverName, orderID).Error
	if err != nil {
		return err
	}

	return nil
}
