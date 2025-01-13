// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0

package repo

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Model struct {
	ID         int32            `json:"id"`
	Field1     string           `json:"field1"`
	Field2     int32            `json:"field2"`
	CreatedBy  pgtype.Int4      `json:"created_by"`
	CreatedAt  pgtype.Timestamp `json:"created_at"`
	ModifiedBy pgtype.Int4      `json:"modified_by"`
	ModifiedAt pgtype.Timestamp `json:"modified_at"`
	DeletedBy  pgtype.Int4      `json:"deleted_by"`
	DeletedAt  pgtype.Timestamp `json:"deleted_at"`
}