package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/jannawro/blog/middleware"
	"github.com/jannawro/blog/static"
)

var (
	port string
)

func main() {
	parseArguments()

	router := http.NewServeMux()

	router.Handle("GET /static/", static.Handler("/static/"))
	router.HandleFunc("GET /", placeholderHandler())
	router.HandleFunc("GET /about", placeholderHandler())

	stack := middleware.CreateStack(
		middleware.Logging,
	)

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

func parseArguments() {
	flag.StringVar(&port, "port", "8888", "The port the server should listen on. The default is 8888.")
}

func placeholderHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		response := "You called a placeholder!"
		log.Println(response)
		_, err := w.Write([]byte(response))
		if err != nil {
			panic(err)
		}
	}
}
