package article

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"
)

const (
	separator             = "==="
	publicationDateFormat = "2006-01-02"
)

var (
	ErrSeparatorNotFound = errors.New("separator between headers and body not found")
	ErrInvalidDateFormat = errors.New("unable to parse date: invalid format")
)

// Article is a representation of a markdown file with specific headers. Example of such a file:
// `title:Fondant recipe
// publicationDate:2005-04-02
// tags:cooking,sweets
// ===
// # Markdown Title
// Markdown contents...`
type Article struct {
	ID              int64
	Title           string
	Content         string
	Tags            []string
	PublicationDate time.Time
}

type Articles []Article

type ArticleRepository interface {
	Create(ctx context.Context, article Article) (*Article, error)
	GetAll(ctx context.Context) (Articles, error)
	GetByID(ctx context.Context, id int64) (*Article, error)
	GetByTitle(ctx context.Context, title string) (*Article, error)
	GetByTags(ctx context.Context, tags []string) (Articles, error)
	Update(ctx context.Context, id int64, updated Article) (*Article, error)
	Delete(ctx context.Context, id int64) error
}

type ArticleHandler http.Handler

// UnmarshalToArticle parses a markdown file with specific headers and stores the result as an article in a
func UnmarshalToArticle(data []byte, a *Article) error {
	headersSection, bodySection, found := strings.Cut(string(data), separator)
	if !found {
		return ErrSeparatorNotFound
	}

	headers := make(map[string]string)
	headerLines := strings.Split(headersSection, "\n")
	for _, line := range headerLines {
		if strings.Contains(line, ":") {
			kv := strings.SplitN(line, ":", 2)
			key := strings.TrimSpace(kv[0])
			value := strings.TrimSpace(kv[1])
			headers[key] = value
		}
	}

	a.Title = strings.ReplaceAll(strings.ToLower(headers["title"]), " ", "-")
	date, err := time.Parse(publicationDateFormat, headers["publicationDate"])
	if err != nil {
		return ErrInvalidDateFormat
	}
	a.PublicationDate = date
	a.Tags = strings.Split(headers["tags"], ",")
	a.Content = strings.TrimSpace(bodySection)

	return nil
}
