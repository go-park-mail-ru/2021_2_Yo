package prometheus

import (
	"backend/pkg/utils"
	"net/http"
	"strconv"
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
		sw := utils.NewModifiedResponse(w)
		path := r.RequestURI[:strings.IndexByte(r.RequestURI, '/')]
		if r.URL.Path != "/metrics" {
			mm.requestNow.With(prometheus.Labels{
				"method": r.Method,
				"path":   path,
			}).Inc()
		}
		start := time.Now()
		next.ServeHTTP(sw, r)
		if r.URL.Path != "/metrics" {
			mm.requestDuration.With(prometheus.Labels{
				"method": r.Method,
				"path":   path,
			}).Observe(float64(int(time.Since(start).Milliseconds())))

			mm.requestNow.With(prometheus.Labels{
				"method": r.Method,
				"path":   path,
			}).Dec()

			mm.opsProcessed.With(prometheus.Labels{
				"method": r.Method,
				"path":   path,
				"status": strconv.Itoa(sw.StatusCode),
			}).Inc()
		}
	})
}
