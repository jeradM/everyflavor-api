package mapper

import (
	"everyflavor/internal/http/api/v1/view"
	"everyflavor/internal/storage/model"
)

func FlavorFromModel(f *model.Flavor, v *model.Vendor) view.Flavor {
	return view.Flavor{
		ID:          f.ID,
		CreatedAt:   f.CreatedAt,
		UpdatedAt:   f.UpdatedAt,
		Name:        f.Name,
		Aliases:     f.Aliases,
		RecipeCount: f.RecipeCount,
		AvgPercent:  f.AvgPercent,
		Vendor: view.Vendor{
			ID:           v.ID,
			Abbreviation: v.Abbreviation,
			Name:         v.Name,
			Aliases:      v.Aliases,
		},
	}
}

func FlavorToModel(flavor view.Flavor) (model.Flavor, model.Vendor) {
	f := model.Flavor{
		ID:          flavor.ID,
		CreatedAt:   flavor.CreatedAt,
		UpdatedAt:   flavor.UpdatedAt,
		VendorID:    flavor.Vendor.ID,
		Name:        flavor.Name,
		Aliases:     flavor.Aliases,
		RecipeCount: flavor.RecipeCount,
		AvgPercent:  flavor.AvgPercent,
	}

	v := model.Vendor{
		ID:           flavor.Vendor.ID,
		Abbreviation: flavor.Vendor.Abbreviation,
		Name:         flavor.Vendor.Name,
		Aliases:      flavor.Vendor.Aliases,
	}

	return f, v
}

func FlavorStashFromModel(m model.FlavorStash) view.FlavorStash {
	return view.FlavorStash{
		ID:         m.ID,
		OnHandM:    m.OnHandM,
		DensityM:   m.DensityM,
		Vg:         m.Vg,
		FlavorID:   m.FlavorID,
		OwnerID:    m.OwnerID,
		NumRecipes: m.NumRecipes,
		Rating:     m.Rating,
	}
}

func FlavorStashToModel(v view.FlavorStash) model.FlavorStash {
	return model.FlavorStash{
		ID:         v.ID,
		OnHandM:    v.OnHandM,
		DensityM:   v.DensityM,
		Vg:         v.Vg,
		FlavorID:   v.FlavorID,
		OwnerID:    v.OwnerID,
		NumRecipes: v.NumRecipes,
		Rating:     v.Rating,
	}
}
