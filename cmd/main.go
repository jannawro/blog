package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jannawro/blog/article"
	"github.com/jannawro/blog/views/components"
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
	component := components.ArticleCard(sampleArticle)

	// Create a context
	ctx := context.Background()

	// Render the component to HTML
	html, err := component.Render(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error rendering component: %v\n", err)
		os.Exit(1)
	}

	// Print the rendered HTML
	fmt.Println(html)
}
