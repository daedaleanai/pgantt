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
	"encoding/json"
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

type PTaskMetadata struct {
	StartDate   string `json:"start_date"`
	Duration    int    `json:"duration"`
	Unscheduled bool   `json:"unscheduled"`
}

type PTask struct {
	Mtime uint64
	Task  Task
}

type Transaction struct {
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

type EditRequest struct {
	requests.Request
	ObjectIdentifier string        `json:"objectIdentifier,omitempty"`
	Transactions     []Transaction `json:"transactions"`
}

type EditResponse struct {
	Object struct {
		Phid string `json:"phid"`
	} `json:"object"`
}

func (r *EditRequest) SetObjectId(phid string) {
	r.ObjectIdentifier = phid
}

func (r *EditRequest) SetParent(phid string) {
	if r.ObjectIdentifier == "" {
		r.Transactions = append(r.Transactions, Transaction{"parent", phid})
	} else {
		r.Transactions = append(r.Transactions, Transaction{"parents.set", []string{phid}})
	}
}

func (r *EditRequest) RemoveParent() {
	r.Transactions = append(r.Transactions, Transaction{"parents.set", []string{}})
}

func (r *EditRequest) SetProject(phid string) {
	r.Transactions = append(r.Transactions, Transaction{"projects.set", []string{phid}})
}

func (r *EditRequest) SetColumn(phid string) {
	r.Transactions = append(r.Transactions, Transaction{"column", []string{phid}})
}

func (r *EditRequest) SetTitle(title string) {
	r.Transactions = append(r.Transactions, Transaction{"title", title})
}

func (r *EditRequest) SetPTaskMetadata(md *PTaskMetadata) {
	bytes, _ := json.Marshal(md)
	r.Transactions = append(r.Transactions, Transaction{"custom.daedalean.pgantt", string(bytes)})
}

func (p *Phabricator) ProjectByName(name string) (*Project, error) {
	req := requests.ProjectQueryRequest{Names: []string{name}}
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

func (p *Phabricator) SyncTasksForProject(phid string, tasks map[string]*PTask) (map[string]*PTask, error) {
	if tasks == nil {
		tasks = make(map[string]*PTask)
	}

	after := ""
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
			taskPhid := el.PHID
			mtime := uint64(el.Fields["dateModified"].(float64))
			update := false

			ptask, ok := tasks[taskPhid]
			if !ok || ptask.Mtime < mtime {
				update = true
			}

			if update {
				log.Debugf("Updating cached task %q", taskPhid)
				ptask := &PTask{}
				tasks[taskPhid] = ptask
				ptask.Mtime = mtime
				ptask.Task.Id = taskPhid
				ptask.Task.Text = el.Fields["name"].(string)
				ptask.Task.Open = el.Fields["status"].(map[string]interface{})["value"].(string) == "open"
				col := el.Attachments["columns"]["boards"].(map[string]interface{})[phid].(map[string]interface{})["columns"].([]interface{})[0].(map[string]interface{})
				ptask.Task.Column = col["phid"].(string)
				ptask.Task.Url = fmt.Sprintf("%sT%d", p.endpoint, el.ID)

				md := &PTaskMetadata{}
				md.Unscheduled = true
				if _, ok := el.Fields["custom.daedalean.pgantt"]; ok && el.Fields["custom.daedalean.pgantt"] != nil {
					if data, ok := el.Fields["custom.daedalean.pgantt"].(string); ok {
						if err := json.Unmarshal([]byte(data), &md); err != nil {
							log.Errorf("Unable to unmarshal PGantt metadata for task %q (%q): %s", taskPhid, data, err)
						}
					}
				}
				ptask.Task.Unscheduled = md.Unscheduled
				ptask.Task.Duration = md.Duration
				ptask.Task.StartDate = md.StartDate

				// Find out who the parent is
				req := requests.SearchRequest{
					Constraints: map[string]interface{}{
						"projects":   []string{phid},
						"subtaskIDs": []int{el.ID},
					},
				}

				var res responses.SearchResponse
				if err := p.c.Call("maniphest.search", &req, &res); err != nil {
					return nil, err
				}

				if len(res.Data) != 0 {
					ptask.Task.Parent = res.Data[0].PHID
				}
			}
		}

		after = res.Cursor.After
		if after == "" {
			break
		}
	}

	return tasks, nil
}

func (p *Phabricator) EditTask(req *EditRequest) (string, error) {
	if req.ObjectIdentifier == "" {
		for _, tr := range req.Transactions {
			if tr.Type == "title" {
				log.Debugf("Creating new task with title: %q", tr.Value)
				break
			}
		}
	} else {
		log.Debugf("Editing task: %q, transactions: %+v", req.ObjectIdentifier, req.Transactions)
	}
	res := EditResponse{}
	if err := p.c.Call("maniphest.edit", req, &res); err != nil {
		return "", err
	}
	log.Debugf("Task %q edited", res.Object.Phid)
	return res.Object.Phid, nil
}

func NewPhabricator(endpoint, key string) (*Phabricator, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	conn, err := gonduit.Dial(endpoint, &core.ClientOptions{
		APIToken: key,
	})
	if err != nil {
		return nil, err
	}

	return &Phabricator{conn, u.String()}, nil
}
