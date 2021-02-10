package main

import (
	"flag"
	"os/user"
	"path"

	"github.com/daedaleanai/pgantt/pkg/pgantt"
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

	log.Info("Starting PGantt...")

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

	configFile := path.Join(usr.HomeDir, ".pgantt")
	opts := pgantt.NewOpts()
	err = opts.LoadYaml(configFile)
	if err != nil {
		log.Fatal(err)
	}

	phab, err := pgantt.NewPhabricator(opts.PhabricatorUri, opts.ApiKey)
	if err != nil {
		log.Fatalf("Cannot make a connection to Phabricator: %s", err)
	}

	tasks, err := phab.GetTasksForProject("Platforms")
	if err != nil {
		log.Fatalf("Unable to fetch tasks: %v", err)
	}
	log.Infof("Top-level tasks: %v", len(tasks))

	pgantt.RunWebServer(opts)
}
