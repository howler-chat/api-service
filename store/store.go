package store

import (
	"github.com/howler-chat/api-service/errors"
	"github.com/howler-chat/api-service/model"
	"golang.org/x/net/context"
)

type contextKey int

const (
	storeKey contextKey = 0
)

type HowlerStore interface {
	InsertMessage(ctx context.Context, msg *model.Message) errors.HttpError
	GetMessage(ctx context.Context, req *model.GetMessageRequest) (*model.Message, errors.HttpError)
	ListMessage(ctx context.Context, req *model.ListMessageRequest) ([]model.Message, errors.HttpError)
}

func AddStore(ctx context.Context, store HowlerStore) context.Context {
	return context.WithValue(ctx, storeKey, store)
}

func GetStore(ctx context.Context) HowlerStore {
	obj, ok := ctx.Value(storeKey).(HowlerStore)
	if !ok {
		panic("No rethink.RethinkContext found in context")
	}
	return obj
}
