package api

import (
	"net/http"
	//"regexp"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (m *APIMod) InitHandlers() {

	r := mux.NewRouter()


	// done for migration to aliyun
	r.Handle("/status", HandlerFunc(m.Health))

	http.Handle("/metrics", promhttp.Handler())

	http.Handle("/", r)
}
