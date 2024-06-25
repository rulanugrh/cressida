package domain

import "gorm.io/gorm"

type Vehicle struct {
	gorm.Model
	Name        string `json:"name" form:"name"`
	Description string `json:"description" form:"description"`
	Transporters []Transporter `json:"transporter" form:"transporter"`
}

type Transporter struct {
	gorm.Model
	VehicleType uint    `json:"type_vehicle" form:"type_vehicle"`
	DriverID    uint    `json:"driver_id" form:"driver_id"`
	MaxWeight   int64   `json:"max_weight" form:"max_weight"`
	MaxDistance int64   `json:"max_distance" form:"max_distance"`
	Vehicle     Vehicle `json:"vehicle" form:"vehicle" gorm:"foreignKey:VehicleType;references:ID"`
	Driver      Driver  `json:"driver" form:"driver" gorm:"foreignKey:DriverID;references:ID"`
}
