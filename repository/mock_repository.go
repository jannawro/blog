package repository

import (
	"context"
	"errors"
	"sync"

	"github.com/jannawro/blog/article"
)

type MockRepository struct {
	articles map[int64]article.Article
	mutex    sync.RWMutex
	nextID   int64
}

func NewMockRepository() *MockRepository {
	return &MockRepository{
		articles: make(map[int64]article.Article),
		nextID:   1,
	}
}

func (m *MockRepository) Create(ctx context.Context, article article.Article) (*article.Article, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	article.ID = m.nextID
	m.articles[article.ID] = article
	m.nextID++

	return &article, nil
}

func (m *MockRepository) GetAll(ctx context.Context) (article.Articles, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	result := make(article.Articles, 0, len(m.articles))
	for _, article := range m.articles {
		result = append(result, article)
	}
	return result, nil
}

func (m *MockRepository) GetByID(ctx context.Context, id int64) (*article.Article, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if article, ok := m.articles[id]; ok {
		return &article, nil
	}
	return nil, errors.New("article not found")
}

func (m *MockRepository) GetByTitle(ctx context.Context, title string) (*article.Article, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	for _, article := range m.articles {
		if article.Title == title {
			return &article, nil
		}
	}
	return nil, errors.New("article not found")
}

func (m *MockRepository) GetByTags(ctx context.Context, tags []string) (article.Articles, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	result := make(article.Articles, 0)
	for _, article := range m.articles {
		if containsAllTags(article.Tags, tags) {
			result = append(result, article)
		}
	}
	return result, nil
}

func (m *MockRepository) Update(ctx context.Context, id int64, updated article.Article) (*article.Article, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, ok := m.articles[id]; !ok {
		return nil, errors.New("article not found")
	}

	updated.ID = id
	m.articles[id] = updated
	return &updated, nil
}

func (m *MockRepository) Delete(ctx context.Context, id int64) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, ok := m.articles[id]; !ok {
		return errors.New("article not found")
	}

	delete(m.articles, id)
	return nil
}

func (m *MockRepository) SetArticles(setArticles []article.Article) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.articles = make(map[int64]article.Article)
	for _, article := range setArticles {
		m.articles[article.ID] = article
		if article.ID >= m.nextID {
			m.nextID = article.ID + 1
		}
	}
}

func (m *MockRepository) Reset() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.articles = make(map[int64]article.Article)
	m.nextID = 1
}

func (m *MockRepository) GetAllTags(ctx context.Context) ([]string, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	tagSet := make(map[string]struct{})
	for _, article := range m.articles {
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
