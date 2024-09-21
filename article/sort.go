package article

import (
	"net/http"
	"sort"
	"strings"
)

// SortOption defines the sorting criteria for Articles
type SortOption string

const (
	SortByTitle           SortOption = "title"
	SortByPublicationDate SortOption = "publication_date"
	SortByID              SortOption = "id"
)

// Sort sorts the Articles slice based on the given SortOption
func (a Articles) Sort(option SortOption) {
	sort.Slice(a, func(i, j int) bool {
		switch option {
		case SortByTitle:
			return strings.ToLower(a[i].Title) < strings.ToLower(a[j].Title)
		case SortByPublicationDate:
			return a[i].PublicationDate.Before(a[j].PublicationDate)
		case SortByID:
			return a[i].ID < a[j].ID
		default:
			return false
		}
	})
}

func GetSortOption(r *http.Request) SortOption {
	sortParam := r.URL.Query().Get("sort")
	switch sortParam {
	case "title":
		return SortByTitle
	case "id":
		return SortByID
	case "date":
		fallthrough
	default:
		return SortByPublicationDate
	}
}
