// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

// Regenerate with `make core_mocks`

package mocks

import (
	view "everyflavor/internal/http/api/v1/view"

	mock "github.com/stretchr/testify/mock"
)

// TagService is an autogenerated mock type for the TagService type
type TagService struct {
	mock.Mock
}

// ListTags provides a mock function with given fields:
func (_m *TagService) ListTags() ([]view.Tag, error) {
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
