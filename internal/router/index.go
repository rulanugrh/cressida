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
	httpSwagger "github.com/swaggo/http-swagger"
)

func RouteEndpoint(user handler.UserHandler, order handler.OrderHandler, vehicle handler.VehicleHandler, nofitication handler.NotificationHandler,registry *prometheus.Registry, observability helper.Metric) {
	cfg := config.GetConfig()


	// depend promhttp for collec with grafana and prometheus
	promHandler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})

	r := mux.NewRouter().StrictSlash(true)

	// handling for metric
	r.Handle("/metrics", promHandler).Methods("GET")
	r.HandleFunc("/api/notification/", nofitication.GetAllNotificationByUserID).Methods("GET")
	r.PathPrefix("/docs/*").Handler(httpSwagger.Handler(
		httpSwagger.URL(cfg.Server.URLDocs),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	)).Methods("GET")

	UserRoute(r, user, observability)
	OrderRoute(r, order, observability)
	VehicleRoute(r, vehicle, observability)
	SSERoute(r, nofitication)

	host := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	server := http.Server{
		Addr: host,
		Handler: r,
	}

	log.Fatal(server.ListenAndServe())
}
