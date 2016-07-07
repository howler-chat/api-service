package main

import (
	"fmt"
	"io/ioutil"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/howler-chat/api-service/http"
	"github.com/thrawn01/args"
)

func main() {
	parser := args.NewParser(args.Name("api-service"), args.EnvPrefix("API_"))

	parser.AddOption("--config").Alias("-c").Env("CONFIG_FILE").
		Help("Specify the location of the config file")
	parser.AddOption("--bind").Alias("-b").IsTrue().Env("BIND").Default("0.0.0.0:8080").
		Help("The interface to bind too")
	parser.AddOption("--debug").Alias("-d").IsTrue().Env("DEBUG").
		Help("Output debug messages")

	rethink := parser.InGroup("rethink")

	rethink.AddOption("--endpoints").Alias("-e").Env("RETHINK_ENDPOINTS").
		Help("comma separated list of rethinkdb cluster endpoints")
	rethink.AddOption("--user").Alias("-u").Env("RETHINK_USER").
		Help("RethinkDB Username")
	rethink.AddOption("--password").Alias("-p").Env("RETHINK_PASSWORD").
		Help("RethinkDB Password")
	rethink.AddOption("--db").Alias("-d").Env("RETHINK_DATABASE").
		Help("RethinkDB Database name")

	opt := parser.ParseArgsSimple(nil)

	// If a config file is provided
	if opt.String("config") != "" {
		reader, err := ioutil.ReadFile(opt.String("config"))
		if err != nil {
			fmt.Printf("Error reading config file - %s\n", err.Error())
			os.Exit(-1)
		}
		// Read the rest of our options from the config file
		if opt, err = parser.ParseIni(reader); err != nil {
			fmt.Printf("Error parsing '%s'  - %s", opt.String("config"), err.Error())
			os.Exit(-1)
		}
	}

	// TODO: Setup a watch for ConfigMap if we decide to use it for credentials

	if opt.Bool("debug") {
		log.Info("Debug Enabled")
		log.SetLevel(log.DebugLevel)
	}

	err := http.Serve(opt)
	if err != nil {
		log.Fatal(err)
	}
}
