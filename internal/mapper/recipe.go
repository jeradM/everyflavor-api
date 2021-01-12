package mapper

import (
	"everyflavor/internal/http/api/v1/view"
	"everyflavor/internal/storage/model"
)

func RecipeFromModel(
	recipe model.Recipe,
	flavors []model.RecipeFlavor,
	collaborators []model.RecipeCollaborator,
	tags []model.RecipeTag,
) (r view.Recipe) {
	r = view.Recipe{
		ID:            recipe.ID,
		CreatedAt:     recipe.CreatedAt,
		UpdatedAt:     recipe.UpdatedAt,
		Current:       recipe.Current,
		Description:   recipe.Description,
		Public:        recipe.Public,
		RemixOfID:     recipe.RemixOfID,
		RemixOfTitle:  recipe.RemixOfTitle,
		Snv:           recipe.Snv,
		SteepDays:     recipe.SteepDays,
		TempF:         recipe.TempF,
		Title:         recipe.Title,
		UUID:          recipe.UUID,
		Version:       recipe.Version,
		VgPercentM:    recipe.VgPercentM,
		Wip:           recipe.Wip,
		AvgRating:     recipe.AvgRating,
		OwnerID:       recipe.OwnerID,
		OwnerUsername: recipe.OwnerUsername,
		Collaborators: []view.Collaborator{},
		Flavors:       []view.RecipeFlavor{},
		Tags:          []view.RecipeTag{},
	}
	//r.RemixOf = &view.Recipe{ID: remixOf.ID, Title: remixOf.Title}
	for _, u := range collaborators {
		r.Collaborators = append(r.Collaborators, view.Collaborator{u.UserID, u.Username})
	}
	for _, f := range flavors {
		r.Flavors = append(r.Flavors, view.RecipeFlavor{ID: f.ID, FlavorID: f.FlavorID, PercentM: f.PercentM})
	}
	for _, t := range tags {
		r.Tags = append(r.Tags, view.RecipeTag{ID: t.TagID, Tag: t.Tag})
	}
	return r
}

func RecipeToModel(recipe view.Recipe) (
	model.Recipe, []model.RecipeFlavor, []model.RecipeCollaborator, []model.RecipeTag) {
	r := model.Recipe{
		ID:            recipe.ID,
		CreatedAt:     recipe.CreatedAt,
		UpdatedAt:     recipe.UpdatedAt,
		Current:       recipe.Current,
		Description:   recipe.Description,
		Public:        recipe.Public,
		RemixOfID:     recipe.RemixOfID,
		RemixOfTitle:  recipe.RemixOfTitle,
		Snv:           recipe.Snv,
		SteepDays:     recipe.SteepDays,
		TempF:         recipe.TempF,
		Title:         recipe.Title,
		UUID:          recipe.UUID,
		Version:       recipe.Version,
		VgPercentM:    recipe.VgPercentM,
		Wip:           recipe.Wip,
		AvgRating:     recipe.AvgRating,
		OwnerID:       recipe.OwnerID,
		OwnerUsername: recipe.OwnerUsername,
	}
	f := []model.RecipeFlavor{}
	for _, flv := range recipe.Flavors {
		f = append(f, model.RecipeFlavor{
			ID:       flv.ID,
			PercentM: flv.PercentM,
			FlavorID: flv.FlavorID,
			RecipeID: r.ID,
		})
	}
	c := []model.RecipeCollaborator{}
	for _, col := range recipe.Collaborators {
		c = append(c, model.RecipeCollaborator{
			UserID:   col.ID,
			Username: col.Username,
			RecipeID: r.ID,
		})
	}
	t := []model.RecipeTag{}
	for _, tag := range recipe.Tags {
		t = append(t, model.RecipeTag{
			TagID:    tag.ID,
			Tag:      tag.Tag,
			RecipeID: r.ID,
		})
	}
	return r, f, c, t
}
