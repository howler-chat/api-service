// Copyright 2016 Derrick J. Wippler. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"encoding/json"

	"github.com/howler-chat/api-service/errors"
	"github.com/howler-chat/api-service/model"
	"github.com/howler-chat/api-service/store"
	"golang.org/x/net/context"
)

/*
The api package provides access to all the public methods for clients to interact with the system. All api interactions
are preformed via json encoded messages. This allows us to transparently use various transport methods to interact with
the system, such as web sockets, http verbs, WebRTC data channels, sockets
*/

// This method posts a message
// Request
//	{ text: "This is a message", "channelId": "A124B343" }
// Response
//	{ id: "AS223SDFS23" }
func PostMessage(ctx context.Context, payload []byte) ([]byte, error) {
	var msg model.Message

	if err := json.Unmarshal(payload, &msg); err != nil {
		return nil, errors.ReceivedInvalidJson(ctx, err)
	}

	// Post the message to the table
	if err := store.SaveMessage(ctx, &msg); err != nil {
		return nil, err
	}

	resp, err := json.Marshal(map[string]interface{}{"id": msg.Id})
	if err != nil {
		return nil, errors.InternalJsonError(ctx, err)
	}
	return resp, nil
}

// This method gets a message
// Request
//	{ "id": "AS223SDFS23", "channelId": "A124B343" }
// Response
//	{ type: "message", text: "This is a message", "channelId": "A124B343" }
func GetMessage(ctx context.Context, payload []byte) ([]byte, error) {
	var request model.GetMessageRequest

	if err := json.Unmarshal(payload, &request); err != nil {
		return nil, errors.ReceivedInvalidJson(ctx, err)
	}

	msg, err := store.GetMessage(ctx, &request)
	if err != nil {
		return nil, err
	}

	resp, err := json.Marshal(msg)
	if err != nil {
		return nil, errors.InternalJsonError(ctx, err)
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
func MessageList(ctx context.Context, payload []byte) {
	var request model.ListMessageRequest

	if err := json.Unmarshal(payload, &request); err != nil {
		return nil, errors.ReceivedInvalidJson(ctx, err)
	}

	msg, err := store.ListMessage(ctx, &request)
	if err != nil {
		return nil, err
	}

	resp, err := json.Marshal(msg)
	if err != nil {
		return nil, errors.InternalJsonError(ctx, err)
	}
	return resp, nil
}
