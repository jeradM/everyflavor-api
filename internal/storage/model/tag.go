package model

import "time"

type Tag struct {
	ID        uint64    `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	Tag       string    `db:"tag"`
}
