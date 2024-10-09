package postgres

import (
	"context"
	"database/sql"
	"errors"
	"io/fs"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jannawro/blog/article"
	"github.com/jannawro/blog/middleware"
	"github.com/jannawro/blog/repository"
	_ "github.com/lib/pq"
)

const DBDriver = "postgres"

type Repository struct {
	db *sql.DB
	q  *Queries
}

func NewDatabase(connString string) (*sql.DB, error) {
	db, err := sql.Open(DBDriver, connString)
	if err != nil {
		return nil, errors.Join(repository.ErrDatabaseConnectionFailed, err)
	}

	if err := db.Ping(); err != nil {
		closeErr := db.Close()
		if closeErr != nil {
			return nil, errors.Join(repository.ErrPingFailed, err, closeErr)
		}
		return nil, errors.Join(repository.ErrPingFailed, err)
	}

	return db, nil
}

func NewRepository(db *sql.DB, migrationFiles fs.FS) (*Repository, error) {
	// Run migration
	if err := runMigration(db, migrationFiles); err != nil {
		closeErr := db.Close()
		if closeErr != nil {
			return nil, errors.Join(repository.ErrMigrationFailed, err, closeErr)
		}
		return nil, errors.Join(repository.ErrMigrationFailed, err)
	}
	return &Repository{
		db: db,
		q:  New(db),
	}, nil
}

func runMigration(db *sql.DB, migrationFiles fs.FS) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return errors.Join(repository.ErrDriverCreationFailed, err)
	}

	// Create an embed source for the migration
	embedSource, err := iofs.New(migrationFiles, ".")
	if err != nil {
		return errors.Join(repository.ErrEmbedFailed, err)
	}

	m, err := migrate.NewWithInstance(
		"iofs", embedSource,
		DBDriver, driver)
	if err != nil {
		return errors.Join(repository.ErrMigrationInstanceFailed, err)
	}

	// Run the migration
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return errors.Join(repository.ErrMigrationRunFailed, err)
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
				slog.Error(errors.Join(repository.ErrTxRollbackFailed, err).Error(), "requestID", middleware.ReqIDFromCtx(ctx))
			}
		}
	}()

	qtx := r.q.WithTx(tx)
	id, err := qtx.CreateArticle(ctx, CreateArticleParams{
		Title:           article.Title,
		Thumbnail:       article.Thumbnail,
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
			Thumbnail:       a.Thumbnail,
			Slug:            a.Slug,
			Content:         a.Content,
			Tags:            a.Tags,
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
		Thumbnail:       dbArticle.Thumbnail,
		Slug:            dbArticle.Slug,
		Content:         dbArticle.Content,
		Tags:            dbArticle.Tags,
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
		Thumbnail:       dbArticle.Thumbnail,
		Slug:            dbArticle.Slug,
		Content:         dbArticle.Content,
		Tags:            dbArticle.Tags,
		PublicationDate: dbArticle.PublicationDate,
	}, nil
}

func (r *Repository) GetByTags(ctx context.Context, tags []string) (article.Articles, error) {
	dbArticles, err := r.q.GetArticlesByTags(ctx, tags)
	if err != nil {
		return nil, err
	}

	articlesSlice := make(article.Articles, len(dbArticles))
	for i, a := range dbArticles {
		articlesSlice[i] = article.Article{
			ID:              a.ID,
			Title:           a.Title,
			Thumbnail:       a.Thumbnail,
			Slug:            a.Slug,
			Content:         a.Content,
			Tags:            a.Tags,
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
				slog.Error(errors.Join(repository.ErrTxRollbackFailed, err).Error(), "requestID", middleware.ReqIDFromCtx(ctx))
			}
		}
	}()

	qtx := r.q.WithTx(tx)
	dbArticle, err := qtx.UpdateArticleByID(ctx, UpdateArticleByIDParams{
		ID:              id,
		Title:           updated.Title,
		Thumbnail:       updated.Thumbnail,
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
		Thumbnail:       dbArticle.Thumbnail,
		Slug:            dbArticle.Slug,
		Content:         dbArticle.Content,
		Tags:            dbArticle.Tags,
		PublicationDate: dbArticle.PublicationDate,
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
				slog.Error(errors.Join(repository.ErrTxRollbackFailed, err).Error(), "requestID", middleware.ReqIDFromCtx(ctx))
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
	return r.q.GetAllTags(ctx)
}
