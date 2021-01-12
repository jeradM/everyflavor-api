package model

import (
	"database/sql"
	"time"
)

type Batch struct {
	ID            uint64       `db:"id"`
	CreatedAt     time.Time    `db:"created_at"`
	UpdatedAt     time.Time    `db:"updated_at"`
	DeletedAt     sql.NullTime `db:"deleted_at"`
	BatchSizeM    uint64       `db:"batch_size_m"`
	BatchStrength uint64       `db:"batch_strength"`
	BatchVgM      uint64       `db:"batch_vg_m"`
	MaxVg         bool         `db:"max_vg"`
	NicStrength   uint64       `db:"nic_strength"`
	NicVgM        uint64       `db:"nic_vg_m"`
	RecipeID      uint64       `db:"recipe_id"`
	OwnerID       uint64       `db:"owner_id"`
	UseNic        bool         `db:"use_nic"`
}

func (b Batch) GetID() uint64 {
	return b.ID
}

func (b Batch) TableName() string {
	return "batches"
}

func (b Batch) SelectFields() []string {
	return []string{
		"batches.id",
		"batches.created_at",
		"batches.updated_at",
		"batches.batch_size_m",
		"batches.batch_strength",
		"batches.batch_vg_m",
		"batches.max_vg",
		"batches.nic_strength",
		"batches.nic_vg_m",
		"batches.recipe_id",
		"batches.owner_id",
		"batches.use_nic",
	}
}

func (b Batch) InsertMap() map[string]interface{} {
	return map[string]interface{}{
		"batch_size_m":   b.BatchSizeM,
		"batch_strength": b.BatchStrength,
		"batch_vg_m":     b.BatchVgM,
		"max_vg":         b.MaxVg,
		"nic_strength":   b.NicStrength,
		"nic_vg_m":       b.NicVgM,
		"recipe_id":      b.RecipeID,
		"owner_id":       b.OwnerID,
		"use_nic":        b.UseNic,
	}
}

func (b Batch) UpdateMap() map[string]interface{} {
	return map[string]interface{}{
		"batch_size_m":   b.BatchSizeM,
		"batch_strength": b.BatchStrength,
		"batch_vg_m":     b.BatchVgM,
		"max_vg":         b.MaxVg,
		"nic_strength":   b.NicStrength,
		"nic_vg_m":       b.NicVgM,
		"use_nic":        b.UseNic,
	}
}

func (b Batch) SelectJoins() []string {
	return []string{}
}

func (b Batch) OwnerField() string {
	return "batches.owner_id"
}

type BatchFlavor struct {
	ID       uint64 `db:"id"`
	PercentM uint64 `db:"percent_m"`
	FlavorID uint64 `db:"flavor_id"`
	Vg       bool   `db:"vg"`
	BatchID  uint64 `db:"batch_id"`
}

func (b BatchFlavor) GetID() uint64 {
	return b.ID
}

func (b BatchFlavor) TableName() string {
	return "batch_flavors"
}

func (b BatchFlavor) SelectFields() []string {
	return []string{
		"batch_flavors.id",
		"batch_flavors.percent_m",
		"batch_flavors.vg",
		"batch_flavors.flavor_id",
		"batch_flavors.batch_id",
	}
}

func (b BatchFlavor) InsertMap() map[string]interface{} {
	return map[string]interface{}{
		"percent_m": b.PercentM,
		"vg":        b.Vg,
		"flavor_id": b.FlavorID,
		"batch_id":  b.BatchID,
	}
}

func (b BatchFlavor) UpdateMap() map[string]interface{} {
	return map[string]interface{}{
		"percent_m": b.PercentM,
		"vg":        b.Vg,
	}
}

func (b BatchFlavor) SelectJoins() []string {
	return []string{}
}
