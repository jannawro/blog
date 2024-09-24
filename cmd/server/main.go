package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"

	"github.com/jannawro/blog/article"
	"github.com/jannawro/blog/handlers/assets"
	"github.com/jannawro/blog/handlers/html"
	"github.com/jannawro/blog/handlers/rest"
	"github.com/jannawro/blog/middleware"
	"github.com/jannawro/blog/repository/postgres"
	"github.com/jannawro/blog/repository/postgres/migrations"
)

var (
	port            string
	apiKey          string
	postgresConnStr string
)

const assetsPath = "/assets/"

func main() {
	parseArguments()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)

	postgresDatabase, err := postgres.NewDatabase(postgresConnStr)
	if err != nil {
		panic(err)
	}
	postgresRepo, err := postgres.NewRepository(postgresDatabase, migrations.Files())
	if err != nil {
		panic(err)
	}
	articleService := article.NewService(postgresRepo)
	htmlHandler := html.NewHandler(articleService, assetsPath)
	restHandler := rest.NewHandler(articleService)

	frontendRouter := http.NewServeMux()
	frontendRouter.Handle("GET "+assetsPath, assets.Serve(assetsPath))
	frontendRouter.Handle("GET /", htmlHandler.ServeBlog())
	frontendRouter.Handle("GET /index", htmlHandler.ServeIndex())
	frontendRouter.Handle("GET /article/{title}", htmlHandler.ServeArticle("title"))
	frontendStack := middleware.CreateStack(
		middleware.Logging(),
	)

	apiRouter := http.NewServeMux()
	apiRouter.Handle("POST /api/articles", restHandler.CreateArticle())
	apiRouter.Handle("GET /api/articles", restHandler.GetAllArticles())
	apiRouter.Handle("GET /api/articles/title/{title}", restHandler.GetArticleByTitle("title"))
	apiRouter.Handle("GET /api/articles/id/{id}", restHandler.GetArticleByID("id"))
	apiRouter.Handle("GET /api/articles/tags", restHandler.GetArticlesByTags())
	apiRouter.Handle("PUT /api/articles/{title}", restHandler.UpdateArticleByTitle("title"))
	apiRouter.Handle("DELETE /api/articles/{title}", restHandler.DeleteArticleByTitle("title"))
	apiRouter.Handle("GET /api/tags", restHandler.GetAllTags())
	apiStack := middleware.CreateStack(
		middleware.Logging(),
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

	slog.Info("Listening on " + port)
	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func parseArguments() {
	flag.StringVar(&port, "port", "8888", "The port the server should listen on. The default is 8888.")
	flag.StringVar(&apiKey, "api-key", "", "API Key for the /api endpoints.")
	flag.StringVar(&postgresConnStr, "postgres-connection-string", "", "Connection string for a postgres database.")
	flag.Parse()
}
