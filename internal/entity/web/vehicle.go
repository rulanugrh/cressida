package web

type VehicleRequest struct {
	Name        string `json:"name" form:"name" validate:"required"`
	Description string `json:"description" form:"description" validate:"required"`
}

type TransporterRequest struct {
	VehicleType uint  `json:"type_vehicle" form:"type_vehicle" validate:"required"`
	DriverID    uint  `json:"driver_id" form:"driver_id" validate:"required"`
	MaxWeight   int64 `json:"max_weight" form:"max_weight" validate:"required"`
	MaxDistance int64 `json:"max_distance" form:"max_distance" validate:"required"`
}

type TransporterResponse struct {
	DriverName string `json:"driver_name" form:"driver_name"`
	Weight     int64  `json:"max_weight" form:"max_weight"`
	Distance   int64  `json:"max_distance" form:"max_distance"`
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