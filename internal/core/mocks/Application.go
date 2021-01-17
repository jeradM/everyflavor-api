// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

// Regenerate with `make core_mocks`

package mocks

import (
	model "everyflavor/internal/storage/model"

	mock "github.com/stretchr/testify/mock"

	view "everyflavor/internal/http/api/v1/view"
)

// Application is an autogenerated mock type for the Application type
type Application struct {
	mock.Mock
}

// AddRecipeRating provides a mock function with given fields: _a0, _a1
func (_m *Application) AddRecipeRating(_a0 uint64, _a1 view.RecipeRating) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(uint64, view.RecipeRating) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CanEditBatch provides a mock function with given fields: userID, batchID
func (_m *Application) CanEditBatch(userID uint64, batchID uint64) (bool, error) {
	ret := _m.Called(userID, batchID)

	var r0 bool
	if rf, ok := ret.Get(0).(func(uint64, uint64) bool); ok {
		r0 = rf(userID, batchID)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint64, uint64) error); ok {
		r1 = rf(userID, batchID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CanEditRecipe provides a mock function with given fields: userID, recipeID
func (_m *Application) CanEditRecipe(userID uint64, recipeID uint64) (bool, error) {
	ret := _m.Called(userID, recipeID)

	var r0 bool
	if rf, ok := ret.Get(0).(func(uint64, uint64) bool); ok {
		r0 = rf(userID, recipeID)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint64, uint64) error); ok {
		r1 = rf(userID, recipeID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CanViewBatch provides a mock function with given fields: userID, batchID
func (_m *Application) CanViewBatch(userID uint64, batchID uint64) (bool, error) {
	ret := _m.Called(userID, batchID)

	var r0 bool
	if rf, ok := ret.Get(0).(func(uint64, uint64) bool); ok {
		r0 = rf(userID, batchID)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint64, uint64) error); ok {
		r1 = rf(userID, batchID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CanViewRecipe provides a mock function with given fields: userID, recipeID
func (_m *Application) CanViewRecipe(userID uint64, recipeID uint64) (bool, error) {
	ret := _m.Called(userID, recipeID)

	var r0 bool
	if rf, ok := ret.Get(0).(func(uint64, uint64) bool); ok {
		r0 = rf(userID, recipeID)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint64, uint64) error); ok {
		r1 = rf(userID, recipeID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// EmailExists provides a mock function with given fields: _a0
func (_m *Application) EmailExists(_a0 string) bool {
	ret := _m.Called(_a0)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// GetBatch provides a mock function with given fields: _a0
func (_m *Application) GetBatch(_a0 uint64) (*view.Batch, error) {
	ret := _m.Called(_a0)

	var r0 *view.Batch
	if rf, ok := ret.Get(0).(func(uint64) *view.Batch); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*view.Batch)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint64) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetBatchesForUser provides a mock function with given fields: ownerId
func (_m *Application) GetBatchesForUser(ownerId uint64) ([]model.Batch, error) {
	ret := _m.Called(ownerId)

	var r0 []model.Batch
	if rf, ok := ret.Get(0).(func(uint64) []model.Batch); ok {
		r0 = rf(ownerId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.Batch)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint64) error); ok {
		r1 = rf(ownerId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetFlavorByID provides a mock function with given fields: _a0
func (_m *Application) GetFlavorByID(_a0 uint64) (view.Flavor, error) {
	ret := _m.Called(_a0)

	var r0 view.Flavor
	if rf, ok := ret.Get(0).(func(uint64) view.Flavor); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(view.Flavor)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint64) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetFlavorsList provides a mock function with given fields:
func (_m *Application) GetFlavorsList() (view.ListResult, error) {
	ret := _m.Called()

	var r0 view.ListResult
	if rf, ok := ret.Get(0).(func() view.ListResult); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(view.ListResult)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetRecipeByID provides a mock function with given fields: _a0
func (_m *Application) GetRecipeByID(_a0 uint64) (*view.Recipe, error) {
	ret := _m.Called(_a0)

	var r0 *view.Recipe
	if rf, ok := ret.Get(0).(func(uint64) *view.Recipe); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*view.Recipe)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint64) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetRecipeByUUID provides a mock function with given fields: _a0
func (_m *Application) GetRecipeByUUID(_a0 string) (*view.Recipe, error) {
	ret := _m.Called(_a0)

	var r0 *view.Recipe
	if rf, ok := ret.Get(0).(func(string) *view.Recipe); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*view.Recipe)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetRecipesList provides a mock function with given fields: _a0
func (_m *Application) GetRecipesList(_a0 *model.RecipeParams) (view.ListResult, error) {
	ret := _m.Called(_a0)

	var r0 view.ListResult
	if rf, ok := ret.Get(0).(func(*model.RecipeParams) view.ListResult); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(view.ListResult)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*model.RecipeParams) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetRolesForUsers provides a mock function with given fields: _a0
func (_m *Application) GetRolesForUsers(_a0 []uint64) ([]model.UserRole, error) {
	ret := _m.Called(_a0)

	var r0 []model.UserRole
	if rf, ok := ret.Get(0).(func([]uint64) []model.UserRole); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.UserRole)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]uint64) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetStashForUser provides a mock function with given fields: _a0
func (_m *Application) GetStashForUser(_a0 uint64) ([]view.FlavorStash, error) {
	ret := _m.Called(_a0)

	var r0 []view.FlavorStash
	if rf, ok := ret.Get(0).(func(uint64) []view.FlavorStash); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]view.FlavorStash)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint64) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserByID provides a mock function with given fields: id
func (_m *Application) GetUserByID(id uint64) (view.User, error) {
	ret := _m.Called(id)

	var r0 view.User
	if rf, ok := ret.Get(0).(func(uint64) view.User); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(view.User)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint64) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserByUsername provides a mock function with given fields: username
func (_m *Application) GetUserByUsername(username string) (view.User, error) {
	ret := _m.Called(username)

	var r0 view.User
	if rf, ok := ret.Get(0).(func(string) view.User); ok {
		r0 = rf(username)
	} else {
		r0 = ret.Get(0).(view.User)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(username)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserList provides a mock function with given fields:
func (_m *Application) GetUserList() ([]view.User, error) {
	ret := _m.Called()

	var r0 []view.User
	if rf, ok := ret.Get(0).(func() []view.User); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]view.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserStatsByID provides a mock function with given fields: _a0
func (_m *Application) GetUserStatsByID(_a0 uint64) (*model.UserStats, error) {
	ret := _m.Called(_a0)

	var r0 *model.UserStats
	if rf, ok := ret.Get(0).(func(uint64) *model.UserStats); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.UserStats)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint64) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetVendorByAbbreviation provides a mock function with given fields: abbrev
func (_m *Application) GetVendorByAbbreviation(abbrev string) (*model.Vendor, error) {
	ret := _m.Called(abbrev)

	var r0 *model.Vendor
	if rf, ok := ret.Get(0).(func(string) *model.Vendor); ok {
		r0 = rf(abbrev)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Vendor)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(abbrev)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetVendorByID provides a mock function with given fields: id
func (_m *Application) GetVendorByID(id uint64) (*model.Vendor, error) {
	ret := _m.Called(id)

	var r0 *model.Vendor
	if rf, ok := ret.Get(0).(func(uint64) *model.Vendor); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Vendor)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint64) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetVendorByName provides a mock function with given fields: name
func (_m *Application) GetVendorByName(name string) (*model.Vendor, error) {
	ret := _m.Called(name)

	var r0 *model.Vendor
	if rf, ok := ret.Get(0).(func(string) *model.Vendor); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Vendor)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetVendorsList provides a mock function with given fields:
func (_m *Application) GetVendorsList() ([]model.Vendor, error) {
	ret := _m.Called()

	var r0 []model.Vendor
	if rf, ok := ret.Get(0).(func() []model.Vendor); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.Vendor)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IsPublic provides a mock function with given fields: _a0, _a1
func (_m *Application) IsPublic(_a0 uint64, _a1 string) (bool, error) {
	ret := _m.Called(_a0, _a1)

	var r0 bool
	if rf, ok := ret.Get(0).(func(uint64, string) bool); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint64, string) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListTags provides a mock function with given fields:
func (_m *Application) ListTags() ([]view.Tag, error) {
	ret := _m.Called()

	var r0 []view.Tag
	if rf, ok := ret.Get(0).(func() []view.Tag); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]view.Tag)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PublishRecipe provides a mock function with given fields: _a0
func (_m *Application) PublishRecipe(_a0 uint64) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(uint64) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SaveBatch provides a mock function with given fields: b
func (_m *Application) SaveBatch(b view.Batch) (view.Batch, error) {
	ret := _m.Called(b)

	var r0 view.Batch
	if rf, ok := ret.Get(0).(func(view.Batch) view.Batch); ok {
		r0 = rf(b)
	} else {
		r0 = ret.Get(0).(view.Batch)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(view.Batch) error); ok {
		r1 = rf(b)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SaveFlavor provides a mock function with given fields: _a0
func (_m *Application) SaveFlavor(_a0 view.Flavor) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(view.Flavor) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SaveRecipe provides a mock function with given fields: _a0
func (_m *Application) SaveRecipe(_a0 view.Recipe) (view.Recipe, error) {
	ret := _m.Called(_a0)

	var r0 view.Recipe
	if rf, ok := ret.Get(0).(func(view.Recipe) view.Recipe); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(view.Recipe)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(view.Recipe) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SaveStash provides a mock function with given fields: _a0
func (_m *Application) SaveStash(_a0 view.FlavorStash) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(view.FlavorStash) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SaveUser provides a mock function with given fields: u
func (_m *Application) SaveUser(u view.User) error {
	ret := _m.Called(u)

	var r0 error
	if rf, ok := ret.Get(0).(func(view.User) error); ok {
		r0 = rf(u)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SaveVendor provides a mock function with given fields: v
func (_m *Application) SaveVendor(v *model.Vendor) error {
	ret := _m.Called(v)

	var r0 error
	if rf, ok := ret.Get(0).(func(*model.Vendor) error); ok {
		r0 = rf(v)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SearchUsersByUsername provides a mock function with given fields: username
func (_m *Application) SearchUsersByUsername(username string) ([]view.User, error) {
	ret := _m.Called(username)

	var r0 []view.User
	if rf, ok := ret.Get(0).(func(string) []view.User); ok {
		r0 = rf(username)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]view.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(username)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateRecipe provides a mock function with given fields: _a0
func (_m *Application) UpdateRecipe(_a0 view.Recipe) (view.Recipe, error) {
	ret := _m.Called(_a0)

	var r0 view.Recipe
	if rf, ok := ret.Get(0).(func(view.Recipe) view.Recipe); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(view.Recipe)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(view.Recipe) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateUserPassword provides a mock function with given fields: id, pw
func (_m *Application) UpdateUserPassword(id uint64, pw string) error {
	ret := _m.Called(id, pw)

	var r0 error
	if rf, ok := ret.Get(0).(func(uint64, string) error); ok {
		r0 = rf(id, pw)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UsernameExists provides a mock function with given fields: _a0
func (_m *Application) UsernameExists(_a0 string) bool {
	ret := _m.Called(_a0)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}
