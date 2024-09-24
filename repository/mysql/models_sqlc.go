// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package mysql

import (
	"database/sql"
	"encoding/json"
	"time"
)

type Article struct {
	ID              int64
	Title           string
	Slug            string
	Content         string
	Tags            json.RawMessage
	PublicationDate time.Time
	CreatedAt       sql.NullTime
	UpdatedAt       sql.NullTime
}
