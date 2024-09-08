package repository

import (
	"context"
	"database/sql"

	. "github.com/jannawro/blog/articles"
)

type PostgresqlRepository struct {
	db *sql.DB
	q  *Queries
}

func NewPostgresqlRepository(db *sql.DB) *PostgresqlRepository {
	return &PostgresqlRepository{
		db: db,
		q:  New(db),
	}
}

func (r *PostgresqlRepository) Create(ctx context.Context, article Article) (*Article, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	qtx := r.q.WithTx(tx)
	id, err := qtx.CreateArticle(ctx, CreateArticleParams{
		Title:           article.Title,
		Content:         article.Content,
		Tags:            article.Tags,
		PublicationDate: article.PublicationDate,
	})
	if err != nil {
		return nil, err
	}

	article.ID = id
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &article, nil
}

func (r *PostgresqlRepository) GetAll(ctx context.Context) (Articles, error) {
	dbArticles, err := r.q.GetAllArticles(ctx)
	if err != nil {
		return nil, err
	}

	articles := make(Articles, len(dbArticles))
	for i, a := range dbArticles {
		articles[i] = Article{
			ID:              a.ID,
			Title:           a.Title,
			Content:         a.Content,
			Tags:            a.Tags,
			PublicationDate: a.PublicationDate,
		}
	}

	return articles, nil
}

func (r *PostgresqlRepository) GetByID(ctx context.Context, id int64) (*Article, error) {
	dbArticle, err := r.q.GetArticleByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &Article{
		ID:              dbArticle.ID,
		Title:           dbArticle.Title,
		Content:         dbArticle.Content,
		Tags:            dbArticle.Tags,
		PublicationDate: dbArticle.PublicationDate,
	}, nil
}

func (r *PostgresqlRepository) GetByTitle(ctx context.Context, title string) (*Article, error) {
	dbArticle, err := r.q.GetArticleByTitle(ctx, title)
	if err != nil {
		return nil, err
	}

	return &Article{
		ID:              dbArticle.ID,
		Title:           dbArticle.Title,
		Content:         dbArticle.Content,
		Tags:            dbArticle.Tags,
		PublicationDate: dbArticle.PublicationDate,
	}, nil
}

func (r *PostgresqlRepository) GetByTags(ctx context.Context, tags []string) (Articles, error) {
	dbArticles, err := r.q.GetArticlesByTags(ctx, tags)
	if err != nil {
		return nil, err
	}

	articles := make(Articles, len(dbArticles))
	for i, a := range dbArticles {
		articles[i] = Article{
			ID:              a.ID,
			Title:           a.Title,
			Content:         a.Content,
			Tags:            a.Tags,
			PublicationDate: a.PublicationDate,
		}
	}

	return articles, nil
}

func (r *PostgresqlRepository) Update(ctx context.Context, id int64, updated Article) (*Article, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	qtx := r.q.WithTx(tx)
	dbArticle, err := qtx.UpdateArticleByID(ctx, UpdateArticleByIDParams{
		ID:              id,
		Title:           updated.Title,
		Content:         updated.Content,
		Tags:            updated.Tags,
		PublicationDate: updated.PublicationDate,
	})
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &Article{
		ID:              dbArticle.ID,
		Title:           dbArticle.Title,
		Content:         dbArticle.Content,
		Tags:            dbArticle.Tags,
		PublicationDate: dbArticle.PublicationDate,
	}, nil
}

func (r *PostgresqlRepository) Delete(ctx context.Context, id int64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	qtx := r.q.WithTx(tx)
	_, err = qtx.DeleteArticleByID(ctx, id)
	if err != nil {
		return err
	}

	return tx.Commit()
}
