package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/rulanugrh/cressida/internal/entity/web"
	"github.com/rulanugrh/cressida/internal/helper"
	"github.com/rulanugrh/cressida/internal/middleware"
	"github.com/rulanugrh/cressida/internal/service"
)

type OrderHandler interface {
	// interface for create order
	CreateOrder(w http.ResponseWriter, r *http.Request)
	// interface for get history by user
	GetHistory(w http.ResponseWriter, r *http.Request)
	// interface for update status
	OrderSuccess(w http.ResponseWriter, r *http.Request)
	// interface for get order with status process
	GetOrderProcess(w http.ResponseWriter, r *http.Request)
	// interface for get order with uuid and userid
	GetOrder(w http.ResponseWriter, r *http.Request)
}

type order struct {
	service service.OrderService
	middleware middleware.InterfaceJWT
	observability helper.Metric
}

func NewOrderHandler(service service.OrderService, observability helper.Metric) OrderHandler {
	return &order{
		service: service,
		middleware: middleware.NewJSONWebToken(),
		observability: observability,
	}
}

// Create Order
// @Summary endpoint for create order
// @ID create_order
// @Tags orders
// @Accept json
// @Produce json
// @Param request body web.OrderRequest true "request body for create order"
// @Router /api/order/create [post]
// @Success 201 {object} web.Response
// @Failure 400 {object} web.Response
// @Failure 401 {object} web.Response
// @Failure 500 {object} web.Response
func(o *order) CreateOrder(w http.ResponseWriter, r *http.Request) {
	// checking user id
	id, err := o.middleware.CheckUserID(r.Header.Get("Authorization"))
	if err != nil {
		// depend counter
		o.observability.CounterOrder("jwt_token", "unauthorized")

		// depend histogram
		o.observability.HistogramOrder("create_order", "401")

		w.WriteHeader(401)
		w.Write(web.Marshalling(web.Unauthorized("cannot get user id by token")))
		return
	}

	// decode struct request
	var request web.OrderRequest
	request.UserID = *id

	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		// depend counter
		o.observability.CounterOrder("create_order", "server_error")

		// depend histogram
		o.observability.HistogramOrder("create_order", "500")

		w.WriteHeader(500)
		w.Write(web.Marshalling(web.InternalServerError("sorry cannot decode request body")))
		return
	}

	// parsing request in to service layer
	data, err := o.service.CreateOrder(request)
	if err != nil {
		// depend counter
		o.observability.CounterOrder("create_order", "failure")

		// depend histogram
		o.observability.HistogramOrder("create_order", "400")

		w.WriteHeader(400)
		w.Write(web.Marshalling(web.BadRequest("sorry cannot create request order")))
		return
	}

	// depend counter
	o.observability.CounterOrder("create_order", "success")

	// depend histogram
	o.observability.HistogramOrder("create_order", "201")
	w.WriteHeader(201)
	w.Write(web.Marshalling(web.Created("success create order", data)))
}

// Get History User
// @Summary endpoint for get history order
// @ID get_history
// @Tags orders
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Router /api/order/history [get]
// @Success 200 {object} web.Response
// @Failure 400 {object} web.Response
// @Failure 401 {object} web.Response
func(o *order) GetHistory(w http.ResponseWriter, r *http.Request) {
	// get user id
	id, err := o.middleware.CheckUserID(r.Header.Get("Authorization"))
	if err != nil {
		// depend counter
		o.observability.CounterOrder("jwt_token", "unauthorized")

		// depend histogram
		o.observability.HistogramOrder("get_history", "401")

		w.WriteHeader(401)
		w.Write(web.Marshalling(web.Unauthorized("sorry cannot get id user")))
		return
	}

	// check history user
	data, err := o.service.GetHistory(*id)
	if err != nil {
		// depend counter
		o.observability.CounterOrder("get_history", "failure")

		// depend histogram
		o.observability.HistogramOrder("get_history", "400")
		w.WriteHeader(400)
		w.Write(web.Marshalling(web.BadRequest("sorry history with this id not found")))
		return
	}

	// depend counter
	o.observability.CounterOrder("get_history", "success")

	// depend histogram
	o.observability.HistogramOrder("get_history", "200")

	w.WriteHeader(200)
	w.Write(web.Marshalling(web.Success("success get history", data)))
}

