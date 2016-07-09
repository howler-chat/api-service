// Copyright 2016 Derrick J. Wippler. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors

import (
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/howler-chat/api-service/metrics"
	"golang.org/x/net/context"
)

type HowlerError interface {
	error
	GetCode() int
	GetMessage() string
	ToJson() []byte
}

func NewHowlerError(msg string) *ErrorResponse {
	return &ErrorResponse{
		Type:    "error",
		Code:    0,
		Message: msg,
	}
}

func Error(ctx context.Context, code int, msg string, stuff ...interface{}) HowlerError {
	return &ErrorResponse{
		Type:    "error",
		Code:    code,
		Message: fmt.Sprintf(msg, stuff...),
		// TODO: RequestId: ctx.GetRequestId()
	}
}

func Internal(ctx context.Context, code int, tags map[string]string, msg string, stuff ...interface{}) HowlerError {
	renderedMsg := fmt.Sprintf(msg, stuff...)
	// Tell metrics about the internal error
	metrics.InternalErrors.With(tags).Inc()
	// Log the detail of the error
	log.WithFields(tags).Error(renderedMsg)

	return &ErrorResponse{
		Type:    "error",
		Code:    code,
		Message: renderedMsg,
		// TODO: RequestId: ctx.GetRequestId()
	}
}

func ReceivedInvalidJson(ctx context.Context, err error) HowlerError {
	return Error(ctx, http.StatusBadRequest, "Received Invalid JSON - %s", err.Error())
}

func InternalJsonError(ctx context.Context, method string, err error) HowlerError {
	tags := map[string]string{
		"type:":  "json",
		"method": method,
	}
	return Internal(ctx, http.StatusInternalServerError, tags, "Marshal JSON Error - %s", err.Error())
}
