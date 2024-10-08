package article_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	a "github.com/jannawro/blog/article"
	"github.com/jannawro/blog/repository/mock"
)

func setupTestService() (*a.Service, *mock.Repository) {
	mockRepo := mock.NewRepository()
	articleService := a.NewService(mockRepo)
	return articleService, mockRepo
}

func TestCreate(t *testing.T) {
	service, _ := setupTestService()
	ctx := context.Background()

	testArticle := a.Article{
		ID:              1,
		Title:           "Article 1",
		Thumbnail:       "Article 1 thumbnail",
		Slug:            "article-1",
		Tags:            []string{"tag1", "tag2"},
		Content:         "This is the content of the test article.",
		PublicationDate: time.Date(2023, 5, 10, 0, 0, 0, 0, time.UTC),
	}

	response, err := service.Create(ctx, testArticle)
	assert.NoError(t, err)

	require.NotNil(t, response)
	assert.NotEmpty(t, response.ID)
	assert.Equal(t, "Article 1", response.Title)
	assert.Equal(t, "Article 1 thumbnail", response.Thumbnail)
	assert.Equal(t, "article-1", response.Slug)
	assert.Equal(t, []string{"tag1", "tag2"}, response.Tags)
	assert.Equal(t, "This is the content of the test article.", response.Content)
	expectedDate, _ := time.Parse("2006-01-02", "2023-05-10")
	assert.Equal(t, expectedDate, response.PublicationDate)
}

func TestGetAllTags(t *testing.T) {
	service, mockRepo := setupTestService()
	ctx := context.Background()

	testArticles := []a.Article{
		{ID: 1, Title: "Article 1", Slug: "article-1", Tags: []string{"tag1", "tag2"}},
		{ID: 2, Title: "Article 2", Slug: "article-2", Tags: []string{"tag2", "tag3"}},
		{ID: 3, Title: "Article 3", Slug: "article-3", Tags: []string{"tag1", "tag3", "tag4"}},
	}
	mockRepo.SetArticles(testArticles)

	tags, err := service.GetAllTags(ctx)

	assert.NoError(t, err)
	assert.ElementsMatch(t, []string{"tag1", "tag2", "tag3", "tag4"}, tags)
}

func TestGetAll(t *testing.T) {
	service, mockRepo := setupTestService()
	ctx := context.Background()

	testArticles := []a.Article{
		{ID: 1, Title: "Article 1", Slug: "article-1", Content: "Content 1", Tags: []string{"tag1", "tag2"}},
		{ID: 2, Title: "Article 2", Slug: "article-2", Content: "Content 2", Tags: []string{"tag2", "tag3"}},
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

	testArticle := a.Article{
		ID:              1,
		Slug:            "test-article",
		Content:         "This is a test article",
		Tags:            []string{"test", "article"},
		PublicationDate: time.Date(2023, 5, 10, 0, 0, 0, 0, time.UTC),
	}
	mockRepo.SetArticles([]a.Article{testArticle})

	tests := []struct {
		name          string
		slug          string
		expectedFound bool
	}{
		{"Existing article", "test-article", true},
		{"Non-existing article", "missing-article", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			article, err := service.GetBySlug(ctx, tt.slug)

			if tt.expectedFound {
				assert.NoError(t, err)
				require.NotNil(t, article)
				assert.Equal(t, tt.slug, article.Slug)
			} else {
				assert.Error(t, err)
				assert.Nil(t, article)
			}
		})
	}
}

func TestGetByID(t *testing.T) {
	service, mockRepo := setupTestService()
	ctx := context.Background()

	testArticle := a.Article{
		ID:              1,
		Slug:            "test-article",
		Content:         "This is a test article",
		Tags:            []string{"test", "article"},
		PublicationDate: time.Date(2023, 5, 10, 0, 0, 0, 0, time.UTC),
	}
	mockRepo.SetArticles([]a.Article{testArticle})

	tests := []struct {
		name          string
		id            int64
		expectedFound bool
	}{
		{"Existing article", 1, true},
		{"Non-existing article", 2, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			article, err := service.GetByID(ctx, tt.id)

			if tt.expectedFound {
				assert.NoError(t, err)
				require.NotNil(t, article)
				assert.Equal(t, tt.id, article.ID)
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

	testArticles := []a.Article{
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

	initialArticle := a.Article{
		ID:      1,
		Slug:    "initial-article",
		Content: "Initial content",
		Tags:    []string{"initial", "tag"},
	}
	mockRepo.SetArticles([]a.Article{initialArticle})

	tests := []struct {
		name        string
		initialSlug string
		updated     a.Article
		expectedErr bool
	}{
		{
			name:        "Update existing article",
			initialSlug: "initial-article",
			updated: a.Article{
				Title:           "Updated Article",
				Slug:            "updated-article",
				Content:         "Updated content",
				Tags:            []string{"updated", "tag"},
				PublicationDate: time.Date(2023, time.May, 10, 0, 0, 0, 0, time.UTC),
			},
			expectedErr: false,
		},
		{
			name:        "Update non-existing article",
			initialSlug: "non-existing-article",
			updated: a.Article{
				Title:           "New Article",
				Slug:            "new-article",
				Content:         "New content",
				Tags:            []string{"new", "tag"},
				PublicationDate: time.Date(2023, time.May, 10, 0, 0, 0, 0, time.UTC),
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updatedArticle, err := service.UpdateBySlug(ctx, tt.initialSlug, tt.updated)

			if tt.expectedErr {
				assert.Error(t, err)
				assert.Nil(t, updatedArticle)
			} else {
				assert.NoError(t, err)
				require.NotNil(t, updatedArticle)
				assert.Equal(t, "Updated Article", updatedArticle.Title)
				assert.Equal(t, "updated-article", updatedArticle.Slug)
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

	initialArticle := a.Article{
		ID:   1,
		Slug: "article-to-delete",
	}
	mockRepo.SetArticles([]a.Article{initialArticle})

	tests := []struct {
		name        string
		slug        string
		expectedErr bool
	}{
		{"Delete existing article", "article-to-delete", false},
		{"Delete non-existing article", "Non-existing Article", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.DeleteBySlug(ctx, tt.slug)

			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				// Verify the article is actually deleted
				_, err := service.GetBySlug(ctx, tt.slug)
				assert.Error(t, err)
			}
		})
	}
}
