package model

import (
	"database/sql"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
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
	ID            uint64       `db:"id"`
	CreatedAt     *time.Time   `db:"created_at"`
	UpdatedAt     *time.Time   `db:"updated_at"`
	DeletedAt     sql.NullTime `deletedAt,omitempty"`
	OwnerID       uint64       `db:"owner_id"`
	Current       bool         `db:"current"`
	Description   string       `db:"description"`
	Public        bool         `db:"public"`
	RemixOfID     *uint64      `db:"remix_of_id"`
	Snv           bool         `db:"snv"`
	SteepDays     uint64       `db:"steep_days"`
	TempF         uint64       `db:"temp_f"`
	Title         string       `db:"title"`
	UUID          string       `db:"uuid"`
	Version       uint64       `db:"version"`
	VgPercentM    uint64       `db:"vg_percent_m"`
	Wip           bool         `db:"wip"`
	AvgRating     uint64       `db:"avg_rating"`
	OwnerUsername string       `db:"owner_username"`
	RemixOfTitle  *string      `db:"remix_of_title"`
}

func (r Recipe) GetID() uint64 {
	return r.ID
}

func (r Recipe) TableName() string {
	return "recipes"
}

func (r Recipe) SelectFields() []string {
	return []string{
		"recipes.id",
		"recipes.created_at",
		"recipes.updated_at",
		"recipes.owner_id",
		"recipes.current",
		"recipes.description",
		"recipes.public",
		"recipes.remix_of_id",
		"recipes.snv",
		"recipes.steep_days",
		"recipes.temp_f",
		"recipes.title",
		"recipes.uuid",
		"recipes.version",
		"recipes.vg_percent_m",
		"recipes.wip",
		"CAST(IFNULL(avg(recipe_ratings.rating) * 1000, 0) as UNSIGNED) as avg_rating",
		"users.username owner_username",
		"remix.title remix_of_title",
	}
}

func (r Recipe) InsertMap() map[string]interface{} {
	return map[string]interface{}{
		"owner_id":     r.OwnerID,
		"current":      1,
		"description":  r.Description,
		"public":       r.Public,
		"remix_of_id":  r.RemixOfID,
		"snv":          r.Snv,
		"steep_days":   r.SteepDays,
		"temp_f":       r.TempF,
		"title":        r.Title,
		"uuid":         r.UUID,
		"version":      r.Version,
		"vg_percent_m": r.VgPercentM,
		"wip":          r.Wip,
	}
}

func (r Recipe) UpdateMap() map[string]interface{} {
	return map[string]interface{}{
		"owner_id":     r.OwnerID,
		"description":  r.Description,
		"public":       r.Public,
		"remix_of_id":  r.RemixOfID,
		"snv":          r.Snv,
		"steep_days":   r.SteepDays,
		"temp_f":       r.TempF,
		"title":        r.Title,
		"vg_percent_m": r.VgPercentM,
		"wip":          r.Wip,
	}
}

func (r Recipe) SelectJoins() []string {
	return []string{
		"LEFT JOIN users ON users.id = recipes.owner_id",
		"LEFT JOIN recipe_ratings ON recipe_ratings.recipe_id = recipes.id",
		"LEFT JOIN recipes remix ON remix.id = recipes.remix_of_id",
	}
}

func (r Recipe) OwnerField() string {
	return "recipes.owner_id"
}

func (r Recipe) PublicField() string {
	return "recipes.public"
}

