// Copyright 2016 Derrick J. Wippler. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors

import (
	"encoding/json"
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/howler-chat/api-service/validate/field"
	"golang.org/x/net/context"
)

const (
	INVALID_JSON = 1 << iota
	INTERNAL_ERROR
	ACCESS_DENIED
	VALIDATION_ERROR
)

type Error struct {
	Type      string `json:"type"`
	Code      int    `json:"code"`
	Message   string `json:"message,omitempty"`
	RequestId string `json:"omit"`
}

func (self *Error) Error() string {
	resp, err := json.Marshal(self)
	if err != nil {
		log.WithField("requestId", self.RequestId).
			Error("json.Marshal() failed on '%+v' with '%s'", self, err.Error())
		return fmt.Sprintf(`{ "type": "error", "code": %d, "message": "Internal Error"}`, INTERNAL_ERROR)
	}
	return string(resp)
}

func Fatal(ctx context.Context, code int, msg string, stuff ...interface{}) Error {
	return Error{
		Type:    "error",
		Code:    code,
		Message: fmt.Sprintf(msg, stuff...),
		// TODO: RequestId: ctx.GetRequestId()
	}
}

func ReceivedInvalidJson(ctx context.Context, err error) error {
	return Fatal(ctx, INVALID_JSON, "Received Invalid JSON - %s", err.Error())
}

func InternalJsonError(ctx context.Context, err error) error {
	return Fatal(ctx, INTERNAL_ERROR, "Marshal JSON Error - %s", err.Error())
}

func ValidationFail(ctx context.Context, msg string, path *field.Path) error {
	return Fatal(ctx, VALIDATION_ERROR, "Validation Failed on '%s' - '%s'", path.String(), msg)
}
