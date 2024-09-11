package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jannawro/blog/article"
	"github.com/jannawro/blog/components"
)

func main() {
	// Create a sample article
	article1 := article.Article{
		ID:              1,
		Title:           "Sample Article 1",
		Content:         "This is a sample article **content**. It's long enough to be truncated in the card view.",
		Tags:            []string{"sample", "test", "article1"},
		PublicationDate: time.Now(),
	}
	article2 := article.Article{
		ID:              2,
		Title:           "Sample Article 2",
		Content:         "This is a sample article **content**. It's long enough to be truncated in the card view.",
		Tags:            []string{"sample", "test", "article2"},
		PublicationDate: time.Now(),
	}
	article3 := article.Article{
		ID:              3,
		Title:           "Sample Article 3",
		Content:         "This is a sample article **content**. It's long enough to be truncated in the card view.",
		Tags:            []string{"sample", "test", "article3"},
		PublicationDate: time.Now(),
	}
	sampleArticles := []article.Article{
		{
			ID:              1,
			Title:           "Sample Article",
			Content:         "This is a sample article content. It's long enough to be truncated in the card view.",
			Tags:            []string{"sample", "test"},
			PublicationDate: time.Now(),
		},
		{
			ID:              2,
			Title:           "Article 2",
			Content:         "**Lorem** ipsum dolor sit amet, consectetur adipiscing elit. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
			Tags:            []string{"lorem", "ipsum"},
			PublicationDate: time.Now(),
		},
		{
			ID:              3,
			Title:           "Article 3",
			Content:         "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
			Tags:            []string{"lorem", "ipsum"},
			PublicationDate: time.Now(),
		},
		{
			ID:              4,
			Title:           "Article 4",
			Content:         "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
			Tags:            []string{"lorem", "ipsum"},
			PublicationDate: time.Now(),
		},
		{
			ID:              5,
			Title:           "Article 5",
			Content:         "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
			Tags:            []string{"lorem", "ipsum"},
			PublicationDate: time.Now(),
		},
	}
	taggedArticles := map[string][]article.Article{
		"sample":   {article1, article2, article3},
		"test":     {article1, article2, article3},
		"article1": {article1},
		"article2": {article2},
		"article3": {article3},
	}

	_ = components.Blog(sampleArticles)
	_ = components.ArticleCard(article1)
	_ = components.ArticlePage(article1)
	indexPage := components.TagIndexPage(taggedArticles)

	// Create a context
	ctx := context.Background()

	// Render the component to HTML
	var html bytes.Buffer
	err := indexPage.Render(ctx, &html)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error rendering component: %v\n", err)
		os.Exit(1)
	}

	// Print the rendered HTML
	fmt.Println(html.String())
}
