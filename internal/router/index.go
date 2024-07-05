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

func RouteEndpoint(user handler.UserHandler, order handler.OrderHandler, vehicle handler.VehicleHandler) {
	cfg := config.GetConfig()
	register := prometheus.NewRegistry()

	r := mux.NewRouter().StrictSlash(true)
	r.Handle("/metric", helper.NewPrometheus(register, nil).WrapHandler("/metrics", promhttp.HandlerFor(register, promhttp.HandlerOpts{})))

	r.Use(middleware.CORS)

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
