package view

import (
	"time"
)

type Vendor struct {
	ID           uint64  `json:"id"`
	Abbreviation string  `json:"abbreviation"`
	Name         string  `json:"name"`
	Aliases      *string `json:"aliases"`
}

type Flavor struct {
	ID          uint64    `json:"id"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
	Name        string    `json:"name"`
	Aliases     *string   `json:"aliases"`
	RecipeCount uint64    `json:"recipeCount"`
	AvgPercent  uint64    `json:"avgPercent"`
	Vendor      Vendor    `json:"vendor"`
}

type FlavorStash struct {
	ID         uint64  `json:"id"`
	OnHandM    *uint64 `json:"onHandM"`
	DensityM   *uint64 `json:"densityM"`
	Vg         bool    `json:"vg"`
	FlavorID   uint64  `json:"flavorId"`
	OwnerID    uint64  `json:"ownerId"`
	NumRecipes uint64  `json:"numRecipes"`
	Rating     *uint64 `json:"rating"`
}
