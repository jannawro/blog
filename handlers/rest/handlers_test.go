package rest_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jannawro/blog/article"
	"github.com/jannawro/blog/handlers/rest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockService is a mock implementation of the article.Service
type MockService struct {
	mock.Mock
}

func (m *MockService) Create(ctx context.Context, articleData []byte) (*article.Article, error) {
	args := m.Called(ctx, articleData)
	return args.Get(0).(*article.Article), args.Error(1)
}

func (m *MockService) GetAll(ctx context.Context, sortOption *article.SortOption) (article.Articles, error) {
	args := m.Called(ctx, sortOption)
	return args.Get(0).(article.Articles), args.Error(1)
}

func (m *MockService) GetBySlug(ctx context.Context, slug string) (*article.Article, error) {
	args := m.Called(ctx, slug)
	return args.Get(0).(*article.Article), args.Error(1)
}

func (m *MockService) GetByID(ctx context.Context, id int64) (*article.Article, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*article.Article), args.Error(1)
}

func (m *MockService) GetByTags(ctx context.Context, tags []string, sortOption *article.SortOption) (article.Articles, error) {
	args := m.Called(ctx, tags, sortOption)
	return args.Get(0).(article.Articles), args.Error(1)
}

func (m *MockService) UpdateBySlug(ctx context.Context, slug string, updatedData []byte) (*article.Article, error) {
	args := m.Called(ctx, slug, updatedData)
	return args.Get(0).(*article.Article), args.Error(1)
}

func (m *MockService) DeleteBySlug(ctx context.Context, slug string) error {
	args := m.Called(ctx, slug)
	return args.Error(0)
}

func (m *MockService) GetAllTags(ctx context.Context) ([]string, error) {
	args := m.Called(ctx)
	return args.Get(0).([]string), args.Error(1)
}

func TestCreateArticle(t *testing.T) {
	mockService := new(MockService)
	handler := rest.NewHandler(mockService)

	articleData := []byte(`{"title":"Test Article","content":"This is a test article."}`)
	expectedArticle := &article.Article{ID: 1, Title: "Test Article", Content: "This is a test article."}

	mockService.On("Create", mock.Anything, articleData).Return(expectedArticle, nil)

	req, err := http.NewRequest("POST", "/articles", bytes.NewBuffer(articleData))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.CreateArticle().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response article.Article
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedArticle, &response)

	mockService.AssertExpectations(t)
}

func TestGetAllArticles(t *testing.T) {
	mockService := new(MockService)
	handler := rest.NewHandler(mockService)

	expectedArticles := article.Articles{
		{ID: 1, Title: "Article 1"},
		{ID: 2, Title: "Article 2"},
	}

	mockService.On("GetAll", mock.Anything, (*article.SortOption)(nil)).Return(expectedArticles, nil)

	req, err := http.NewRequest("GET", "/articles", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.GetAllArticles().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response article.Articles
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedArticles, response)

	mockService.AssertExpectations(t)
}

func TestGetArticleByTitle(t *testing.T) {
	mockService := new(MockService)
	handler := rest.NewHandler(mockService)

	expectedArticle := &article.Article{ID: 1, Title: "Test Article", Slug: "test-article"}

	mockService.On("GetBySlug", mock.Anything, "test-article").Return(expectedArticle, nil)

	req, err := http.NewRequest("GET", "/articles/test-article", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.GetArticleByTitle().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response article.Article
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedArticle, &response)

	mockService.AssertExpectations(t)
}

// Add more tests for other handler methods...
