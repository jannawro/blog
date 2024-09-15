package repository

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jannawro/blog/article"
	_ "github.com/lib/pq"
)

type PostgresqlRepository struct {
	db *sql.DB
	q  *Queries
}

func NewPostgresqlRepository(connString string) (*PostgresqlRepository, error) {
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Run migration
	if err := runMigration(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to run migration: %w", err)
	}

	return &PostgresqlRepository{
		db: db,
		q:  New(db),
	}, nil
}

func runMigration(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create database driver: %w", err)
	}

	// Create a file source for the migration
	fileSource, err := (&file.File{}).Open("file://repository/sqlc/migrations")
	if err != nil {
		return fmt.Errorf("failed to create file source: %w", err)
	}

	m, err := migrate.NewWithInstance(
		"file", fileSource,
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	// Run the migration
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migration: %w", err)
	}

	return nil
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

func (r *PostgresqlRepository) GetBySlug(ctx context.Context, slug string) (*article.Article, error) {
	dbArticle, err := r.q.GetArticleBySlug(ctx, slug)
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
