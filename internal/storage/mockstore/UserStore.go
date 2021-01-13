// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

// Regenerate using `make store_mocks`

package mocks

import (
	model "everyflavor/internal/storage/model"

	mock "github.com/stretchr/testify/mock"

	sqlx "github.com/jmoiron/sqlx"
)

// UserStore is an autogenerated mock type for the UserStore type
type UserStore struct {
	mock.Mock
}

// EmailExists provides a mock function with given fields: _a0
func (_m *UserStore) EmailExists(_a0 string) bool {
	ret := _m.Called(_a0)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// FindAllByUsernameLike provides a mock function with given fields: _a0
func (_m *UserStore) FindAllByUsernameLike(_a0 string) ([]model.User, error) {
	ret := _m.Called(_a0)

	var r0 []model.User
	if rf, ok := ret.Get(0).(func(string) []model.User); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.User)
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

// FindByUsername provides a mock function with given fields: _a0
func (_m *UserStore) FindByUsername(_a0 string) (*model.User, error) {
	ret := _m.Called(_a0)

	var r0 *model.User
	if rf, ok := ret.Get(0).(func(string) *model.User); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.User)
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

// Get provides a mock function with given fields: _a0
func (_m *UserStore) Get(_a0 uint64) (*model.User, error) {
	ret := _m.Called(_a0)

	var r0 *model.User
	if rf, ok := ret.Get(0).(func(uint64) *model.User); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.User)
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

// GetStats provides a mock function with given fields: _a0
func (_m *UserStore) GetStats(_a0 uint64) (*model.UserStats, error) {
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

// Insert provides a mock function with given fields: _a0, _a1
func (_m *UserStore) Insert(_a0 *model.User, _a1 sqlx.Execer) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(*model.User, sqlx.Execer) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// List provides a mock function with given fields:
func (_m *UserStore) List() ([]model.User, error) {
	ret := _m.Called()

	var r0 []model.User
	if rf, ok := ret.Get(0).(func() []model.User); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.User)
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

// ListRoles provides a mock function with given fields: _a0
func (_m *UserStore) ListRoles(_a0 []uint64) ([]model.UserRole, error) {
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

// Update provides a mock function with given fields: _a0, _a1
func (_m *UserStore) Update(_a0 *model.User, _a1 sqlx.Execer) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(*model.User, sqlx.Execer) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdatePassword provides a mock function with given fields: _a0, _a1
func (_m *UserStore) UpdatePassword(_a0 uint64, _a1 string) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(uint64, string) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UsernameExists provides a mock function with given fields: _a0
func (_m *UserStore) UsernameExists(_a0 string) bool {
	ret := _m.Called(_a0)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}
