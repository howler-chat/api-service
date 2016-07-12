package rethink

import (
	"net/http"

	"github.com/dancannon/gorethink"
	"github.com/howler-chat/api-service/errors"
	"golang.org/x/net/context"
)

var runOpts = gorethink.RunOpts{
	Durability: "hard",
}

var runOptsHard = gorethink.RunOpts{
	Durability: "hard",
}

var execOpts = gorethink.ExecOpts{
	Durability: "hard",
}

func Error(ctx context.Context, method string, msg string) errors.HttpError {
	tags := map[string]string{"type": "rethink", "method": method}
	return errors.NewHttpError(ctx, http.StatusServiceUnavailable, tags, "Rethinkdb Error - %s", msg)
}
