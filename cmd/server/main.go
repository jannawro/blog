package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/jannawro/blog/article"
	"github.com/jannawro/blog/handlers/assets"
	"github.com/jannawro/blog/handlers/html"
	"github.com/jannawro/blog/handlers/rest"
	"github.com/jannawro/blog/middleware"
	"github.com/jannawro/blog/repository/postgres"
	"github.com/jannawro/blog/repository/postgres/migrations"
)

var (
	port        string
	apiKey      string
	databaseURL string
	logLevel    string
)

const assetsPath = "/assets/"

func main() {
	parseArguments()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: parseLogLevel(logLevel),
	}))
	slog.SetDefault(logger)

	postgresDatabase, err := postgres.NewDatabase(databaseURL)
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

	assetsRouter := http.NewServeMux()
	assetsRouter.Handle("GET "+assetsPath, assets.Serve(assetsPath))

	frontendRouter := http.NewServeMux()
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
	mainRouter.Handle(assetsPath, assetsRouter)
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
	flag.StringVar(&port, "port", os.Getenv("PORT"), "The port the server should listen on. The default is 8888.")
	flag.StringVar(&apiKey, "api-key", os.Getenv("API_KEY"), "API Key for the /api endpoints.")
	flag.StringVar(&databaseURL,
		"database-url",
		os.Getenv("DATABASE_URL"),
		"Database URL. Should be a connection string containing auth credentials and configuration.",
	)
	flag.StringVar(&logLevel, "log-level", os.Getenv("LOG_LEVEL"), "Set the log level (debug, info, warn, error)")
	flag.Parse()
}

func parseLogLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
