package register

import (
	"github.com/gorilla/mux"
	"testing"
)

func TestRegister(t *testing.T) {
	r := mux.NewRouter()
	AuthHTTPEndpoints(r, nil, nil)
	UserHTTPEndpoints(r, nil, nil, nil)
	EventHTTPEndpoints(r, nil, nil)
}
