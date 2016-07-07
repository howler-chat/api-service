package http

import (
	"net/http"
	"runtime/debug"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/pressly/chi"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/net/context"
)

// Recovers from panics, logs the panic stack and returns a 500 to the caller
func PanicRecoverer(next chi.Handler) chi.Handler {
	return chi.HandlerFunc(func(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Errorf("panic: %+v", err)
				debug.PrintStack()
				http.Error(resp, http.StatusText(500), 500)
			}
		}()
		next.ServeHTTPC(ctx, resp, req)
	})
}

// Sets the 'Content-Type' to 'application/json'
func MimeJson(next chi.Handler) chi.Handler {
	return chi.HandlerFunc(func(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
		// TODO: Also set 'Accepts' header?
		resp.Header().Set("Content-Type", "application/json; charset=utf-8")
		next.ServeHTTPC(ctx, resp, req)
	})
}

var PrometheusHTTPRequestCount = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Namespace: "howler-api",
		Name:      "http_request_count",
		Help:      "The number of HTTP requests.",
	},
	[]string{"method", "endpoint"},
)

var PrometheusHTTPRequestLatency = prometheus.NewSummaryVec(
	prometheus.SummaryOpts{
		Namespace: "howler-api",
		Name:      "http_request_latency",
		Help:      "The latency of HTTP requests.",
	},
	[]string{"method", "endpoint"},
)

// Must call before using the RecordMetrics() middleware
func InitMetrics() {
	prometheus.MustRegister(PrometheusHTTPRequestCount)
	prometheus.MustRegister(PrometheusHTTPRequestLatency)
}

// Records request count, and latency
func RecordMetrics(next chi.Handler) chi.Handler {
	return chi.HandlerFunc(func(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
		// TODO: Record some metrics about the request

		PrometheusHTTPRequestCount.WithLabelValues(req.Method, req.URL.Path).Inc()
		startTime := time.Now()
		next.ServeHTTPC(ctx, resp, req)
		stopTime := time.Now()

		elapsed := stopTime.Sub(startTime)
		PrometheusHTTPRequestLatency.WithLabelValues(req.Method, req.URL.Path).
			Observe(float64(elapsed) / float64(time.Millisecond))
	})
}
