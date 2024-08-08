package helper

import (
	"net/http"
	"runtime"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metric interface {
	WrapHandler(name string, handler http.Handler) http.HandlerFunc
	CounterUser(name string, conclusion string)
	HistogramUser(name string, code string)
	CounterOrder(name string, conclusion string)
	HistogramOrder(name string, code string)
	CounterVehicle(name string, conclusion string)
	HistogramVehicle(name string, code string)
	CounterTransporter(name string, conclusion string)
	HistogramTransporter(name string, code string)
}

type metric struct {
	buckets              []float64
	reg                  prometheus.Registerer
	userCounter          *prometheus.CounterVec
	userHistogram        *prometheus.HistogramVec
	orderCounter         *prometheus.CounterVec
	orderHistogram       *prometheus.HistogramVec
	vehicleCounter       *prometheus.CounterVec
	vehicleHistogram     *prometheus.HistogramVec
	transporterCounter   *prometheus.CounterVec
	transporterHistogram *prometheus.HistogramVec
}

func NewPrometheus(register prometheus.Registerer, buckets []float64) Metric {
	totalMemory := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "total_used_memory",
		Help: "Total memory used in server",
	})

	totalCPU := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "total_used_cpu",
		Help: "Total CPU used in server",
	})

	userCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "user_",
		Name:      "counter_request_user_per_function",
		Help:      "Counter request user per function",
	}, []string{"function", "type"})

	userHistory := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "user_",
		Name:      "histogram_request_user_per_function",
		Help:      "Histogram request user per function",
	}, []string{"function", "code"})

	orderCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "order_",
		Name:      "counter_request_order_per_function",
		Help:      "Counter request order per function",
	}, []string{"function", "type"})

	orderHistory := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "order_",
		Name:      "histogram_request_order_per_function",
		Help:      "Histogram request order per function",
	}, []string{"function", "code"})

	vehicleCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "vehicle_",
		Name:      "counter_request_vehicle_per_function",
		Help:      "Counter request vehicle per function",
	}, []string{"function", "type"})

	vehicleHistogram := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "vehicle_",
		Name:      "histogram_request_vehicle_per_function",
		Help:      "Histogram request vehicle per function",
	}, []string{"function", "code"})

	transporterCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "vehicle_",
		Name:      "counter_request_transporter_per_function",
		Help:      "Counter request transporter per function",
	}, []string{"function", "type"})

	transporterHistogram := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "vehicle_",
		Name:      "histogram_request_transporter_per_function",
		Help:      "Histogram request transporter per function",
	}, []string{"function", "code"})

	// Set total Memory
	var memory runtime.MemStats
	runtime.ReadMemStats(&memory)
	totalMemory.Set(float64(memory.Alloc))

	// Set Total CPU
	totalCPU.Set(float64(runtime.NumCPU()))

	// Register our gauge with the registry.
	register.MustRegister(totalCPU, userCounter, userHistory, totalMemory, collectors.NewGoCollector(), collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))

	if buckets == nil {
		buckets = append(buckets, 0.1, 1.5, 5)
	}

	return &metric{
		buckets:              buckets,
		reg:                  register,
		userCounter:          userCounter,
		userHistogram:        userHistory,
		orderCounter:         orderCounter,
		orderHistogram:       orderHistory,
		vehicleCounter:       vehicleCounter,
		vehicleHistogram:     vehicleHistogram,
		transporterCounter:   transporterCounter,
		transporterHistogram: transporterHistogram,
	}

}

func (m *metric) CounterUser(name string, conclusion string) {
	m.userCounter.WithLabelValues(name, conclusion).Inc()
}

func (m *metric) HistogramUser(name string, code string) {
	m.userHistogram.WithLabelValues(name, code).Observe(time.Since(time.Now()).Seconds())
}

func (m *metric) CounterOrder(name string, conclusion string) {
	m.orderCounter.WithLabelValues(name, conclusion).Inc()
}

func (m *metric) HistogramOrder(name string, code string) {
	m.orderHistogram.WithLabelValues(name, code).Observe(time.Since(time.Now()).Seconds())
}

func (m *metric) CounterVehicle(name string, conclusion string) {
	m.vehicleCounter.WithLabelValues(name, conclusion).Inc()
}

func (m *metric) HistogramVehicle(name string, code string) {
	m.vehicleHistogram.WithLabelValues(name, code).Observe(time.Since(time.Now()).Seconds())
}

func (m *metric) CounterTransporter(name string, conclusion string) {
	m.transporterCounter.WithLabelValues(name, conclusion).Inc()
}

func (m *metric) HistogramTransporter(name string, code string) {
	m.transporterHistogram.WithLabelValues(name, code).Observe(time.Since(time.Now()).Seconds())
}

func (m *metric) WrapHandler(name string, handler http.Handler) http.HandlerFunc {
	reg := prometheus.WrapRegistererWith(prometheus.Labels{"handler": name}, m.reg)

	requestTotal := promauto.With(reg).NewCounterVec(prometheus.CounterOpts{
		Name:      "http_request_total",
		Help:      "Total number of HTTP requests",
	}, []string{"method", "code"})

	requestDuration := promauto.With(reg).NewHistogramVec(prometheus.HistogramOpts{
		Name:      "http_request_duration_seconds",
		Help:      "Duration of HTTP requests in seconds",
		Buckets:   m.buckets,
	}, []string{"method", "code"})

	requestSize := promauto.With(reg).NewSummaryVec(prometheus.SummaryOpts{
		Name:      "http_request_size_bytes",
		Help:      "Size of HTTP requests in bytes",
	}, []string{"method", "code"})

	responseSize := promauto.With(reg).NewSummaryVec(prometheus.SummaryOpts{
		Name:      "http_response_size_bytes",
		Help:      "Size of HTTP responses in bytes",
	}, []string{"method", "code"})

	// wrap handling for http
	base := promhttp.InstrumentHandlerCounter(
		requestTotal,
		promhttp.InstrumentHandlerDuration(
			requestDuration,
			promhttp.InstrumentHandlerRequestSize(
				requestSize,
				promhttp.InstrumentHandlerResponseSize(
					responseSize,
					handler,
				),
			),
		),
	)

	return base.ServeHTTP
}
