package model

import (
	"database/sql"
	"fmt"
	"time"
)

type RecipeParams struct {
	Limit          uint64     `json:"limit" form:"limit"`
	Offset         uint64     `json:"offset" form:"offset"`
	Sort           string     `json:"sort" form:"sort"`
	Order          string     `json:"order" form:"order"`
	Title          string     `json:"title" form:"title"`
	UserID         *uint64    `json:"userId" form:"userId"`
	IncludeFlavors []uint64   `json:"includeFlavor" form:"includeFlavors"`
	ExcludeFlavors []uint64   `json:"excludeFlavors" form:"excludeFlavors"`
	CreatedFrom    *time.Time `json:"createdFrom" form:"createdFrom"`
	CreatedTo      *time.Time `json:"createdTo" form:"createdTo"`
	Snv            *int       `json:"snv" form:"snv"`
	Public         bool       `json:"-"`
	Current        bool       `json:"current"`
	Group          string
}

func (r RecipeParams) GetLimit() uint64 {
	limit := r.Limit
	if limit < 0 {
		limit = 25
	} else if limit > 100 {
		limit = 100
	}
	return limit
}

func (r RecipeParams) GetOffset() uint64 {
	if r.Offset < 0 {
		return 0
	}
	return r.Offset
}

func (r RecipeParams) GetSort() string {
	if r.Sort == "" {
		return r.Sort
	}
	var sort string
	switch r.Sort {
	case "rating":
		sort = "avg_rating"
	case "owner":
		sort = "owner_username"
	default:
		sort = fmt.Sprintf("recipes.%s", r.Sort)
	}
	order := r.Order
	if order != "desc" {
		order = "asc"
	}
	return fmt.Sprintf("%s %s", sort, order)
}

func (r RecipeParams) GetGroup() string {
	return r.Group
}

type Recipe struct {
	ID            uint64     `db:"id"`
	CreatedAt     *time.Time `db:"created_at"`
	UpdatedAt     *time.Time `db:"updated_at"`
	OwnerID       uint64     `db:"owner_id"`
	Current       bool       `db:"current"`
	Description   string     `db:"description"`
	Public        bool       `db:"public"`
	PublishedAt   *time.Time `db:"published_at"`
	RemixOfID     *uint64    `db:"remix_of_id"`
	Snv           bool       `db:"snv"`
	SteepDays     uint64     `db:"steep_days"`
	TempF         uint64     `db:"temp_f"`
	Title         string     `db:"title"`
	UUID          string     `db:"uuid"`
	Version       uint64     `db:"version"`
	VgPercentM    uint64     `db:"vg_percent_m"`
	Wip           bool       `db:"wip"`
	AvgRating     uint64     `db:"avg_rating"`
	OwnerUsername string     `db:"owner_username"`
	RemixOfTitle  *string    `db:"remix_of_title"`
	DeletedAt     sql.NullTime
}

type RecipeCollaboratorParams struct {
	RecipeID  uint64
	RecipeIDs []uint64
}

type RecipeCollaborator struct {
	UserID   uint64 `db:"user_id"`
	RecipeID uint64 `db:"recipe_id"`
	Username string `db:"username"`
}

type RecipeFlavorParams struct {
	RecipeID  uint64
	RecipeIDs []uint64
}

type RecipeFlavor struct {
	ID       uint64 `db:"id"`
	PercentM uint64 `db:"percent_m"`
	FlavorID uint64 `db:"flavor_id"`
	RecipeID uint64 `db:"recipe_id"`
}

type RecipeFlavorSubstitution struct {
	RecipeFlavorID uint64
	FlavorID       uint64
}

type RecipeRating struct {
	ID        uint64     `db:"id"`
	CreatedAt *time.Time `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
	OwnerID   uint64     `db:"owner_id"`
	Rating    uint64     `db:"rating"`
	RecipeID  uint64     `db:"recipe_id"`
}

type RecipeComment struct {
	ID        uint64    `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	DeletedAt sql.NullTime
	Content   string
	OwnerID   uint64
	RecipeID  uint64
	ReplyToID uint64
}

type RecipeTagParams struct {
	RecipeID  uint64
	RecipeIDs []uint64
}

type RecipeTag struct {
	TagID    uint64 `db:"tag_id"`
	Tag      string `db:"tag"`
	RecipeID uint64 `db:"recipe_id"`
}
