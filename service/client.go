// Copyright 2016 Derrick J. Wippler. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package service

import (
	"path"

	"net/url"

	. "github.com/howler-chat/api-service/model"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

var contentType = "application/json"

type ServiceClient struct {
	Endpoint string
}

func NewServiceClient(endpoint string) (*ServiceClient, error) {
	_, err := url.Parse(endpoint)
	if err != nil {
		return nil, errors.Wrap(err, "Invalid Endpoint")
	}
	return &ServiceClient{endpoint}, nil
}

func (self *ServiceClient) buildUrl(slug string) string {
	parts, _ := url.Parse(self.Endpoint)
	parts.Path = path.Join(parts.Path, slug)
	return parts.String()
}

func (self *ServiceClient) PostMessage(ctx context.Context, msg *Message) (*MessageResponse, error) {
	resp, err := Post(ctx, self.buildUrl("message.post"), &msg)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, FromErrorResponse(resp.Body)
	}

	entity := MessageResponse{}
	if err := FromJson(resp.Body, &entity); err != nil {
		return nil, err
	}
	return &entity, nil
}

func (self *ServiceClient) GetMessage(ctx context.Context, msgId, chanId string) (*Message, error) {
	request := GetMessageRequest{MessageId: msgId, ChannelId: chanId}
	//fmt.Printf("Endpoint: %s\n", self.Endpoint)
	//fmt.Printf("Url: %s\n", self.buildUrl("/message.get"))
	resp, err := Post(ctx, self.buildUrl("/api/message.get"), &request)
	//fmt.Printf("Resp: %+v\n", resp)
	//fmt.Printf("Err: %+v\n", err)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, FromErrorResponse(resp.Body)
	}

	var entity Message
	if err := FromJson(resp.Body, &entity); err != nil {
		return nil, err
	}
	return &entity, nil
}
