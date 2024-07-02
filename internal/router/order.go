package router

import (
	"github.com/gorilla/mux"
	handler "github.com/rulanugrh/cressida/internal/http"
)

func OrderRoute(r *mux.Router, handler handler.OrderHandler) {
	subrouter := r.PathPrefix("/api/order").Subrouter()
	subrouter.HandleFunc("/create",handler.CreateOrder).Methods("POST")
	subrouter.HandleFunc("/history", handler.GetHistory).Methods("GET")
	subrouter.HandleFunc("/update/status", handler.UpdateStatus).Methods("PUT")
	subrouter.HandleFunc("/process", handler.GetOrderProcess).Methods("GET")
	subrouter.HandleFunc("/find/{uuid}", handler.GetOrder).Methods("GET")
}
