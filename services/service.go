package service

import (
	"context"
	"errors"

	"github.com/jannawro/blog/articles"
)

type ArticleService struct {
	repo articles.ArticleRepository
}

func NewArticleService(repo articles.ArticleRepository) *ArticleService {
	return &ArticleService{repo: repo}
}

func (s *ArticleService) Create(ctx context.Context, articleData []byte) (*articles.Article, error) {
	var article articles.Article
	if err := articles.UnmarshalToArticle(articleData, &article); err != nil {
		return nil, err
	}
	return s.repo.Create(ctx, article)
}

func (s *ArticleService) GetAll(ctx context.Context, sortBy *articles.SortOption) (articles.Articles, error) {
	articles, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	if sortBy != nil {
		articles.Sort(*sortBy)
	}
	return articles, nil
}

func (s *ArticleService) GetByTitle(ctx context.Context, title string) (*articles.Article, error) {
	return s.repo.GetByTitle(ctx, title)
}

func (s *ArticleService) GetByTags(ctx context.Context, tags []string, sortBy *articles.SortOption) (articles.Articles, error) {
	articles, err := s.repo.GetByTags(ctx, tags)
	if err != nil {
		return nil, err
	}
	if sortBy != nil {
		articles.Sort(*sortBy)
	}
	return articles, nil
}

func (s *ArticleService) UpdateByTitle(
	ctx context.Context,
	title string,
	updatedData []byte,
) (*articles.Article, error) {
	existingArticle, err := s.repo.GetByTitle(ctx, title)
	if err != nil {
		return nil, err
	}

	var updatedArticle articles.Article
	if err := articles.UnmarshalToArticle(updatedData, &updatedArticle); err != nil {
		return nil, err
	}

	updatedArticle.ID = existingArticle.ID
	return s.repo.Update(ctx, existingArticle.ID, updatedArticle)
}

func (s *ArticleService) DeleteByTitle(ctx context.Context, title string) error {
	article, err := s.repo.GetByTitle(ctx, title)
	if err != nil {
		return err
	}
	if article == nil {
		return errors.New("article not found")
	}
	return s.repo.Delete(ctx, article.ID)
}
