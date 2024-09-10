package main

import (
	"bytes"
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
		Content:         "This is a sample article **content**. It's long enough to be truncated in the card view.",
		Tags:            []string{"sample", "test"},
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
	// Render the ArticleCard component
	homePage := components.Home(sampleArticles)
	_ = components.ArticleCard(sampleArticle)
	_ = components.ArticleContentPage(sampleArticle)

	// Create a context
	ctx := context.Background()

	// Render the component to HTML
	var html bytes.Buffer
	err := homePage.Render(ctx, &html)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error rendering component: %v\n", err)
		os.Exit(1)
	}

	// Print the rendered HTML
	fmt.Println(html.String())
}
