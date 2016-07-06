// Copyright 2016 Derrick J. Wippler. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors

import (
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/howler-chat/api-service/validate/field"
	"golang.org/x/net/context"
)

type HowlerError interface {
	error
	Code() int
	ToJson() []byte
}

type howlerError struct {
	Type      string `json:"type"`
	Code      int    `json:"code"`
	Message   string `json:"message,omitempty"`
	RequestId string `json:"omit"`
}

func (self *howlerError) Error() string {
	return string(self.ToJson())
}

func (self *howlerError) ToJson() []byte {
	resp, err := json.Marshal(self)
	if err != nil {
		log.WithField("requestId", self.RequestId).
			Error("json.Marshal() failed on '%+v' with '%s'", self, err.Error())
		return []byte(fmt.Sprintf(`{ "type": "error", "code": %d, "message": "Internal Error"}`,
			http.StatusInternalServerError))
	}
	return resp
}

func Fatal(ctx context.Context, code int, msg string, stuff ...interface{}) HowlerError {
	return howlerError{
		Type:    "error",
		Code:    code,
		Message: fmt.Sprintf(msg, stuff...),
		// TODO: RequestId: ctx.GetRequestId()
	}
}

func ReceivedInvalidJson(ctx context.Context, err error) HowlerError {
	return Fatal(ctx, http.StatusBadRequest, "Received Invalid JSON - %s", err.Error())
}

func InternalJsonError(ctx context.Context, err error) HowlerError {
	return Fatal(ctx, http.StatusInternalServerError, "Marshal JSON Error - %s", err.Error())
}

func ValidationFail(ctx context.Context, msg string, path *field.Path) HowlerError {
	return Fatal(ctx, http.StatusNotAcceptable, "Validation Failed on '%s' - '%s'", path.String(), msg)
}
