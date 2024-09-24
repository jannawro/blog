package postgres_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/jannawro/blog/article"
	"github.com/jannawro/blog/repository/postgres"
	"github.com/jannawro/blog/repository/postgres/migrations"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	pqtest "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupTestDatabase(t *testing.T) (*sql.DB, func()) {
	ctx := context.Background()

	dbName := "postgres"
	dbUser := "postgres"
	dbPassword := "password"

	postgresContainer, err := pqtest.Run(ctx,
		"docker.io/postgres:16-alpine",
		pqtest.WithDatabase(dbName),
		pqtest.WithUsername(dbUser),
		pqtest.WithPassword(dbPassword),
		pqtest.WithSQLDriver("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	assert.NoError(t, err)

	connString, err := postgresContainer.ConnectionString(ctx, "sslmode=disable", "parseTime=true")
	assert.NoError(t, err)

	db, err := sql.Open("postgres", connString)
	require.NoError(t, err)

	cleanup := func() {
		db.Close()
		postgresContainer.Terminate(ctx)
	}

	return db, cleanup
}

func TestRepository(t *testing.T) {
	db, cleanup := setupTestDatabase(t)
	defer cleanup()

	repo, err := postgres.NewRepository(db, migrations.Files())
	require.NoError(t, err)

	ctx := context.Background()

	t.Run("Create and GetByID", func(t *testing.T) {
		article := article.Article{
			Title:           "Test Article",
			Slug:            "test-article",
			Content:         "This is a test article",
			Tags:            []string{"test", "golang"},
			PublicationDate: time.Now().UTC().Truncate(time.Second),
		}

		createdArticle, err := repo.Create(ctx, article)
		require.NoError(t, err)
		assert.NotZero(t, createdArticle.ID)

		fetchedArticle, err := repo.GetByID(ctx, createdArticle.ID)
		require.NoError(t, err)
		assert.Equal(t, createdArticle, fetchedArticle)
	})

	t.Run("GetAll", func(t *testing.T) {
		articles, err := repo.GetAll(ctx)
		require.NoError(t, err)
		assert.NotEmpty(t, articles)
	})

	t.Run("GetBySlug", func(t *testing.T) {
		article, err := repo.GetBySlug(ctx, "test-article")
		require.NoError(t, err)
		assert.Equal(t, "Test Article", article.Title)
	})

	t.Run("GetByTags", func(t *testing.T) {
		articles, err := repo.GetByTags(ctx, []string{"test"})
		require.NoError(t, err)
		assert.NotEmpty(t, articles)
	})

	t.Run("Update", func(t *testing.T) {
		article, err := repo.GetBySlug(ctx, "test-article")
		require.NoError(t, err)

		article.Title = "Updated Test Article"
		updatedArticle, err := repo.Update(ctx, article.ID, *article)
		require.NoError(t, err)
		assert.Equal(t, "Updated Test Article", updatedArticle.Title)
	})

	t.Run("GetAllTags", func(t *testing.T) {
		tags, err := repo.GetAllTags(ctx)
		require.NoError(t, err)
		assert.Contains(t, tags, "test")
		assert.Contains(t, tags, "golang")
	})

	t.Run("Delete", func(t *testing.T) {
		article, err := repo.GetBySlug(ctx, "test-article")
		require.NoError(t, err)

		err = repo.Delete(ctx, article.ID)
		require.NoError(t, err)

		_, err = repo.GetByID(ctx, article.ID)
		assert.Error(t, err)
	})
}
