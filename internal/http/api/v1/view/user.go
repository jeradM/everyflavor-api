package view

import (
	"time"
)

type User struct {
	ID        uint64     `json:"id"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	Password  string     `json:"-"`
	Roles     []UserRole `json:"roles"`
	Stats     UserStats  `json:"stats"`
}

type UserRole struct {
	ID        uint64 `json:"id"`
	Authority string `json:"authority"`
}

type UserStats struct {
	NumPublic  uint64 `json:"numPublic"`
	NumPrivate uint64 `json:"numPrivate"`
	AvgRating  uint64 `json:"avgRating"`
}