func (r Recipe) CollaboratorIDsQuery(id uint64) (string, []interface{}) {
	query, args, _ := sq.Select("user_id").
		From("recipe_collaborators").
		Where(sq.Eq{"recipe_id": id}).
		ToSql()
	return query, args
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

func (r RecipeCollaborator) GetID() uint64 {
	return 0
}

func (r RecipeCollaborator) TableName() string {
	return "recipe_collaborators"
}

func (r RecipeCollaborator) SelectFields() []string {
	return []string{
		"recipe_collaborators.user_id",
		"users.username",
		"recipe_collaborators.recipe_id",
	}
}

func (r RecipeCollaborator) InsertMap() map[string]interface{} {
	return map[string]interface{}{
		"recipe_id": r.RecipeID,
		"flavor_id": r.UserID,
	}
}

func (r RecipeCollaborator) UpdateMap() map[string]interface{} {
	return r.InsertMap()
}

func (r RecipeCollaborator) SelectJoins() []string {
	return []string{"JOIN users ON users.id = recipe_collaborators.user_id"}
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

func (r RecipeFlavor) GetID() uint64 {
	return r.ID
}

func (r RecipeFlavor) TableName() string {
	return "recipe_flavors"
}

func (r RecipeFlavor) SelectFields() []string {
	return []string{
		"recipe_flavors.id",
		"recipe_flavors.percent_m",
		"recipe_flavors.flavor_id",
		"recipe_flavors.recipe_id",
	}
}

func (r RecipeFlavor) InsertMap() map[string]interface{} {
	return map[string]interface{}{
		"percent_m": r.PercentM,
		"flavor_id": r.FlavorID,
		"recipe_id": r.RecipeID,
	}
}

func (r RecipeFlavor) UpdateMap() map[string]interface{} {
	return r.InsertMap()
}

func (r RecipeFlavor) SelectJoins() []string {
	return []string{}
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

func (r RecipeRating) GetID() uint64 {
	return r.ID
}

func (r RecipeRating) TableName() string {
	return "recipe_ratings"
}

func (r RecipeRating) SelectFields() []string {
	return []string{
		"recipe_ratings.id",
		"recipe_ratings.created_at",
		"recipe_ratings.updated_at",
		"recipe_ratings.rating",
		"recipe_ratings.recipe_id",
		"recipe_ratings.owner_id",
	}
}

func (r RecipeRating) InsertMap() map[string]interface{} {
	return map[string]interface{}{
		"rating":    r.Rating,
		"recipe_id": r.RecipeID,
		"owner_id":  r.OwnerID,
	}
}

func (r RecipeRating) UpdateMap() map[string]interface{} {
	return map[string]interface{}{
		"rating": r.Rating,
	}
}

func (r RecipeRating) SelectJoins() []string {
	return []string{}
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

func (r RecipeComment) GetID() uint64 {
	return r.ID
}

func (r RecipeComment) TableName() string {
	return "recipe_comments"
}

func (r RecipeComment) SelectFields() []string {
	return []string{
		"recipe_comments.id",
		"recipe_comments.created_at",
		"recipe_comments.updated_at",
		"recipe_comments.content",
		"recipe_comments.owner_id",
		"recipe_comments.recipe_id",
		"recipe_comments.reply_to_id",
	}
}

func (r RecipeComment) InsertMap() map[string]interface{} {
	return map[string]interface{}{
		"content":     r.Content,
		"recipe_id":   r.RecipeID,
		"owner_id":    r.OwnerID,
		"reply_to_id": r.ReplyToID,
	}
}

func (r RecipeComment) UpdateMap() map[string]interface{} {
	return map[string]interface{}{
		"content": r.Content,
	}
}

func (r RecipeComment) SelectJoins() []string {
	return []string{}
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

func (r RecipeTag) GetID() uint64 {
	return 0
}

func (r RecipeTag) TableName() string {
	return "recipe_tags"
}

func (r RecipeTag) SelectFields() []string {
	return []string{
		"recipe_tags.recipe_id",
		"recipe_tags.tag_id",
		"tags.tag",
	}
}

func (r RecipeTag) InsertMap() map[string]interface{} {
	return map[string]interface{}{
		"recipe_id": r.RecipeID,
		"tag_id":    r.TagID,
	}
}

func (r RecipeTag) UpdateMap() map[string]interface{} {
	return r.InsertMap()
}

func (r RecipeTag) SelectJoins() []string {
	return []string{"JOIN tags ON tags.id = recipe_tags.tag_id"}
}
