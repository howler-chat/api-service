package service

import (
	"net/http"
	"runtime/debug"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/howler-chat/api-service/api"
	"github.com/howler-chat/api-service/metrics"
	"github.com/howler-chat/api-service/store"
	"github.com/howler-chat/api-service/store/rethink"
	"github.com/pressly/chi"
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

// Records request count, and latency
func RecordMetrics(next chi.Handler) chi.Handler {
	return chi.HandlerFunc(func(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
		// TODO: Record some metrics about the request

		metrics.HTTPRequestCount.WithLabelValues(req.Method, req.URL.Path).Inc()
		startTime := time.Now()
		next.ServeHTTPC(ctx, resp, req)
		stopTime := time.Now()

		elapsed := stopTime.Sub(startTime)
		metrics.HTTPRequestLatency.WithLabelValues(req.Method, req.URL.Path).
			Observe(float64(elapsed) / float64(time.Millisecond))
	})
}

func SetupContext(serviceCtx *ServiceContext) func(chi.Handler) chi.Handler {
	return func(next chi.Handler) chi.Handler {
		return chi.HandlerFunc(func(ctx context.Context, resp http.ResponseWriter, req *http.Request) {

			// TODO: At some point we will have some logic here to decide what rethink session should be
			// associated with this request, (probably based on user or team)
			ctx = rethink.AddRethinkSession(ctx, serviceCtx.RethinkContext.GetRethinkSession())

			// TODO: If in the future we migrate teams to a different store (mongodb?) we would have logic
			// here to decide what Store interface to use for this request, right now we always use rethink
			ctx = store.AddStore(ctx, serviceCtx.Store)

			// Same for API
			ctx = api.AddApi(ctx, serviceCtx.Api)

			next.ServeHTTPC(ctx, resp, req)
		})
	}
}
