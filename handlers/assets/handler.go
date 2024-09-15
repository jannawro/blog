package assets

import (
	"embed"
	"net/http"
)

//go:embed styles.css favicon.png red_door_nobg.png red_door_cropped.png
var assets embed.FS

// Serve returns an http.Handler that serves static assets on "path" endpoint
func Serve(path string) http.Handler {
	fileServer := http.FileServer(http.FS(assets))
	return http.StripPrefix(path, fileServer)
}
