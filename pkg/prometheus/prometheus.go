package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"net/http"
	"strconv"
	"strings"
	"time"
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
		pathArr := strings.Split(path, "/")
		var resultPath string
		for index, _ := range pathArr {
			_, err := strconv.Atoi(pathArr[index])
			if err != nil {
				resultPath += pathArr[index]
				resultPath += "/"
			}
		}
		if r.URL.Path != "/metrics" {
			mm.requestNow.With(prometheus.Labels{
				"method": r.Method,
				"path":   resultPath,
			}).Inc()
		}
		start := time.Now()
		next.ServeHTTP(w, r)
		elapsed := time.Since(start)
		if r.URL.Path != "/metrics" {
			mm.requestDuration.With(prometheus.Labels{
				"method": r.Method,
				"path":   resultPath,
			}).Observe(float64(elapsed) / float64(time.Second))

			mm.requestNow.With(prometheus.Labels{
				"method": r.Method,
				"path":   resultPath,
			}).Dec()

			mm.opsProcessed.With(prometheus.Labels{
				"method": r.Method,
				"path":   resultPath,
				"status": "200",
			}).Inc()
		}
	})
}
