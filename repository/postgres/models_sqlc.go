// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package postgres

import (
	"database/sql"
	"time"
)

type Article struct {
	ID              int64
	Title           string
	Thumbnail       string
	Slug            string
	Content         string
	Tags            []string
	PublicationDate time.Time
	CreatedAt       sql.NullTime
	UpdatedAt       sql.NullTime
}
