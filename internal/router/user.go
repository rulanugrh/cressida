package router

import (
	"github.com/gorilla/mux"
	handler "github.com/rulanugrh/cressida/internal/http"
)
func UserRoute(r *mux.Router, handler handler.UserHandler) {
	subrouter := r.PathPrefix("/api/user").Subrouter()
	subrouter.HandleFunc("/register", handler.Register).Methods("POST")
	subrouter.HandleFunc("/login", handler.Login).Methods("POST")
	subrouter.HandleFunc("/getme", handler.GetMe).Methods("GET")
}
