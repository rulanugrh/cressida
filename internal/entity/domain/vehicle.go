package domain

import "gorm.io/gorm"

type Vehicle struct {
	gorm.Model
	Name         string        `json:"name" form:"name"`
	Description  string        `json:"description" form:"description"`
	Transporters []Transporter `json:"transporter" form:"transporter"`
}

type Transporter struct {
	gorm.Model
	DriverID    uint    `json:"driver_id" form:"driver_id"`
	Price       float64 `json:"price" form:"price"`
	VehicleType uint    `json:"type_vehicle" form:"type_vehicle"`
	MaxWeight   int64   `json:"max_weight" form:"max_weight"`
	MaxDistance int64   `json:"max_distance" form:"max_distance"`
	Driver      Driver  `json:"driver" form:"driver" gorm:"foreignKey:DriverID;references:ID"`
	Vehicle     Vehicle `json:"vehicle" form:"vehicle" gorm:"foreignKey:VehicleType;references:ID"`
}
