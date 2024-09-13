package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/jannawro/blog/article"
	"github.com/jannawro/blog/handlers/html"
	"github.com/jannawro/blog/middleware"
	"github.com/jannawro/blog/repository"
	"github.com/jannawro/blog/static"
)

var (
	port string
)

func main() {
	parseArguments()

	mockRepo := initMockRepository()
	articleService := article.NewService(mockRepo)
	htmlHandler := html.NewHandler(articleService)

	router := http.NewServeMux()

	router.Handle("GET /static/", static.Handler("/static/"))
	router.HandleFunc("GET /", htmlHandler.ServeBlog())
	router.HandleFunc("GET /index", htmlHandler.ServeIndex())
	router.HandleFunc("GET /article/{id}", htmlHandler.ServeArticle())

	stack := middleware.CreateStack(
		middleware.Logging,
	)

	server := http.Server{
		Addr:    ":" + port,
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
	flag.Parse()
}

func initMockRepository() *repository.MockRepository {
	mockRepo := repository.NewMockRepository()

	sampleArticles := []article.Article{
		{
			ID:              1,
			Title:           "Getting Started with Go",
			Content:         "Go is a statically typed, compiled programming language designed at Google...",
			Tags:            []string{"go", "programming", "beginner"},
			PublicationDate: time.Now().AddDate(0, 0, -7),
		},
		{
			ID:              2,
			Title:           "Web Development with Go",
			Content:         "Go is an excellent choice for web development due to its simplicity and performance...",
			Tags:            []string{"go", "web-development", "backend"},
			PublicationDate: time.Now().AddDate(0, 0, -3),
		},
		{
			ID:              3,
			Title:           "Concurrency in Go",
			Content:         "One of Go's standout features is its built-in support for concurrency...",
			Tags:            []string{"go", "concurrency", "advanced"},
			PublicationDate: time.Now().AddDate(0, 0, -1),
		},
	}

	mockRepo.SetArticles(sampleArticles)
	return mockRepo
}
