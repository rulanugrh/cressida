package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rulanugrh/cressida/internal/helper"
	handler "github.com/rulanugrh/cressida/internal/http"
	"github.com/rulanugrh/cressida/internal/middleware"
)

func OrderRoute(r *mux.Router, handler handler.OrderHandler, observ helper.Metric) {
	subrouter := r.PathPrefix("/api/order").Subrouter()
	subrouter.Use(middleware.CORSHTTP)

	subrouter.HandleFunc("/create", observ.WrapHandler("CreateOrder", http.HandlerFunc(handler.CreateOrder))).Methods("POST")
	subrouter.HandleFunc("/history", observ.WrapHandler("GetHistoryByUserID", http.HandlerFunc(handler.GetHistory))).Methods("GET")
	subrouter.HandleFunc("/update/status", observ.WrapHandler("UpdateOrderStatus", http.HandlerFunc(handler.OrderSuccess))).Methods("PUT")
	subrouter.HandleFunc("/process", observ.WrapHandler("GetOrderProcess", http.HandlerFunc(handler.GetOrderProcess))).Methods("GET")
	subrouter.HandleFunc("/find/{uuid}", observ.WrapHandler("GetAllOrder", http.HandlerFunc(handler.GetOrder))).Methods("GET")
}
