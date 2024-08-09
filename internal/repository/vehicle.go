package repository

import (
	"fmt"

	"github.com/rulanugrh/cressida/config"
	"github.com/rulanugrh/cressida/internal/entity/domain"
	"github.com/rulanugrh/cressida/internal/entity/web"
	"github.com/rulanugrh/cressida/internal/helper"
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
	log helper.ILog
}

func NewVehicleRepository(conn *config.SDatabase) VehicleRepository {
	return &vehicle{conn: conn, log: helper.NewLogger()}
}

func(v *vehicle) CreateVehicle(request web.VehicleRequest) (*domain.Vehicle, error) {
	req := domain.Vehicle{
		Name: request.Name,
		Description: request.Description,
	}

	err := v.conn.DB.Create(&req).Error
	if err != nil {
		v.log.Error(fmt.Sprintf("[REPOSITORY] - [CreateVehicle] %s", err.Error()))
		return nil, err
	}

	v.log.Info(fmt.Sprintf("[REPOSITORY] - [CreateVehicle] new data has been append: %d", req.Name))
	return &req, nil
}

func(v *vehicle) FindByID(id uint) (*domain.Vehicle, error) {
	var response domain.Vehicle
	err := v.conn.DB.Exec("SELECT * FROM vehicles WHERE id = ?", id).Preload("Transporters").Preload("Transporters.Driver").Find(&response).Error
	if err != nil {
		v.log.Error(fmt.Sprintf("[REPOSITORY] - [FindByIDVehicle] %s", err.Error()))
		return nil, err
	}

	v.log.Info(fmt.Sprintf("[REPOSITORY] - [FindByIDVehicle] id: %d, success found", id))
	return &response, nil
}

func(v *vehicle) FindAll(perPage int, page int) (*[]domain.Vehicle, error) {
	var response []domain.Vehicle
	err := v.conn.DB.Exec("SELECT * FROM vehicles").Scopes(helper.ScopesPagination(page, perPage)).Preload("Transporters").Preload("Transporters.Driver").Find(&response).Error
	if err != nil {
		v.log.Error(fmt.Sprintf("[REPOSITORY] - [FindAllVehicle] %s", err.Error()))
		return nil, err
	}

	v.log.Info("[REPOSITORY] - [FindAllVehicle] success find all vehicle")
	return &response, nil
}

func(v *vehicle) CreateTransporter(request web.TransporterRequest) (*domain.Transporter, error) {
	var response domain.Transporter
	err := v.conn.DB.Exec("INSERT INTO transporters(driver_id, vehicle_type, max_weight, max_distance, price) VALUES (?,?,?,?,?)",
		request.DriverID,
		request.VehicleType,
		request.MaxWeight,
		request.MaxDistance,
		request.Price,
	).Scan(&response).Error

	if err != nil {
		v.log.Error(fmt.Sprintf("[REPOSITORY] - [CreateTransporter] %s", err.Error()))
		return nil, err
	}

	v.log.Info(fmt.Sprintf("[REPOSITORY] - [CreateTransporter] new data has been append: %d", response.ID))
	return &response, nil
}

func(v *vehicle) FindByIDTransporter(id uint) (*domain.Transporter, error) {
	var response domain.Transporter
	err := v.conn.DB.Exec("SELECT * FROM transporter WHERE id = ?", id).Preload("Vehicle").Preload("Driver").Find(&response).Error
	if err != nil {
		v.log.Error(fmt.Sprintf("[REPOSITORY] - [FindTransporterByID] %s", err.Error()))
		return nil, err
	}

	v.log.Info(fmt.Sprintf("[REPOSITORY] - [FindTransporterByID] id: %d, success found", id))
	return &response, nil
}

func(v *vehicle) FindAllTransporter(perPage int, page int) (*[]domain.Transporter, error) {
	var response []domain.Transporter
	err := v.conn.DB.Exec("SELECT * FROM transporter").Scopes(helper.ScopesPagination(page, perPage)).Preload("Vehicle").Preload("Driver").Find(&response).Error
	if err != nil {
		v.log.Error(fmt.Sprintf("[REPOSITORY] - [FindAllTransporter] %s", err.Error()))
		return nil, err
	}

	v.log.Info("[REPOSITORY] - [FindAllTransporter] success get all data")
	return &response, nil
}
