package mysql

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/fs"
	"log/slog"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jannawro/blog/article"
	"github.com/jannawro/blog/middleware"
)

type Repository struct {
	db *sql.DB
	q  *Queries
}

func NewDatabase(connString string) (*sql.DB, error) {
	db, err := sql.Open("mysql", connString)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

func NewRepository(database *sql.DB, migrationFiles fs.FS) (*Repository, error) {
	// Run migration
	if err := runMigration(database, migrationFiles); err != nil {
		database.Close()
		return nil, fmt.Errorf("failed to run migration: %w", err)
	}
	return &Repository{
		db: database,
		q:  New(database),
	}, nil
}

func runMigration(db *sql.DB, migrationFiles fs.FS) error {
	driver, err := mysql.WithInstance(db, &mysql.Config{})
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

func (r *Repository) Create(ctx context.Context, article article.Article) (*article.Article, error) {
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
	result, err := qtx.CreateArticle(ctx, CreateArticleParams{
		Title:           article.Title,
		Slug:            article.Slug,
		Content:         article.Content,
		Tags:            tagsToJSON(article.Tags),
		PublicationDate: article.PublicationDate,
	})
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	article.ID = id

	return &article, nil
}

func (r *Repository) GetAll(ctx context.Context) (article.Articles, error) {
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
			Tags:            jsonToTags(a.Tags),
			PublicationDate: a.PublicationDate,
		}
	}

	return articlesSlice, nil
}

func (r *Repository) GetByID(ctx context.Context, id int64) (*article.Article, error) {
	dbArticle, err := r.q.GetArticleByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &article.Article{
		ID:              dbArticle.ID,
		Title:           dbArticle.Title,
		Slug:            dbArticle.Slug,
		Content:         dbArticle.Content,
		Tags:            jsonToTags(dbArticle.Tags),
		PublicationDate: dbArticle.PublicationDate,
	}, nil
}

func (r *Repository) GetBySlug(ctx context.Context, slug string) (*article.Article, error) {
	dbArticle, err := r.q.GetArticleBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	return &article.Article{
		ID:              dbArticle.ID,
		Title:           dbArticle.Title,
		Slug:            dbArticle.Slug,
		Content:         dbArticle.Content,
		Tags:            jsonToTags(dbArticle.Tags),
		PublicationDate: dbArticle.PublicationDate,
	}, nil
}

func (r *Repository) GetByTags(ctx context.Context, tags []string) (article.Articles, error) {
	dbArticles, err := r.q.GetArticlesByTags(ctx, tagsToJSON(tags))
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
			Tags:            jsonToTags(a.Tags),
			PublicationDate: a.PublicationDate,
		}
	}

	return articlesSlice, nil
}

func (r *Repository) Update(
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
	id, err = qtx.UpdateArticleByID(ctx, UpdateArticleByIDParams{
		ID:              id,
		Title:           updated.Title,
		Slug:            updated.Slug,
		Content:         updated.Content,
		Tags:            tagsToJSON(updated.Tags),
		PublicationDate: updated.PublicationDate,
	})
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	a, err := r.q.GetArticleByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &article.Article{
		ID:              a.ID,
		Title:           a.Title,
		Slug:            a.Slug,
		Content:         a.Content,
		Tags:            jsonToTags(a.Tags),
		PublicationDate: a.PublicationDate,
	}, nil
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
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

func (r *Repository) GetAllTags(ctx context.Context) ([]string, error) {
	tags, err := r.q.GetAllTags(ctx)
	if err != nil {
		return nil, err
	}

	// Clean up the tags
	var cleanTags []string
	for _, tag := range tags {
		// Remove any leading or trailing whitespace and quotes
		cleanTag := strings.Trim(tag, " \"\n")
		if cleanTag != "" {
			cleanTags = append(cleanTags, cleanTag)
		}
	}

	return cleanTags, nil
}

func tagsToJSON(tags []string) json.RawMessage {
	jsonTags, err := json.Marshal(tags)
	if err != nil {
		panic(err)
	}
	return jsonTags
}

func jsonToTags(j json.RawMessage) []string {
	var tags []string
	err := json.Unmarshal(j, &tags)
	if err != nil {
		panic(err)
	}
	return tags
}
