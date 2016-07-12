package api

import "golang.org/x/net/context"

type contextKey int

const (
	apiContextKey contextKey = 0
)

func AddApi(ctx context.Context, api HowlerApi) context.Context {
	return context.WithValue(ctx, apiContextKey, api)
}

func GetApi(ctx context.Context) HowlerApi {
	obj, ok := ctx.Value(apiContextKey).(HowlerApi)
	if !ok {
		panic("No rethink.RethinkContext found in context")
	}
	return obj
}
