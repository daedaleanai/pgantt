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

	"github.com/ljanyst/go-srvutils/fs"
	log "github.com/sirupsen/logrus"
)

type Response struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

type StateHandler struct {
	s *StateManager
}

func writeError(w http.ResponseWriter, code int, err error) {
	resp := Response{
		"ERROR",
		err.Error(),
	}

	bytes, err := json.Marshal(resp)
	if err != nil {
		log.Errorf("Cannot serialize error respense: %s", err)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache, private, max-age=0")
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
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache, private, max-age=0")
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}

type ProjectsHandler StateHandler
type PlanHandler StateHandler

func (h ProjectsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	projects := h.s.Projects()
	writeData(w, projects)
}

func (h PlanHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	planning := h.s.PlanningData(r.URL.Path)
	if planning == nil {
		writeError(w, 404, fmt.Errorf("Unknown project %s", r.URL.Path))
		return
	}
	writeData(w, planning)
}

func RunWebServer(sm *StateManager, opts *Opts) {
	assets := &fs.Index404Fs{Assets}
	ui := http.FileServer(assets)
	http.Handle("/", ui)
	http.Handle("/api/projects", ProjectsHandler{sm})
	http.Handle("/api/plan/", http.StripPrefix("/api/plan/", PlanHandler{sm}))
	addressString := fmt.Sprintf("localhost:%d", opts.Port)
	log.Infof("Serving at: http://%s", addressString)
	log.Fatal("Server failure: ", http.ListenAndServe(addressString, nil))
}
