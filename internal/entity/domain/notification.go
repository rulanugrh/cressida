package domain

import "gorm.io/gorm"

type Notification struct {
	gorm.Model
	UserID     uint   `json:"user_id" form:"user_id"`
	Content    string `json:"content" form:"content"`
	Status     string `json:"status" form:"status"`
	OrderID    string `json:"order_id" form:"order_id"`
	DriverName string `json:"driver_name" form:"driver_name"`
}

type NotificationTakeOrder struct {
	Content    string `json:"content" form:"content"`
	DriverName string `json:"driver_name" form:"driver_name"`
	Status     string `json:"status" form:"status"`
	OrderID    string `json:"order_id" form:"order_id"`
}

type NotificationUpdateOrder struct {
	Content    string `json:"content" form:"content"`
	Status     string `json:"status" form:"status"`
	OrderID    string `json:"order_id" form:"order_id"`
}

type NotificationStreamAfterCreateOrder struct {
	UserID map[uint]chan Notification
}

type NotificationStreamAfterTakeOrder struct {
	UserID map[uint]chan NotificationTakeOrder
}

type NotificationStreamAfterUpdateOrder struct {
	UserID map[uint]chan NotificationUpdateOrder
}
