package model

import (
	"database/sql"
	"time"
)

type User struct {
	ID        uint64       `db:"id"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt time.Time    `db:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at"`
	Username  string       `db:"username"`
	Email     string       `db:"email"`
	Password  string
	Roles     []Role
}

type Role struct {
	ID        uint64 `db:"id"`
	Authority string `db:"authority"`
}

type UserRole struct {
	UserID    uint64 `db:"user_id"`
	RoleID    uint64 `db:"role_id"`
	Authority string `db:"authority"`
}

type UserStats struct {
	UserID     uint64 `db:"user_id"`
	NumPublic  uint64 `db:"num_public"`
	NumPrivate uint64 `db:"num_private"`
	AvgRating  uint64 `db:"avg_rating"`
}
