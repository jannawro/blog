package mysql

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"io/fs"
	"log/slog"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jannawro/blog/article"
	"github.com/jannawro/blog/middleware"
	"github.com/jannawro/blog/repository"
)

const DBDriver = "mysql"

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
	driver, err := mysql.WithInstance(db, &mysql.Config{})
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
				slog.Error(errors.Join(repository.ErrTxRollbackFailed, err).Error(), "requestID", middleware.ReqIDFromCtx(ctx))
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
	tags, err := r.q.GetAllTags(ctx)
	if err != nil {
		return nil, err
	}

	// Clean up the tags
	var cleanTags []string
	for _, tag := range tags {
		// Remove any leading or trailing whitespace, quotes, and square brackets
		cleanTag := strings.Trim(tag, " \"\n[]")
		// Split the tag in case it contains multiple tags
		splitTags := strings.Split(cleanTag, ",")
		for _, t := range splitTags {
			t = strings.TrimSpace(t)
			if t != "" {
				cleanTags = append(cleanTags, t)
			}
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
