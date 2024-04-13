package main

import (
	"log"
	"net/http"

	"github.com/jannawro/blog/middleware"
	"github.com/jannawro/blog/static"
	pages "github.com/jannawro/blog/views/pages"
)

func main() {
	router := http.NewServeMux()
	static.Mount(router)

	router.HandleFunc("GET /hello/{name}", func(w http.ResponseWriter, r *http.Request) {
		name := r.PathValue("name")
		err := pages.Index(name).Render(r.Context(), w)
		if err != nil {
			log.Println("An error occured: ", err)
		}
	})

	stack := middleware.CreateStack(
		middleware.Logging,
	)

	port := ":8888"
	server := http.Server{
		Addr:    port,
		Handler: stack(router),
	}

	log.Println("Listening on", port)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
