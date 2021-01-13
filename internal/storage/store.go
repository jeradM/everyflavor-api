package storage

import (
	"database/sql"
	"database/sql/driver"
	"everyflavor/internal/storage/model"

	"github.com/jmoiron/sqlx"
)

type DB interface {
	Beginx() (*sqlx.Tx, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	Queryx(query string, args ...interface{}) (*sqlx.Rows, error)
	QueryRowx(query string, args ...interface{}) *sqlx.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
}

type Store interface {
	Connection() DB
	DB() DBQueryExecer
	Auth() AuthStore
	Batch() BatchStore
	Flavor() FlavorStore
	Recipe() RecipeStore
	Stash() FlavorStashStore
	Tag() TagStore
	User() UserStore
	Vendor() VendorStore
}

type AuthStore interface {
	IsPublic(uint64, model.PublishableEntity) (bool, error)
	IsOwner(uint64, uint64, model.OwnedEntity) (bool, error)
	IsCollaborator(uint64, uint64, model.SharedEntity) (bool, error)
}

type UserStore interface {
	Get(uint64) (*model.User, error)
	List() ([]model.User, error)
	Insert(*model.User, sqlx.Execer) error
	Update(*model.User, sqlx.Execer) error
	UpdatePassword(uint64, string) error
	FindByUsername(string) (*model.User, error)
	FindAllByUsernameLike(string) ([]model.User, error)
	GetStats(uint64) (*model.UserStats, error)
	ListRoles([]uint64) ([]model.UserRole, error)
	UsernameExists(string) bool
	EmailExists(string) bool
}

type VendorStore interface {
	Get(uint64) (*model.Vendor, error)
	List() ([]model.Vendor, uint64, error)
	Insert(*model.Vendor, sqlx.Execer) error
	Update(*model.Vendor, sqlx.Execer) error
	FindByAbbreviation(string) (*model.Vendor, error)
	FindByName(string) (*model.Vendor, error)
}

type FlavorStore interface {
	Get(uint64) (model.Flavor, error)
	List() ([]model.Flavor, uint64, error)
	Insert(model.Flavor, sqlx.Execer) error
	Update(model.Flavor, sqlx.Execer) error
}

type RecipeStore interface {
	Get(uint64) (*model.Recipe, error)
	ByUUID(string) (*model.Recipe, error)
	List(*model.RecipeParams) ([]model.Recipe, uint64, error)
	Insert(*model.Recipe, sqlx.Ext) error
	Update(*model.Recipe, sqlx.Ext) error
	AddComment(*model.RecipeComment) error
	ListFlavors(model.RecipeFlavorParams) ([]model.RecipeFlavor, error)
	ReplaceFlavors(recipeID uint64, flavors []model.RecipeFlavor, tx sqlx.Execer) error
	ListCollaborators(model.RecipeCollaboratorParams) ([]model.RecipeCollaborator, error)
	ReplaceCollaborators(recipeID uint64, collabs []model.RecipeCollaborator, tx sqlx.Execer) error
	ListTags(model.RecipeTagParams) ([]model.RecipeTag, error)
	ReplaceTags(recipeID uint64, tags []model.RecipeTag, tx sqlx.Execer) error
	InsertRating(*model.RecipeRating, sqlx.Execer) error
}

type RecipeFlavorStore interface {
	Get(uint64) (*model.RecipeFlavor, error)
	List(model.RecipeFlavorParams) ([]model.RecipeFlavor, error)
	Insert(*model.RecipeFlavor, sqlx.Execer) error
	Update(*model.RecipeFlavor, sqlx.Execer) error
	DeleteAllByRecipeID(uint64, sqlx.Execer) (int64, error)
	Replace(uint64, []model.RecipeFlavor, sqlx.Execer) error
}

type BatchStore interface {
	Get(uint64) (model.Batch, error)
	List(uint64) ([]model.Batch, error)
	Insert(*model.Batch, sqlx.Execer) error
	Update(*model.Batch, sqlx.Execer) error
	ListFlavorsForBatches([]uint64) ([]model.BatchFlavor, error)
	InsertFlavor(*model.BatchFlavor, sqlx.Execer) error
}

type RecipeCollaboratorStore interface {
	Replace(uint64, []model.RecipeCollaborator, sqlx.Execer) error
}

type RecipeTagStore interface {
	List(model.RecipeTagParams) ([]model.RecipeTag, error)
	Insert(*model.RecipeTag, sqlx.Execer) error
	DeleteAllByRecipeID(uint64, sqlx.Execer) (int64, error)
	Replace(uint64, []model.RecipeTag, sqlx.Execer) error
}

type TagStore interface {
	List() ([]model.Tag, error)
}

type FlavorStashStore interface {
	List(uint64) ([]model.FlavorStash, error)
	Insert(*model.FlavorStash, sqlx.Execer) error
	Delete(*model.FlavorStash, sqlx.Execer) (int64, error)
}

type DBQueryExecer interface {
	Enable(bool)
	Get(interface{}, string, []interface{}) error
	GetWithTX(sqlx.Queryer, interface{}, string, []interface{}) error
	Select(interface{}, string, []interface{}) error
	SelectWithTX(sqlx.Queryer, interface{}, string, []interface{}) error
	Exec(string, []interface{}) (driver.Result, error)
	ExecWithTX(sqlx.Execer, string, []interface{}) (driver.Result, error)
}
