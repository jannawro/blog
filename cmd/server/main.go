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

func articleHandler(repo article.ArticleRepository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		articleID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			http.Error(w, "Invalid article ID", http.StatusBadRequest)
			return
		}

		ctx := r.Context()
		a, err := repo.GetByID(ctx, articleID)
		if err != nil {
			http.Error(w, "Article not found", http.StatusNotFound)
			return
		}

		component := pages.ArticlePage(*a)
		err = component.Render(ctx, w)
		if err != nil {
			http.Error(w, "Error rendering page", http.StatusInternalServerError)
			return
		}
	}
}
