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
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type StateManager struct {
	phab     *Phabricator
	m        sync.Mutex
	projects []Project
	tasks    map[string]map[string]*PTask
	users    []User
}

func NewStateManager(opts *Opts) (*StateManager, error) {
	sm := new(StateManager)
	var err error

	sm.phab, err = NewPhabricator(opts.PhabricatorUri, opts.ApiKey)
	if err != nil {
		return nil, fmt.Errorf("Cannot make a connection to Phabricator: %s", err)
	}

	projects := opts.PGantt.Projects
	if len(projects) == 0 {
		projects, err = sm.phab.MyProjectNames()
		if err != nil {
			return nil, fmt.Errorf("Cannot fetch project names: %s", err)
		}
	}

	sm.tasks = make(map[string]map[string]*PTask)
	for _, projName := range projects {
		log.Debugf("Attempting to fetch project info for: %s", projName)
		proj, err := sm.phab.ProjectByName(projName)
		if err != nil {
			return nil, err
		}
		sm.projects = append(sm.projects, *proj)
		sm.tasks[proj.Phid] = make(map[string]*PTask)
	}

	if sm.users, err = sm.phab.Users(); err != nil {
		return nil, err
	}

	log.Infof("Syncing tasks, it may take a while...")
	if err := sm.SyncTasks(); err != nil {
		return nil, err
	}

	go func() {
		for {
			time.Sleep(time.Duration(opts.PGantt.PollInterval) * time.Second)
			if err := sm.SyncTasks(); err != nil {
				log.Errorf("Failed to sync tasks: %s", err)
			}
		}
	}()

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
	plan.Data = make([]Task, 0)
	plan.Links = make([]Link, 0)

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
		for _, link := range t.Links {
			plan.Links = append(plan.Links, *link)
		}
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

	sort.Slice(plan.Links[:], func(i, j int) bool {
		return plan.Links[i].Id < plan.Links[j].Id
	})

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

	var tm time.Time
	var err error
	if task.StartDate != "" {
		tm, err = time.Parse("2006-01-02", task.StartDate)
		if err != nil {
			return "", fmt.Errorf("Malformed start date: %s", err)
		}
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
		if task.StartDate != "" && tm.Unix() != 0 {
			req.SetStartDate(tm.Unix())
			req.SetDuration(task.Duration)
		}
		req.SetProgress(task.Progress)
		req.SetType(task.Type)

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

	if task.Parent != "" && ptask.Task.Parent != task.Parent {
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
		if task.StartDate == "" || tm.Unix() == 0 {
			req.RemoveStartDate()
			task.Duration = 0
		} else {
			req.SetStartDate(tm.Unix())
		}
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

	if ptask.Task.Type != task.Type {
		req.SetType(task.Type)
		numEds++
	}

	if numEds > 0 {
		return s.phab.EditTask(&req)
	}

	return task.Id, nil
}

func (s *StateManager) DeleteLink(projPhid, id string) error {
	s.m.Lock()
	defer s.m.Unlock()

	tasks, ok := s.tasks[projPhid]
	if !ok {
		return fmt.Errorf("No such project: %q", projPhid)
	}

	fragments := strings.Split(id, "#")
	if len(fragments) != 3 {
		return fmt.Errorf("Unable to decode link ID: %s", id)
	}

	ptask, ok := tasks[fragments[0]]
	if !ok {
		return fmt.Errorf("No such source task: %q", fragments[0])
	}

	if _, ok := ptask.Links[id]; !ok {
		return fmt.Errorf("No such link: %q", id)
	}

	delete(ptask.Links, id)

	linkData := getLinkSlice(ptask.Links)
	req := EditRequest{}
	req.SetObjectId(fragments[0])
	req.SetSuccessors(linkData)
	if _, err := s.phab.EditTask(&req); err != nil {
		return err
	}
	return nil
}

func getLinkSlice(links map[string]*Link) []PLinkData {
	linkIds := make([]string, 0, len(links))
	for id := range links {
		linkIds = append(linkIds, id)
	}

	sort.Strings(linkIds)

	lSlice := make([]PLinkData, 0, len(links))
	for _, id := range linkIds {
		lSlice = append(lSlice, PLinkData{
			Target: links[id].Target,
			Type:   links[id].Type,
		})
	}
	return lSlice
}

func generateLinkId(link *Link) string {
	return fmt.Sprintf("%s#%s#%s", link.Source, link.Target, link.Type)
}

func (s *StateManager) CreateLink(projPhid string, link *Link) (string, error) {
	s.m.Lock()
	defer s.m.Unlock()

	tasks, ok := s.tasks[projPhid]
	if !ok {
		return "", fmt.Errorf("No such project: %q", projPhid)
	}

	ptask, ok := tasks[link.Source]
	if !ok {
		return "", fmt.Errorf("No such source task: %q", link.Source)
	}

	_, ok = tasks[link.Source]
	if !ok {
		return "", fmt.Errorf("No such target task: %q", link.Target)
	}

	id := generateLinkId(link)
	link.Id = id

	// This actually happrens because of a bug in the front end. It's fine to assume
	// success because the ID encapsulates the complete link data
	if _, ok = ptask.Links[id]; ok {
		return id, nil
	}

	ptask.Links[id] = link
	linkData := getLinkSlice(ptask.Links)
	req := EditRequest{}
	req.SetObjectId(link.Source)
	req.SetSuccessors(linkData)
	if _, err := s.phab.EditTask(&req); err != nil {
		return "", err
	}
	return id, nil
}
