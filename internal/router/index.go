package router

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rulanugrh/cressida/config"
	"github.com/rulanugrh/cressida/internal/helper"
	handler "github.com/rulanugrh/cressida/internal/http"
	"github.com/rulanugrh/cressida/internal/middleware"
)

func RouteEndpoint(user handler.UserHandler, order handler.OrderHandler, vehicle handler.VehicleHandler, registry *prometheus.Registry, observability helper.Metric) {
	cfg := config.GetConfig()


	// depend promhttp for collec with grafana and prometheus
	promHandler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})

	r := mux.NewRouter().StrictSlash(true)
	r.Use(middleware.CORS)

	// handling for metric
	r.Handle("/metric", promHandler).Methods("GET")

	UserRoute(r, user, observability)
	OrderRoute(r, order, observability)
	VehicleRoute(r, vehicle, observability)

	host := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	server := http.Server{
		Addr: host,
		Handler: r,
	}

	log.Fatal(server.ListenAndServe())
}
