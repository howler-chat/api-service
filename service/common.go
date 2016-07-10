// Copyright 2016 Derrick J. Wippler. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	. "github.com/howler-chat/api-service/errors"
	"github.com/thrawn01/args"
	"golang.org/x/net/context"
	"golang.org/x/net/context/ctxhttp"
)

func FromErrorResponse(body io.ReadCloser) HowlerError {
	payload, err := ioutil.ReadAll(body)
	fmt.Printf("resp: %s\n", string(payload))
	if err != nil {
		return NewHowlerError(0, err.Error(), payload)
	}

	var entity ErrorResponse
	if err := json.Unmarshal(payload, &entity); err != nil {
		return NewHowlerError(0, fmt.Sprintf("Invalid JSON from server - %s", err.Error()), payload)
	}
	entity.Raw = payload
	return &entity
}

func FromJson(body io.ReadCloser, value interface{}) error {
	payload, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(payload, &value); err != nil {
		return err
	}
	return nil
}

func CurlString(req *http.Request, payload *[]byte) string {
	parts := []string{"curl", "-i", "-X", req.Method, req.URL.String()}
	for key, value := range req.Header {
		parts = append(parts, fmt.Sprintf("-H \"%s: %s\"", key, value[0]))
	}

	if payload != nil {
		parts = append(parts, fmt.Sprintf(" -d '%s'", string(*payload)))
	}

	return strings.Join(parts, " ")
}

func Post(ctx context.Context, url string, value interface{}) (*http.Response, error) {
	payload, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	return ctxhttp.Post(ctx, nil, url, contentType, bytes.NewReader(payload))
}

// Return the message associated with this error
func GetErrorMsg(err error) string {
	obj, ok := err.(HowlerError)
	if ok {
		return obj.GetMessage()
	}
	return err.Error()
}

// Return the error code associated with this error, if Error Code is 0, no JSON is associated with this error
func GetErrorCode(err error) int {
	obj, ok := err.(HowlerError)
	if ok {
		return obj.GetCode()
	}
	return 0
}

// Return the RAW un-parsed JSON, returns nil if no JSON is associated with this error
func GetErrorRaw(err error) []byte {
	obj, ok := err.(HowlerError)
	if ok {
		return obj.GetRaw()
	}
	return nil
}

func ParseRethinkArgs(argv *[]string) *args.ArgParser {
	parser := args.NewParser()
	rethink := parser.InGroup("rethink")
	rethink.AddOption("--endpoints").Env("RETHINK_ENDPOINTS")
	rethink.AddOption("--user").Env("RETHINK_USER")
	rethink.AddOption("--password").Env("RETHINK_PASSWORD")
	rethink.AddOption("--db").Env("RETHINK_DATABASE")
	parser.ParseArgs(argv)
	return parser
}
