package html

import (
	"errors"
	"net/http"

	a "github.com/jannawro/blog/article"
	"github.com/jannawro/blog/components"
)

type Handler struct {
	service *a.Service
}

func NewHandler(service *a.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) ServeArticle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Extract title from URL path parameter
		title := r.PathValue("title")
		if title == "" {
			http.Error(w, "Article title is required", http.StatusBadRequest)
			return
		}

		// Fetch article by title
		article, err := h.service.GetBySlug(ctx, title)
		if err != nil {
			if errors.Is(err, a.ErrArticleNotFound) {
				http.Error(w, "Article not found", http.StatusNotFound)
			} else {
				http.Error(w, "Failed to fetch article", http.StatusInternalServerError)
			}
			return
		}

		// Generate HTML using ArticlePage component
		articlePage := components.ArticlePage(*article)

		// Render the HTML
		err = articlePage.Render(ctx, w)
		if err != nil {
			http.Error(w, "Failed to render article page", http.StatusInternalServerError)
			return
		}
	}
}

func (h *Handler) ServeBlog() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		sortOption := getSortOption(r)
		tags := r.URL.Query()["tag"]

		var articles a.Articles
		var err error

		if len(tags) > 0 {
			articles, err = h.service.GetByTags(ctx, tags, &sortOption)
		} else {
			articles, err = h.service.GetAll(ctx, &sortOption)
		}

		if err != nil {
			http.Error(w, "Failed to fetch articles", http.StatusInternalServerError)
			return
		}

		blogComponent := components.Blog(articles)

		err = blogComponent.Render(ctx, w)
		if err != nil {
			http.Error(w, "Failed to render blog", http.StatusInternalServerError)
			return
		}
	}
}

func (h *Handler) ServeIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Fetch all tags
		tags, err := h.service.GetAllTags(ctx)
		if err != nil {
			http.Error(w, "Failed to fetch tags", http.StatusInternalServerError)
			return
		}

		// Initialize map for tagged articles
		taggedArticles := make(map[string][]a.Article)

		// Fetch articles for each tag
		for _, tag := range tags {
			articles, err := h.service.GetByTags(ctx, []string{tag}, nil)
			if err != nil {
				http.Error(w, "Failed to fetch articles for tag: "+tag, http.StatusInternalServerError)
				return
			}
			taggedArticles[tag] = articles
		}

		// Create and render TagIndexPage component
		indexPage := components.TagIndexPage(taggedArticles)
		err = indexPage.Render(ctx, w)
		if err != nil {
			http.Error(w, "Failed to render index page", http.StatusInternalServerError)
			return
		}
	}
}

func getSortOption(r *http.Request) a.SortOption {
	sortParam := r.URL.Query().Get("sort")
	switch sortParam {
	case "title":
		return a.SortByTitle
	case "id":
		return a.SortByID
	case "date":
		fallthrough
	default:
		return a.SortByPublicationDate
	}
}
