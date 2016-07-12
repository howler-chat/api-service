// Copyright 2016 Derrick J. Wippler. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package auth

import (
	"golang.org/x/net/context"
	"github.com/howler-chat/api-service/errors"
)

func CanAccessChannel(ctx context.Context, channelId string) errors.HttpError {
	//errors.Fatal(ctx, errors.ACCESS_DENIED, "You do not have access to channel '%s'", channelId)
	return nil
}
