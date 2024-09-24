package mock

import (
	"context"
	"errors"
	"sync"

	"github.com/jannawro/blog/article"
)

type Repository struct {
	articles map[int64]article.Article
	mutex    sync.RWMutex
	nextID   int64
}

func NewRepository() *Repository {
	return &Repository{
		articles: make(map[int64]article.Article),
		nextID:   1,
	}
}

func (r *Repository) Create(ctx context.Context, article article.Article) (*article.Article, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	article.ID = r.nextID
	r.articles[article.ID] = article
	r.nextID++

	return &article, nil
}

func (r *Repository) GetAll(ctx context.Context) (article.Articles, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	result := make(article.Articles, 0, len(r.articles))
	for _, article := range r.articles {
		result = append(result, article)
	}
	return result, nil
}

func (r *Repository) GetByID(ctx context.Context, id int64) (*article.Article, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if article, ok := r.articles[id]; ok {
		return &article, nil
	}
	return nil, errors.New("article not found")
}

func (r *Repository) GetBySlug(ctx context.Context, slug string) (*article.Article, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, article := range r.articles {
		if article.Slug == slug {
			return &article, nil
		}
	}
	return nil, errors.New("article not found")
}

func (r *Repository) GetByTags(ctx context.Context, tags []string) (article.Articles, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	result := make(article.Articles, 0)
	for _, article := range r.articles {
		if containsAllTags(article.Tags, tags) {
			result = append(result, article)
		}
	}
	return result, nil
}

func (r *Repository) Update(ctx context.Context, id int64, updated article.Article) (*article.Article, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, ok := r.articles[id]; !ok {
		return nil, errors.New("article not found")
	}

	updated.ID = id
	r.articles[id] = updated
	return &updated, nil
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, ok := r.articles[id]; !ok {
		return errors.New("article not found")
	}

	delete(r.articles, id)
	return nil
}

func (r *Repository) SetArticles(setArticles []article.Article) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.articles = make(map[int64]article.Article)
	for _, article := range setArticles {
		r.articles[article.ID] = article
		if article.ID >= r.nextID {
			r.nextID = article.ID + 1
		}
	}
}

func (r *Repository) Reset() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.articles = make(map[int64]article.Article)
	r.nextID = 1
}

func (r *Repository) GetAllTags(ctx context.Context) ([]string, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	tagSet := make(map[string]struct{})
	for _, article := range r.articles {
		for _, tag := range article.Tags {
			tagSet[tag] = struct{}{}
		}
	}

	tags := make([]string, 0, len(tagSet))
	for tag := range tagSet {
		tags = append(tags, tag)
	}
	return tags, nil
}

func containsAllTags(articleTags, searchTags []string) bool {
	tagSet := make(map[string]struct{})
	for _, tag := range articleTags {
		tagSet[tag] = struct{}{}
	}

	for _, tag := range searchTags {
		if _, ok := tagSet[tag]; !ok {
			return false
		}
	}
	return true
}
