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

	"github.com/ghodss/yaml"
)

type Opts struct {
	Port           int      // Port to serve the on
	PhabricatorUri string   // Rendez-vous point with phabricator
	ApiKey         string   // Phabricator API key
	Projects       []string // List of projects to be handled
}

func NewOpts() (opts *Opts) {
	opts = new(Opts)
	opts.Port = 9999
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

	return nil
}
