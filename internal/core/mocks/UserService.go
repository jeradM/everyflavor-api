// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

// Regenerate with `make core_mocks`

package mocks

import (
	model "everyflavor/internal/storage/model"

	mock "github.com/stretchr/testify/mock"

	view "everyflavor/internal/http/api/v1/view"
)

// UserService is an autogenerated mock type for the UserService type
type UserService struct {
	mock.Mock
}

// EmailExists provides a mock function with given fields: _a0
func (_m *UserService) EmailExists(_a0 string) bool {
	ret := _m.Called(_a0)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// GetRolesForUsers provides a mock function with given fields: _a0
func (_m *UserService) GetRolesForUsers(_a0 []uint64) ([]model.UserRole, error) {
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

// GetUserByID provides a mock function with given fields: id
func (_m *UserService) GetUserByID(id uint64) (view.User, error) {
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
func (_m *UserService) GetUserByUsername(username string) (view.User, error) {
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
func (_m *UserService) GetUserList() ([]view.User, error) {
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
func (_m *UserService) GetUserStatsByID(_a0 uint64) (*model.UserStats, error) {
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

// SaveUser provides a mock function with given fields: u
func (_m *UserService) SaveUser(u view.User) error {
	ret := _m.Called(u)

	var r0 error
	if rf, ok := ret.Get(0).(func(view.User) error); ok {
		r0 = rf(u)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SearchUsersByUsername provides a mock function with given fields: username
func (_m *UserService) SearchUsersByUsername(username string) ([]view.User, error) {
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

// UpdateUserPassword provides a mock function with given fields: id, pw
func (_m *UserService) UpdateUserPassword(id uint64, pw string) error {
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
func (_m *UserService) UsernameExists(_a0 string) bool {
	ret := _m.Called(_a0)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}
