// Copyright 2016 Derrick J. Wippler. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package store

import (
	"github.com/howler-chat/api-service/auth"
	"github.com/howler-chat/api-service/model"
	"golang.org/x/net/context"
)

// Save and send the message to the requested channel, will return non nil error if the message is invalid or the client
// does not have access to the requested channel
func SaveMessage(ctx context.Context, msg *model.Message) error {
	// Validate the Model
	if err := msg.Validate(ctx); err != nil {
		return err
	}

	// Does client have access to the channel?
	if err := auth.CanAccessChannel(ctx, msg.ChannelId); err != nil {
		return nil, err
	}

	// TODO: Save the message
	// TODO: get the new message id and set msg.Id
	return nil
}

// Get a message, will return non nil error if the message doesn't exist, or the client does not have accees to the
// requested channel
func GetMessage(ctx context.Context, req *model.GetMessageRequest) (*model.Message, error) {
	// Validate the Model
	if err := req.Validate(ctx); err != nil {
		return err
	}

	// Does client have access to the channel?
	if err := auth.CanAccessChannel(ctx, req.ChannelId); err != nil {
		return nil, err
	}
	// TODO: Get the message
	return model.Message{}, nil
}

func ListMessage(ctx context.Context, req *model.ListMessageRequest) ([]model.Message, error) {
	// Validate the Model
	if err := req.Validate(ctx); err != nil {
		return err
	}

	// Does client have access to the channel?
	if err := auth.CanAccessChannel(ctx, req.ChannelId); err != nil {
		return nil, err
	}
	// TODO: Get all the messages for the requested channel
	return
}
