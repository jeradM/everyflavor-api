package mapper

import (
	"database/sql"
	"everyflavor/internal/http/api/v1/view"
	"everyflavor/internal/storage/model"
)

func BatchFromModel(batch model.Batch, flavors []model.BatchFlavor) view.Batch {
	b := view.Batch{
		ID:            batch.ID,
		CreatedAt:     batch.CreatedAt,
		UpdatedAt:     batch.UpdatedAt,
		BatchSizeM:    batch.BatchSizeM,
		BatchStrength: batch.BatchStrength,
		BatchVgM:      batch.BatchVgM,
		MaxVg:         batch.MaxVg,
		NicStrength:   batch.NicStrength,
		NicVgM:        batch.NicVgM,
		RecipeID:      batch.RecipeID,
		OwnerID:       batch.OwnerID,
		UseNic:        batch.UseNic,
	}
	for _, f := range flavors {
		b.Flavors = append(b.Flavors, view.BatchFlavor{
			PercentM: f.PercentM,
			FlavorID: f.FlavorID,
			Vg:       f.Vg,
		})
	}
	return b
}

func BatchToModel(batch view.Batch) (model.Batch, []model.BatchFlavor) {
	b := model.Batch{
		ID:            batch.ID,
		CreatedAt:     batch.CreatedAt,
		UpdatedAt:     batch.UpdatedAt,
		DeletedAt:     sql.NullTime{},
		BatchSizeM:    batch.BatchSizeM,
		BatchStrength: batch.BatchStrength,
		BatchVgM:      batch.BatchVgM,
		MaxVg:         batch.MaxVg,
		NicStrength:   batch.NicStrength,
		NicVgM:        batch.NicVgM,
		RecipeID:      batch.RecipeID,
		OwnerID:       batch.OwnerID,
		UseNic:        batch.UseNic,
	}
	f := make([]model.BatchFlavor, len(batch.Flavors))
	for idx, flv := range batch.Flavors {
		f[idx] = model.BatchFlavor{
			PercentM: flv.PercentM,
			FlavorID: flv.FlavorID,
			Vg:       flv.Vg,
			BatchID:  batch.ID,
		}
	}
	return b, f
}
