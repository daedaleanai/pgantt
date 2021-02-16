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

// +build ignore

package main

import (
	"github.com/daedaleanai/pgantt/pkg/pgantt"
	"github.com/ljanyst/go-srvutils/gen"
	"log"
)

func main() {
	err := gen.GenerateNodeProject(gen.Options{
		ProjectPath:  "../../ui",
		BuildProject: true,
		Assets:       pgantt.Assets,
		PackageName:  "pgantt",
		BuildTags:    "!dev",
		VariableName: "Assets",
		Filename:     "assets_prod.go",
	})

	if err != nil {
		log.Fatal(err)
	}
}
