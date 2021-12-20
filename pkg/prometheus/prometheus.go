package prometheus

import (
	"net/http"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type metricsMiddleware struct {
	opsProcessed    *prometheus.CounterVec
	requestNow      *prometheus.GaugeVec
	requestDuration *prometheus.HistogramVec
}

func NewMetricsMiddleware() *metricsMiddleware {

	opsProcessed := promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "bmstusa_processed_ops_total",
		Help: "The total number of processed ops",
	}, []string{"method", "path", "status"})

	requestNow := promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "bmstusa_req_status",
		Help: "Diagram of total requests Now",
	}, []string{"method", "path"})

	requestDuration := promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "bmstusa_req_duration",
		Help: "Request Duration in seconds",
	}, []string{"method", "path"})

	return &metricsMiddleware{
		opsProcessed:    opsProcessed,
		requestNow:      requestNow,
		requestDuration: requestDuration,
	}
}

func (mm *metricsMiddleware) Metrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.RequestURI
		pathArr := strings.Split(path, "?")
		if r.URL.Path != "/metrics" {
			mm.requestNow.With(prometheus.Labels{
				"method": r.Method,
				"path":   pathArr[0],
			}).Inc()
		}
		start := time.Now()
		next.ServeHTTP(w, r)
		elapsed := time.Since(start)
		if r.URL.Path != "/metrics" {
			mm.requestDuration.With(prometheus.Labels{
				"method": r.Method,
				"path":   pathArr[0],
			}).Observe(float64(elapsed)/float64(time.Second))

			mm.requestNow.With(prometheus.Labels{
				"method": r.Method,
				"path":   pathArr[0],
			}).Dec()

			mm.opsProcessed.With(prometheus.Labels{
				"method": r.Method,
				"path":   pathArr[0],
				"status": "200",
			}).Inc()
		}
	})
}