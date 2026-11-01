package web

import (
	"embed"
	"io/fs"
	"net/http"
)

// DistFS embeds the distribution build of the web frontend.
//
//go:embed all:dist
var DistFS embed.FS

// SubFS returns the sub-filesystem pointing to the bundled dist directory.
func SubFS() (http.FileSystem, fs.FS, error) {
	sub, err := fs.Sub(DistFS, "dist")
	if err != nil {
		return nil, nil, err
	}
	return http.FS(sub), sub, nil
}
