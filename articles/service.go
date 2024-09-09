package articles

import (
	"context"
	"errors"
)

type Service struct {
	repo ArticleRepository
}

func NewService(repo ArticleRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, articleData []byte) (*Article, error) {
	var article Article
	if err := UnmarshalToArticle(articleData, &article); err != nil {
		return nil, err
	}
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

func (s *Service) GetByTitle(ctx context.Context, title string) (*Article, error) {
	return s.repo.GetByTitle(ctx, title)
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

func (s *Service) UpdateByTitle(
	ctx context.Context,
	title string,
	updatedData []byte,
) (*Article, error) {
	existingArticle, err := s.repo.GetByTitle(ctx, title)
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

func (s *Service) DeleteByTitle(ctx context.Context, title string) error {
	article, err := s.repo.GetByTitle(ctx, title)
	if err != nil {
		return err
	}
	if article == nil {
		return errors.New("article not found")
	}
	return s.repo.Delete(ctx, article.ID)
}
