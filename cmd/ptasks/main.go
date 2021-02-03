package main

import (
	"flag"
	"os/user"
	"path"

	"github.com/daedaleanai/ptasks/pkg/ptasks"
	log "github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

func main() {
	// Commandline
	logLevel := flag.String("log-level", "Info", "verbosity of the diagnostic information")
	flag.Parse()

	// Logging
	log.SetFormatter(&prefixed.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
		ForceFormatting: true,
	})

	log.Info("Starting PTask...")

	level := log.InfoLevel
	if *logLevel != "" {
		l, err := log.ParseLevel(*logLevel)
		if err != nil {
			log.Fatal(err)
		}
		level = l
	}
	log.SetLevel(level)

	// Configuration
	usr, err := user.Current()
	if err != nil {
		log.Fatalf("Cannot get user info for current user: %s", err)
	}

	configFile := path.Join(usr.HomeDir, ".ptasks")
	opts := ptasks.NewOpts()
	err = opts.LoadYaml(configFile)
	if err != nil {
		log.Fatal(err)
	}

	phab, err := ptasks.NewPhabricator(opts.PhabricatorUri, opts.ApiKey)
	if err != nil {
		log.Fatalf("Cannot make a connection to Phabricator: %s", err)
	}

	phab.GetTasksForProject("Platforms")
}
