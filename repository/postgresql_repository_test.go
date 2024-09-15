package repository

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/jannawro/blog/article"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupTestDatabase(t *testing.T) (*sql.DB, func()) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:13",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_DB":       "testdb",
			"POSTGRES_USER":     "testuser",
			"POSTGRES_PASSWORD": "testpass",
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections"),
	}

	postgresC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)

	host, err := postgresC.Host(ctx)
	require.NoError(t, err)

	port, err := postgresC.MappedPort(ctx, "5432")
	require.NoError(t, err)

	connString := fmt.Sprintf("postgres://testuser:testpass@%s:%s/testdb?sslmode=disable", host, port.Port())

	db, err := sql.Open("postgres", connString)
	require.NoError(t, err)

	return db, func() {
		db.Close()
		postgresC.Terminate(ctx)
	}
}

func TestPostgresqlRepository(t *testing.T) {
	db, cleanup := setupTestDatabase(t)
	defer cleanup()

	repo, err := NewPostgresqlRepository(db)
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

	t.Run("Delete", func(t *testing.T) {
		article, err := repo.GetBySlug(ctx, "test-article")
		require.NoError(t, err)

		err = repo.Delete(ctx, article.ID)
		require.NoError(t, err)

		_, err = repo.GetByID(ctx, article.ID)
		assert.Error(t, err)
	})

	t.Run("GetAllTags", func(t *testing.T) {
		tags, err := repo.GetAllTags(ctx)
		require.NoError(t, err)
		assert.Contains(t, tags, "test")
		assert.Contains(t, tags, "golang")
	})
}
