package html

import (
	"net/http"

	"github.com/jannawro/blog/article"
	"github.com/jannawro/blog/components"
)

type HTMLHandler struct {
	service *article.Service
}

func NewHTMLHandler(service *article.Service) *HTMLHandler {
	return &HTMLHandler{service: service}
}

func (h *HTMLHandler) ServeBlog() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		articles, err := h.service.GetAll(ctx, nil) // Assuming no sorting for now
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
