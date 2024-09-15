-- name: CreateArticle :one
INSERT INTO articles (title, slug, content, tags, publication_date)
VALUES ($1, $2, $3, $4, $5)
RETURNING id;

-- name: GetAllArticles :many
SELECT * FROM articles;

-- name: GetArticleByID :one
SELECT * FROM articles
WHERE id = $1 LIMIT 1;

-- name: GetArticleBySlug :one
SELECT * FROM articles
WHERE slug = $1 LIMIT 1;

-- name: GetArticlesByTags :many
SELECT * FROM articles
WHERE tags && sqlc.arg(tags)::text[];

-- name: GetAllTags :many
SELECT DISTINCT unnest(tags)::TEXT AS unique_tag
FROM articles
WHERE tags IS NOT NULL
ORDER BY unique_tag ASC;

-- name: UpdateArticleByID :one
UPDATE articles
SET title = $1,
    content = $2,
    tags = $3,
    publication_date = $4
WHERE id = $5
RETURNING id, title, content, tags, publication_date;

-- name: DeleteArticleByID :one
DELETE FROM articles
WHERE id = $1
RETURNING id, title, content, tags, publication_date;
