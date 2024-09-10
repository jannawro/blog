package repository

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/jannawro/blog/article"
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

func (r *PostgresqlRepository) Create(ctx context.Context, article article.Article) (*article.Article, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			if err != sql.ErrTxDone {
				fmt.Fprintf(os.Stderr, "Error rolling back transaction: %v\n", err)
			}
		}
	}()

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

func (r *PostgresqlRepository) GetAll(ctx context.Context) (article.Articles, error) {
	dbArticles, err := r.q.GetAllArticles(ctx)
	if err != nil {
		return nil, err
	}

	articlesSlice := make(article.Articles, len(dbArticles))
	for i, a := range dbArticles {
		articlesSlice[i] = article.Article{
			ID:              a.ID,
			Title:           a.Title,
			Content:         a.Content,
			Tags:            a.Tags,
			PublicationDate: a.PublicationDate,
		}
	}

	return articlesSlice, nil
}

func (r *PostgresqlRepository) GetByID(ctx context.Context, id int64) (*article.Article, error) {
	dbArticle, err := r.q.GetArticleByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &article.Article{
		ID:              dbArticle.ID,
		Title:           dbArticle.Title,
		Content:         dbArticle.Content,
		Tags:            dbArticle.Tags,
		PublicationDate: dbArticle.PublicationDate,
	}, nil
}

func (r *PostgresqlRepository) GetByTitle(ctx context.Context, title string) (*article.Article, error) {
	dbArticle, err := r.q.GetArticleByTitle(ctx, title)
	if err != nil {
		return nil, err
	}

	return &article.Article{
		ID:              dbArticle.ID,
		Title:           dbArticle.Title,
		Content:         dbArticle.Content,
		Tags:            dbArticle.Tags,
		PublicationDate: dbArticle.PublicationDate,
	}, nil
}

func (r *PostgresqlRepository) GetByTags(ctx context.Context, tags []string) (article.Articles, error) {
	dbArticles, err := r.q.GetArticlesByTags(ctx, tags)
	if err != nil {
		return nil, err
	}

	articlesSlice := make(article.Articles, len(dbArticles))
	for i, a := range dbArticles {
		articlesSlice[i] = article.Article{
			ID:              a.ID,
			Title:           a.Title,
			Content:         a.Content,
			Tags:            a.Tags,
			PublicationDate: a.PublicationDate,
		}
	}

	return articlesSlice, nil
}

func (r *PostgresqlRepository) Update(
	ctx context.Context,
	id int64,
	updated article.Article,
) (*article.Article, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			if err != sql.ErrTxDone {
				fmt.Fprintf(os.Stderr, "Error rolling back transaction: %v\n", err)
			}
		}
	}()

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

	return &article.Article{
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
	defer func() {
		if err := tx.Rollback(); err != nil {
			if err != sql.ErrTxDone {
				fmt.Fprintf(os.Stderr, "Error rolling back transaction: %v\n", err)
			}
		}
	}()

	qtx := r.q.WithTx(tx)
	_, err = qtx.DeleteArticleByID(ctx, id)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *PostgresqlRepository) GetAllTags(ctx context.Context) ([]string, error) {
	return r.q.GetAllTags(ctx)
}
