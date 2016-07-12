// Copyright 2016 Derrick J. Wippler. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package service

import (
	"fmt"
	"net/http"
	"time"

	"golang.org/x/net/context"

	"github.com/howler-chat/api-service/api"
	"github.com/howler-chat/api-service/errors"
	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/thrawn01/args"
)

func Serve(parser *args.ArgParser) error {
	ctx := NewServiceContext(parser)
	defer ctx.Stop()

	// Start Context Services
	ctx.Start()

	// Listen on our selected interface
	return http.ListenAndServe(parser.GetOpts().String("bind"), NewService(ctx))
}

func NewRouter() chi.Router {
	router := chi.NewRouter()

	// Add NotFound Handler
	router.NotFound(func(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
		err := errors.NewClientError(404, fmt.Sprintf("Path '%s' Not Found", req.URL.RequestURI()), nil)
		resp.WriteHeader(404)
		resp.Write(err.ToJson())
	})

	return router
}

func NewService(ctx *ServiceContext) http.Handler {
	router := NewRouter()

	// Capture any panics
	router.Use(PanicRecoverer)
	// Stop processing if client disconnects
	//router.Use(middleware.CloseNotify)
	// Log Requests
	router.Use(Logger)
	// Stop processing after 2.5 seconds.
	router.Use(middleware.Timeout(2500 * time.Millisecond))
	// Inject the correct rethink session into our current context
	router.Use(SetupContext(ctx))

	router.Route("/api", func(router chi.Router) {
		// Set JSON headers for every request
		router.Use(MimeJson)
		// Record Metrics for every request
		router.Use(RecordMetrics)

		// Use '.' dot to indicate to our users this is not a rest endpoint
		router.Post("/message.post", MessagePost)
		router.Post("/message.get", MessageGet)
		router.Post("/message.list", MessageList)
	})

	// Expose the metrics we have collected
	router.Get("/metrics", prometheus.Handler())

	return router
}

func MessagePost(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
	chatApi := api.GetApi(ctx)
	payload, err := chatApi.PostMessage(ctx, req.Body)
	if err != nil {
		resp.WriteHeader(err.GetCode())
	}
	resp.Write(payload)
	req.Body.Close()
}

func MessageGet(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
	chatApi := api.GetApi(ctx)
	payload, err := chatApi.GetMessage(ctx, req.Body)
	if err != nil {
		resp.WriteHeader(err.GetCode())
	}
	resp.Write(payload)
	req.Body.Close()
}

func MessageList(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
	chatApi := api.GetApi(ctx)
	payload, err := chatApi.MessageList(ctx, req.Body)
	if err != nil {
		resp.WriteHeader(err.GetCode())
	}
	resp.Write(payload)
	req.Body.Close()
}
