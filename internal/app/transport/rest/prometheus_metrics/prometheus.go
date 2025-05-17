package prometheusMetrics

import "github.com/prometheus/client_golang/prometheus"

var RequestCount = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "http_requests_total",
	Help: "Total number of HTTP requests",
}, []string{"method", "path", "status"})

var RequestDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "Histogram of response durations for HTTP requests",
		Buckets: prometheus.DefBuckets,
	},
	[]string{"method", "path"},
)

var PvzCreated = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "pvz_created_total",
	Help: "Total number of pvz successfully created",
})

var ReceptionsCreated = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "receptions_created_total",
	Help: "Total number of successfully created receptions",
})

var ProductsCreated = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "products_created_total",
	Help: "Total number of successfully created products",
})

func Init() {
	prometheus.MustRegister(RequestCount)
	prometheus.MustRegister(RequestDuration)
	prometheus.MustRegister(PvzCreated)
	prometheus.MustRegister(ReceptionsCreated)
	prometheus.MustRegister(ProductsCreated)
}
