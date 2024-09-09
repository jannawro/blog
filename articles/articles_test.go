package articles_test

import (
	"testing"
	"time"

	"github.com/jannawro/blog/articles"
	"github.com/stretchr/testify/assert"
)

var article = []byte(`title:Fondant recipe
publicationDate:2005-04-02
tags:cooking,sweets
===
# Markdown Title
Markdown contents...`)

func TestUnmarshalToArticle(t *testing.T) {
	a := articles.Article{}
	err := articles.UnmarshalToArticle(article, &a)
	if err != nil {
		t.Fatal("expected no error but got:", err)
	}

	assert := assert.New(t)

	assert.Equal("fondant-recipe", a.Title)

	assert.Equal([]string{"cooking", "sweets"}, a.Tags)

	date, err := time.Parse("2006-01-02", "2005-04-02")
	if err != nil {
		t.Fatal("expected no error but got:", err)
	}

	assert.Equal(date, a.PublicationDate)

	assert.Equal(`# Markdown Title
Markdown contents...`, a.Content)
}

func TestArticlesSort(t *testing.T) {
	testArticles := articles.Articles{
		{ID: 3, Title: "C Article", PublicationDate: time.Date(2023, 1, 3, 0, 0, 0, 0, time.UTC)},
		{ID: 1, Title: "A Article", PublicationDate: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)},
		{ID: 2, Title: "B Article", PublicationDate: time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC)},
	}

	t.Run("Sort by Title", func(t *testing.T) {
		sorted := make(articles.Articles, len(testArticles))
		copy(sorted, testArticles)
		sorted.Sort(articles.SortByTitle)
		
		assert.Equal(t, "A Article", sorted[0].Title)
		assert.Equal(t, "B Article", sorted[1].Title)
		assert.Equal(t, "C Article", sorted[2].Title)
	})

	t.Run("Sort by PublicationDate", func(t *testing.T) {
		sorted := make(articles.Articles, len(testArticles))
		copy(sorted, testArticles)
		sorted.Sort(articles.SortByPublicationDate)
		
		assert.Equal(t, time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), sorted[0].PublicationDate)
		assert.Equal(t, time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC), sorted[1].PublicationDate)
		assert.Equal(t, time.Date(2023, 1, 3, 0, 0, 0, 0, time.UTC), sorted[2].PublicationDate)
	})

	t.Run("Sort by ID", func(t *testing.T) {
		sorted := make(articles.Articles, len(testArticles))
		copy(sorted, testArticles)
		sorted.Sort(articles.SortByID)
		
		assert.Equal(t, int64(1), sorted[0].ID)
		assert.Equal(t, int64(2), sorted[1].ID)
		assert.Equal(t, int64(3), sorted[2].ID)
	})
}
