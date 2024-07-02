package router

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rulanugrh/cressida/config"
	handler "github.com/rulanugrh/cressida/internal/http"
)

func RouteEndpoint(user handler.UserHandler, order handler.OrderHandler, vehicle handler.VehicleHandler) {
	cfg := config.GetConfig()
	r := mux.NewRouter().StrictSlash(true)

	UserRoute(r, user)
	OrderRoute(r, order)
	VehicleRoute(r, vehicle)

	host := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	server := http.Server{
		Addr: host,
		Handler: r,
	}

	log.Fatal(server.ListenAndServe())
}
