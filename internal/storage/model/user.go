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

func (u *User) GetID() uint64 {
	return u.ID
}

func (u *User) TableName() string {
	return "users"
}

func (u *User) SelectFields() []string {
	return []string{
		"users.id",
		"users.created_at",
		"users.updated_at",
		"users.username",
		"users.email",
		"users.password",
	}
}

func (u *User) InsertMap() map[string]interface{} {
	return map[string]interface{}{
		"username": u.Username,
		"email":    u.Email,
		"password": u.Password,
	}
}

func (u *User) UpdateMap() map[string]interface{} {
	return u.InsertMap()
}

func (u *User) SelectJoins() []string {
	return []string{}
}

type Role struct {
	ID        uint64 `db:"id"`
	Authority string `db:"authority"`
}

func (r *Role) GetID() uint64 {
	return r.ID
}

func (r Role) TableName() string {
	return "roles"
}

func (r Role) SelectFields() []string {
	return []string{"roles.id", "roles.authority"}
}

func (r Role) InsertMap() map[string]interface{} {
	return map[string]interface{}{"authority": r.Authority}
}

func (r Role) UpdateMap() map[string]interface{} {
	return r.InsertMap()
}

func (r Role) SelectJoins() []string {
	return []string{}
}

type UserRole struct {
	UserID    uint64 `db:"user_id"`
	RoleID    uint64 `db:"role_id"`
	Authority string `db:"authority"`
}

func (u UserRole) GetID() uint64 {
	return 0
}

func (u UserRole) TableName() string {
	return "user_roles"
}

func (u UserRole) SelectFields() []string {
	return []string{
		"user_roles.user_id",
		"user_roles.role_id",
		"roles.authority",
	}
}

func (u UserRole) InsertMap() map[string]interface{} {
	return map[string]interface{}{
		"user_id": u.UserID,
		"role_id": u.RoleID,
	}
}

func (u UserRole) UpdateMap() map[string]interface{} {
	return u.InsertMap()
}

func (u UserRole) SelectJoins() []string {
	return []string{"JOIN roles ON roles.id = user_roles.role_id"}
}

type UserStats struct {
	UserID     uint64 `db:"user_id"`
	NumPublic  uint64 `db:"num_public"`
	NumPrivate uint64 `db:"num_private"`
	AvgRating  uint64 `db:"avg_rating"`
}
