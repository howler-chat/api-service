package store

import (
	"time"

	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/dancannon/gorethink"
	"github.com/golang/go/src/pkg/sync"
	"github.com/howler-chat/api-service/errors"
	"github.com/thrawn01/args"
	"golang.org/x/net/context"
)

// The rethink Session object
var Session *gorethink.Session

var runOpts = gorethink.RunOpts{
	Durability: "hard",
}

var runOptsHard = gorethink.RunOpts{
	Durability: "hard",
}

var execOpts = gorethink.ExecOpts{
	Durability: "hard",
}

// Used for deciding if we are currently reconnecting
var mutex sync.Mutex
var connecting bool

func notConnecting() bool {
	mutex.Lock()
	result := !connecting
	mutex.Unlock()
	return result
}

func setConnecting(isConnecting bool) {
	mutex.Lock()
	connecting = isConnecting
	mutex.Unlock()
}

var parser *args.ArgParser

// Attempt to reconnect to rethinkdb after disconnecting
func Reconnect() {
	Session = nil
	Init(parser)
}

// Connect to rethink, you may call this function multiple times and it will only connect once
func Init(newParser *args.ArgParser) {
	// Do not attempt to connect if we already have a session
	if Session != nil {
		return
	}

	// If not other goroutine is already trying to connect
	if notConnecting() {
		parser = newParser
		setConnecting(true)
		go func() {
			for {
				// Always fetch the latest version of the config
				config := parser.GetOpts()

				// Attempt to connect to rethinkdb
				session, err := gorethink.Connect(gorethink.ConnectOpts{
					Addresses: config.Group("rethink").StringSlice("endpoints"),
					Database:  config.Group("rethink").String("database"),
					Username:  config.Group("rethink").String("user"),
					Password:  config.Group("rethink").String("password"),
				})
				if err != nil {
					log.Errorf("Rethinkdb Connect Failed - %s", err.Error())
					time.Sleep(time.Second)
					continue
				}
				Session = session
				setConnecting(false)
				return
			}
		}()
	}
}

func RethinkError(ctx context.Context, functionName string, msg string) errors.HowlerError {
	// TODO: If 'err' reports Session is a not connected to rethink, signal to start reconnecting
	// TODO: Call Reconnect()
	// TODO: If we need to reconnect, add the 'reconnect' tag
	tags := []string{"rethink"}
	return errors.Internal(ctx, http.StatusServiceUnavailable, tags, "%s Rethinkdb Error - %s", functionName, msg)
}
