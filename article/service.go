package article

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

func (s *Service) Create(ctx context.Context, article Article) (*Article, error) {
	a, err := s.repo.Create(ctx, article)
	if err != nil {
		return nil, errors.Join(ErrArticleCreationFailed, err)
	}
	return a, nil
}

func (s *Service) GetAll(ctx context.Context, sortBy *SortOption) (Articles, error) {
	articles, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, errors.Join(ErrArticlesNotFound, err)
	}
	if sortBy != nil {
		articles.Sort(*sortBy)
	}
	return articles, nil
}

func (s *Service) GetBySlug(ctx context.Context, slug string) (*Article, error) {
	article, err := s.repo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, errors.Join(ErrArticleNotFound, err)
	}
	return article, nil
}

func (s *Service) GetByID(ctx context.Context, id int64) (*Article, error) {
	article, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.Join(ErrArticleNotFound, err)
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
		return nil, errors.Join(ErrArticlesNotFound, err)
	}
	if sortBy != nil {
		articles.Sort(*sortBy)
	}
	return articles, nil
}

func (s *Service) UpdateBySlug(
	ctx context.Context,
	slug string,
	updatedArticle Article,
) (*Article, error) {
	existingArticle, err := s.repo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, errors.Join(ErrArticleNotFound, err)
	}

	updatedArticle.ID = existingArticle.ID
	a, err := s.repo.Update(ctx, existingArticle.ID, updatedArticle)
	if err != nil {
		return nil, errors.Join(ErrArticleUpdateFailed, err)
	}
	return a, nil
}

func (s *Service) DeleteBySlug(ctx context.Context, slug string) error {
	article, err := s.repo.GetBySlug(ctx, slug)
	if err != nil {
		return errors.Join(ErrArticleDeletionFailed, err)
	}
	if article == nil {
		return ErrArticleNotFound
	}
	err = s.repo.Delete(ctx, article.ID)
	if err != nil {
		return errors.Join(ErrArticleDeletionFailed, err)
	}
	return nil
}

func (s *Service) GetAllTags(ctx context.Context) ([]string, error) {
	tags, err := s.repo.GetAllTags(ctx)
	if err != nil {
		return nil, errors.Join(ErrArticleDeletionFailed, err)
	}
	return tags, nil
}
