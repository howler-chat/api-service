package service

import (
	"github.com/howler-chat/api-service/api"
	"github.com/howler-chat/api-service/store"
	"github.com/howler-chat/api-service/store/rethink"
	"github.com/thrawn01/args"
)

// This handles all the context for the service, including hot reloading of objects and config changes
type ServiceContext struct {
	RethinkContext *rethink.RethinkContext
	Api            api.HowlerApi
	Store          store.HowlerStore
}

// This should create a new context based on the config passed in via the parser
func NewServiceContext(parser *args.ArgParser) *ServiceContext {
	return &ServiceContext{
		RethinkContext: rethink.NewRethinkContext(parser),
		Api:            api.NewApi(),
		Store:          rethink.NewStore(),
	}
}

func (self *ServiceContext) Start() {
	self.RethinkContext.Start()
}

func (self *ServiceContext) Stop() {
	self.RethinkContext.Stop()
}
