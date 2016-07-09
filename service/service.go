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

func Serve(opt *args.Options) error {
	handler := NewService()
	return http.ListenAndServe(opt.String("bind"), handler)
}

func NewRouter() chi.Router {
	router := chi.NewRouter()

	// Add NotFound Handler
	router.NotFound(func(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
		err := errors.NewHowlerError(404, fmt.Sprintf("Path '%s' Not Found", req.URL.RequestURI()), nil)
		resp.WriteHeader(404)
		resp.Write(err.ToJson())
	})

	return router
}

func NewService() http.Handler {
	router := NewRouter()

	// Capture any panics
	router.Use(PanicRecoverer)
	// Stop processing if client disconnects
	//router.Use(middleware.CloseNotify)
	// Log Requests
	router.Use(Logger)
	// Stop processing after 2.5 seconds.
	router.Use(middleware.Timeout(2500 * time.Millisecond))

	router.Route("/api", func(router chi.Router) {
		// Set JSON headers for every request
		router.Use(MimeJson)
		// Record Metrics for every request
		router.Use(RecordMetrics)

		// Use '.' dot to indicate to our users this is not a rest endpoint
		router.Post("/message.post", messagePost)
		router.Post("/message.get", messageGet)
		router.Post("/message.list", messageList)
	})

	// Expose the metrics we have collected
	router.Get("/metrics", prometheus.Handler())

	return router
}

func messagePost(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
	payload, err := api.PostMessage(ctx, req.Body)
	if err != nil {
		resp.WriteHeader(err.GetCode())
	}
	resp.Write(payload)
	req.Body.Close()
}

func messageGet(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
	payload, err := api.GetMessage(ctx, req.Body)
	if err != nil {
		resp.WriteHeader(err.GetCode())
	}
	resp.Write(payload)
	req.Body.Close()
}

func messageList(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
	payload, err := api.MessageList(ctx, req.Body)
	if err != nil {
		resp.WriteHeader(err.GetCode())
	}
	resp.Write(payload)
	req.Body.Close()
}
