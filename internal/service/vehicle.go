package service

import (
	"github.com/rulanugrh/cressida/internal/entity/web"
	"github.com/rulanugrh/cressida/internal/middleware"
	"github.com/rulanugrh/cressida/internal/repository"
)

type VehicleService interface {
	CreateVehicle(request web.VehicleRequest) (*web.VehicleResponseCreate, error)
	FindByID(id uint) (*web.VehicleResponseGet, error)
	FindAll(perPage int, page int) (*[]web.VehicleResponseGet, error)
	CreateTransporter(request web.TransporterRequest) (*web.TransporterResponse, error)
	FindTransporterByID(id uint) (*web.TransporterResponse, error)
	FindAllTransporter(perPage int, page int) (*[]web.TransporterResponse, error) 
}

type vehicle struct {
	repository repository.VehicleRepository
	validate middleware.IValidation
}

func NewVehicleService(repository repository.VehicleRepository) VehicleService {
	return &vehicle{
		repository: repository,
		validate: middleware.NewValidation(),
	}
}