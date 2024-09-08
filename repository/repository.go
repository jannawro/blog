package repository

import (
	"context"

	. "github.com/jannawro/blog/articles"
)

type ArticleRepositoryInterface interface {
	Create(ctx context.Context, article Article) (*Article, error)
	GetAll(ctx context.Context) (Articles, error)
	GetByID(ctx context.Context, id int64) (*Article, error)
	GetByTitle(ctx context.Context, title string) (*Article, error)
	GetByTags(ctx context.Context, tags []string) (Articles, error)
	Update(ctx context.Context, id int64, updated Article) (*Article, error)
	Delete(ctx context.Context, id int64) error
}

type PostgresqlRepository struct{}
