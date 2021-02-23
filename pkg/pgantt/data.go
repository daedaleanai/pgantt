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

type Column struct {
	Name string `json:"name"`
	Phid string `json:"phid"`
}

type Project struct {
	Name    string   `json:"name"`
	Phid    string   `json:"phid"`
	Columns []Column `json:"columns"`
}

type Task struct {
	Id          string  `json:"id"`
	Parent      string  `json:"parent,omitempty"`
	Text        string  `json:"text"`
	Type        string  `json:"type,omitempty"`
	StartDate   string  `json:"start_date,omitempty"`
	Duration    string  `json:"duration,omitempty"`
	Progress    float32 `json:"progress"`
	Open        bool    `json:"open"`
	Unscheduled bool    `json:"unscheduled"`
	Column      string  `json:"column"`
	Url         string  `json:"url"`
}

type Link struct {
	Id     string `json:"id"`
	Source string `json:"source"`
	Target string `json:"target"`
	Type   string `json:"type"`
}

type PlanningData struct {
	Data  []Task `json:"data"`
	Links []Link `json:"links"`
}
