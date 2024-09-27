package rest

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	a "github.com/jannawro/blog/article"
	"github.com/jannawro/blog/middleware"
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
			slog.Error(err.Error(), "requestID", middleware.ReqIDFromCtx(r.Context()))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Parse the article data to get the title
		parsedArticle, err := a.ParseArticle([]byte(articleData.Article))
		if err != nil {
			slog.Error("Failed to parse article", "requestID", middleware.ReqIDFromCtx(r.Context()), "error", err)
			http.Error(w, "Invalid article format", http.StatusBadRequest)
			return
		}

		// Check if an article with the same title already exists
		existingArticle, err := h.service.GetBySlug(r.Context(), a.Slugify(parsedArticle.Title))
		if err == nil {
			// Article with the same title already exists
			slog.Info("Article with this title already exists", "requestID", middleware.ReqIDFromCtx(r.Context()), "title", parsedArticle.Title)
			http.Error(w, "An article with this title already exists", http.StatusConflict)
			return
		} else if err != a.ErrArticleNotFound {
			// An unexpected error occurred
			slog.Error(err.Error(), "requestID", middleware.ReqIDFromCtx(r.Context()))
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		slog.Debug("Creating article", "requestID", middleware.ReqIDFromCtx(r.Context()), "articleDetails", articleData.Article)
		article, err := h.service.Create(r.Context(), []byte(articleData.Article))
		if err != nil {
			slog.Error(err.Error(), "requestID", middleware.ReqIDFromCtx(r.Context()))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(article)
		if err != nil {
			slog.Error(err.Error(), "requestID", middleware.ReqIDFromCtx(r.Context()))
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}

func (h *Handler) GetAllArticles() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sortOption := a.GetSortOption(r)

		slog.Debug("Fetching all articles", "requestID", middleware.ReqIDFromCtx(r.Context()), "sortOption", sortOption)
		articles, err := h.service.GetAll(r.Context(), &sortOption)
		if err != nil {
			slog.Error(err.Error(), "requestID", middleware.ReqIDFromCtx(r.Context()))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(articles)
		if err != nil {
			slog.Error(err.Error(), "requestID", middleware.ReqIDFromCtx(r.Context()))
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}

func (h *Handler) GetArticleByTitle(slugPathParam string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue(slugPathParam)

		slog.Debug("Fetching article by slug", "requestID", middleware.ReqIDFromCtx(r.Context()), "slug", slug)
		article, err := h.service.GetBySlug(r.Context(), slug)
		if err != nil {
			if err == a.ErrArticleNotFound {
				slog.Error(err.Error(), "requestID", middleware.ReqIDFromCtx(r.Context()))
				http.Error(w, err.Error(), http.StatusNotFound)
			} else {
				slog.Error(err.Error(), "requestID", middleware.ReqIDFromCtx(r.Context()))
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		err = json.NewEncoder(w).Encode(article)
		if err != nil {
			slog.Error(err.Error(), "requestID", middleware.ReqIDFromCtx(r.Context()))
			http.Error(w, err.Error(), http.StatusInternalServerError)
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
			if err == a.ErrArticleNotFound {
				slog.Error(err.Error(), "requestID", middleware.ReqIDFromCtx(r.Context()))
				http.Error(w, err.Error(), http.StatusNotFound)
			} else {
				slog.Error(err.Error(), "requestID", middleware.ReqIDFromCtx(r.Context()))
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		err = json.NewEncoder(w).Encode(article)
		if err != nil {
			slog.Error(err.Error(), "requestID", middleware.ReqIDFromCtx(r.Context()))
			http.Error(w, err.Error(), http.StatusInternalServerError)
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
			slog.Error(err.Error(), "requestID", middleware.ReqIDFromCtx(r.Context()))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(articles)
		if err != nil {
			slog.Error(err.Error(), "requestID", middleware.ReqIDFromCtx(r.Context()))
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}

func (h *Handler) UpdateArticleByTitle(slugPathParam string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue(slugPathParam)
		var updatedArticleDate struct {
			Article string `json:"article"`
		}

		slog.Debug("Updating article", "requestID", middleware.ReqIDFromCtx(r.Context()),
			"slug", slug,
			"requestBody", r.Body,
		)
		if err := json.NewDecoder(r.Body).Decode(&updatedArticleDate); err != nil {
			slog.Error(err.Error(), "requestID", middleware.ReqIDFromCtx(r.Context()))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		article, err := h.service.UpdateBySlug(r.Context(), slug, []byte(updatedArticleDate.Article))
		if err != nil {
			if err == a.ErrArticleNotFound {
				slog.Error(err.Error(), "requestID", middleware.ReqIDFromCtx(r.Context()))
				http.Error(w, err.Error(), http.StatusNotFound)
			} else {
				slog.Error(err.Error(), "requestID", middleware.ReqIDFromCtx(r.Context()))
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		err = json.NewEncoder(w).Encode(article)
		if err != nil {
			slog.Error(err.Error(), "requestID", middleware.ReqIDFromCtx(r.Context()))
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}

func (h *Handler) DeleteArticleByTitle(slugPathParam string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue(slugPathParam)

		slog.Debug("Deleting article", "requestID", middleware.ReqIDFromCtx(r.Context()), "slug", slug)
		err := h.service.DeleteBySlug(r.Context(), slug)
		if err != nil {
			if err == a.ErrArticleNotFound {
				slog.Error(err.Error(), "requestID", middleware.ReqIDFromCtx(r.Context()))
				http.Error(w, err.Error(), http.StatusNotFound)
			} else {
				slog.Error(err.Error(), "requestID", middleware.ReqIDFromCtx(r.Context()))
				http.Error(w, err.Error(), http.StatusInternalServerError)
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
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(tags)
		if err != nil {
			slog.Error(err.Error(), "requestID", middleware.ReqIDFromCtx(r.Context()))
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}
func ParseArticle(data []byte) (*Article, error) {
	lines := strings.Split(string(data), "\n")
	article := &Article{}
	var content []string
	inContent := false

	for _, line := range lines {
		if inContent {
			content = append(content, line)
		} else if strings.TrimSpace(line) == "===" {
			inContent = true
		} else {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) != 2 {
				continue
			}
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			switch key {
			case "title":
				article.Title = value
			case "publicationDate":
				date, err := time.Parse("2006-01-02", value)
				if err != nil {
					return nil, fmt.Errorf("invalid date format: %v", err)
				}
				article.PublicationDate = date
			case "tags":
				article.Tags = strings.Split(value, ",")
			}
		}
	}

	article.Content = strings.TrimSpace(strings.Join(content, "\n"))
	article.Slug = Slugify(article.Title)

	return article, nil
}
