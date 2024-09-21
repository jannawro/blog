package repository

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jannawro/blog/article"
	"github.com/jannawro/blog/middleware"
	_ "github.com/lib/pq"
)

type PostgresqlRepository struct {
	db *sql.DB
	q  *Queries
}

func NewPostgresDatabase(connString string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

func NewPostgresqlRepository(database *sql.DB, migrationFiles fs.FS) (*PostgresqlRepository, error) {
	// Run migration
	if err := runMigration(database, migrationFiles); err != nil {
		database.Close()
		return nil, fmt.Errorf("failed to run migration: %w", err)
	}
	return &PostgresqlRepository{
		db: database,
		q:  New(database),
	}, nil
}

func runMigration(db *sql.DB, migrationFiles fs.FS) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create database driver: %w", err)
	}

	// Create an embed source for the migration
	embedSource, err := iofs.New(migrationFiles, ".")
	if err != nil {
		return fmt.Errorf("failed to create embed source: %w", err)
	}

	m, err := migrate.NewWithInstance(
		"iofs", embedSource,
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
				slog.Error(fmt.Sprintf("Error rolling back transaction: %v", err), "requestID", middleware.ReqIDFromCtx(ctx))
			}
		}
	}()

	qtx := r.q.WithTx(tx)
	id, err := qtx.CreateArticle(ctx, CreateArticleParams{
		Title:           article.Title,
		Slug:            article.Slug,
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
			Slug:            a.Slug,
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
		Slug:            dbArticle.Slug,
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
		Slug:            dbArticle.Slug,
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
			Slug:            a.Slug,
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
				slog.Error(fmt.Sprintf("Error rolling back transaction: %v", err), "requestID", middleware.ReqIDFromCtx(ctx))
			}
		}
	}()

	qtx := r.q.WithTx(tx)
	dbArticle, err := qtx.UpdateArticleByID(ctx, UpdateArticleByIDParams{
		ID:              id,
		Title:           updated.Title,
		Slug:            updated.Slug,
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
		Slug:            dbArticle.Slug,
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
				slog.Error(fmt.Sprintf("Error rolling back transaction: %v", err), "requestID", middleware.ReqIDFromCtx(ctx))
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
