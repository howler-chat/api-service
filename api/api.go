// Copyright 2016 Derrick J. Wippler. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"encoding/json"
	"io"

	"github.com/howler-chat/api-service/auth"
	. "github.com/howler-chat/api-service/errors"
	"github.com/howler-chat/api-service/model"
	"github.com/howler-chat/api-service/rethink"
	"golang.org/x/net/context"
)

/*
The api interface provides access to all the public methods for clients to interact with the system. All api
interactions are preformed via json encoded messages. This allows us to transparently use various transport methods to
interact with the system, such as web sockets, http verbs, WebRTC data channels, sockets
*/

type HowlerApi interface {
	PostMessage(ctx context.Context, payload io.Reader) ([]byte, HowlerError)
	GetMessage(ctx context.Context, payload io.Reader) ([]byte, HowlerError)
	MessageList(ctx context.Context, payload io.Reader) ([]byte, HowlerError)
}

type Api struct {
	store rethink.Store
}

func NewApi() HowlerApi {
	return &Api{
		store: rethink.NewStore(),
	}
}

// This method posts a message
// Request
//	{ text: "This is a message", "channelId": "A124B343" }
// Response
//	{ id: "AS223SDFS23" }
func (self *Api) PostMessage(ctx context.Context, payload io.Reader) ([]byte, HowlerError) {
	var msg model.Message

	// TODO: Test how this reacts to multiple json bodies in a single reader
	decoder := json.NewDecoder(payload)
	if err := decoder.Decode(&msg); err != nil {
		err := ReceivedInvalidJson(ctx, err)
		return err.ToJson(), err
	}

	// Validate the Model
	if err := msg.Validate(ctx); err != nil {
		return err.ToJson(), err
	}

	// Does client have access to the channel?
	if err := auth.CanAccessChannel(ctx, msg.ChannelId); err != nil {
		return err.ToJson(), err
	}

	if err := self.store.InsertMessage(ctx, &msg); err != nil {
		return err.ToJson(), err
	}

	resp, err := json.Marshal(map[string]interface{}{"id": msg.Id})
	if err != nil {
		err := InternalJsonError(ctx, "api.PostMessage()", err)
		return err.ToJson(), err
	}
	return resp, nil
}

// This method gets a message
// Request
//	{ "id": "AS223SDFS23", "channelId": "A124B343" }
// Response
//	{ type: "message", text: "This is a message", "channelId": "A124B343" }
func (self *Api) GetMessage(ctx context.Context, payload io.Reader) ([]byte, HowlerError) {
	var request model.GetMessageRequest

	decoder := json.NewDecoder(payload)
	if err := decoder.Decode(&request); err != nil {
		err := ReceivedInvalidJson(ctx, err)
		return err.ToJson(), err
	}

	// Validate the Model
	if err := request.Validate(ctx); err != nil {
		return err.ToJson(), err
	}

	// Does client have access to the channel?
	if err := auth.CanAccessChannel(ctx, request.ChannelId); err != nil {
		return err.ToJson(), err
	}

	msg, err := self.store.GetMessage(ctx, &request)
	if err != nil {
		return err.ToJson(), err
	}

	resp, jsonErr := json.Marshal(msg)
	if jsonErr != nil {
		err := InternalJsonError(ctx, "api.GetMessage()", jsonErr)
		return err.ToJson(), err
	}
	return resp, nil
}

// This method lists all messages for a channel
// Request
//	{ "channelId": "A124B343" }
// Response
//	[
// 		{ type: "message", text: "This is a message", "channelId": "A124B343" }
//		...
//	]
func (self *Api) MessageList(ctx context.Context, payload io.Reader) ([]byte, HowlerError) {
	var request model.ListMessageRequest

	decoder := json.NewDecoder(payload)
	if err := decoder.Decode(&request); err != nil {
		err := ReceivedInvalidJson(ctx, err)
		return err.ToJson(), err
	}

	// Validate the Model
	if err := request.Validate(ctx); err != nil {
		return err.ToJson(), err
	}

	// Does client have access to the channel?
	if err := auth.CanAccessChannel(ctx, request.ChannelId); err != nil {
		return err.ToJson(), err
	}

	msg, err := self.store.ListMessage(ctx, &request)
	if err != nil {
		return err.ToJson(), err
	}

	resp, jsonErr := json.Marshal(msg)
	if jsonErr != nil {
		err := InternalJsonError(ctx, "api.MessageList()", jsonErr)
		return err.ToJson(), err
	}
	return resp, nil
}
