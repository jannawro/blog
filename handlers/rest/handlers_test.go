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

	articleData := []byte(`{"title":"Test Article","content":"This is a test article.","tags":["test"]}`)
	expectedArticle := &article.Article{
		Title:   "Test Article",
		Content: "This is a test article.",
		Tags:    []string{"test"},
	}

	req, err := http.NewRequest("POST", "/articles", bytes.NewBuffer(articleData))
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
	assert.Equal(t, expectedArticle.Content, response.Content)
	assert.Equal(t, expectedArticle.Tags, response.Tags)

	// Verify the article was added to the mock repository
	articles, _ := mockRepo.GetAll(req.Context())
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

	expectedArticle := &article.Article{ID: 1, Title: "Test Article", Slug: "test-article"}
	mockRepo.SetArticles([]article.Article{*expectedArticle})

	req, err := http.NewRequest("GET", "/articles/test-article", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.GetArticleByTitle().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response article.Article
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedArticle.Title, response.Title)
	assert.Equal(t, expectedArticle.Slug, response.Slug)
}

// Add more tests for other handler methods...
