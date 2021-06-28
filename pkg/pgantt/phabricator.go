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
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/thought-machine/gonduit"
	"github.com/thought-machine/gonduit/core"
	"github.com/thought-machine/gonduit/entities"
	"github.com/thought-machine/gonduit/requests"
	"github.com/thought-machine/gonduit/responses"
)

type Phabricator struct {
	c              *gonduit.Conn
	endpoint       string
	fieldsVerified bool
}

type PTask struct {
	Mtime  uint64
	IsLeaf bool
	Links  map[string]*Link
	Task   Task
}

type PLinkData struct {
	Target string `json:"target"`
	Type   string `json:"type"`
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

func (r *EditRequest) SetStartDate(date int64) {
	r.Transactions = append(r.Transactions, Transaction{"custom.daedalean.start_date", float64(date)})
}

func (r *EditRequest) RemoveStartDate() {
	r.Transactions = append(r.Transactions, Transaction{"custom.daedalean.start_date", nil})
}

func (r *EditRequest) SetScheduled(scheduled bool) {
	r.Transactions = append(r.Transactions, Transaction{"custom.daedalean.scheduled", scheduled})
}

func (r *EditRequest) SetDuration(duration int) {
	r.Transactions = append(r.Transactions, Transaction{"custom.daedalean.duration", float64(duration)})
}

func (r *EditRequest) SetProgress(progress float32) {
	r.Transactions = append(r.Transactions, Transaction{"custom.daedalean.progress", float64(int(progress * 100))})
}

func (r *EditRequest) SetSuccessors(links []PLinkData) {
	data, _ := json.Marshal(links)
	r.Transactions = append(r.Transactions, Transaction{"custom.daedalean.successors", string(data)})
}

func (r *EditRequest) SetType(typ string) {
	phTyp := "daedalean:task"
	if typ == "milestone" {
		phTyp = "daedalean:milestone"
	} else if typ == "project" {
		phTyp = "daedalean:project"
	}
	r.Transactions = append(r.Transactions, Transaction{"custom.daedalean.type", phTyp})
}

func (p *Phabricator) MyProjectNames() ([]string, error) {
	userReq := requests.Request{}
	var userRes entities.User
	if err := p.c.Call("user.whoami", &userReq, &userRes); err != nil {
		return nil, err
	}

	log.Debugf("User: %s (%s)", userRes.UserName, userRes.RealName)

	projects := []string{}
	after := ""
	for {
		req := requests.SearchRequest{
			Attachments: map[string]bool{
				"members": true,
			},
			After: after,
		}
		var res responses.SearchResponse
		if err := p.c.Call("project.search", &req, &res); err != nil {
			return nil, err
		}

		for _, el := range res.Data {
			if el.Fields["icon"].(map[string]interface{})["key"].(string) != "project" {
				continue
			}

			if membersIface, ok := el.Attachments["members"]["members"]; ok && membersIface != nil {
				for _, mem := range membersIface.([]interface{}) {
					if mem.(map[string]interface{})["phid"].(string) == userRes.PHID {
						name := el.Fields["name"].(string)
						log.Debugf("Found project: %s", name)
						projects = append(projects, name)
					}
				}
			}
		}

		after = res.Cursor.After
		if after == "" {
			break
		}
	}

	return projects, nil
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

func (p *Phabricator) verifyCustomFields(fields map[string]interface{}) {
	if p.fieldsVerified {
		return
	}

	fieldNames := []string{
		"custom.daedalean.scheduled",
		"custom.daedalean.start_date",
		"custom.daedalean.duration",
		"custom.daedalean.progress",
		"custom.daedalean.type",
		"custom.daedalean.successors",
	}
	for _, name := range fieldNames {
		if _, ok := fields[name]; !ok {
			log.Fatalf("Task field %q missing. Please go to " +
				"https://github.com/daedaleanai/pgantt for instructions on how" +
				"to configure Phabricator")
		}
	}
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
			p.verifyCustomFields(el.Fields)
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
				ptask.IsLeaf = true
				ptask.Links = make(map[string]*Link)
				ptask.Mtime = mtime
				ptask.Task.Id = taskPhid
				ptask.Task.Text = el.Fields["name"].(string)
				ptask.Task.Open = el.Fields["status"].(map[string]interface{})["value"].(string) == "open"
				board := el.Attachments["columns"]["boards"].(map[string]interface{})[phid]
				if board == nil {
					// The task is on a different board than `phid`.
					// Can happen with old tasks.
					continue
				}
				col := board.(map[string]interface{})["columns"].([]interface{})[0].(map[string]interface{})
				ptask.Task.Column = col["phid"].(string)
				ptask.Task.Url = fmt.Sprintf("%s/T%d", p.endpoint, el.ID)

				ptask.Task.Unscheduled = true
				if el.Fields["custom.daedalean.scheduled"] != nil {
					ptask.Task.Unscheduled = !el.Fields["custom.daedalean.scheduled"].(bool)
				}

				if el.Fields["custom.daedalean.duration"] != nil {
					ptask.Task.Duration = int(el.Fields["custom.daedalean.duration"].(float64))
				}

				if el.Fields["custom.daedalean.progress"] != nil {
					ptask.Task.Progress = float32(el.Fields["custom.daedalean.progress"].(float64) / 100)
				}

				if el.Fields["custom.daedalean.start_date"] != nil {
					tm := time.Unix(int64(el.Fields["custom.daedalean.start_date"].(float64)), 0)
					ptask.Task.StartDate = tm.Format("2006-01-02")
				} else {
					ptask.Task.Unscheduled = true
				}

				if el.Fields["custom.daedalean.successors"] != nil {
					data := el.Fields["custom.daedalean.successors"].(string)
					linkData := []PLinkData{}
					if err := json.Unmarshal([]byte(data), &linkData); err != nil {
						log.Errorf("Cannot unmarshal successors in task %q titled %q: %s", taskPhid, ptask.Task.Text, err)
					} else {
						for _, ld := range linkData {
							link := &Link{}
							link.Source = taskPhid
							link.Target = ld.Target
							link.Type = ld.Type
							link.Id = generateLinkId(link)
							ptask.Links[link.Id] = link
						}
					}
				}

				if el.Fields["custom.daedalean.type"] != nil {
					val := el.Fields["custom.daedalean.type"].(string)
					if val == "daedalean:milestone" {
						ptask.Task.Type = "milestone"
					} else if val == "daedalean:project" {
						ptask.Task.Type = "project"
					} else if val == "daedalean:task" {
						ptask.Task.Type = "task"
					}
				}

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

	for _, task := range tasks {
		if task.Task.Parent != "" {
			tasks[task.Task.Parent].IsLeaf = false
		}
	}

	for _, task := range tasks {
		if task.Task.Type == "" {
			if task.IsLeaf {
				task.Task.Type = "task"
			} else {
				task.Task.Type = "project"
			}
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

func (p *Phabricator) Users() ([]User, error) {
	after := ""
	users := []User{}
	for {
		req := requests.SearchRequest{
			After: after,
		}
		var res responses.SearchResponse
		if err := p.c.Call("user.search", &req, &res); err != nil {
			return nil, err
		}

	UserLoop:
		for _, el := range res.Data {
			for _, role := range el.Fields["roles"].([]interface{}) {
				if role.(string) == "disabled" {
					continue UserLoop
				}
			}
			user := User{
				Phid:     el.PHID,
				Name:     el.Fields["username"].(string),
				RealName: el.Fields["realName"].(string),
			}
			log.Debugf("Found active user: %s (%q)", user.Name, user.RealName)
			users = append(users, user)
		}

		after = res.Cursor.After
		if after == "" {
			break
		}
	}
	return users, nil
}

func NewPhabricator(endpoint, key string) (*Phabricator, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	endpointUri := u.Scheme + "://" + u.Host
	if u.Port() != "" {
		endpointUri += ":" + u.Port()
	}

	log.Debugf("Attempting to connect to Phabricator at %q", endpointUri)

	conn, err := gonduit.Dial(endpointUri, &core.ClientOptions{
		APIToken: key,
	})
	if err != nil {
		return nil, err
	}

	log.Debugf("Created connection to Phabricator at %q", endpointUri)
	return &Phabricator{conn, endpointUri, false}, nil
}
