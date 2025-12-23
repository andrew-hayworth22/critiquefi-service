package main

import (
	"io/fs"
	"net/http"
	"path"
	"strings"
)

func frontendHandler(fileSystem fs.FS) func(http.ResponseWriter, *http.Request) {
	fileServer := http.FileServer(http.FS(fileSystem))

	return func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/") {
			http.NotFound(w, r)
			return
		}

		r.URL.Path = path.Clean(r.URL.Path)

		fileServer.ServeHTTP(w, r)
	}
}
