package api

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"regexp"
	"syscall"
	"time"

	"github.com/eapache/go-resiliency/breaker"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	tlog "github.com/opentracing/opentracing-go/log"
	"gopkg.in/tokopedia/logging.v1"

	"github.com/Somesh/go-boilerplate/common/constant"
	"github.com/Somesh/go-boilerplate/lib"
)

type Base struct {
	Status            string      `json:"status"`
	Config            interface{} `json:"config"`
	ServerProcessTime string      `json:"server_process_time"`
	ErrorMessage      []string    `json:"message_error,omitempty"`
	StatusMessage     []string    `json:"message_status,omitempty"`
}

type Response struct {
	Base
	Data interface{} `json:"data"`
}

// each handler can return the data and error, and serveHTTP can chose how to convert this
type HandlerFunc func(rw http.ResponseWriter, r *http.Request) (interface{}, error)

// FIXME: hate keeping these global
var circuitBreaker *breaker.Breaker
var excludeRegex *regexp.Regexp

func shouldTraceURL(url string) bool {
	if excludeRegex == nil || !excludeRegex.Match([]byte(url)) {
		return true
	}
	return false
}

func (fn HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	response := Response{}
	response.Base.Status = "OK"

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET,PUT,OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "X-Tkpd-UserId,Tkpd-UserId,Authorization,Origin")

	if r.Method == "OPTIONS" {
		w.WriteHeader(200)
		return
	}

	var span opentracing.Span

	ctx := r.Context()

	if shouldTraceURL(r.URL.Path) {
		span, ctx = opentracing.StartSpanFromContext(ctx, r.URL.Path)
		defer span.Finish()
	}

	// SS: TODO
	// stats.reqs.Add(1)

	// set a 10 second timeout on responses. TODO: make this come from config
	ctx, cancelFn := context.WithTimeout(ctx, 40*time.Second)
	defer cancelFn()

	r = r.WithContext(ctx)

	start := time.Now()

	var data interface{}
	var err error

	errStatus := http.StatusInternalServerError

	// open the circuit breaker and return 503 if we see too many client timeouts
	if circuitBreaker != nil {
		res := circuitBreaker.Run(func() error {
			data, err = fn(w, r)
			if err, ok := err.(net.Error); ok && err.Timeout() {
				logging.Debug.Println("network error", err)
				if span != nil {
					ext.Error.Set(span, true)
				}
				return err
			}
			return nil
		})

		if res == breaker.ErrBreakerOpen {
			err = res
			errStatus = http.StatusServiceUnavailable
		}
	} else {
		data, err = fn(w, r)
	}

	response.Base.ServerProcessTime = time.Since(start).String()

	var buf []byte

	w.Header().Set("Content-Type", "application/json")

	if data != nil && err == nil {
		response.Data = data
		if buf, err = json.Marshal(response); err == nil {
			w.Write(buf)
			return
		}
	}

	if err != nil {
		// stats.errs.Add(1)
		// stats.perr.With(prometheus.Labels{"env": environ}).Inc()
		response.Base.ErrorMessage = []string{
			err.Error(),
		}

		switch t := err.(type) {
		case lib.APIError:
			errStatus = t.Status
			response.Data = t
		case net.Error:
			if t.Timeout() {
				response.Base.ErrorMessage = []string{
					constant.ErrConnectivity.Error(),
				}
			}
		case *net.OpError:
			// t.Op == "dial",  "read" , "write"
			response.Base.ErrorMessage = []string{
				constant.ErrConnectivity.Error(),
			}
		case syscall.Errno:
			// t == syscall.ECONNREFUSED
			response.Base.ErrorMessage = []string{
				constant.ErrConnectivity.Error(),
			}
		}

		log.Println("handler error", err.Error(), r.URL.Path)
		if span != nil {
			ext.HTTPStatusCode.Set(span, uint16(errStatus))
			ext.Error.Set(span, true)
			span.LogFields(
				tlog.String("error_msg", err.Error()),
			)
		}
		w.WriteHeader(errStatus)
	}

	buf, _ = json.Marshal(response)
	logging.Debug.Println(string(buf[:]))
	w.Write(buf)
	return
}

// status check url for the eventapp
func (m *APIMod) Ping(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	//TODO: Update PING Function
	// ctx := r.Context()
	return "OK", nil
}
