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
	"sync"

	log "github.com/sirupsen/logrus"
)

type StateManager struct {
	phab     *Phabricator
	m        sync.Mutex
	projects []Project
}

func NewStateManager(opts *Opts) (*StateManager, error) {
	sm := new(StateManager)
	var err error

	sm.phab, err = NewPhabricator(opts.PhabricatorUri, opts.ApiKey)
	if err != nil {
		return nil, fmt.Errorf("Cannot make a connection to Phabricator: %s", err)
	}

	log.Infof("Created a connection to Phabricator at: %s", opts.PhabricatorUri)

	for _, projName := range opts.Projects {
		log.Debugf("Attempting to fetch project info for: %s", projName)
		proj, err := sm.phab.ProjectByName(projName)
		if err != nil {
			return nil, err
		}
		sm.projects = append(sm.projects, *proj)
	}

	return sm, nil
}

func (s *StateManager) Projects() []Project {
	s.m.Lock()
	defer s.m.Unlock()
	return s.projects
}
