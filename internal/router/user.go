package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rulanugrh/cressida/internal/helper"
	handler "github.com/rulanugrh/cressida/internal/http"
)
func UserRoute(r *mux.Router, handler handler.UserHandler, observ helper.Metric) {
	subrouter := r.PathPrefix("/api/user").Subrouter()
	subrouter.HandleFunc("/register", observ.WrapHandler("RegisterAccount", http.HandlerFunc(handler.Register))).Methods("POST")
	subrouter.HandleFunc("/login", observ.WrapHandler("LoginAccount", http.HandlerFunc(handler.Login))).Methods("POST")
	subrouter.HandleFunc("/getme", observ.WrapHandler("GetAccount", http.HandlerFunc(handler.GetMe))).Methods("GET")
}
