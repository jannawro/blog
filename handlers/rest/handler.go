package rest

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	a "github.com/jannawro/blog/article"
	"github.com/jannawro/blog/middleware"
)

const internalServerErrorMsg = "Internal server error"

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
			slog.Error(err.Error(), "requestID", middleware.ReqIDFromCtx(r.Context()))
			http.Error(w, internalServerErrorMsg, http.StatusBadRequest)
			return
		}

		var unmarshaledArticle a.Article
		err := a.UnmarshalToArticle([]byte(articleData.Article), &unmarshaledArticle)
		if err != nil {
			slog.Error("Failed to unmarshal article", "requestID", middleware.ReqIDFromCtx(r.Context()), "error", err)
			http.Error(w, "Invalid article format", http.StatusBadRequest)
			return
		}

		// Check if an article with the same title already exists
		_, err = h.service.GetBySlug(r.Context(), unmarshaledArticle.Slug)
		if err == nil {
			// Article with the same title already exists
			slog.Info("Article with this title already exists",
				"requestID", middleware.ReqIDFromCtx(r.Context()),
				"title", unmarshaledArticle.Title,
			)
			http.Error(w, "An article with this title already exists", http.StatusConflict)
			return
		} else if !errors.Is(err, a.ErrArticleNotFound) {
			// An unexpected error occurred
			slog.Error(err.Error(), "requestID", middleware.ReqIDFromCtx(r.Context()))
			http.Error(w, internalServerErrorMsg, http.StatusInternalServerError)
			return
		}

		slog.Debug("Creating article",
			"requestID", middleware.ReqIDFromCtx(r.Context()),
			"articleDetails", articleData.Article,
		)
		createdArticle, err := h.service.Create(r.Context(), unmarshaledArticle)
		if err != nil {
			slog.Error(err.Error(), "requestID", middleware.ReqIDFromCtx(r.Context()))
			http.Error(w, internalServerErrorMsg, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(createdArticle)
		if err != nil {
			slog.Error(err.Error(), "requestID", middleware.ReqIDFromCtx(r.Context()))
			http.Error(w, internalServerErrorMsg, http.StatusInternalServerError)
		}
	})
}

func (h *Handler) GetAllArticles() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sortOption := a.GetSortOption(r)

		slog.Debug("Fetching all articles", "requestID", middleware.ReqIDFromCtx(r.Context()), "sortOption", sortOption)
		articles, err := h.service.GetAll(r.Context(), &sortOption)
		if err != nil {
			if errors.Is(err, a.ErrArticlesNotFound) {
				slog.Error(err.Error(), "requestID", middleware.ReqIDFromCtx(r.Context()))
				http.Error(w, a.ErrArticleNotFound.Error(), http.StatusNotFound)
			} else {
				slog.Error(err.Error(), "requestID", middleware.ReqIDFromCtx(r.Context()))
				http.Error(w, internalServerErrorMsg, http.StatusInternalServerError)
			}
			return
		}

		err = json.NewEncoder(w).Encode(articles)
		if err != nil {
			slog.Error(err.Error(), "requestID", middleware.ReqIDFromCtx(r.Context()))
			http.Error(w, internalServerErrorMsg, http.StatusInternalServerError)
		}
	})
}

func (h *Handler) GetArticleByTitle(slugPathParam string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue(slugPathParam)

		slog.Debug("Fetching article by slug", "requestID", middleware.ReqIDFromCtx(r.Context()), "slug", slug)
		article, err := h.service.GetBySlug(r.Context(), slug)
		if err != nil {
			if errors.Is(err, a.ErrArticleNotFound) {
				slog.Error(err.Error(), "requestID", middleware.ReqIDFromCtx(r.Context()))
				http.Error(w, a.ErrArticleNotFound.Error(), http.StatusNotFound)
			} else {
				slog.Error(err.Error(), "requestID", middleware.ReqIDFromCtx(r.Context()))
				http.Error(w, internalServerErrorMsg, http.StatusInternalServerError)
			}
			return
		}

		err = json.NewEncoder(w).Encode(article)
		if err != nil {
			slog.Error(err.Error(), "requestID", middleware.ReqIDFromCtx(r.Context()))
			http.Error(w, internalServerErrorMsg, http.StatusInternalServerError)
		}
	})
}

func (h *Handler) GetArticleByID(idPathParam string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue(idPathParam)

		slog.Debug("Fetching article by ID", "requestID", middleware.ReqIDFromCtx(r.Context()), "ID", idStr)
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			slog.Error(err.Error(), "requestID", middleware.ReqIDFromCtx(r.Context()))
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		article, err := h.service.GetByID(r.Context(), id)
		if err != nil {
			if errors.Is(err, a.ErrArticleNotFound) {
				slog.Error(err.Error(), "requestID", middleware.ReqIDFromCtx(r.Context()))
				http.Error(w, a.ErrArticleNotFound.Error(), http.StatusNotFound)
			} else {
				slog.Error(err.Error(), "requestID", middleware.ReqIDFromCtx(r.Context()))
				http.Error(w, internalServerErrorMsg, http.StatusInternalServerError)
			}
			return
		}

		err = json.NewEncoder(w).Encode(article)
		if err != nil {
			slog.Error(err.Error(), "requestID", middleware.ReqIDFromCtx(r.Context()))
			http.Error(w, internalServerErrorMsg, http.StatusInternalServerError)
		}
	})
}

