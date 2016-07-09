// Copyright 2016 Derrick J. Wippler. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package service

import (
	"strings"

	. "github.com/howler-chat/api-service/errors"
	. "github.com/howler-chat/api-service/model"
	"golang.org/x/net/context"
)

var contentType = "application/json"

type ServiceClient struct {
	Endpoint string
}

func NewServiceClient(endpoint string) *ServiceClient {
	return &ServiceClient{endpoint}
}

func (self *ServiceClient) buildUrl(slug string) string {
	return strings.Join([]string{self.Endpoint, slug}, "/")
}

func (self *ServiceClient) PostMessage(ctx context.Context, msg *Message) (MessageResponse, error) {
	resp, err := Post(ctx, self.buildUrl("/message.post"), &msg)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return FromErrorResponse(resp.Body)
	}

	entity := MessageResponse{}
	if err := FromJson(resp.Body, &entity); err != nil {
		return nil, err
	}
	return entity
}

func (self *ServiceClient) GetMessage(ctx context.Context, msgId, chanId string) (MessageResponse, HowlerError) {
	request := GetMessageRequest{MessageId: msgId, ChannelId: chanId}
	resp, err := Post(ctx, self.buildUrl("/message.get"), &request)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return FromErrorResponse(resp.Body)
	}

	var entity Message
	if err := FromJson(resp.Body, &entity); err != nil {
		return nil, err
	}
	return entity
}
