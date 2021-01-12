package core

import (
	"everyflavor/internal/storage/model"
)

func (a *App) IsPublic(resourceID uint64, t model.PublishableEntity) (bool, error) {
	return a.Store.Auth().IsPublic(resourceID, t)
}

func (a *App) CanViewRecipe(userID, recipeID uint64) (bool, error) {
	var m model.Recipe
	p, err := a.Store.Auth().IsPublic(recipeID, m)
	if err != nil || p {
		return p, err
	}
	o, err := a.Store.Auth().IsOwner(userID, recipeID, m)
	if err != nil || o {
		return o, err
	}
	return a.Store.Auth().IsCollaborator(userID, recipeID, m)
}

func (a *App) CanEditRecipe(userID, recipeID uint64) (bool, error) {
	var m model.Recipe
	o, err := a.Store.Auth().IsOwner(userID, recipeID, m)
	if err != nil || o {
		return o, err
	}
	return a.Store.Auth().IsCollaborator(userID, recipeID, m)
}

func (a *App) CanViewBatch(userID, batchID uint64) (bool, error) {
	var m model.Batch
	return a.Store.Auth().IsOwner(userID, batchID, m)
}

func (a *App) CanEditBatch(userID, batchID uint64) (bool, error) {
	var m model.Batch
	return a.Store.Auth().IsOwner(userID, batchID, m)
}
