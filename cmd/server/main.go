package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/jannawro/blog/article"
	"github.com/jannawro/blog/handlers/html"
	"github.com/jannawro/blog/handlers/rest"
	"github.com/jannawro/blog/middleware"
	"github.com/jannawro/blog/repository"
	"github.com/jannawro/blog/static"
)

var (
	port   string
	apiKey string
)

func main() {
	parseArguments()

	mockRepo := initMockRepository()
	articleService := article.NewService(mockRepo)
	htmlHandler := html.NewHandler(articleService)
	restHandler := rest.NewHandler(articleService)

	frontendRouter := http.NewServeMux()
	frontendRouter.Handle("GET /static/", static.Handler("/static/"))
	frontendRouter.Handle("GET /", htmlHandler.ServeBlog())
	frontendRouter.Handle("GET /index", htmlHandler.ServeIndex())
	frontendRouter.Handle("GET /article/{title}", htmlHandler.ServeArticle())
	frontendStack := middleware.CreateStack(
		middleware.Logging,
	)

	apiRouter := http.NewServeMux()
	apiRouter.Handle("POST /api/articles", restHandler.CreateArticle())
	apiRouter.Handle("GET /api/articles", restHandler.GetAllArticles())
	apiRouter.Handle("GET /api/articles/title/{title}", restHandler.GetArticleByTitle())
	apiRouter.Handle("GET /api/articles/id/{id}", restHandler.GetArticleByID())
	apiRouter.Handle("GET /api/articles/tags", restHandler.GetArticlesByTags())
	apiRouter.Handle("PUT /api/articles/{title}", restHandler.UpdateArticleByTitle())
	apiRouter.Handle("DELETE /api/articles/{title}", restHandler.DeleteArticleByTitle())
	apiRouter.Handle("GET /api/tags", restHandler.GetAllTags())
	apiStack := middleware.CreateStack(
		middleware.Logging,
		middleware.APIKeyAuth(middleware.APIKeyConfig{
			KeyName: "X-API-Key",
			Keys: map[string]bool{
				apiKey: true,
			},
		}),
	)

	mainRouter := http.NewServeMux()
	mainRouter.Handle("/", frontendStack(frontendRouter))
	mainRouter.Handle("/api/", apiStack(apiRouter))

	server := http.Server{
		Addr:    ":" + port,
		Handler: mainRouter,
	}

	log.Println("Listening on", port)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func parseArguments() {
	flag.StringVar(&port, "port", "8888", "The port the server should listen on. The default is 8888.")
	flag.StringVar(&apiKey, "api-key", "", "API Key for the /api endpoints.")
	flag.Parse()
}

func initMockRepository() *repository.MockRepository {
	mockRepo := repository.NewMockRepository()

	sampleArticles := []article.Article{
		{
			ID:              1,
			Title:           "Getting Started with Go",
			Slug:            "getting-started-with-go",
			Content:         "Go is a statically typed, compiled programming language designed at Google...",
			Tags:            []string{"go", "programming", "beginner"},
			PublicationDate: time.Now().AddDate(0, 0, -7),
		},
		{
			ID:              2,
			Title:           "Web Development with Go",
			Slug:            "web-development-with-go",
			Content:         "Go is an excellent choice for web development due to its simplicity and performance...",
			Tags:            []string{"go", "web-development", "backend"},
			PublicationDate: time.Now().AddDate(0, 0, -3),
		},
		{
			ID:              3,
			Title:           "Concurrency in Go",
			Slug:            "concurrency-in-go",
			Content:         "One of Go's standout features is its built-in support for concurrency...",
			Tags:            []string{"go", "concurrency", "advanced"},
			PublicationDate: time.Now().AddDate(0, 0, -1),
		},
	}

	mockRepo.SetArticles(sampleArticles)
	return mockRepo
}
