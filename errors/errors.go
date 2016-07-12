// Copyright 2016 Derrick J. Wippler. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors

import (
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/howler-chat/api-service/metrics"
	"github.com/howler-chat/api-service/utils"
	"golang.org/x/net/context"
)

// Used internally by the http service
type HttpError interface {
	error
	GetCode() int
	GetMessage() string
	ToJson() []byte
}

// Returned by the http client
type ClientError interface {
	HttpError
	GetRaw() []byte
}

// Used by a client to indicate some error occurred while communicating with the server
func NewClientError(code int, msg string, body []byte) HttpError {
	return &ErrorResponse{
		Type:    "error",
		Code:    code,
		Message: msg,
		Raw:     body,
	}
}

// Used by a Service to create a new HttpError, suitable for JSON encoding as a response to a request. If tags is != nil
// the error will be logged, and included in the our metrics
func NewHttpError(ctx context.Context, code int, tags map[string]string, msg string, stuff ...interface{}) HttpError {
	renderedMsg := fmt.Sprintf(msg, stuff...)

	// If we included tags, assume we want metrics on this error
	if tags != nil {
		// Tell metrics about the internal error
		metrics.InternalErrors.With(tags).Inc()
		// Log the detail of the error
		log.WithFields(utils.ToFields(tags)).Error(renderedMsg)
	}

	return &ErrorResponse{
		Type:    "error",
		Code:    code,
		Message: renderedMsg,
		// TODO: RequestId: ctx.GetRequestId()
	}
}

// Tell the client, it sent invalid json
func HttpErrorInvalidJson(ctx context.Context, err error) HttpError {
	return NewHttpError(ctx, http.StatusBadRequest, nil, "Received Invalid JSON - %s", err.Error())
}

// Tell the client we had some issue un-marshalling json internally
func HttpErrorInternalJson(ctx context.Context, method string, err error) HttpError {
	tags := map[string]string{
		"type:":  "json",
		"method": method,
	}
	return NewHttpError(ctx, http.StatusInternalServerError, tags, "Marshal JSON Error - %s", err.Error())
}
