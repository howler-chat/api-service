// Copyright 2016 Derrick J. Wippler. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package store

import (
	"fmt"

	"github.com/dancannon/gorethink"
	"github.com/howler-chat/api-service/errors"
	"github.com/howler-chat/api-service/model"
	"golang.org/x/net/context"
)

// Insert the message on the requested channel
func InsertMessage(ctx context.Context, msg *model.Message) errors.HowlerError {
	changed, err := gorethink.Table("Message").Insert(msg).RunWrite(Session, runOpts)
	if err != nil {
		return RethinkError(ctx, "InsertMessage()", err.Error())
	} else if changed.Errors != 0 {
		return RethinkError(ctx, "InsertMessage()", changed.FirstError)
	}
	if len(changed.GeneratedKeys) == 0 {
		return RethinkError(ctx, "InsertMessage()",
			fmt.Sprintf("GeneratedKeys Empty after insert %+v", changed))
	}
	msg.Id = changed.GeneratedKeys[0]
	return nil
}

// Get a message, will return non nil error if the message doesn't exist, or the client does not have accees to the
// requested channel
func GetMessage(ctx context.Context, req *model.GetMessageRequest) (*model.Message, errors.HowlerError) {
	var message model.Message
	cursor, err := gorethink.Table("Message").
		Filter(gorethink.Row.Field("ChannelId").Eq(req.ChannelId).
			And(gorethink.Row.Field("MessageId").Eq(req.MessageId))).Run(Session, runOpts)

	if err != nil {
		return RethinkError(ctx, "GetMessage()", err.Error())
	} else if err := cursor.One(&message); err != nil {
		return RethinkError(ctx, "GetMessage().One()", err.Error())
	}
	return &message, nil
}

func ListMessage(ctx context.Context, req *model.ListMessageRequest) ([]model.Message, errors.HowlerError) {
	var messages []model.Message
	cursor, err := gorethink.Table("Message").
		Filter(gorethink.Row.Field("ChannelId").Eq(req.ChannelId)).Run(Session, runOpts)

	if err != nil {
		return RethinkError(ctx, "ListMessage()", err.Error())
	} else if err := cursor.All(&messages); err != nil {
		return RethinkError(ctx, "ListMessage().All()", err.Error())
	}
	return &messages, nil
}
