package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rulanugrh/cressida/internal/helper"
	handler "github.com/rulanugrh/cressida/internal/http"
)

func VehicleRoute(r *mux.Router, handler handler.VehicleHandler, observ helper.Metric) {
	subrouter := r.PathPrefix("/api/vehicles").Subrouter()
	subrouter.HandleFunc("/add", observ.WrapHandler("AddVehicle", http.HandlerFunc(handler.CreateVehicle))).Methods("POST")
	subrouter.HandleFunc("/find/{id}", observ.WrapHandler("GetVehicleByID", http.HandlerFunc(handler.GetVehicleByID))).Methods("GET")
	subrouter.HandleFunc("/get", observ.WrapHandler("GetAllVehicle", http.HandlerFunc(handler.GetAllVehicle))).Methods("GET")

	subrouterT := r.PathPrefix("/api/transporters").Subrouter()
	subrouterT.HandleFunc("/add", observ.WrapHandler("AddTransporter", http.HandlerFunc(handler.CreateTransporter))).Methods("POST")
	subrouterT.HandleFunc("/find/{id}", observ.WrapHandler("GetTransporterByID", http.HandlerFunc(handler.GetTransporterByID))).Methods("GET")
	subrouterT.HandleFunc("/get", observ.WrapHandler("GetAllTransporter", http.HandlerFunc(handler.GetAllTransporter))).Methods("GET")
}
