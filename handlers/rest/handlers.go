package rest

import (
	"encoding/json"
	"net/http"
	"strconv"

	a "github.com/jannawro/blog/article"
)

type Handler struct {
	service *a.Service
}

func NewHandler(service *a.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) CreateArticle() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var articleData struct {
			Article string `json:"article"`
		}
		if err := json.NewDecoder(r.Body).Decode(&articleData); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		article, err := h.service.Create(r.Context(), []byte(articleData.Article))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(article)
	})
}

func (h *Handler) GetAllArticles() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sortBy := r.URL.Query().Get("sort")
		var sortOption *a.SortOption
		if sortBy != "" {
			so := a.SortOption(sortBy)
			sortOption = &so
		}

		articles, err := h.service.GetAll(r.Context(), sortOption)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(articles)
	})
}

func (h *Handler) GetArticleByTitle() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		title := r.PathValue("title")
		article, err := h.service.GetBySlug(r.Context(), title)
		if err != nil {
			if err == a.ErrArticleNotFound {
				http.Error(w, err.Error(), http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		json.NewEncoder(w).Encode(article)
	})
}

func (h *Handler) GetArticleByID() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		article, err := h.service.GetByID(r.Context(), id)
		if err != nil {
			if err == a.ErrArticleNotFound {
				http.Error(w, err.Error(), http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		json.NewEncoder(w).Encode(article)
	})
}

func (h *Handler) GetArticlesByTags() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tags := r.URL.Query()["tag"]
		sortBy := r.URL.Query().Get("sort")
		var sortOption *a.SortOption
		if sortBy != "" {
			so := a.SortOption(sortBy)
			sortOption = &so
		}

		articles, err := h.service.GetByTags(r.Context(), tags, sortOption)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(articles)
	})
}

func (h *Handler) UpdateArticleByTitle() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		title := r.PathValue("title")
		var updatedData []byte
		if err := json.NewDecoder(r.Body).Decode(&updatedData); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		article, err := h.service.UpdateBySlug(r.Context(), title, updatedData)
		if err != nil {
			if err == a.ErrArticleNotFound {
				http.Error(w, err.Error(), http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		json.NewEncoder(w).Encode(article)
	})
}

func (h *Handler) DeleteArticleByTitle() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		title := r.PathValue("title")
		err := h.service.DeleteBySlug(r.Context(), title)
		if err != nil {
			if err == a.ErrArticleNotFound {
				http.Error(w, err.Error(), http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})
}

func (h *Handler) GetAllTags() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tags, err := h.service.GetAllTags(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(tags)
	})
}
