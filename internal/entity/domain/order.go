package domain

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Order struct {
	ID               uuid.UUID   `json:"id" form:"id"`
	TransporterID    uint        `json:"transporter_id" form:"transporter_id"`
	UserID           uint        `json:"user_id" form:"user_id"`
	Weight           int64       `json:"weight" form:"weight"`
	Distance         int64       `json:"distance" form:"distance"`
	PickupLocation   string      `json:"pickup_location" form:"pickup_location"`
	DeliveryLocation string      `json:"delivery_location" form:"delivery_location"`
	Description      string      `json:"description" form:"description"`
	Status           string      `json:"status" form:"status"`
	Transporter      Transporter `json:"transporter" form:"transporter" gorm:"foreignKey:TransporterID;references:ID"`
	User             User        `json:"user" form:"user" gorm:"foreignKey:UserID;references:ID"`
}

type Transaction struct {
	gorm.Model
	OrderID     uuid.UUID `json:"order_id" form:"order_id"`
	TypePayment string    `json:"type_payment" form:"type_payment"`
	UserEmail   string    `json:"user_id" form:"user_id"`
	Weight      int64     `json:"weight" form:"weight"`
	Distance    int64     `json:"distance" form:"distance"`
	Subtotal    int64     `json:"subtotal" form:"subtotal"`
}
