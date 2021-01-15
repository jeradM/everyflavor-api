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

type FlavorRating struct {
	ID        uint64        `db:"id"`
	CreatedAt *time.Time    `db:"created_at"`
	UpdatedAt *time.Time    `db:"updated_at"`
	DeletedAt *sql.NullTime `db:"deleted_at"`
	FlavorID  uint64        `db:"flavor_id"`
	Rating    uint64        `db:"rating"`
	OwnerID   uint64        `db:"owner_id"`
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
