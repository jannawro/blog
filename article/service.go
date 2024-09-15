package article

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var ErrArticleNotFound = errors.New("article not found")

type Service struct {
	repo ArticleRepository
}

func NewService(repo ArticleRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, article Article) (*Article, error) {
	return s.repo.Create(ctx, article)
}

func (s *Service) GetAll(ctx context.Context, sortBy *SortOption) (Articles, error) {
	articles, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	if sortBy != nil {
		articles.Sort(*sortBy)
	}
	return articles, nil
}

func (s *Service) GetBySlug(ctx context.Context, slug string) (*Article, error) {
	article, err := s.repo.GetBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%w: %v", ErrArticleNotFound, err)
		}
		return nil, err
	}
	return article, nil
}

func (s *Service) GetByID(ctx context.Context, id int64) (*Article, error) {
	article, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%w: %v", ErrArticleNotFound, err)
		}
		return nil, err
	}
	return article, nil
}

func (s *Service) GetByTags(
	ctx context.Context,
	tags []string,
	sortBy *SortOption,
) (Articles, error) {
	articles, err := s.repo.GetByTags(ctx, tags)
	if err != nil {
		return nil, err
	}
	if sortBy != nil {
		articles.Sort(*sortBy)
	}
	return articles, nil
}

func (s *Service) UpdateBySlug(
	ctx context.Context,
	slug string,
	updatedData []byte,
) (*Article, error) {
	existingArticle, err := s.repo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	var updatedArticle Article
	if err := UnmarshalToArticle(updatedData, &updatedArticle); err != nil {
		return nil, err
	}

	updatedArticle.ID = existingArticle.ID
	return s.repo.Update(ctx, existingArticle.ID, updatedArticle)
}

func (s *Service) DeleteBySlug(ctx context.Context, slug string) error {
	article, err := s.repo.GetBySlug(ctx, slug)
	if err != nil {
		return err
	}
	if article == nil {
		return errors.New("article not found")
	}
	return s.repo.Delete(ctx, article.ID)
}

func (s *Service) GetAllTags(ctx context.Context) ([]string, error) {
	return s.repo.GetAllTags(ctx)
}
