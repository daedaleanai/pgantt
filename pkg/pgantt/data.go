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

type User struct {
	Name     string `json:"name"`
	Phid     string `json:"phid"`
	RealName string `json:"real_name"`
}

type Task struct {
	// The json name needs to be "id" otherwise the DHTMLX Gantt library says
	// "Task not found id=...".
	Phid        string  `json:"id"`
	Parent      string  `json:"parent"`
	Text        string  `json:"text"`
	Type        string  `json:"type"`
	StartDate   string  `json:"start_date"`
	Duration    int     `json:"duration"`
	Progress    float32 `json:"progress"`
	Open        bool    `json:"open"`
	Unscheduled bool    `json:"unscheduled"`
	Column      string  `json:"column"`
	Url         string  `json:"url"`
}

type Link struct {
	// "Source#Target#Type"
	Id string `json:"id"`
	// Source task Phid.
	Source string `json:"source"`
	// Target task Phid.
	Target string `json:"target"`
	Type   string `json:"type"`
}

type PlanningData struct {
	Data  []Task `json:"data"`
	Links []Link `json:"links"`
}
