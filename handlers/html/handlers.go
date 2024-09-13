package html

import (
	"net/http"

	"github.com/jannawro/blog/article"
	"github.com/jannawro/blog/components"
)

type Handler struct {
	service *article.Service
}

func NewHandler(service *article.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) ServeBlog() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		sortOption := getSortOption(r)
		tag := r.URL.Query().Get("tag")

		var articles article.Articles
		var err error

		if tag != "" {
			articles, err = h.service.GetByTags(ctx, []string{tag}, &sortOption)
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

func getSortOption(r *http.Request) article.SortOption {
	sortParam := r.URL.Query().Get("sort")
	switch sortParam {
	case "title":
		return article.SortByTitle
	case "id":
		return article.SortByID
	case "date":
		fallthrough
	default:
		return article.SortByPublicationDate
	}
}
