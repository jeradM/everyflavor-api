package view

import (
	"time"
)

type Recipe struct {
	ID            uint64         `json:"id"`
	CreatedAt     *time.Time     `json:"createdAt"`
	UpdatedAt     *time.Time     `json:"updatedAt"`
	Current       bool           `json:"current"`
	Description   string         `json:"description"`
	Public        bool           `json:"public"`
	PublishedAt   *time.Time     `json:"published_at"`
	RemixOfID     *uint64        `json:"remixOfId"`
	RemixOfTitle  *string        `json:"remixOfTitle"`
	Snv           bool           `json:"snv"`
	SteepDays     uint64         `json:"steepDays"`
	TempF         uint64         `json:"temp"`
	Title         string         `json:"title"`
	UUID          string         `json:"uuid"`
	Version       uint64         `json:"version"`
	VgPercentM    uint64         `json:"vgPercentM"`
	Wip           bool           `json:"wip"`
	AvgRating     uint64         `json:"avgRating"`
	OwnerID       uint64         `json:"ownerId"`
	OwnerUsername string         `json:"ownerUsername"`
	Collaborators []Collaborator `json:"collaborators"`
	Flavors       []RecipeFlavor `json:"flavors"`
	Ratings       []RecipeRating `json:"ratings"`
	Tags          []RecipeTag    `json:"tags"`
}

type Collaborator struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`
}

type RecipeFlavor struct {
	ID       uint64 `json:"id"`
	FlavorID uint64 `json:"flavorId"`
	PercentM uint64 `json:"percentM"`
	RecipeID uint64 `json:"recipeId"`
}

type RecipeTag struct {
	ID  uint64 `json:"id"`
	Tag string `json:"tag"`
}

type RecipeRating struct {
	ID      uint64 `json:"id"`
	OwnerID uint64 `json:"owner_id"`
	Rating  uint64 `json:"rating"`
}
