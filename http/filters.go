package http

import (
	"net/http"
	"runtime/debug"

	log "github.com/Sirupsen/logrus"
	"github.com/pressly/chi"
	"golang.org/x/net/context"
)

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

func MimeJson(next chi.Handler) chi.Handler {
	return chi.HandlerFunc(func(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
		// TODO: Also set 'Accepts' header?
		resp.Header().Set("Content-Type", "application/json; charset=utf-8")
		next.ServeHTTPC(ctx, resp, req)
	})
}

func RecordMetrics(next chi.Handler) chi.Handler {
	return chi.HandlerFunc(func(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
		// TODO: Record some metrics about the request
		next.ServeHTTPC(ctx, resp, req)
	})
}
