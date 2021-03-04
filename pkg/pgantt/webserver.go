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

//go:generate go run --tags=dev assets_generate.go

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"

	"github.com/ljanyst/go-srvutils/fs"
	log "github.com/sirupsen/logrus"
)

type Response struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

type ActionStatus struct {
	Action string `json:"action"`
	Tid    string `json:"tid,omitempty"`
}

type StateHandler struct {
	s *StateManager
}

func setupHeader(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache, private, max-age=0")
}

func writeError(w http.ResponseWriter, code int, err error) {
	resp := Response{
		"ERROR",
		err.Error(),
	}
	log.Errorf("Writing an error message to the client: %s", err)

	bytes, err := json.Marshal(resp)
	if err != nil {
		log.Errorf("Cannot serialize error respense: %s", err)
		return
	}

	setupHeader(w)
	w.WriteHeader(code)
	w.Write(bytes)
}

func writeData(w http.ResponseWriter, data interface{}) {
	resp := Response{
		"SUCCESS",
		data,
	}

	bytes, err := json.Marshal(resp)
	if err != nil {
		log.Errorf("Cannot serialize data respense: %s", err)
		return
	}
	setupHeader(w)
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}

type ProjectsHandler StateHandler
type PlanProvider StateHandler
type PlanEditor StateHandler

func (h ProjectsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	projects := h.s.Projects()
	writeData(w, projects)
}

func (h PlanProvider) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	planning := h.s.PlanningData(r.URL.Path)
	if planning == nil {
		writeError(w, 404, fmt.Errorf("Unknown project %s", r.URL.Path))
		return
	}
	writeData(w, planning)
}

func (h PlanEditor) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		setupHeader(w)
		w.WriteHeader(http.StatusOK)
		return
	}

	typ := path.Base(r.URL.Path)
	phid := path.Dir(r.URL.Path)

	if typ != "task" && typ != "link" {
		writeError(w, 400, fmt.Errorf("Unsupported %s request for %q", r.Method, typ))
		return
	}

	defer h.s.SyncTasks()
	status := ActionStatus{}
	var err error
	var id string

	if typ == "task" {
		var task Task
		err = json.NewDecoder(r.Body).Decode(&task)
		if err != nil {
			writeError(w, 400, err)
			return
		}

		if r.Method == "DELETE" {
			writeError(w, 400, fmt.Errorf("Task deletion is not supported"))
			return
		}

		id, err = h.s.EditTask(phid, &task)
		if err != nil {
			writeError(w, 400, err)
			return
		}
	}

	if typ == "link" {
		if r.Method == "DELETE" {
			var linkId string
			err = json.NewDecoder(r.Body).Decode(&linkId)
			if err != nil {
				writeError(w, 400, err)
				return
			}

			err = h.s.DeleteLink(phid, linkId)
			if err != nil {
				writeError(w, 400, err)
				return
			}

			writeData(w, ActionStatus{"deleted", ""})
			return
		}

		var link Link
		err = json.NewDecoder(r.Body).Decode(&link)
		if err != nil {
			writeError(w, 400, err)
			return
		}

		if r.Method == "PUT" {
			writeError(w, 400, fmt.Errorf("Link edition is not supported"))
			return
		}

		id, err = h.s.CreateLink(phid, &link)
		if err != nil {
			writeError(w, 400, err)
			return
		}
	}

	if r.Method == "POST" {
		status.Action = "inserted"
		status.Tid = id
	} else {
		status.Action = "updated"
	}
	writeData(w, status)
}

func RunWebServer(sm *StateManager, opts *Opts) {
	assets := &fs.Index404Fs{Assets}
	ui := http.FileServer(assets)
	http.Handle("/", ui)
	http.Handle("/api/projects", ProjectsHandler{sm})
	http.Handle("/api/plan/", http.StripPrefix("/api/plan/", PlanProvider{sm}))
	http.Handle("/api/edit/", http.StripPrefix("/api/edit/", PlanEditor{sm}))
	addressString := fmt.Sprintf("localhost:%d", opts.PGantt.Port)
	log.Infof("Serving at: http://%s", addressString)
	log.Fatal("Server failure: ", http.ListenAndServe(addressString, nil))
}
