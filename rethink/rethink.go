package rethink

import (
	"net/http"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/dancannon/gorethink"
	"github.com/howler-chat/api-service/errors"
	"github.com/thrawn01/args"
	"golang.org/x/net/context"
)

// The key type is unexported to prevent collisions with context keys defined in
// other packages.
type contextKey int

// rethinkSessionKey is the context key for the RethinkSession object.  Its value of zero is
// arbitrary.  If this package defined other context keys, they would have
// different integer values.
const rethinkSessionKey contextKey = 0

func NewContext(ctx context.Context, cluster *RethinkCluster) context.Context {
	return context.WithValue(ctx, rethinkSessionKey, cluster)
}

// Get the rethinkSession from our current context
func FromContext(ctx context.Context) *RethinkCluster {
	cluster, ok := ctx.Value(rethinkSessionKey).(*RethinkCluster)
	if !ok {
		panic("No rethink.RethinkCluster found in context")
	}
	return cluster
}

var runOpts = gorethink.RunOpts{
	Durability: "hard",
}

var runOptsHard = gorethink.RunOpts{
	Durability: "hard",
}

var execOpts = gorethink.ExecOpts{
	Durability: "hard",
}

// A singleton that hands out RethinkCluster sessions and handles connect and reconnect to clusters
type Factory struct {
	session         *gorethink.Session
	sessionMutex    sync.Mutex
	parser          *args.ArgParser
	connectingMutex sync.Mutex
	connecting      bool
	stopConnecting  chan bool
}

// TODO: If/When we want to support a different database Factory should become an interface
func NewFactory(parser *args.ArgParser) *Factory {
	factory := &Factory{
		parser:         parser,
		stopConnecting: make(chan bool, 1),
	}
	factory.Connect()
	return factory
}

// Returns true if we are NOT in the process of connecting to rethink
func (self *Factory) notConnecting() bool {
	self.connectingMutex.Lock()
	result := !self.connecting
	self.connectingMutex.Unlock()
	return result
}

// Indicate if we are in the process of connecting to rethink
func (self *Factory) setConnecting(isConnecting bool) {
	self.connectingMutex.Lock()
	self.connecting = isConnecting
	self.connectingMutex.Unlock()
}

// Thread safe session set
func (self *Factory) SetSession(session *gorethink.Session) {
	self.sessionMutex.Lock()
	self.session = session
	self.sessionMutex.Unlock()
}

// Thread safe session get
func (self *Factory) GetSession() *gorethink.Session {
	self.sessionMutex.Lock()
	session := self.session
	self.sessionMutex.Unlock()
	return session
}

// Attempt to reconnect to rethinkdb after disconnecting
func (self *Factory) Reconnect() {
	self.SetSession(nil)
	self.Connect()
}

// Connect to rethink, you may call this function multiple times and it will only connect once
func (self *Factory) Connect() {
	// TODO: This needs to be redone.... perhaps using a complete channel only solution
	// Do not attempt to connect if we already have a session
	if self.GetSession() != nil {
		return
	}

	// If no other goroutine is already trying to connect
	if self.notConnecting() {
		var isRunning sync.WaitGroup
		var once sync.Once
		self.setConnecting(true)
		isRunning.Add(1)
		go func() {
			for {
				// Always fetch the latest version of the config
				config := self.parser.GetOpts()

				// Attempt to connect to rethinkdb
				session, err := gorethink.Connect(gorethink.ConnectOpts{
					Addresses: config.Group("rethink").StringSlice("endpoints"),
					Database:  config.Group("rethink").String("database"),
					Username:  config.Group("rethink").String("user"),
					Password:  config.Group("rethink").String("password"),
				})
				once.Do(func() { isRunning.Done() }) // Notify we attempted to connect at least once
				if err != nil {
					logrus.WithFields(logrus.Fields{
						"type":   "rethink",
						"method": "Connect()",
					}).Errorf("Rethinkdb Connect Failed - %s", err.Error())
					timer := time.NewTimer(time.Second).C
					select {
					case <-timer:
						continue
					case <-self.stopConnecting:
						return
					}
					time.Sleep(time.Second)
					continue
				}
				self.SetSession(session)
				self.setConnecting(false)
				return
			}
		}()
		// Wait until we attempt to connect at least once before continuing
		isRunning.Wait()
	}
}

// Close all rethink sessions and or stop attempting to reconnect
func (self *Factory) Close() {
	// Set Connecting to True, so when we kill the session no one tries to reconnect
	self.setConnecting(true)
	// Kill our current session
	session := self.GetSession()
	if session != nil {
		session.Close()
		self.SetSession(nil)
	}
	// If we are still attempting to connect, stop doing that
	self.stopConnecting <- true
}

// Gets the current connected cluster, Eventually you should be able to ask for different clusters, but currently only
// supports one
func (self *Factory) GetCluster() *RethinkCluster {
	return &RethinkCluster{
		Session: self.GetSession(),
		Factory: self,
	}
}

// This Struct represents a single session connected to a rethink cluster
type RethinkCluster struct {
	Factory *Factory
	Session *gorethink.Session
}

func Error(ctx context.Context, method string, msg string) errors.HowlerError {
	// TODO: If 'err' reports Session is a not connected to rethink, signal to start reconnecting
	// cluster := rethink.FromContext(ctx)
	// cluster.Factory.Reconnect()
	tags := map[string]string{"type": "rethink", "method": method}
	return errors.Internal(ctx, http.StatusServiceUnavailable, tags, "Rethinkdb Error - %s", msg)
}
