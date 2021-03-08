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
	"fmt"
	"os"
	"os/exec"

	"github.com/daedaleanai/pgantt/pkg/pgantt"
	"github.com/ljanyst/go-srvutils/gen"
	log "github.com/sirupsen/logrus"
)

func main() {
	projectPath := "../../ui"

	hook := func() error {
		ganttUrl := os.Getenv("DDLN_GANTT")
		if ganttUrl != "" {
			log.Infof("Installing the DDLN version of dhtml-gantt...")
			cmd := exec.Command("npm", "install", ganttUrl)
			cmd.Dir = projectPath
			if err := cmd.Run(); err != nil {
				log.Errorf("See: https://docs.google.com/document/d/1Vqt_ojK2kxSVHi_-lziuE_lvtgTddpvbKt3t-3x-qBc")
				output, _ := cmd.CombinedOutput()
				return fmt.Errorf("Cannot install the DDLN version of dhtml-pgantt:\n%s", string(output))
			}
		}
		return nil
	}

	err := gen.GenerateNodeProject(gen.Options{
		ProjectPath:     projectPath,
		PostInstallHook: hook,
		BuildProject:    true,
		Assets:          pgantt.Assets,
		PackageName:     "pgantt",
		BuildTags:       "!dev",
		VariableName:    "Assets",
		Filename:        "assets_prod.go",
	})

	if err != nil {
		log.Fatal(err)
	}
}
