// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

// Regenerate using `make store_mocks`

package mocks

import (
	model "everyflavor/internal/storage/model"

	mock "github.com/stretchr/testify/mock"

	sqlx "github.com/jmoiron/sqlx"
)

// RecipeCollaboratorStore is an autogenerated mock type for the RecipeCollaboratorStore type
type RecipeCollaboratorStore struct {
	mock.Mock
}

// Replace provides a mock function with given fields: _a0, _a1, _a2
func (_m *RecipeCollaboratorStore) Replace(_a0 uint64, _a1 []model.RecipeCollaborator, _a2 sqlx.Execer) error {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 error
	if rf, ok := ret.Get(0).(func(uint64, []model.RecipeCollaborator, sqlx.Execer) error); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}