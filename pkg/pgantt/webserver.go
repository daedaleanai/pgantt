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
