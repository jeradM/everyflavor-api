package model

import (
	"database/sql"
	"time"
)

type Vendor struct {
	ID           uint64       `db:"id"`
	CreatedAt    time.Time    `db:"created_at"`
	UpdatedAt    time.Time    `db:"updated_at"`
	DeletedAt    sql.NullTime `db:"deleted_at"`
	Abbreviation string       `db:"abbreviation"`
	Name         string       `db:"name"`
	Aliases      *string      `db:"aliases"`
}

func (v Vendor) GetID() uint64 {
	return v.ID
}

func (Vendor) TableName() string {
	return "vendors"
}

func (Vendor) SelectFields() []string {
	return []string{"vendors.id", "vendors.abbreviation", "vendors.name", "vendors.aliases"}
}

func (v Vendor) InsertMap() map[string]interface{} {
	return map[string]interface{}{
		"vendors.abbreviation": v.Abbreviation,
		"vendors.name":         v.Name,
		"vendors.aliases":      v.Aliases,
	}
}

func (v Vendor) UpdateMap() map[string]interface{} {
	return v.InsertMap()
}

func (Vendor) SelectJoins() []string {
	return []string{}
}

type Flavor struct {
	ID          uint64       `db:"id"`
	CreatedAt   time.Time    `db:"created_at"`
	UpdatedAt   time.Time    `db:"updated_at"`
	DeletedAt   sql.NullTime `db:"deleted_at"`
	VendorID    uint64       `db:"vendor_id"`
	Name        string       `db:"name"`
	Aliases     *string      `db:"aliases"`
	RecipeCount uint64       `db:"recipe_count"`
	AvgPercent  uint64       `db:"avg_percent"`
}

func (f Flavor) GetID() uint64 {
	return f.ID
}

func (f Flavor) TableName() string {
	return "flavors"
}

func (f Flavor) SelectFields() []string {
	return []string{
		"flavors.id",
		"flavors.created_at",
		"flavors.updated_at",
		"flavors.vendor_id",
		"flavors.name",
		"flavors.aliases",
		"count(recipes.id) recipe_count",
		"IFNULL(CAST(avg(recipe_flavors.percent_m) as UNSIGNED), 0) avg_percent",
	}
}

func (f Flavor) SelectJoins() []string {
	return []string{
		"LEFT JOIN recipe_flavors ON recipe_flavors.flavor_id = flavors.id",
		"LEFT JOIN recipes ON recipes.id = recipe_flavors.recipe_id",
	}
}

func (f Flavor) InsertMap() map[string]interface{} {
	return map[string]interface{}{
		"flavors.vendor_id": f.VendorID,
		"flavors.name":      f.Name,
		"flavors.aliases":   f.Aliases,
	}
}

func (f Flavor) UpdateMap() map[string]interface{} {
	return f.InsertMap()
}

type FlavorRating struct {
	ID        uint64        `db:"id"`
	CreatedAt *time.Time    `db:"created_at"`
	UpdatedAt *time.Time    `db:"updated_at"`
	DeletedAt *sql.NullTime `db:"deleted_at"`
	FlavorID  uint64        `db:"flavor_id"`
	Rating    uint64        `db:"rating"`
	OwnerID   uint64        `db:"owner_id"`
}

func (f *FlavorRating) GetID() uint64 {
	return f.ID
}

func (f *FlavorRating) TableName() string {
	return "flavor_ratings"
}

func (f *FlavorRating) SelectFields() []string {
	return []string{
		"flavor_ratings.id",
		"flavor_ratings.created_at",
		"flavor_ratings.updated_at",
		"flavor_ratings.flavor_id",
		"flavor_ratings.rating",
		"flavor_ratings.owner_id",
	}
}

func (f *FlavorRating) SelectJoins() []string {
	return []string{}
}

func (f *FlavorRating) InsertMap() map[string]interface{} {
	return map[string]interface{}{
		"flavor_ratings.flavor_id": f.FlavorID,
		"flavor_ratings.rating":    f.Rating,
		"flavor_ratings.owner_id":  f.OwnerID,
	}
}

func (f *FlavorRating) UpdateMap() map[string]interface{} {
	return f.InsertMap()
}

type FlavorReview struct {
	ID        uint64        `db:"id"`
	CreatedAt *time.Time    `db:"created_at"`
	UpdatedAt *time.Time    `db:"updated_at"`
	DeletedAt *sql.NullTime `db:"deleted_at"`
	Rating    FlavorRating  `db:"-"`
	RatingID  uint64        `db:"ratingId"`
	Content   string        `db:"content"`
}

func (f *FlavorReview) GetID() uint64 {
	return f.ID
}

func (f *FlavorReview) TableName() string {
	return "flavor_reviews"
}

func (f *FlavorReview) SelectFields() []string {
	return []string{
		"flavor_reviews.id",
		"flavor_reviews.created_at",
		"flavor_reviews.updated_at",
		"flavor_reviews.rating_id",
		"flavor_reviews.content",
	}
}

func (f *FlavorReview) InsertMap() map[string]interface{} {
	return map[string]interface{}{
		"rating_id": f.RatingID,
		"content":   f.Content,
	}
}

func (f *FlavorReview) UpdateMap() map[string]interface{} {
	return f.InsertMap()
}

func (f *FlavorReview) SelectJoins() []string {
	return []string{}
}

type FlavorStash struct {
	ID         uint64        `db:"id"`
	CreatedAt  *time.Time    `db:"created_at"`
	UpdatedAt  *time.Time    `db:"updated_at"`
	DeletedAt  *sql.NullTime `db:"deleted_at"`
	OnHandM    *uint64       `db:"on_hand_m"`
	DensityM   *uint64       `db:"density_m"`
	Vg         bool          `db:"vg"`
	FlavorID   uint64        `db:"flavor_id"`
	OwnerID    uint64        `db:"owner_id"`
	NumRecipes uint64        `db:"num_recipes"`
	Rating     *uint64       `db:"rating"`
}

func (f *FlavorStash) GetID() uint64 {
	return f.ID
}

func (f *FlavorStash) SetID(id uint64) {
	f.ID = id
}

func (f *FlavorStash) TableName() string {
	return "flavor_stashes"
}

func (f *FlavorStash) SelectFields() []string {
	return []string{
		"flavor_stashes.id",
		"flavor_stashes.created_at",
		"flavor_stashes.updated_at",
		"flavor_stashes.on_hand_m",
		"flavor_stashes.density_m",
		"flavor_stashes.vg",
		"flavor_stashes.flavor_id",
		"flavor_stashes.owner_id",
	}
}

func (f *FlavorStash) InsertMap() map[string]interface{} {
	return map[string]interface{}{
		"on_hand_m": f.OnHandM,
		"density_m": f.DensityM,
		"vg":        f.Vg,
		"flavor_id": f.FlavorID,
		"owner_id":  f.OwnerID,
	}
}

func (f *FlavorStash) UpdateMap() map[string]interface{} {
	return f.InsertMap()
}

func (f *FlavorStash) SelectJoins() []string {
	return []string{
		"LEFT JOIN recipe_flavors ON recipe_flavors.flavor_id = flavor_stashes.flavor_id",
		"LEFT JOIN recipes ON recipes.id = recipe_flavors.recipe_id AND recipes.owner_id = flavor_stashes.owner_id",
	}
}
