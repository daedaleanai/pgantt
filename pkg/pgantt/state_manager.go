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
	"sort"
	"sync"
	"time"

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

	sm.tasks = make(map[string]map[string]*PTask)
	for _, projName := range opts.Projects {
		log.Debugf("Attempting to fetch project info for: %s", projName)
		proj, err := sm.phab.ProjectByName(projName)
		if err != nil {
			return nil, err
		}
		sm.projects = append(sm.projects, *proj)
		sm.tasks[proj.Phid] = make(map[string]*PTask)
	}

	log.Infof("Syncing tasks, it may take a while...")
	if err := sm.SyncTasks(); err != nil {
		return nil, err
	}
	return sm, nil
}

func (s *StateManager) SyncTasks() error {
	s.m.Lock()
	defer s.m.Unlock()

	var err error
	for _, proj := range s.projects {
		s.tasks[proj.Phid], err = s.phab.SyncTasksForProject(proj.Phid, s.tasks[proj.Phid])
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *StateManager) Projects() []Project {
	s.m.Lock()
	defer s.m.Unlock()
	return s.projects
}

func (s *StateManager) PlanningData(phid string) *PlanningData {
	s.m.Lock()
	defer s.m.Unlock()

	tasks, ok := s.tasks[phid]
	if !ok {
		return nil
	}

	plan := &PlanningData{}
	taskMap := make(map[string]bool)
	var add func(t *PTask)
	add = func(t *PTask) {
		if _, added := taskMap[t.Task.Id]; added {
			return
		}

		if t.Task.Parent != "" {
			add(tasks[t.Task.Parent])
		}

		taskMap[t.Task.Id] = true
		plan.Data = append(plan.Data, t.Task)
	}

	taskPhids := make([]string, 0, len(tasks))
	for phid := range tasks {
		taskPhids = append(taskPhids, phid)
	}

	sort.Strings(taskPhids)

	for _, phid := range taskPhids {
		task := tasks[phid]
		add(task)
	}

	plan.Links = make([]Link, 0)
	return plan
}

func (s *StateManager) EditTask(projPhid string, task *Task) (string, error) {
	s.m.Lock()
	defer s.m.Unlock()

	tasks, ok := s.tasks[projPhid]
	if !ok {
		return "", fmt.Errorf("No such project: %q", projPhid)
	}

	ptask, ok := tasks[task.Id]
	tm, err := time.Parse("2006-01-02", task.StartDate)
	if err != nil {
		return "", fmt.Errorf("Malformed start date: %s", err)
	}

	// New task
	if !ok {
		req := EditRequest{}
		req.SetProject(projPhid)
		if task.Parent != "0" {
			req.SetParent(task.Parent)
		}
		req.SetColumn(task.Column)
		req.SetTitle(task.Text)

		req.SetScheduled(!task.Unscheduled)
		req.SetStartDate(tm.Unix())
		req.SetDuration(task.Duration)
		req.SetProgress(task.Progress)

		return s.phab.EditTask(&req)
	}

	// Edit task
	numEds := 0
	req := EditRequest{}
	req.SetObjectId(task.Id)
	if ptask.Task.Column != task.Column {
		req.SetColumn(task.Column)
		numEds++
	}

	if ptask.Task.Text != task.Text {
		req.SetTitle(task.Text)
		numEds++
	}

	if ptask.Task.Parent != "" && ptask.Task.Parent != task.Parent {
		if task.Parent == "0" {
			req.RemoveParent()
		} else {
			req.SetParent(task.Parent)
		}
		numEds++
	}

	if ptask.Task.Unscheduled != task.Unscheduled {
		req.SetScheduled(!task.Unscheduled)
		numEds++
	}

	if ptask.Task.StartDate != task.StartDate {
		req.SetStartDate(tm.Unix())
		numEds++
	}

	if ptask.Task.Duration != task.Duration {
		req.SetDuration(task.Duration)
		numEds++
	}

	if ptask.Task.Progress != task.Progress {
		req.SetProgress(task.Progress)
		numEds++
	}

	if numEds > 0 {
		return s.phab.EditTask(&req)
	}

	return task.Id, nil
}
