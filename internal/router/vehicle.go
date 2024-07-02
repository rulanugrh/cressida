package router

import (
	"github.com/gorilla/mux"
	handler "github.com/rulanugrh/cressida/internal/http"
)

func VehicleRoute(r *mux.Router, handler handler.VehicleHandler) {
	subrouter := r.PathPrefix("/api/vehicles").Subrouter()
	subrouter.HandleFunc("/add", handler.CreateVehicle).Methods("POST")
	subrouter.HandleFunc("/find/{id}", handler.GetVehicleByID).Methods("GET")
	subrouter.HandleFunc("/get", handler.GetAllVehicle).Methods("GET")

	subrouterT := r.PathPrefix("/api/transporters").Subrouter()
	subrouterT.HandleFunc("/add", handler.CreateTransporter).Methods("POST")
	subrouterT.HandleFunc("/find/{id}", handler.GetTransporterByID).Methods("GET")
	subrouterT.HandleFunc("/get", handler.GetAllTransporter).Methods("GET")
}
