// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

// Regenerate using `make store_mocks`

package mocks

import (
	model "everyflavor/internal/storage/model"

	mock "github.com/stretchr/testify/mock"

	sqlx "github.com/jmoiron/sqlx"
)

// FlavorStashStore is an autogenerated mock type for the FlavorStashStore type
type FlavorStashStore struct {
	mock.Mock
}

// Delete provides a mock function with given fields: _a0, _a1
func (_m *FlavorStashStore) Delete(_a0 *model.FlavorStash, _a1 sqlx.Execer) (int64, error) {
	ret := _m.Called(_a0, _a1)

	var r0 int64
	if rf, ok := ret.Get(0).(func(*model.FlavorStash, sqlx.Execer) int64); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*model.FlavorStash, sqlx.Execer) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Insert provides a mock function with given fields: _a0, _a1
func (_m *FlavorStashStore) Insert(_a0 *model.FlavorStash, _a1 sqlx.Execer) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(*model.FlavorStash, sqlx.Execer) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// List provides a mock function with given fields: _a0
func (_m *FlavorStashStore) List(_a0 uint64) ([]model.FlavorStash, error) {
	ret := _m.Called(_a0)

	var r0 []model.FlavorStash
	if rf, ok := ret.Get(0).(func(uint64) []model.FlavorStash); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.FlavorStash)
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
