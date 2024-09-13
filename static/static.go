package static

import (
	"embed"
	"net/http"
)

//go:embed styles.css favicon.png red_door.png red_door_cropped.png
var static embed.FS

// Handler returns an http.Handler that serves static assets on "path" endpoint
func Handler(path string) http.Handler {
	fileServer := http.FileServer(http.FS(static))
	return http.StripPrefix(path, fileServer)
}
