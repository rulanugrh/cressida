package repository

import (
	"github.com/rulanugrh/cressida/config"
	"github.com/rulanugrh/cressida/internal/entity/domain"
	"github.com/rulanugrh/cressida/internal/entity/web"
)

type VehicleRepository interface {
	CreateVehicle(request web.VehicleRequest) (*domain.Vehicle, error)
	FindByID(id uint) (*domain.Vehicle, error)
	FindAll(perPage int, page int) (*[]domain.Vehicle, error)
	CreateTransporter(request web.TransporterRequest) (*domain.Transporter, error)
	FindByIDTransporter(id uint) (*domain.Transporter, error)
	FindAllTransporter(perPage int, page int) (*[]domain.Transporter, error)
}

type vehicle struct {
	conn *config.SDatabase
}

func NewVehicleRepository(conn *config.SDatabase) VehicleRepository {
	return &vehicle{conn: conn}
}

func(v *vehicle) CreateVehicle(request web.VehicleRequest) (*domain.Vehicle, error) {
	var response domain.Vehicle
	err := v.conn.DB.Exec("INSERT INTO vehicles(name, description) VALUES(?,?)", request.Name, request.Description).Scan(&response).Error
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func(v *vehicle) FindByID(id uint) (*domain.Vehicle, error) {
	var response domain.Vehicle
	err := v.conn.DB.Exec("SELECT * FROM vehicles WHERE id = ?", id).Preload("Transporters").Find(&response).Error
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func(v *vehicle) FindAll(perPage int, page int) (*[]domain.Vehicle, error) {
	var response []domain.Vehicle
	err := v.conn.DB.Exec("SELECT * FROM vehicles").Offset((page - 1) * perPage).Limit(perPage).Find(&response).Error
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func(v *vehicle) CreateTransporter(request web.TransporterRequest) (*domain.Transporter, error) {
	var response domain.Transporter
	err := v.conn.DB.Exec("INSERT INTO transporters(vehicle_type, max_weight, max_distance) VALUES (?,?,?,?)",
		request.VehicleType,
		request.MaxWeight,
		request.MaxDistance,
	).Scan(&response).Error

	if err != nil {
		return nil, err
	}

	return &response, nil
}

func(v *vehicle) FindByIDTransporter(id uint) (*domain.Transporter, error) {
	var response domain.Transporter
	err := v.conn.DB.Exec("SELECT * FROM transporter WHERE id = ?", id).Preload("Vehicle").Preload("Driver").Find(&response).Error
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func(v *vehicle) FindAllTransporter(perPage int, page int) (*[]domain.Transporter, error) {
	var response []domain.Transporter
	err := v.conn.DB.Exec("SELECT * FROM transporter").Offset((page - 1) * page).Limit(perPage).Preload("Vehicle").Preload("Driver").Find(&response).Error
	if err != nil {
		return nil, err
	}

	return &response, nil
}
