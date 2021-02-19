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
	"net/url"
	"reflect"

	log "github.com/sirupsen/logrus"
	"github.com/thought-machine/gonduit"
	"github.com/thought-machine/gonduit/core"
	"github.com/thought-machine/gonduit/requests"
	"github.com/thought-machine/gonduit/responses"
)

type Phabricator struct {
	c        *gonduit.Conn
	endpoint string
}

type Task struct {
	ID     int
	PHID   string
	Title  string
	Column struct {
		ID   int
		PHID string
		Name string
	}
	Owner struct {
		PHID string
		Name string
	}
	Priority struct {
		Name  string
		Value int
	}
	Status struct {
		Name  string
		Value string
	}
	URL      string
	Children []Task
}

func (p *Phabricator) ProjectByName(name string) (*Project, error) {
	req := requests.ProjectQueryRequest{Names: []string{"Platforms"}}
	res, err := p.c.ProjectQuery(req)
	if err != nil {
		return nil, err
	}

	keys := reflect.ValueOf(res.Data).MapKeys()

	if len(keys) == 0 {
		return nil, fmt.Errorf("Project not found: %s", name)
	}

	phid, ok := keys[0].Interface().(string)
	if !ok {
		return nil, fmt.Errorf("Malformed project query response")
	}

	log.Debugf("Located PHID for %q: %s", name, phid)

	var proj Project
	proj.Name = name
	proj.Phid = phid

	after := ""
	for {
		req := requests.SearchRequest{
			Constraints: map[string]interface{}{
				"projects": []string{phid},
			},
			After: after,
		}
		var res responses.SearchResponse
		if err := p.c.Call("project.column.search", &req, &res); err != nil {
			return nil, err
		}

		for _, el := range res.Data {
			col := Column{}
			col.Name = el.Fields["name"].(string)
			col.Phid = el.PHID
			log.Debugf("Found column in %q: %q (%s)", name, col.Name, col.Phid)
			proj.Columns = append(proj.Columns, col)
		}

		after = res.Cursor.After
		if after == "" {
			break
		}
	}

	return &proj, nil
}

func (p *Phabricator) getProjectPhid(name string) (string, error) {
	req := requests.ProjectQueryRequest{Names: []string{"Platforms"}}
	res, err := p.c.ProjectQuery(req)
	if err != nil {
		return "", err
	}

	keys := reflect.ValueOf(res.Data).MapKeys()

	if len(keys) == 0 {
		return "", fmt.Errorf("Project not found: %s", name)
	}

	return keys[0].Interface().(string), err
}

func (p *Phabricator) GetTasksForProject(name string) ([]Task, error) {
	phid, err := p.getProjectPhid(name)
	if err != nil {
		return nil, err
	}

	after := ""
	type TaskInfo struct {
		top  bool
		task *Task
	}
	taskMap := make(map[string]TaskInfo)

	for {
		req := requests.SearchRequest{
			Constraints: map[string]interface{}{
				"projects": []string{phid},
			},
			Attachments: map[string]bool{
				"columns": true,
			},
			After: after,
		}
		var res responses.SearchResponse
		if err := p.c.Call("maniphest.search", &req, &res); err != nil {
			return nil, err
		}

		for _, el := range res.Data {
			t := Task{}
			t.ID = el.ID
			t.PHID = el.PHID
			t.Title = el.Fields["name"].(string)
			col := el.Attachments["columns"]["boards"].(map[string]interface{})[phid].(map[string]interface{})["columns"].([]interface{})[0].(map[string]interface{})
			t.Column.ID = int(col["id"].(float64))
			t.Column.PHID = col["phid"].(string)
			t.Column.Name = col["name"].(string)
			if ophid, ok := el.Fields["ownerPHID"]; ok && ophid != nil {
				t.Owner.PHID = ophid.(string)
			}
			t.Priority.Name = el.Fields["priority"].(map[string]interface{})["name"].(string)
			t.Priority.Value = int(el.Fields["priority"].(map[string]interface{})["value"].(float64))
			t.Status.Name = el.Fields["status"].(map[string]interface{})["name"].(string)
			t.Status.Value = el.Fields["status"].(map[string]interface{})["value"].(string)
			t.URL = fmt.Sprintf("%sT%d", p.endpoint, t.ID)
			taskMap[t.PHID] = TaskInfo{
				top:  true,
				task: &t,
			}
		}

		after = res.Cursor.After
		if after == "" {
			break
		}
	}

	for _, t := range taskMap {
		req := requests.SearchRequest{
			Constraints: map[string]interface{}{
				"projects":  []string{phid},
				"parentIDs": []int{t.task.ID},
			},
			Attachments: map[string]bool{
				"columns": true,
			},
			After: after,
		}

		var res responses.SearchResponse
		if err := p.c.Call("maniphest.search", &req, &res); err != nil {
			return nil, err
		}

		for _, el := range res.Data {
			if c, ok := taskMap[el.PHID]; ok {
				t.task.Children = append(t.task.Children, *c.task)
				c.top = false
			}
		}
	}

	var tasks []Task
	for _, t := range taskMap {
		if t.top {
			tasks = append(tasks, *t.task)
		}
	}

	return tasks, nil
}

func NewPhabricator(endpoint, key string) (*Phabricator, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}
	conn, err := gonduit.Dial(endpoint, &core.ClientOptions{APIToken: key})
	if err != nil {
		return nil, err
	}

	return &Phabricator{conn, u.String()}, nil
}
