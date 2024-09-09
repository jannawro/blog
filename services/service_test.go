package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jannawro/blog/articles"
	"github.com/jannawro/blog/repository"
	"github.com/jannawro/blog/services/service"
)

func setupTestService() (*service.ArticleService, *repository.MockRepository) {
	mockRepo := repository.NewMockRepository()
	articleService := service.NewArticleService(mockRepo)
	return service, mockRepo
}

func TestCreate(t *testing.T) {
	service, _ := setupTestService()
	ctx := context.Background()

	tests := []struct {
		name        string
		input       []byte
		expectedErr bool
	}{
		{
			name: "Valid article",
			input: []byte(`title:Test Article
publicationDate:2023-05-10
tags:test,article
===
This is the content of the test article.`),
			expectedErr: false,
		},
		{
			name:        "Invalid article (missing separator)",
			input:       []byte("Title: Invalid Article\nTags: test\nThis is invalid content without separator."),
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			article, err := service.Create(ctx, tt.input)

			if tt.expectedErr {
				assert.Error(t, err)
				assert.Nil(t, article)
			} else {
				assert.NoError(t, err)
				require.NotNil(t, article)
				assert.NotEmpty(t, article.ID)
				assert.Equal(t, "test-article", article.Title)
				assert.Equal(t, []string{"test", "article"}, article.Tags)
				expectedDate, _ := time.Parse("2006-01-02", "2023-05-10")
				assert.Equal(t, expectedDate, article.PublicationDate)
				assert.Equal(t, "This is the content of the test article.", article.Content)
			}
		})
	}
}

func TestGetAll(t *testing.T) {
	service, mockRepo := setupTestService()
	ctx := context.Background()

	testArticles := []articles.Article{
		{ID: 1, Title: "Article 1", Content: "Content 1", Tags: []string{"tag1", "tag2"}},
		{ID: 2, Title: "Article 2", Content: "Content 2", Tags: []string{"tag2", "tag3"}},
	}
	mockRepo.SetArticles(testArticles)

	result, err := service.GetAll(ctx, nil)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Contains(t, result, testArticles[0])
	assert.Contains(t, result, testArticles[1])
}

func TestGetByTitle(t *testing.T) {
	service, mockRepo := setupTestService()
	ctx := context.Background()

	testArticle := articles.Article{
		ID:              1,
		Title:           "test-article",
		Content:         "This is a test article",
		Tags:            []string{"test", "article"},
		PublicationDate: time.Date(2023, 5, 10, 0, 0, 0, 0, time.UTC),
	}
	mockRepo.SetArticles([]articles.Article{testArticle})

	tests := []struct {
		name          string
		title         string
		expectedFound bool
	}{
		{"Existing article", "test-article", true},
		{"Non-existing article", "missing-article", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			article, err := service.GetByTitle(ctx, tt.title)

			if tt.expectedFound {
				assert.NoError(t, err)
				require.NotNil(t, article)
				assert.Equal(t, tt.title, article.Title)
			} else {
				assert.Error(t, err)
				assert.Nil(t, article)
			}
		})
	}
}

func TestGetByTags(t *testing.T) {
	service, mockRepo := setupTestService()
	ctx := context.Background()

	testArticles := []articles.Article{
		{ID: 1, Title: "Article 1", Content: "Content 1", Tags: []string{"tag1", "tag2"}},
		{ID: 2, Title: "Article 2", Content: "Content 2", Tags: []string{"tag2", "tag3"}},
		{ID: 3, Title: "Article 3", Content: "Content 3", Tags: []string{"tag3", "tag4"}},
	}
	mockRepo.SetArticles(testArticles)

	tests := []struct {
		name          string
		tags          []string
		expectedCount int
		expectedIDs   []int64
	}{
		{"Single tag", []string{"tag1"}, 1, []int64{1}},
		{"Multiple tags", []string{"tag2"}, 2, []int64{1, 2}},
		{"No matching tags", []string{"tag5"}, 0, []int64{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			arts, err := service.GetByTags(ctx, tt.tags, nil)

			assert.NoError(t, err)
			assert.Len(t, arts, tt.expectedCount)

			actualIDs := make([]int64, len(arts))
			for i, art := range arts {
				actualIDs[i] = art.ID
			}
			assert.ElementsMatch(t, tt.expectedIDs, actualIDs)
		})
	}
}

func TestUpdateByTitle(t *testing.T) {
	service, mockRepo := setupTestService()
	ctx := context.Background()

	initialArticle := articles.Article{
		ID:      1,
		Title:   "Initial Article",
		Content: "Initial content",
		Tags:    []string{"initial", "tag"},
	}
	mockRepo.SetArticles([]articles.Article{initialArticle})

	tests := []struct {
		name        string
		title       string
		updatedData []byte
		expectedErr bool
	}{
		{
			name:  "Update existing article",
			title: "Initial Article",
			updatedData: []byte(`title:Updated Article
publicationDate:2023-05-10
tags:updated,tag
===
Updated content`),
			expectedErr: false,
		},
		{
			name:        "Update non-existing article",
			title:       "Non-existing Article",
			updatedData: []byte("Title: New Article\nTags: new, tag\nPublication Date: 2023-05-10\n\nNew content"),
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updatedArticle, err := service.UpdateByTitle(ctx, tt.title, tt.updatedData)

			if tt.expectedErr {
				assert.Error(t, err)
				assert.Nil(t, updatedArticle)
			} else {
				assert.NoError(t, err)
				require.NotNil(t, updatedArticle)
				assert.Equal(t, "updated-article", updatedArticle.Title)
				assert.Equal(t, []string{"updated", "tag"}, updatedArticle.Tags)
				expectedDate, _ := time.Parse("2006-01-02", "2023-05-10")
				assert.Equal(t, expectedDate, updatedArticle.PublicationDate)
				assert.Equal(t, "Updated content", updatedArticle.Content)
			}
		})
	}
}

func TestDeleteByTitle(t *testing.T) {
	service, mockRepo := setupTestService()
	ctx := context.Background()

	initialArticle := articles.Article{
		ID:    1,
		Title: "Article to Delete",
	}
	mockRepo.SetArticles([]articles.Article{initialArticle})

	tests := []struct {
		name        string
		title       string
		expectedErr bool
	}{
		{"Delete existing article", "Article to Delete", false},
		{"Delete non-existing article", "Non-existing Article", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.DeleteByTitle(ctx, tt.title)

			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				// Verify the article is actually deleted
				_, err := service.GetByTitle(ctx, tt.title)
				assert.Error(t, err)
			}
		})
	}
}
