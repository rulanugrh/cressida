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


func(v *vehicle) CreateVehicle(request web.VehicleRequest) (*web.VehicleResponseCreate, error) {
	// validate request struct data
	err := v.validate.Validate(request)
	if err != nil {
		return nil, v.validate.ValidationMessage(err)
	}

	// insert data into repository
	data, err := v.repository.CreateVehicle(request)
	if err != nil {
		return nil, web.BadRequest("cannot insert vehicle")
	}

	// parsing rsponse into new struct
	response := web.VehicleResponseCreate{
		Name: data.Name,
		Description: data.Description,
	}

	return &response, nil
}

func(v *vehicle) FindByID(id uint) (*web.VehicleResponseGet, error) {
	// find data by id
	data, err := v.repository.FindByID(id)
	if err != nil {
		return nil, web.NotFound("sorry data with this id not found")
	}

	// looping data transporter
	var transporters []web.TransporterResponse
	for _, v := range data.Transporters {
		transporters = append(transporters, web.TransporterResponse{
			VehicleName: data.Name,
			DriverName: v.Driver.FName + " " + v.Driver.LName,
			Weight: v.MaxWeight,
			Distance: v.MaxDistance,
		})
	}

	// parsing response into new struct
	response := web.VehicleResponseGet{
		Name: data.Name,
		Description: data.Description,
		Transporter: transporters,
	}

	return &response, nil
}

func(v *vehicle) FindAll(perPage int, page int) (*[]web.VehicleResponseGet, error) {
	// find all data vehicle
	data, err := v.repository.FindAll(perPage, page)
	if err != nil {
		return nil, web.BadRequest("cannot find all data vehicles")
	}

	// looping data for parsing data
	var response []web.VehicleResponseGet

	// variabel for parsing data transporter
	var transporters []web.TransporterResponse

	for _, result := range *data {
		// looping data transporter
		for _, v := range result.Transporters {
			transporters = append(transporters, web.TransporterResponse{
				VehicleName: v.Vehicle.Name,
				DriverName: v.Driver.FName + " " + v.Driver.LName,
				Weight: v.MaxWeight,
				Distance: v.MaxDistance,
			})
		}

		vehicle := web.VehicleResponseGet{
			Name: result.Name,
			Description: result.Description,
			Transporter: transporters,
		}

		response = append(response, vehicle)
	}

	return &response, nil
}

func(v *vehicle) CreateTransporter(request web.TransporterRequest) (*web.TransporterResponse, error) {
	// validate request struct data
	err := v.validate.Validate(request)
	if err != nil {
		return nil, v.validate.ValidationMessage(err)
	}

	// insert into repository layer
	data, err := v.repository.CreateTransporter(request)
	if err != nil {
		return nil, web.BadRequest("cannot insert data")
	}

	// parsing response into new struct
	response := web.TransporterResponse{
		VehicleName: data.Vehicle.Name,
		DriverName: data.Driver.FName + " " + data.Driver.LName,
		Weight: data.MaxWeight,
		Distance: data.MaxDistance,
	}

	return &response, nil
}

func(v *vehicle) FindTransporterByID(id uint) (*web.TransporterResponse, error) {
	// find transporter by id
	data, err := v.repository.FindByIDTransporter(id)
	if err != nil {
		return nil, web.NotFound("sorry transporter not found")
	}

	// parsing response into new struct
	response := web.TransporterResponse{
		VehicleName: data.Vehicle.Name,
		DriverName: data.Driver.FName + " " + data.Driver.LName,
		Weight: data.MaxWeight,
		Distance: data.MaxDistance,
	}

	return &response, nil
}

func(v *vehicle) FindAllTransporter(perPage int, page int) (*[]web.TransporterResponse, error)  {
	// find all data transporter
	data, err := v.repository.FindAllTransporter(perPage, page)
	if err != nil {
		return nil, web.BadRequest("data cannot response")
	}

	// new variable for parsing value
	var response []web.TransporterResponse
	
	// looping data for parsing
	for _, val := range *data {
		response = append(response, web.TransporterResponse{
			VehicleName: val.Vehicle.Name,
			DriverName: val.Driver.FName + " " + val.Driver.LName,
			Weight: val.MaxWeight,
			Distance: val.MaxDistance,
		})
	}

	return &response, nil
}
