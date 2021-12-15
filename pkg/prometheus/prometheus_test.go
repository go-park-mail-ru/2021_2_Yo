package prometheus

import (
	"bytes"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func testHandlerFunc(w http.ResponseWriter, r *http.Request) {}

func TestUtils(t *testing.T) {
	metricsMW := NewMetricsMiddleware()
	r := mux.NewRouter()
	r.Use(metricsMW.Metrics)
	r.HandleFunc("/metrics", testHandlerFunc).Methods("GET")
	r.HandleFunc("/not_metrics", testHandlerFunc).Methods("GET")
	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "test/metrics", bytes.NewBuffer(nil))
	require.NoError(t, err)
	r.ServeHTTP(w, req)
}
