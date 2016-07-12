package rethink

import (
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/dancannon/gorethink"
	"github.com/thrawn01/args"
	"golang.org/x/net/context"
)

type contextKey int

const (
	rethinkContextKey contextKey = 0
	rethinkStoreKey   contextKey = 1
	rethinkSessionKey contextKey = 2
)

type RethinkContext struct {
	rethinkChan chan *gorethink.Session
	done        chan struct{}
	parser      *args.ArgParser
}

func NewRethinkContext(parser *args.ArgParser) *RethinkContext {
	return &RethinkContext{
		parser: parser,
	}
}

func (self *RethinkContext) Stop() {
	close(self.done)
}

func (self *RethinkContext) Start() {
	self.rethinkChan = make(chan *gorethink.Session)

	go func() {
		var session *gorethink.Session
		defer close(self.rethinkChan)
		var err error

		for {
			// If we already have a session object, and we are already connected
			if session != nil && session.IsConnected() {
				// Feed a rethink session into the channel, until done channel is closed
				select {
				case self.rethinkChan <- session:
				case <-self.done:
					return
				}
				continue
			}

			// Always fetch the latest version of the config
			config := self.parser.GetOpts()

			// Attempt to connect to rethinkdb
			session, err = gorethink.Connect(gorethink.ConnectOpts{
				Addresses: config.Group("rethink").StringSlice("endpoints"),
				Database:  config.Group("rethink").String("database"),
				Username:  config.Group("rethink").String("user"),
				Password:  config.Group("rethink").String("password"),
			})
			// If something went wrong
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"type":   "rethink",
					"method": "RethinkConnectJob()",
				}).Errorf("Rethinkdb Connect Failed - %s", err.Error())

				// Sleep for 1 second, or until the done channel is closed
				timer := time.NewTimer(time.Second).C
				select {
				case <-timer:
					continue
				case <-self.done:
					return
				}
			}
		}
	}()
}

func (self *RethinkContext) GetRethinkSession() *gorethink.Session {
	return <-self.rethinkChan
}

func AddRethinkContext(ctx context.Context, service *RethinkContext) context.Context {
	return context.WithValue(ctx, rethinkContextKey, service)
}

func AddRethinkSession(ctx context.Context, session *gorethink.Session) context.Context {
	return context.WithValue(ctx, rethinkSessionKey, session)
}

func GetRethinkSession(ctx context.Context) *gorethink.Session {
	obj, ok := ctx.Value(rethinkSessionKey).(*gorethink.Session)
	if !ok {
		panic("No rethink.RethinkContext found in context")
	}
	return obj
}

func GetRethinkContext(ctx context.Context) *RethinkContext {
	obj, ok := ctx.Value(rethinkContextKey).(*RethinkContext)
	if !ok {
		panic("No rethink.RethinkContext found in context")
	}
	return obj
}
