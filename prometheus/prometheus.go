package prometheus

import (
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var RequestCount = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "request_count",
	Help: "Requests count",
},[]string{"method","path","status"})

func RegisterPrometheus(r* mux.Router) {
	r.Handle("/metrics", promhttp.Handler())
	prometheus.MustRegister(RequestCount)
}