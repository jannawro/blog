package html

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"

	a "github.com/jannawro/blog/article"
	"github.com/jannawro/blog/components"
	"github.com/jannawro/blog/middleware"
)

type Handler struct {
	service    *a.Service
	assetsPath string
}

func NewHandler(service *a.Service, assetsPath string) *Handler {
	return &Handler{
		service:    service,
		assetsPath: assetsPath,
	}
}

func (h *Handler) ServeArticle(slugPathParam string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Extract slug from URL path parameter
		slug := r.PathValue(slugPathParam)
		if slug == "" {
			http.Error(w, "Article title is required", http.StatusBadRequest)
			return
		}

		slog.Debug("Serving article", "requestID", middleware.ReqIDFromCtx(r.Context()), "slug", slug)

		slog.Debug("Fetching article by slug", "requestID", middleware.ReqIDFromCtx(r.Context()), "slug", slug)
		article, err := h.service.GetBySlug(ctx, slug)
		if err != nil {
			if errors.Is(err, a.ErrArticleNotFound) {
				slog.Error(err.Error(), "requestID", middleware.ReqIDFromCtx(r.Context()), "slug", slug)
				http.Error(w, err.Error(), http.StatusNotFound)
			} else {
				slog.Error(err.Error(), "requestID", middleware.ReqIDFromCtx(r.Context()), "slug", slug)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		// Generate HTML using ArticlePage component
		articlePage := components.ArticlePage(*article, h.assetsPath)

		// Render the HTML
		err = articlePage.Render(ctx, w)
		if err != nil {
			http.Error(w, "Failed to render article page", http.StatusInternalServerError)
			return
		}
	})
}

func (h *Handler) ServeBlog() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		slog.Debug("Serving blog", "requestID", middleware.ReqIDFromCtx(r.Context()))

		sortOption := a.GetSortOption(r)
		tags := r.URL.Query()["tag"]

		var articles a.Articles
		var err error

		if len(tags) > 0 {
			slog.Debug("Fetching articles by tags: "+strings.Join(tags, ", "),
				"requestID", middleware.ReqIDFromCtx(r.Context()),
				"sortOption", sortOption,
			)
			articles, err = h.service.GetByTags(ctx, tags, &sortOption)
		} else {
			slog.Debug("Fetching all articles", "requestID", middleware.ReqIDFromCtx(r.Context()), "sortOption", sortOption)
			articles, err = h.service.GetAll(ctx, &sortOption)
		}

		if err != nil {
			slog.Error(err.Error(), "requestID", middleware.ReqIDFromCtx(r.Context()))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		blogComponent := components.Blog(articles, h.assetsPath)

		err = blogComponent.Render(ctx, w)
		if err != nil {
			slog.Error(err.Error(), "requestID", middleware.ReqIDFromCtx(r.Context()))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

func (h *Handler) ServeIndex() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		slog.Debug("Serving index", "requestID", middleware.ReqIDFromCtx(r.Context()))

		slog.Debug("Fetching all tags", "requestID", middleware.ReqIDFromCtx(r.Context()))
		tags, err := h.service.GetAllTags(ctx)
		if err != nil {
			http.Error(w, "Failed to fetch tags", http.StatusInternalServerError)
			return
		}

		// Initialize map for tagged articles
		taggedArticles := make(map[string][]a.Article)

		// Fetch articles for each tag
		for _, tag := range tags {
			slog.Debug("Fetching articles by tags: "+strings.Join(tags, ", "),
				"requestID", middleware.ReqIDFromCtx(r.Context()),
			)
			articles, err := h.service.GetByTags(ctx, []string{tag}, nil)
			if err != nil {
				http.Error(w, "Failed to fetch articles for tag: "+tag, http.StatusInternalServerError)
				return
			}
			taggedArticles[tag] = articles
		}

		// Create and render TagIndexPage component
		indexPage := components.TagIndexPage(taggedArticles, h.assetsPath)
		err = indexPage.Render(ctx, w)
		if err != nil {
			http.Error(w, "Failed to render index page", http.StatusInternalServerError)
			return
		}
	})
}
