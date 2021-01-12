// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

// Regenerate using `make store_mocks`

package mocks

import (
	model "everyflavor/internal/storage/model"

	mock "github.com/stretchr/testify/mock"

	sqlx "github.com/jmoiron/sqlx"
)

// BatchStore is an autogenerated mock type for the BatchStore type
type BatchStore struct {
	mock.Mock
}

// Get provides a mock function with given fields: _a0
func (_m *BatchStore) Get(_a0 uint64) (model.Batch, error) {
	ret := _m.Called(_a0)

	var r0 model.Batch
	if rf, ok := ret.Get(0).(func(uint64) model.Batch); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(model.Batch)
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
func (_m *BatchStore) Insert(_a0 *model.Batch, _a1 sqlx.Execer) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(*model.Batch, sqlx.Execer) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// InsertFlavor provides a mock function with given fields: _a0, _a1
func (_m *BatchStore) InsertFlavor(_a0 *model.BatchFlavor, _a1 sqlx.Execer) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(*model.BatchFlavor, sqlx.Execer) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// List provides a mock function with given fields: _a0
func (_m *BatchStore) List(_a0 uint64) ([]model.Batch, error) {
	ret := _m.Called(_a0)

	var r0 []model.Batch
	if rf, ok := ret.Get(0).(func(uint64) []model.Batch); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.Batch)
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

// ListFlavorsForBatches provides a mock function with given fields: _a0
func (_m *BatchStore) ListFlavorsForBatches(_a0 []uint64) ([]model.BatchFlavor, error) {
	ret := _m.Called(_a0)

	var r0 []model.BatchFlavor
	if rf, ok := ret.Get(0).(func([]uint64) []model.BatchFlavor); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.BatchFlavor)
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
func (_m *BatchStore) Update(_a0 *model.Batch, _a1 sqlx.Execer) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(*model.Batch, sqlx.Execer) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
