package core

import (
	"errors"
	"everyflavor/internal/http/api/v1/view"
	"everyflavor/internal/mapper"
	"everyflavor/internal/storage/model"
)

func (a *App) GetRecipeByID(id uint64) (*view.Recipe, error) {
	r, err := a.Store.Recipe().Get(id)
	if err != nil {
		return nil, err
	}
	return a.loadRecipeRelationships(r)
}

func (a *App) GetRecipeByUUID(uuid string) (*view.Recipe, error) {
	r, err := a.Store.Recipe().ByUUID(uuid)
	if err != nil {
		return nil, err
	}
	return a.loadRecipeRelationships(r)
}

func (a *App) GetRecipesList(p *model.RecipeParams) (view.ListResult, error) {
	var result view.ListResult
	m, count, err := a.Store.Recipe().List(p)
	if err != nil {
		return result, err
	}

	ids := make([]uint64, len(m))
	for idx, r := range m {
		ids[idx] = r.ID
	}
	flavors, err := a.Store.Recipe().
		ListFlavors(model.RecipeFlavorParams{RecipeIDs: ids})
	if err != nil {
		return result, err
	}

	collaborators, err := a.Store.Recipe().
		ListCollaborators(model.RecipeCollaboratorParams{RecipeIDs: ids})
	if err != nil {
		return result, err
	}

	tags, err := a.Store.Recipe().
		ListTags(model.RecipeTagParams{RecipeIDs: ids})
	if err != nil {
		return result, err
	}

	recipes := []view.Recipe{}
	for _, recipe := range m {
		rf := []model.RecipeFlavor{}
		rc := []model.RecipeCollaborator{}
		rt := []model.RecipeTag{}
		for _, flavor := range flavors {
			if flavor.RecipeID == recipe.ID {
				rf = append(rf, flavor)
			}
		}
		for _, collaborator := range collaborators {
			if collaborator.RecipeID == recipe.ID {
				rc = append(rc, collaborator)
			}
		}
		for _, tag := range tags {
			if tag.RecipeID == recipe.ID {
				rt = append(rt, tag)
			}
		}
		r := mapper.RecipeFromModel(recipe, rf, rc, rt)
		recipes = append(recipes, r)
	}
	result.Results = recipes
	result.Count = count
	return result, nil
}

func (a *App) SaveRecipe(r view.Recipe) (view.Recipe, error) {
	recipe, flavors, collaborators, tags := mapper.RecipeToModel(r)
	if r.ID != 0 {
		return r, errors.New("can't insert existing recipe")
	}
	tx, err := a.Store.Connection().Beginx()
	if err != nil {
		return r, err
	}
	err = a.Store.Recipe().Insert(&recipe, tx)
	if err != nil {
		_ = tx.Rollback()
		return r, err
	}
	err = a.Store.Recipe().ReplaceFlavors(recipe.ID, flavors, tx)
	if err != nil {
		_ = tx.Rollback()
		return r, err
	}
	err = a.Store.Recipe().ReplaceCollaborators(recipe.ID, collaborators, tx)
	if err != nil {
		_ = tx.Rollback()
		return r, err
	}
	err = a.Store.Recipe().ReplaceTags(recipe.ID, tags, tx)
	if err != nil {
		_ = tx.Rollback()
		return r, err
	}
	return mapper.RecipeFromModel(recipe, flavors, collaborators, tags), tx.Commit()
}

func (a *App) UpdateRecipe(r view.Recipe) (view.Recipe, error) {
	if r.ID == 0 {
		return r, errors.New("no recipe ID provided")
	}
	tx, err := a.Store.Connection().Beginx()
	if err != nil {
		return r, err
	}
	recipe, flavors, collaborators, tags := mapper.RecipeToModel(r)
	err = a.Store.Recipe().Update(&recipe, tx)
	if err != nil {
		_ = tx.Rollback()
		return r, err
	}
	err = a.Store.Recipe().ReplaceFlavors(recipe.ID, flavors, tx)
	if err != nil {
		_ = tx.Rollback()
		return r, err
	}
	err = a.Store.Recipe().ReplaceCollaborators(recipe.ID, collaborators, tx)
	if err != nil {
		_ = tx.Rollback()
		return r, err
	}
	err = a.Store.Recipe().ReplaceTags(recipe.ID, tags, tx)
	if err != nil {
		_ = tx.Rollback()
		return r, err
	}
	v := mapper.RecipeFromModel(recipe, flavors, collaborators, tags)
	return v, tx.Commit()
}

func (a *App) AddRecipeRating(recipeID uint64, v view.RecipeRating) error {
	rr := model.RecipeRating{
		OwnerID:  v.OwnerID,
		Rating:   v.Rating,
		RecipeID: recipeID,
	}
	return a.Store.Recipe().InsertRating(&rr, nil)
}

func (a *App) PublishRecipe(recipeID uint64) error {
	return a.Store.Recipe().Publish(recipeID, nil)
}

func (a *App) loadRecipeRelationships(r *model.Recipe) (*view.Recipe, error) {
	flavors, err := a.Store.Recipe().
		ListFlavors(model.RecipeFlavorParams{RecipeID: r.ID})
	if err != nil {
		return nil, err
	}

	collaborators, err := a.Store.Recipe().
		ListCollaborators(model.RecipeCollaboratorParams{RecipeID: r.ID})
	if err != nil {
		return nil, err
	}

	tags, err := a.Store.Recipe().
		ListTags(model.RecipeTagParams{RecipeID: r.ID})
	if err != nil {
		return nil, err
	}

	recipe := mapper.RecipeFromModel(*r, flavors, collaborators, tags)
	return &recipe, nil
}
