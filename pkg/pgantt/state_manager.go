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
	tasks    map[string]map[string]*PTask
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

	sm.tasks = make(map[string]map[string]*PTask)
	for _, proj := range sm.projects {
		log.Infof("Syncing tasks for %q, it may take a while...", proj.Name)
		sm.tasks[proj.Phid] = make(map[string]*PTask)
		sm.tasks[proj.Phid], err = sm.phab.SyncTasksForProject(proj.Phid, sm.tasks[proj.Phid])
		if err != nil {
			return nil, err
		}
	}

	return sm, nil
}

func (s *StateManager) Projects() []Project {
	s.m.Lock()
	defer s.m.Unlock()
	return s.projects
}

func (s *StateManager) PlanningData(phid string) *PlanningData {
	s.m.Lock()
	defer s.m.Unlock()

	if _, ok := s.tasks[phid]; !ok {
		return nil
	}

	plan := &PlanningData{}
	tasks := s.tasks[phid]
	taskMap := make(map[string]bool)
	var add func(t *PTask)
	add = func(t *PTask) {
		if _, added := taskMap[t.Task.Id]; added {
			return
		}

		if t.Task.Parent != "" {
			add(tasks[t.Task.Parent])
		}

		plan.Data = append(plan.Data, t.Task)
	}

	for _, task := range s.tasks[phid] {
		add(task)
	}

	plan.Links = make([]Link, 0)
	return plan
}
