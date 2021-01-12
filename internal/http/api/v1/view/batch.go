package view

import (
	"time"
)

type Batch struct {
	ID            uint64        `json:"id"`
	CreatedAt     time.Time     `json:"createdAt"`
	UpdatedAt     time.Time     `json:"updatedAt"`
	BatchSizeM    uint64        `json:"batchSizeM"`
	BatchStrength uint64        `json:"batchStrength"`
	BatchVgM      uint64        `json:"batchVgM"`
	MaxVg         bool          `json:"maxVg"`
	NicStrength   uint64        `json:"nicStrength"`
	NicVgM        uint64        `json:"nicVgM"`
	RecipeID      uint64        `json:"recipeId"`
	OwnerID       uint64        `json:"ownerId"`
	UseNic        bool          `json:"useNic"`
	Flavors       []BatchFlavor `json:"flavors"`
}

type BatchFlavor struct {
	PercentM uint64 `json:"percentM"`
	FlavorID uint64 `json:"flavorId"`
	Vg       bool   `json:"vg"`
}
