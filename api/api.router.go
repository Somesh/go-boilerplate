package api

import (
	"net/http"
	//"regexp"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (m *APIMod) InitHandlers() {

	r := mux.NewRouter()


	r.Handle("/v1/config", HandlerFunc(m.ConfigHandler))


	// Health and system metrics
	r.Handle("/status", HandlerFunc(m.Health))

	http.Handle("/metrics", promhttp.Handler())

	http.Handle("/", r)
}