// Update Status Order
// @Summary endpoint for update status order
// @ID order_success
// @Tags orders
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param uuid query string true "query for uuid"
// @Router /api/order/update/status [put]
// @Success 200 {object} web.Response
// @Failure 400 {object} web.Response
// @Failure 403 {object} web.Response
func(o *order) OrderSuccess(w http.ResponseWriter, r *http.Request) {
	// create query for uuid
	uuid := r.URL.Query().Get("uuid")

	// checking is driver or admin
	valid := o.middleware.ValidateAdminOrDriver(r.Header.Get("Authorization"))
	if !valid {
		// depend counter
		o.observability.CounterOrder("jwt_token", "forbidden")

		// depend histogram
		o.observability.HistogramOrder("order_success", "403")
		w.WriteHeader(403)
		w.Write(web.Marshalling(web.Forbidden("sorry you are not admin or driver")))
		return
	}


	// parsing value into service layer
	data, err := o.service.OrderSuccess(uuid)
	if err != nil {
		// depend counter
		o.observability.CounterOrder("order_success", "failure")

		// depend histogram
		o.observability.HistogramOrder("order_success", "400")
		w.WriteHeader(400)
		w.Write(web.Marshalling(web.BadRequest(fmt.Sprintf("cannot update status something error: %s", err.Error()))))
		return
	}

	// depend counter
	o.observability.CounterOrder("order_success", "success")

	// depend histogram
	o.observability.HistogramOrder("order_success", "200")
	w.WriteHeader(200)
	w.Write(web.Marshalling(web.Success("success update status", data)))
}

// Get Order Process
// @Summary endpoint for get order process
// @ID get_order_process
// @Tags orders
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param per_page query int true "Per page for get data"
// @Param page query int true "Page for get data"
// @Router /api/order/process [get]
// @Success 200 {object} web.Response
// @Failure 400 {object} web.Response
// @Failure 403 {object} web.Response
func(o *order) GetOrderProcess(w http.ResponseWriter, r *http.Request) {
	// create query for per_page
	per_page, _ := strconv.Atoi(r.URL.Query().Get("per_page"))
	// create query for page
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))

	if per_page == 0 && page == 0 {
		page = 1
		per_page = 10
	}

	// checking is driver or admin
	valid := o.middleware.ValidateAdminOrDriver(r.Header.Get("Authorization"))
	if !valid {
		// depend histogram
		o.observability.HistogramOrder("jwt_token", "403")

		w.WriteHeader(403)
		w.Write(web.Marshalling(web.Forbidden("sorry you are not admin or driver")))
		return
	}

	// check from service layer
	data, err := o.service.GetOrderProcess(per_page, page)
	if err != nil {
		// depend histogram
		o.observability.HistogramOrder("get_order_process", "400")

		w.WriteHeader(400)
		w.Write(web.Marshalling(web.BadRequest("sorry cannot get order with status process")))
		return
	}

	// depend histogram
	o.observability.HistogramOrder("get_order_process", "200")

	w.WriteHeader(200)
	w.Write(web.Marshalling(web.Success("success get order with status process", data)))
}

// Get Order by UUID
// @Summary endpoint for get order process
// @ID get_order_by_uuid
// @Tags orders
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param uuid query string true "UUID order data"
// @Router /api/order/find/{uuid} [get]
// @Success 200 {object} web.Response
// @Failure 400 {object} web.Response
// @Failure 401 {object} web.Response
func(o *order) GetOrder(w http.ResponseWriter, r *http.Request) {
	// create query for uuid
	uuid := r.URL.Query().Get("uuid")

	// checking user id
	id, err := o.middleware.CheckUserID(r.Header.Get("Authorization"))
	if err != nil {
		// depend histogram
		o.observability.HistogramOrder("jwt_token", "401")

		w.WriteHeader(401)
		w.Write(web.Marshalling(web.Unauthorized("sorry cannot get user id, you must login first")))
		return
	}

	// checking data from service layer
	data, err := o.service.GetOrder(uuid, *id)
	if err != nil {
		// depend histogram
		o.observability.HistogramOrder("get_order", "400")

		w.WriteHeader(400)
		w.Write(web.Marshalling(web.BadRequest("cannot get order with this uuid")))
		return
	}

	// depend histogram
	o.observability.HistogramOrder("get_order", "200")

	w.WriteHeader(200)
	w.Write(web.Marshalling(web.Success("success get order", data)))
}
