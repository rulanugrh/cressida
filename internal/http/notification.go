package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rulanugrh/cressida/internal/entity/domain"
	"github.com/rulanugrh/cressida/internal/entity/web"
	"github.com/rulanugrh/cressida/internal/middleware"
	"github.com/rulanugrh/cressida/internal/service"
)

type NotificationHandler interface {
	// GetAllNotificationByUserID
	GetAllNotificationByUserID(w http.ResponseWriter, r *http.Request)
	// send streaming event while after order created
	StreamNotificationAfterCreate(w http.ResponseWriter, r *http.Request)
	// send streaming event while after order taking
	StreamNotificationAfterTaking(w http.ResponseWriter, r *http.Request)
	// send streaming event after update
	StreamNotificationAfterUpdateOrder(w http.ResponseWriter, r *http.Request)
}

type notification struct {
	service service.NotificationService
	middleware middleware.InterfaceJWT
}

func NewNotificiationHandler(service service.NotificationService) NotificationHandler {
	return &notification{service: service, middleware: middleware.NewJSONWebToken()}
}

// Get Notificiation By UserID
// @Summary endpoint for get notification by userID
// @ID get_notification_by_user_id
// @Tags notifications
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Router /api/notification/ [get]
// @Success 200 {object} web.Response
// @Failure 401 {object} web.Response
// @Failure 404 {object} web.Response
func(n *notification) GetAllNotificationByUserID(w http.ResponseWriter, r *http.Request) {
	// get user id
	userID, err := n.middleware.CheckUserID(r.Header.Get("Authorization"))
	if err != nil {
		w.WriteHeader(401)
		w.Write(web.Marshalling(web.Unauthorized("you must login first")))
		return
	}

	// parsing into service
	data, err := n.service.GetNotificationByUserID(*userID)
	if err != nil {
		w.WriteHeader(404)
		w.Write(web.Marshalling(web.NotFound("sorry data with this id not found")))
		return
	}

	w.WriteHeader(200)
	w.Write(web.Marshalling(web.Success("success get notification", data)))
}

// Server Send Event New Order
// @Summary endpoint for sse new order
// @ID sse_new_order
// @Tags notifications
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Router /api/notif-stream/new-order [get]
// @Success 200 {string} string "notification-create-order"
func(n *notification) StreamNotificationAfterCreate(w http.ResponseWriter, r *http.Request) {
	var data domain.NotificationStreamAfterCreateOrder


	// for initial streaming
	event := fmt.Sprintf("event: %s\n" + "data: \n\n", "initial")
	fmt.Fprintf(w, "%s", event)
	w.(http.Flusher).Flush()

	checkID, _ := n.middleware.CheckUserID(r.Header.Get("Authorization"))
	data.UserID[*checkID] = make(chan domain.Notification)

	// for get from data Notif
	for result := range data.UserID[*checkID] {
		response, _ := json.Marshal(result)
		event := fmt.Sprintf("event: %s\n" + "data: %s\n\n", "notification-create-order", response)

		fmt.Fprintf(w, "%s", event)
		w.(http.Flusher).Flush()
	}
}

// Server Send Event Take Order
// @Summary endpoint for sse take order
// @ID sse_take_order
// @Tags notifications
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Router /api/notif-stream/take-order [get]
// @Success 200 {string} string "notification-take-order"
func(n *notification) StreamNotificationAfterTaking(w http.ResponseWriter, r *http.Request) {
	var data domain.NotificationStreamAfterTakeOrder


	// for initial streaming
	checkID, _ := n.middleware.CheckUserID(r.Header.Get("Authorization"))
	data.UserID[*checkID] = make(chan domain.NotificationTakeOrder)
	event := fmt.Sprintf("event: %s\n" + "data: \n\n", "initial")
	fmt.Fprintf(w, "%s", event)
	w.(http.Flusher).Flush()

	// for get from data Notif
	for result := range data.UserID[*checkID] {
		response, _ := json.Marshal(result)
		event := fmt.Sprintf("event: %s\n" + "data: %s\n\n", "notification-take-order", response)

		fmt.Fprintf(w, "%s", event)
		w.(http.Flusher).Flush()
	}
}

// Server Send Event Order Success
// @Summary endpoint for sse order success
// @ID sse_order_success
// @Tags notifications
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Router /api/notif-steream/order-success [get]
// @Success 200 {string} string "notification-order-order"
func(n *notification) StreamNotificationAfterUpdateOrder(w http.ResponseWriter, r *http.Request) {
	var data domain.NotificationStreamAfterUpdateOrder


	// for initial streaming
	checkID, _ := n.middleware.CheckUserID(r.Header.Get("Authorization"))
	data.UserID[*checkID] = make(chan domain.NotificationUpdateOrder)
	event := fmt.Sprintf("event: %s\n" + "data: \n\n", "initial")
	fmt.Fprintf(w, "%s", event)
	w.(http.Flusher).Flush()

	// for get from data Notif
	for result := range data.UserID[*checkID] {
		response, _ := json.Marshal(result)
		event := fmt.Sprintf("event: %s\n" + "data: %s\n\n", "notification-update-order", response)

		fmt.Fprintf(w, "%s", event)
		w.(http.Flusher).Flush()
	}
}
