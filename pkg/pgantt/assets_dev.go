// +build dev

package pgantt

import "net/http"

var Assets http.FileSystem = http.Dir("../../ui/build")
