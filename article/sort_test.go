package article_test

import (
	"testing"
	"time"

	article1 "github.com/jannawro/blog/article"
	"github.com/stretchr/testify/assert"
)

func TestArticlesSort(t *testing.T) {
	testArticles := article1.Articles{
		{ID: 3, Title: "C Article", PublicationDate: time.Date(2023, 1, 3, 0, 0, 0, 0, time.UTC)},
		{ID: 1, Title: "A Article", PublicationDate: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)},
		{ID: 2, Title: "B Article", PublicationDate: time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC)},
	}

	t.Run("Sort by Title", func(t *testing.T) {
		sorted := make(article1.Articles, len(testArticles))
		copy(sorted, testArticles)
		sorted.Sort(article1.SortByTitle)

		assert.Equal(t, "A Article", sorted[0].Title)
		assert.Equal(t, "B Article", sorted[1].Title)
		assert.Equal(t, "C Article", sorted[2].Title)
	})

	t.Run("Sort by PublicationDate", func(t *testing.T) {
		sorted := make(article1.Articles, len(testArticles))
		copy(sorted, testArticles)
		sorted.Sort(article1.SortByPublicationDate)

		assert.Equal(t, time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), sorted[0].PublicationDate)
		assert.Equal(t, time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC), sorted[1].PublicationDate)
		assert.Equal(t, time.Date(2023, 1, 3, 0, 0, 0, 0, time.UTC), sorted[2].PublicationDate)
	})

	t.Run("Sort by ID", func(t *testing.T) {
		sorted := make(article1.Articles, len(testArticles))
		copy(sorted, testArticles)
		sorted.Sort(article1.SortByID)

		assert.Equal(t, int64(1), sorted[0].ID)
		assert.Equal(t, int64(2), sorted[1].ID)
		assert.Equal(t, int64(3), sorted[2].ID)
	})
}
