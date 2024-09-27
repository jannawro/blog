package article

import (
	"errors"
	"fmt"
)

var (
	ErrSeparatorNotFound         = fmt.Errorf("headers and body separator '%s' not found", separator)
	ErrDateFormatFailed          = errors.New("date formatting failed")
	ErrArticleUnmarshalingFailed = errors.New("article unmarshaling failed")
	ErrArticleNotFound           = errors.New("article not found")
	ErrArticlesNotFound          = errors.New("articles not found")
	ErrArticleCreationFailed     = errors.New("article creation failed")
	ErrArticleUpdateFailed       = errors.New("updating article failed")
	ErrArticleDeletionFailed     = errors.New("deleting article failed")
)
