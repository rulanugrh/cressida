package handler

import "net/http"

type VehicleService interface {
	// save vehicle into database
	CreateVehicle(w http.ResponseWriter, r *http.Request)
	// get vehicle by id
	GetVehicleByID(w http.ResponseWriter, r *http.Request)
	// get all vehicle
	GetAllVehicle(w http.ResponseWriter, r *http.Request)
	// save transporter into database
	CreateTransporter(w http.ResponseWriter, r *http.Request)
}
