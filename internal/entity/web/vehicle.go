package web

type VehicleRequest struct {
	Name        string `json:"name" form:"name" validate:"required"`
	Description string `json:"description" form:"description" validate:"required"`
}

type TransporterRequest struct {
	DriverID    uint    `json:"driver_id" form:"driver_id" validate:"required"`
	VehicleType uint    `json:"type_vehicle" form:"type_vehicle" validate:"required"`
	MaxWeight   int64   `json:"max_weight" form:"max_weight" validate:"required"`
	MaxDistance int64   `json:"max_distance" form:"max_distance" validate:"required"`
	Price       float64 `json:"price" form:"price" validate:"required"`
}

type TransporterResponse struct {
	VehicleName string  `json:"vehicle_name" form:"vehicle_name"`
	Weight      int64   `json:"max_weight" form:"max_weight"`
	Distance    int64   `json:"max_distance" form:"max_distance"`
	DriverName  string  `json:"driver_name" form:"driver_name"`
	Price       float64 `json:"price" form:"price"`
}

type VehicleResponseGet struct {
	Name        string                `json:"name" form:"name"`
	Description string                `json:"description" form:"description"`
	Transporter []TransporterResponse `json:"transporter" form:"transporter"`
}

type VehicleResponseCreate struct {
	Name        string `json:"name" form:"name"`
	Description string `json:"description" form:"description"`
}