func (h *Handler) GetArticlesByTags() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tags := r.URL.Query()["tag"]
		sortOption := a.GetSortOption(r)

		slog.Debug("Fetching articles by tags: "+strings.Join(tags, ", "),
			"requestID", middleware.ReqIDFromCtx(r.Context()),
			"sortOption", sortOption,
		)
		articles, err := h.service.GetByTags(r.Context(), tags, &sortOption)
		if err != nil {
			if errors.Is(err, a.ErrArticlesNotFound) {
				slog.Error(err.Error(), "requestID", middleware.ReqIDFromCtx(r.Context()))
				http.Error(w, a.ErrArticleNotFound.Error(), http.StatusNotFound)
			} else {
				slog.Error(err.Error(), "requestID", middleware.ReqIDFromCtx(r.Context()))
				http.Error(w, internalServerErrorMsg, http.StatusInternalServerError)
			}
			return
		}

		err = json.NewEncoder(w).Encode(articles)
		if err != nil {
			slog.Error(err.Error(), "requestID", middleware.ReqIDFromCtx(r.Context()))
			http.Error(w, internalServerErrorMsg, http.StatusInternalServerError)
		}
	})
}

func (h *Handler) UpdateArticleByTitle(slugPathParam string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue(slugPathParam)
		var updatedArticleData struct {
			Article string `json:"article"`
		}

		if err := json.NewDecoder(r.Body).Decode(&updatedArticleData); err != nil {
			slog.Error(err.Error(), "requestID", middleware.ReqIDFromCtx(r.Context()))
			http.Error(w, internalServerErrorMsg, http.StatusBadRequest)
			return
		}
		slog.Debug("Updating article", "requestID", middleware.ReqIDFromCtx(r.Context()),
			"slug", slug,
			"articleDetails", updatedArticleData.Article,
		)

		var unmarshaledArticle a.Article
		err := a.UnmarshalToArticle([]byte(updatedArticleData.Article), &unmarshaledArticle)
		if err != nil {
			slog.Error("Failed to unmarshal article", "requestID", middleware.ReqIDFromCtx(r.Context()), "error", err)
			http.Error(w, "Invalid article format", http.StatusBadRequest)
			return
		}

		updatedArticle, err := h.service.UpdateBySlug(r.Context(), slug, unmarshaledArticle)
		if err != nil {
			if errors.Is(err, a.ErrArticleNotFound) {
				slog.Error(err.Error(), "requestID", middleware.ReqIDFromCtx(r.Context()))
				http.Error(w, a.ErrArticleNotFound.Error(), http.StatusNotFound)
			} else {
				slog.Error(err.Error(), "requestID", middleware.ReqIDFromCtx(r.Context()))
				http.Error(w, internalServerErrorMsg, http.StatusInternalServerError)
			}
			return
		}

		err = json.NewEncoder(w).Encode(updatedArticle)
		if err != nil {
			slog.Error(err.Error(), "requestID", middleware.ReqIDFromCtx(r.Context()))
			http.Error(w, internalServerErrorMsg, http.StatusInternalServerError)
		}
	})
}

func (h *Handler) DeleteArticleByTitle(slugPathParam string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue(slugPathParam)

		slog.Debug("Deleting article", "requestID", middleware.ReqIDFromCtx(r.Context()), "slug", slug)
		err := h.service.DeleteBySlug(r.Context(), slug)
		if err != nil {
			if errors.Is(err, a.ErrArticleNotFound) {
				slog.Error(err.Error(), "requestID", middleware.ReqIDFromCtx(r.Context()))
				http.Error(w, a.ErrArticleNotFound.Error(), http.StatusNotFound)
			} else {
				slog.Error(err.Error(), "requestID", middleware.ReqIDFromCtx(r.Context()))
				http.Error(w, internalServerErrorMsg, http.StatusInternalServerError)
			}
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})
}

func (h *Handler) GetAllTags() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Debug("Fetching all tags", "requestID", middleware.ReqIDFromCtx(r.Context()))
		tags, err := h.service.GetAllTags(r.Context())
		if err != nil {
			slog.Error(err.Error(), "requestID", middleware.ReqIDFromCtx(r.Context()))
			http.Error(w, internalServerErrorMsg, http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(tags)
		if err != nil {
			slog.Error(err.Error(), "requestID", middleware.ReqIDFromCtx(r.Context()))
			http.Error(w, internalServerErrorMsg, http.StatusInternalServerError)
		}
	})
}
