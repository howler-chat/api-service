// Copyright 2016 Derrick J. Wippler. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rethink

import (
	"fmt"

	"github.com/dancannon/gorethink"
	"github.com/howler-chat/api-service/errors"
	"github.com/howler-chat/api-service/model"
	"golang.org/x/net/context"
)

type Store interface {
	InsertMessage(ctx context.Context, msg *model.Message) errors.HowlerError
	GetMessage(ctx context.Context, req *model.GetMessageRequest) (*model.Message, errors.HowlerError)
	ListMessage(ctx context.Context, req *model.ListMessageRequest) ([]model.Message, errors.HowlerError)
}

type RethinkStore struct{}

func NewStore() Store {
	return &RethinkStore{}
}

// Insert the message on the requested channel
func (self *RethinkStore) InsertMessage(ctx context.Context, msg *model.Message) errors.HowlerError {
	cluster := FromContext(ctx)

	changed, err := gorethink.Table("Message").Insert(msg).RunWrite(cluster.Session, runOpts)
	if err != nil {
		return Error(ctx, "InsertMessage()", err.Error())
	} else if changed.Errors != 0 {
		return Error(ctx, "InsertMessage()", changed.FirstError)
	}
	if len(changed.GeneratedKeys) == 0 {
		return Error(ctx, "InsertMessage()",
			fmt.Sprintf("GeneratedKeys Empty after insert %+v", changed))
	}
	msg.Id = changed.GeneratedKeys[0]
	return nil
}

// Get a message, will return non nil error if the message doesn't exist
func (self *RethinkStore) GetMessage(ctx context.Context, req *model.GetMessageRequest) (*model.Message, errors.HowlerError) {
	cluster := FromContext(ctx)

	var message model.Message
	cursor, err := gorethink.Table("Message").
		Filter(gorethink.Row.Field("ChannelId").Eq(req.ChannelId).
			And(gorethink.Row.Field("MessageId").Eq(req.MessageId))).Run(cluster.Session, runOpts)

	if err != nil {
		return nil, Error(ctx, "GetMessage()", err.Error())
	} else if err := cursor.One(&message); err != nil {
		return nil, Error(ctx, "GetMessage().One()", err.Error())
	}
	return &message, nil
}

func (self *RethinkStore) ListMessage(ctx context.Context, req *model.ListMessageRequest) ([]model.Message, errors.HowlerError) {
	cluster := FromContext(ctx)

	var messages []model.Message
	cursor, err := gorethink.Table("Message").
		Filter(gorethink.Row.Field("ChannelId").Eq(req.ChannelId)).Run(cluster.Session, runOpts)

	if err != nil {
		return nil, Error(ctx, "ListMessage()", err.Error())
	} else if err := cursor.All(&messages); err != nil {
		return nil, Error(ctx, "ListMessage().All()", err.Error())
	}
	return messages, nil
}
