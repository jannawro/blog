package rest_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jannawro/blog/article"
	"github.com/jannawro/blog/handlers/rest"
	"github.com/jannawro/blog/middleware"
	"github.com/jannawro/blog/repository"
	"github.com/stretchr/testify/assert"
)

func setupTest() (*rest.Handler, *repository.MockRepository) {
	mockRepo := repository.NewMockRepository()
	service := article.NewService(mockRepo)
	handler := rest.NewHandler(service)
	return handler, mockRepo
}

func TestCreateArticle(t *testing.T) {
	handler, mockRepo := setupTest()

	articleData := []byte(`{
		"article": "title:Test Article\npublicationDate:2023-05-15\ntags:test,article\n===\nThis is a test article."
	}`)
	expectedArticle := &article.Article{
		Title:           "Test Article",
		Slug:            "test-article",
		Content:         "This is a test article.",
		Tags:            []string{"test", "article"},
		PublicationDate: time.Date(2023, 5, 15, 0, 0, 0, 0, time.UTC),
	}

	req, err := http.NewRequest("POST", "/articles", bytes.NewBuffer(articleData))
	req = middleware.SetReqID(req)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.CreateArticle().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response article.Article
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Compare relevant fields
	assert.Equal(t, expectedArticle.Title, response.Title)
	assert.Equal(t, expectedArticle.Slug, response.Slug)
	assert.Equal(t, expectedArticle.Content, response.Content)
	assert.Equal(t, expectedArticle.Tags, response.Tags)
	assert.Equal(t, expectedArticle.PublicationDate, response.PublicationDate)

	// Verify the article was added to the mock repository
	articles, err := mockRepo.GetAll(req.Context())
	assert.NoError(t, err)
	assert.Len(t, articles, 1)
	assert.Equal(t, expectedArticle.Title, articles[0].Title)
}

func TestGetAllArticles(t *testing.T) {
	handler, mockRepo := setupTest()

	// Add some test articles to the mock repository
	mockRepo.SetArticles([]article.Article{
		{ID: 1, Title: "Article 1", Slug: "article-1"},
		{ID: 2, Title: "Article 2", Slug: "article-2"},
	})

	req, err := http.NewRequest("GET", "/articles", nil)
	req = middleware.SetReqID(req)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.GetAllArticles().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response article.Articles
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 2)
	assert.Equal(t, "Article 1", response[0].Title)
	assert.Equal(t, "Article 2", response[1].Title)
}

func TestGetArticleByTitle(t *testing.T) {
	handler, mockRepo := setupTest()

	expectedArticle := article.Article{ID: 1, Title: "Test Article", Slug: "test-article"}
	mockRepo.SetArticles([]article.Article{expectedArticle})

	req := httptest.NewRequest("GET", "/articles/test-article", nil)
	req = middleware.SetReqID(req)
	pathParam := "title"
	req.SetPathValue(pathParam, "test-article")

	rr := httptest.NewRecorder()
	handler.GetArticleByTitle(pathParam).ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response article.Article
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, expectedArticle.Title, response.Title)
	assert.Equal(t, expectedArticle.Slug, response.Slug)
}

func TestGetArticleByID(t *testing.T) {
	handler, mockRepo := setupTest()

	expectedArticle := article.Article{ID: 1, Title: "Test Article", Slug: "test-article"}
	mockRepo.SetArticles([]article.Article{expectedArticle})

	req := httptest.NewRequest("GET", "/articles/1", nil)
	req = middleware.SetReqID(req)
	pathParam := "id"
	req.SetPathValue(pathParam, "1")

	rr := httptest.NewRecorder()
	handler.GetArticleByID(pathParam).ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response article.Article
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, expectedArticle.ID, response.ID)
	assert.Equal(t, expectedArticle.Title, response.Title)
}

func TestGetArticlesByTags(t *testing.T) {
	handler, mockRepo := setupTest()

	articles := []article.Article{
		{ID: 1, Title: "Article 1", Tags: []string{"tag1", "tag2"}},
		{ID: 2, Title: "Article 2", Tags: []string{"tag2", "tag3"}},
		{ID: 3, Title: "Article 3", Tags: []string{"tag1", "tag3"}},
	}
	mockRepo.SetArticles(articles)

	req := httptest.NewRequest("GET", "/articles?tag=tag1&tag=tag2", nil)
	req = middleware.SetReqID(req)

	rr := httptest.NewRecorder()
	handler.GetArticlesByTags().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response article.Articles
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Len(t, response, 1)
	assert.Equal(t, "Article 1", response[0].Title)
}

func TestUpdateArticleByTitle(t *testing.T) {
	handler, mockRepo := setupTest()

	originalArticle := article.Article{ID: 1, Title: "Original Title", Slug: "original-title", Content: "Original content"}
	mockRepo.SetArticles([]article.Article{originalArticle})

	updatedData := []byte(`{
		"article": "title:Updated Title\npublicationDate:2023-05-16\ntags:updated,article\n===\nThis is updated content."
	}`)

	req := httptest.NewRequest("PUT", "/articles/original-title", bytes.NewBuffer(updatedData))
	req = middleware.SetReqID(req)
	pathParam := "title"
	req.SetPathValue(pathParam, "original-title")
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.UpdateArticleByTitle(pathParam).ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response article.Article
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), response.ID)
	assert.Equal(t, "Updated Title", response.Title)
	assert.Equal(t, "This is updated content.", response.Content)
}

func TestDeleteArticleByTitle(t *testing.T) {
	handler, mockRepo := setupTest()

	articleToDelete := article.Article{ID: 1, Title: "Article to Delete", Slug: "article-to-delete"}
	mockRepo.SetArticles([]article.Article{articleToDelete})

	req := httptest.NewRequest("DELETE", "/articles/article-to-delete", nil)
	req = middleware.SetReqID(req)
	pathParam := "title"
	req.SetPathValue(pathParam, "article-to-delete")

	rr := httptest.NewRecorder()
	handler.DeleteArticleByTitle(pathParam).ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)

	// Verify the article was deleted
	articles, err := mockRepo.GetAll(req.Context())
	assert.NoError(t, err)
	assert.Len(t, articles, 0)
}

func TestGetAllTags(t *testing.T) {
	handler, mockRepo := setupTest()

	articles := []article.Article{
		{ID: 1, Title: "Article 1", Tags: []string{"tag1", "tag2"}},
		{ID: 2, Title: "Article 2", Tags: []string{"tag2", "tag3"}},
		{ID: 3, Title: "Article 3", Tags: []string{"tag1", "tag3"}},
	}
	mockRepo.SetArticles(articles)

	req := httptest.NewRequest("GET", "/tags", nil)
	req = middleware.SetReqID(req)

	rr := httptest.NewRecorder()
	handler.GetAllTags().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response []string
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Len(t, response, 3)
	assert.ElementsMatch(t, []string{"tag1", "tag2", "tag3"}, response)
}
