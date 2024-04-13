package static

import (
	"embed"
	"net/http"
)

//go:embed styles.css favicon.svg
var static embed.FS

// Mount adds a handler for the /static/ path that serves static assets
func Mount(router *http.ServeMux) {
	fileServer := http.FileServer(http.FS(static))
	router.Handle("/static/", http.StripPrefix("/static/", fileServer))
}
