//------------------------------------------------------------------------------
// Copyright (C) 2021 Daedalean AG
//
// This file is part of PGantt.
//
// PGantt is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 2 of the License, or
// (at your option) any later version.
//
// PGantt is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with PGantt.  If not, see <https://www.gnu.org/licenses/>.
//------------------------------------------------------------------------------

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
