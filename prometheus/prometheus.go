package prometheus

import (
	"backend/pkg/utils"
	"net/http"
	"strconv"
	log "backend/pkg/logger"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"time"
)

type metricsMiddleware struct {
	opsProcessed *prometheus.CounterVec
	requestNow *prometheus.GaugeVec
	requestDuration *prometheus.HistogramVec
}

func NewMetricsMiddleware() *metricsMiddleware {

	opsProcessed := promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "bmstusa_processed_ops_total",
		Help: "The total number of processed ops",
	}, []string{"method","path","status"})
	
	requestNow := promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "bmstusa_req_status",
		Help: "Diagram of total requests Now",
	}, []string{"method", "path"})

	requestDuration := promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "bmstusa_req_duration",
		Help: "Request Duration in seconds",
	}, []string{"method", "path"})

	return &metricsMiddleware{
		opsProcessed: opsProcessed,
		requestNow: requestNow,
		requestDuration: requestDuration,
	}
}

func (mm *metricsMiddleware) Metrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sw := utils.NewModifiedResponse(w)
		if r.URL.Path != "/metrics" {
			mm.requestNow.With(prometheus.Labels{
				"method": r.Method, 
				"path": r.RequestURI,
				}).Inc()
		}
		start := time.Now()
		next.ServeHTTP(sw,r)
		if r.URL.Path != "/metrics" {
			mm.requestDuration.With(prometheus.Labels{
				"method": r.Method, 
				"path": r.RequestURI,
				}).Observe(float64(int(time.Since(start).Milliseconds())))
			
			mm.requestNow.With(prometheus.Labels{
				"method": r.Method, 
				"path": r.RequestURI,
				}).Dec()
				
			mm.opsProcessed.With(prometheus.Labels{
				"method": r.Method, 
				"path": r.RequestURI, 
				"status": strconv.Itoa(sw.StatusCode),
				}).Inc()
		}
	})
}