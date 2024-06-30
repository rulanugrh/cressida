package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/rulanugrh/cressida/internal/entity/web"
	"github.com/rulanugrh/cressida/internal/middleware"
	"github.com/rulanugrh/cressida/internal/service"
)

type VehicleHandler interface {
	// save vehicle into database
	CreateVehicle(w http.ResponseWriter, r *http.Request)
	// get vehicle by id
	GetVehicleByID(w http.ResponseWriter, r *http.Request)
	// get all vehicle
	GetAllVehicle(w http.ResponseWriter, r *http.Request)
	// save transporter into database
	CreateTransporter(w http.ResponseWriter, r *http.Request)
	// get transporter by id
	GetTransporterByID(w http.ResponseWriter, r *http.Request)
	// get all transporter
	GetAllTransporter(w http.ResponseWriter, r *http.Request)
}

type vehicle struct {
	service    service.VehicleService
	middleware middleware.InterfaceJWT
}

func NewVehicleHandler(service service.VehicleService) VehicleHandler {
	return &vehicle{
		service:    service,
		middleware: middleware.NewJSONWebToken(),
	}
}

// @ID create_vehicle
// @Summary save vehicle into db
// @Tags vehicles
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body web.VehicleRequest true "request body for add vehicle"
// @Router /api/vehicles/add [post]
// @Success 201 {object} web.Response
// @Failure 400 {object} web.Response
// @Failure 401 {object} web.Response
// @Failure 403 {object} web.Response
// @Failure 500 {object} web.Response
func (v *vehicle) CreateVehicle(w http.ResponseWriter, r *http.Request) {
	// check token
	if r.Header.Get("Authorization") == "" {
		w.WriteHeader(401)
		w.Write(web.Marshalling(web.Unauthorized("required token for this page")))
		return
	}

	// validate for role admin
	admin := v.middleware.ValidateAdmin(r.Header.Get("Authorization"))
	if !admin {
		w.WriteHeader(403)
		w.Write(web.Marshalling(web.Forbidden("sorry you're not admin")))
		return
	}

	// decode request body
	var request web.VehicleRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(500)
		w.Write(web.Marshalling(web.InternalServerError("error while decode requesst body")))
		return
	}

	// parsing request to service layer
	data, err := v.service.CreateVehicle(request)
	if err != nil {
		w.WriteHeader(400)
		w.Write(web.Marshalling(web.BadRequest(err.Error())))
		return
	}

	w.WriteHeader(201)
	w.Write(web.Marshalling(web.Created("success add new vehicles", data)))
}

// @ID get_vehicle_by_id
// @Summary get vehicle by id
// @Tags vehicles
// @Accept json
// @Produce json
// @Param id path int true "Vehicle ID"
// @Router /api/vehicles/find/{id} [get]
// @Success 200 {object} web.Response
// @Failure 404 {object} web.Response
func (v *vehicle) GetVehicleByID(w http.ResponseWriter, r *http.Request) {
	// parsing path request into int
	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/api/vehicles/find/"))
	if err != nil {
		w.WriteHeader(500)
		return
	}

	// get data from service layer
	data, err := v.service.FindByID(uint(id))
	if err != nil {
		w.WriteHeader(404)
		w.Write(web.Marshalling(web.NotFound("sorry vehicle with this id not found")))
		return
	}

	w.WriteHeader(200)
	w.Write(web.Marshalling(web.Success("success found vehicle", data)))
}

// @ID get_all_vehicle
// @Summary get all vehicle
// @Tags vehicles
// @Accept json
// @Produce json
// @Param per_page query int true "Per page for get data"
// @Param page query int true "Page for get data"
// @Router /api/vehicles/get [get]
// @Success 200 {object} web.Response
// @Failure 400 {object} web.Response
func (v *vehicle) GetAllVehicle(w http.ResponseWriter, r *http.Request) {
	// create query for per_page
	per_page, _ := strconv.Atoi(r.URL.Query().Get("per_page"))
	// create query for page
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))

	// get data from service layer
	data, err := v.service.FindAll(per_page, page)
	if err != nil {
		w.WriteHeader(400)
		w.Write(web.Marshalling(web.BadRequest(err.Error())))
		return
	}

	w.WriteHeader(200)
	w.Write(web.Marshalling(web.Success("success found vehicle", data)))
}

// @ID create_transporter
// @Summary save transporter into db
// @Tags transporters
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body web.TranspoterRequest true "request body for add transporter"
// @Router /api/transporters/add [post]
// @Success 201 {object} web.Response
// @Failure 401 {object} web.Response
// @Failure 403 {object} web.Response
// @Failure 400 {object} web.Response
// @Failure 500 {object} web.Response
func (v *vehicle) CreateTransporter(w http.ResponseWriter, r *http.Request) {
	// check token
	if r.Header.Get("Authorization") == "" {
		w.WriteHeader(401)
		w.Write(web.Marshalling(web.Unauthorized("required token for this page")))
		return
	}

	// validate for role admin
	driver := v.middleware.ValidateDriver(r.Header.Get("Authorization"))
	if !driver {
		w.WriteHeader(403)
		w.Write(web.Marshalling(web.Forbidden("sorry you're not driver")))
		return
	}
	// decode request body
	var request web.TransporterRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(500)
		w.Write(web.Marshalling(web.InternalServerError("error while decode requesst body")))
		return
	}

	// parsing request to service layer
	data, err := v.service.CreateTransporter(request)
	if err != nil {
		w.WriteHeader(400)
		w.Write(web.Marshalling(web.BadRequest(err.Error())))
		return
	}

	w.WriteHeader(201)
	w.Write(web.Marshalling(web.Created("success add new transporter", data)))
}

// @ID get_transporter_byd_id
// @Summary get transporter by id
// @Tags transporters
// @Accept json
// @Produce json
// @Param id path int true "Vehicle ID"
// @Router /api/transporters/find/{id} [get]
// @Success 200 {object} web.Response
// @Failure 404 {object} web.Response
func (v *vehicle) GetTransporterByID(w http.ResponseWriter, r *http.Request) {
	// parsing path request into int
	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/api/transporter/find/"))
	if err != nil {
		w.WriteHeader(500)
		return
	}

	// get data from service layer
	data, err := v.service.FindTransporterByID(uint(id))
	if err != nil {
		w.WriteHeader(404)
		w.Write(web.Marshalling(web.NotFound("sorry transporter with this id not found")))
		return
	}

	w.WriteHeader(200)
	w.Write(web.Marshalling(web.Success("success found transporter", data)))
}

// @ID get_all_transporter
// @Summary get all transporter
// @Tags transporters
// @Accept json
// @Produce json
// @Param per_page query int true "Per page for get data"
// @Param page query int true "Page for get data"
// @Router /api/transporters/get [get]
// @Success 200 {object} web.Response
// @Failure 400 {object} web.Response
func (v *vehicle) GetAllTransporter(w http.ResponseWriter, r *http.Request) {
	// create query for per_page
	per_page, _ := strconv.Atoi(r.URL.Query().Get("per_page"))
	// create query for page
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))

	// get data from service layer
	data, err := v.service.FindAllTransporter(per_page, page)
	if err != nil {
		w.WriteHeader(400)
		w.Write(web.Marshalling(web.BadRequest(err.Error())))
		return
	}

	w.WriteHeader(200)
	w.Write(web.Marshalling(web.Success("success found transporters", data)))
}
