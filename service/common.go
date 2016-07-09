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
	"golang.org/x/net/context"
	"golang.org/x/net/context/ctxhttp"
)

func FromErrorResponse(body io.ReadCloser) HowlerError {
	var entity ErrorResponse
	if err := FromJson(body, &entity); err != nil {
		return NewHowlerError(fmt.Sprintf("Error Marshalling ErrorResponse from server - %s", err.Error()))
	}
	return &entity
}

func FromJson(body io.ReadCloser, value interface{}) error {
	body, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, &value); err != nil {
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

func Post(ctx context.Context, url, value interface{}) (*http.Response, error) {
	payload, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	return ctxhttp.Post(ctx, nil, url, contentType, bytes.NewReader(payload))
}
