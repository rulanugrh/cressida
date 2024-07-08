package router

import (
	"github.com/gorilla/mux"
	handler "github.com/rulanugrh/cressida/internal/http"
	"github.com/rulanugrh/cressida/internal/middleware"
)

func SSERoute(r *mux.Router, handler handler.NotificationHandler) {
	subrouter := r.PathPrefix("/api/notif-stream").Subrouter()
	subrouter.Use(middleware.CORSStreaming)

	subrouter.HandleFunc("/new-order", handler.StreamNotificationAfterCreate).Methods("GET")
	subrouter.HandleFunc("/take-order", handler.StreamNotificationAfterTaking).Methods("GET")
	subrouter.HandleFunc("/order-success", handler.StreamNotificationAfterUpdateOrder).Methods("GET")
}
