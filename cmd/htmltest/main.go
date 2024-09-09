package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jannawro/blog/article"
	"github.com/jannawro/blog/views/components"
	"github.com/jannawro/blog/views/pages"
)

func main() {
	// Create a sample article
	sampleArticle := article.Article{
		ID:              1,
		Title:           "Sample Article",
		Content:         "This is a sample article content. It's long enough to be truncated in the card view.",
		Tags:            []string{"sample", "test"},
		PublicationDate: time.Now(),
	}

	// Render the ArticleCard component
	indexPage := pages.Index()
	component := components.ArticleCard(sampleArticle)

	// Create a context
	ctx := context.Background()

	// Render the component to HTML
	var html bytes.Buffer
	err := component.Render(ctx, &html)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error rendering component: %v\n", err)
		os.Exit(1)
	}

	// Print the rendered HTML
	fmt.Println(html.String())
}
