package rest

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/jannawro/blog/article"
)

type Handler struct {
	service *article.Service
}

func NewHandler(service *article.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) CreateArticle(w http.ResponseWriter, r *http.Request) {
	var articleData []byte
	if err := json.NewDecoder(r.Body).Decode(&articleData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	article, err := h.service.Create(r.Context(), articleData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(article)
}

func (h *Handler) GetAllArticles(w http.ResponseWriter, r *http.Request) {
	sortBy := r.URL.Query().Get("sort")
	var sortOption *article.SortOption
	if sortBy != "" {
		so := article.SortOption(sortBy)
		sortOption = &so
	}

	articles, err := h.service.GetAll(r.Context(), sortOption)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(articles)
}

func (h *Handler) GetArticleByTitle(w http.ResponseWriter, r *http.Request) {
	title := r.PathValue("title")
	article, err := h.service.GetByTitle(r.Context(), title)
	if err != nil {
		if err == article.ErrArticleNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	json.NewEncoder(w).Encode(article)
}

func (h *Handler) GetArticleByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	article, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		if err == article.ErrArticleNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	json.NewEncoder(w).Encode(article)
}

func (h *Handler) GetArticlesByTags(w http.ResponseWriter, r *http.Request) {
	tags := r.URL.Query()["tag"]
	sortBy := r.URL.Query().Get("sort")
	var sortOption *article.SortOption
	if sortBy != "" {
		so := article.SortOption(sortBy)
		sortOption = &so
	}

	articles, err := h.service.GetByTags(r.Context(), tags, sortOption)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(articles)
}

func (h *Handler) UpdateArticleByTitle(w http.ResponseWriter, r *http.Request) {
	title := r.PathValue("title")
	var updatedData []byte
	if err := json.NewDecoder(r.Body).Decode(&updatedData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	article, err := h.service.UpdateByTitle(r.Context(), title, updatedData)
	if err != nil {
		if err == article.ErrArticleNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	json.NewEncoder(w).Encode(article)
}

func (h *Handler) DeleteArticleByTitle(w http.ResponseWriter, r *http.Request) {
	title := r.PathValue("title")
	err := h.service.DeleteByTitle(r.Context(), title)
	if err != nil {
		if err == article.ErrArticleNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetAllTags(w http.ResponseWriter, r *http.Request) {
	tags, err := h.service.GetAllTags(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(tags)
}
