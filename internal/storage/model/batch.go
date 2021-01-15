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

type BatchFlavor struct {
	ID       uint64 `db:"id"`
	PercentM uint64 `db:"percent_m"`
	FlavorID uint64 `db:"flavor_id"`
	Vg       bool   `db:"vg"`
	BatchID  uint64 `db:"batch_id"`
}
