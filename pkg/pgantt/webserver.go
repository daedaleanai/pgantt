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
	"fmt"
	"net/http"

	"github.com/ljanyst/go-srvutils/fs"
	log "github.com/sirupsen/logrus"
)

func RunWebServer(opts *Opts) {
	assets := &fs.Index404Fs{Assets}
	ui := http.FileServer(assets)
	http.Handle("/", ui)
	addressString := fmt.Sprintf("localhost:%d", opts.Port)
	log.Infof("Serving at: http://%s", addressString)
	log.Fatal("Server failure: ", http.ListenAndServe(addressString, nil))
}
