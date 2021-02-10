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
