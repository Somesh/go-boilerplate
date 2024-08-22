package api

import (
	"html/template"
	"log"
	"net/http"

	"github.com/Somesh/go-boilerplate/common/config"
)

type APIMod struct {
	cfg    *config.Config
	client *http.Client
	log    *log.Logger

	// Health Check Http Client.
	//Using Seperate because Health Check client needs to have very low Timeout
	pingClient *http.Client
}

var templateMap map[string]*template.Template = make(map[string]*template.Template)

func InitAPIMod(cfg *config.Config) *APIMod {

	w := APIMod{
		client:     newHttpClient(cfg, false),
		pingClient: newHttpClient(cfg, true),
		cfg:        cfg,
		log:        config.GetLogger(),
	}

	return &w
}

func newHttpClient(cfg *config.Config, isPingClient bool) *http.Client {
	timeout := cfg.Server.Timeout
	if isPingClient {
		timeout = cfg.Server.PingClientTimeout
	}
	client := &http.Client{
		Timeout: timeout,
	}
	return client
}
