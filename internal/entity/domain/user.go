package domain

import "gorm.io/gorm"

type User struct {
	gorm.Model
	FName    string `json:"f_name" form:"f_name"`
	LName    string `json:"l_name" form:"l_name"`
	Email    string `json:"email" form:"email" gorm:"type:unique"`
	Password string `json:"password" form:"password"`
	Address  string `json:"address" form:"address"`
	Phone    string `json:"phone" form:"phone"`
	RoleID   uint   `json:"role_id" form:"role_id"`
	Role     Role   `json:"role" form:"role" gorm:"foreignKey:RoleID;references:ID"`
}

type Role struct {
	gorm.Model
	Name        string `json:"name" form:"name"`
	Description string `json:"description" form:"description"`
	User        []User `json:"user" form:"user" gorm:"many2many:user_role"`
}

type Driver struct {
	gorm.Model
	FName        string        `json:"f_name" form:"f_name"`
	LName        string        `json:"l_name" form:"l_name"`
	Email        string        `json:"email" form:"email" gorm:"type:unique"`
	Password     string        `json:"password" form:"password"`
	Address      string        `json:"address" form:"address"`
	Phone        string        `json:"phone" form:"phone"`
	KTP          string        `json:"ktp" form:"ktp"`
	SIM          string        `json:"sim" form:"sim"`
	Profile      string        `json:"profile" form:"profile"`
	Transporters []Transporter `json:"transporter" gorm:"many2many:driver_transporter"`
}
