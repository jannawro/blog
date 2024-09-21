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
	ErrArticleCreationFailed     = errors.New("article creation failed")
	ErrFetchAllFailed            = errors.New("fetching all articles failed")
	ErrFetchBySlugFailed         = errors.New("fetching article by slug failed")
	ErrFetchByTagsFailed         = errors.New("fetching articles by tags failed")
	ErrUpdateBySlugFailed        = errors.New("updating article by slug failed")
	ErrDeleteBySlugFailed        = errors.New("deleting article by slug failed")
	ErrFetchAllTagsFailed        = errors.New("fetching all tags failed")
)
