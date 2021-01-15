package core

func (a *App) IsPublic(resourceID uint64, table string) (bool, error) {
	return a.Store.Auth().IsPublic(resourceID, table)
}

func (a *App) CanViewRecipe(userID, recipeID uint64) (bool, error) {
	p, err := a.Store.Auth().IsPublic(recipeID, "recipes")
	if err != nil {
		return false, err
	}
	if p {
		return true, nil
	}
	o, err := a.Store.Auth().IsOwner(userID, recipeID, "recipes")
	if err != nil {
		return false, err
	}
	if o {
		return true, nil
	}
	return a.Store.Auth().IsCollaborator(userID, recipeID, "recipe")
}

func (a *App) CanEditRecipe(userID, recipeID uint64) (bool, error) {
	o, err := a.Store.Auth().IsOwner(userID, recipeID, "recipes")
	if err != nil {
		return false, err
	}
	if o {
		return true, nil
	}
	return a.Store.Auth().IsCollaborator(userID, recipeID, "recipe")
}

func (a *App) CanViewBatch(userID, batchID uint64) (bool, error) {
	v, err := a.Store.Auth().IsOwner(userID, batchID, "batches")
	if err != nil {
		return false, err
	}
	return v, nil
}

func (a *App) CanEditBatch(userID, batchID uint64) (bool, error) {
	e, err := a.Store.Auth().IsOwner(userID, batchID, "batches")
	if err != nil {
		return false, err
	}
	return e, nil
}
