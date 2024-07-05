package helper

import (
	"net/http"
	"runtime"

	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metric interface {
	WrapHandler(name string, handler http.Handler) http.HandlerFunc
}

type metric struct {
	buckets []float64
	reg prometheus.Registerer
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

	// Set total Memory
	var memory runtime.MemStats
	runtime.ReadMemStats(&memory)
	totalMemory.Set(float64(memory.Alloc))

	// Set Total CPU
	totalCPU.Set(float64(runtime.NumCPU()))

	// Register our gauge with the registry.
	register.MustRegister(totalCPU, totalMemory, collectors.NewGoCollector(), collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))

	if buckets == nil {
		buckets = append(buckets, 0.1, 1.5, 5)
	}

	return &metric{
		buckets: buckets,
		reg: register,
	}

}

func (m *metric) WrapHandler(name string, handler http.Handler) http.HandlerFunc {
	reg := prometheus.WrapRegistererWith(prometheus.Labels{"handler": name}, m.reg)

	requestTotal := promauto.With(reg).NewCounterVec(prometheus.CounterOpts{
		Name: "http_request_total",
		Help: "Total number of HTTP requests",
	}, []string{"method", "code"})

	requestDuration := promauto.With(reg).NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "Duration of HTTP requests in seconds",
		Buckets: m.buckets,
	}, []string{"method", "code"})

	requestSize := promauto.With(reg).NewSummaryVec(prometheus.SummaryOpts{
		Name: "http_request_size_bytes",
		Help: "Size of HTTP requests in bytes",
	}, []string{"method", "code"})

	responseSize := promauto.With(reg).NewSummaryVec(prometheus.SummaryOpts{
		Name: "http_response_size_bytes",
		Help: "Size of HTTP responses in bytes",
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
