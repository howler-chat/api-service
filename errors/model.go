// Copyright 2016 Derrick J. Wippler. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errors

import (
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
)

type ErrorResponse struct {
	Type      string `json:"type"`
	Code      int    `json:"code"`
	Message   string `json:"message,omitempty"`
	RequestId string `json:"omit"`
}

func (self *ErrorResponse) Error() string {
	return string(self.ToJson())
}

func (self *ErrorResponse) GetCode() int {
	return self.Code
}

func (self *ErrorResponse) GetMessage() string {
	return self.Message
}

func (self *ErrorResponse) ToJson() []byte {
	resp, err := json.Marshal(self)
	if err != nil {
		log.WithField("requestId", self.RequestId).
			Error("json.Marshal() failed on '%+v' with '%s'", self, err.Error())
		return []byte(fmt.Sprintf(`{ "type": "error", "code": %d, "message": "Internal Error"}`,
			http.StatusInternalServerError))
	}
	return resp
}
