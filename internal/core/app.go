//go:generate mockery --all --output/mocks
package core

import (
	"everyflavor/internal/http/api/v1/view"
	"everyflavor/internal/storage"
	"everyflavor/internal/storage/model"
)

type AuthService interface {
	IsPublic(uint64, model.PublishableEntity) (bool, error)
	CanViewRecipe(userID, recipeID uint64) (bool, error)
	CanEditRecipe(userID, recipeID uint64) (bool, error)
	CanViewBatch(userID, batchID uint64) (bool, error)
	CanEditBatch(userID, batchID uint64) (bool, error)
}

type BatchService interface {
	GetBatch(uint64) (*view.Batch, error)
	GetBatchesForUser(ownerId uint64) ([]model.Batch, error)
	SaveBatch(b view.Batch) (view.Batch, error)
}

type FlavorService interface {
	GetFlavorByID(uint64) (view.Flavor, error)
	GetFlavorsList() (view.ListResult, error)
	SaveFlavor(view.Flavor) error
	GetStashForUser(uint64) ([]view.FlavorStash, error)
	SaveStash(view.FlavorStash) error
}

type RecipeService interface {
	GetRecipeByID(uint64) (*view.Recipe, error)
	GetRecipeByUUID(string) (*view.Recipe, error)
	GetRecipesList(*model.RecipeParams) (view.ListResult, error)
	SaveRecipe(view.Recipe) (view.Recipe, error)
	UpdateRecipe(view.Recipe) (view.Recipe, error)
	AddRecipeRating(uint64, view.RecipeRating) error
}

type TagService interface {
	ListTags() ([]view.Tag, error)
}

type UserService interface {
	GetUserByID(id uint64) (view.User, error)
	GetUserList() ([]model.User, error)
	SaveUser(u view.User) error
	UpdateUserPassword(id uint64, pw string) error
	GetUserByUsername(username string) (view.User, error)
	SearchUsersByUsername(username string) ([]view.User, error)
	GetUserStatsByID(uint64) (*model.UserStats, error)
	GetRolesForUsers([]uint64) ([]model.UserRole, error)
}

type VendorService interface {
	GetVendorByID(id uint64) (*model.Vendor, error)
	GetVendorsList() ([]model.Vendor, error)
	SaveVendor(v *model.Vendor) error
	GetVendorByAbbreviation(abbrev string) (*model.Vendor, error)
	GetVendorByName(name string) (*model.Vendor, error)
}

type Application interface {
	AuthService
	BatchService
	FlavorService
	RecipeService
	TagService
	UserService
	VendorService
}

type App struct {
	Config AppConfig
	Store  storage.Store
}

func NewApp(c AppConfig, r storage.Store) *App {
	return &App{c, r}
}
