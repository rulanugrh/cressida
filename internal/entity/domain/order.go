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
	PickupLat        string      `json:"pickup_lat" form:"pickup_lat"`
	PickupLang       string      `json:"pickup_lang" form:"pickup_lang"`
	PickupCoordinate string      `json:"pickup_coordinate" form:"pickup_coordinate"`
	PickupAddress    string      `json:"pickup_address" form:"pickup_address"`
	DropLat          string      `json:"drop_lat" from:"drop_lat"`
	DropLang         string      `json:"drop_lang" from:"drop_lang"`
	DropCoordinate   string      `json:"drop_coordinate" from:"drop_coordinate"`
	DropAddress      string      `json:"drop_address" from:"drop_address"`
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
