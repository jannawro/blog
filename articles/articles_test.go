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
