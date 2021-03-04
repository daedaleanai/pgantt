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

package pgantt

import (
	"fmt"
	"io/ioutil"
	"reflect"

	"github.com/ghodss/yaml"
)

type HostOpts struct {
	Token string `json:"token"`
}

type PGanttOpts struct {
	Port         int      `json:"port"`          // Port to serve the on
	Projects     []string `json:"projects"`      // List of projects to be handled
	PollInterval int      `json:"poll_interval"` // How often to pool Phabricator for changes in seconds
}

type Opts struct {
	Hosts          map[string]HostOpts `json:"hosts"`
	PGantt         PGanttOpts          `json:"pgantt"`
	PhabricatorUri string
	ApiKey         string
}

func NewOpts() (opts *Opts) {
	opts = new(Opts)
	opts.PGantt.Port = 9999
	opts.PGantt.PollInterval = 10
	opts.PGantt.Projects = []string{}
	return
}

// Load the configuration data from a Yaml file
func (opts *Opts) LoadYaml(fileName string) error {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf("Unable to read the configuration file %s: %s", fileName, err)
	}

	err = yaml.Unmarshal(data, opts)
	if err != nil {
		return fmt.Errorf("Malformed config %s: %s", fileName, err)
	}

	hosts := reflect.ValueOf(opts.Hosts).MapKeys()
	if len(hosts) == 0 {
		return fmt.Errorf("No host definitions found in %s", fileName)
	}

	host := hosts[0].Interface().(string)
	token := opts.Hosts[host].Token
	if token == "" {
		return fmt.Errorf("Token for host %q missing in %s", host, fileName)
	}

	opts.PhabricatorUri = host
	opts.ApiKey = token

	return nil
}
