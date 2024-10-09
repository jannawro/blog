-- name: CreateArticle :execresult
INSERT INTO articles (title, thumbnail, slug, content, tags, publication_date)
VALUES (?, ?, ?, ?, ?, ?);

-- name: GetAllArticles :many
SELECT * FROM articles;

-- name: GetArticleByID :one
SELECT * FROM articles
WHERE id = ? LIMIT 1;

-- name: GetArticleBySlug :one
SELECT * FROM articles
WHERE slug = ? LIMIT 1;

-- name: GetArticlesByTags :many
SELECT * FROM articles
WHERE JSON_OVERLAPS(tags, CAST(sqlc.arg(tags) AS JSON));

-- name: GetAllTags :many
SELECT DISTINCT SUBSTRING_INDEX(SUBSTRING_INDEX(tags, ',', numbers.n), ',', -1) AS unique_tag
FROM articles
CROSS JOIN (
    SELECT 1 AS n UNION ALL SELECT 2 UNION ALL SELECT 3 UNION ALL SELECT 4
) numbers
WHERE tags IS NOT NULL
  AND CHAR_LENGTH(tags) - CHAR_LENGTH(REPLACE(tags, ',', '')) >= numbers.n - 1
ORDER BY unique_tag ASC;

-- name: UpdateArticleByID :execrows
UPDATE articles
SET title = ?,
    thumbnail = ?,
    slug = ?,
    content = ?,
    tags = ?,
    publication_date = ?
WHERE id = ?;

-- name: DeleteArticleByID :execrows
DELETE FROM articles
WHERE id = ?;

SELECT id, title, thumbnail, slug, content, tags, publication_date
FROM articles
WHERE id = ?;

