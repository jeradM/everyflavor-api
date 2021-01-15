package mysql

import (
	mocks "everyflavor/internal/storage/mockstore"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthStore_IsPublic_True(t *testing.T) {
	db := new(mocks.DBQueryExecer)
	db.On("Get", mock.AnythingOfType("*uint64"), mock.AnythingOfType("string"), mock.Anything).
		Return(nil).
		Run(func(args mock.Arguments) {
			p := args.Get(0).(*uint64)
			*p = uint64(1)
		})
	store := NewMySQLStore(new(mocks.DB), db)
	authStore := NewAuthStore(store)
	public, err := authStore.IsPublic(0, "recipes")
	assert.NoError(t, err)
	assert.True(t, public)
}

func TestAuthStore_IsPublic_False(t *testing.T) {
	db := new(mocks.DBQueryExecer)
	db.On("Get", mock.AnythingOfType("*uint64"), mock.AnythingOfType("string"), mock.Anything).
		Return(nil).
		Run(func(args mock.Arguments) {
			p := args.Get(0).(*uint64)
			*p = uint64(0)
		})
	store := NewMySQLStore(new(mocks.DB), db)
	authStore := NewAuthStore(store)
	public, err := authStore.IsPublic(0, "recipes")
	assert.NoError(t, err)
	assert.False(t, public)
}

func TestAuthStore_IsPublic_Error(t *testing.T) {
	db := new(mocks.DBQueryExecer)
	db.On("Get", mock.AnythingOfType("*uint64"), mock.AnythingOfType("string"), mock.Anything).
		Return(errors.New(""))
	store := NewMySQLStore(new(mocks.DB), db)
	authStore := NewAuthStore(store)
	public, err := authStore.IsPublic(0, "recipes")
	assert.Error(t, err)
	assert.False(t, public)
}

func TestAuthStore_IsOwner_True(t *testing.T) {
	db := new(mocks.DBQueryExecer)
	db.On("Get", mock.AnythingOfType("*uint64"), mock.AnythingOfType("string"), mock.Anything).
		Return(nil).
		Run(func(args mock.Arguments) {
			p := args.Get(0).(*uint64)
			*p = uint64(1)
		})
	store := NewMySQLStore(new(mocks.DB), db)
	authStore := NewAuthStore(store)
	owner, err := authStore.IsOwner(1, 0, "recipes")
	assert.NoError(t, err)
	assert.True(t, owner)
}

func TestAuthStore_IsOwner_False(t *testing.T) {
	db := new(mocks.DBQueryExecer)
	db.On("Get", mock.AnythingOfType("*uint64"), mock.Anything, mock.Anything).
		Return(nil).
		Run(func(args mock.Arguments) {
			p := args.Get(0).(*uint64)
			*p = uint64(1)
		})
	store := NewMySQLStore(new(mocks.DB), db)
	authStore := NewAuthStore(store)
	owner, err := authStore.IsOwner(2, 0, "recipes")
	assert.NoError(t, err)
	assert.False(t, owner)
}

func TestAuthStore_IsOwner_Error(t *testing.T) {
	db := new(mocks.DBQueryExecer)
	db.On("Get", mock.AnythingOfType("*uint64"), mock.Anything, mock.Anything).
		Return(errors.New("error"))
	store := NewMySQLStore(new(mocks.DB), db)
	authStore := NewAuthStore(store)
	owner, err := authStore.IsOwner(2, 0, "recipes")
	assert.Error(t, err)
	assert.False(t, owner)
}

func TestAuthStore_IsCollaborator_True(t *testing.T) {
	db := new(mocks.DBQueryExecer)
	db.On("Select", mock.AnythingOfType("*[]uint64"), mock.Anything, mock.Anything).
		Return(nil).
		Run(func(args mock.Arguments) {
			p := args.Get(0).(*[]uint64)
			*p = []uint64{1, 2, 3}
		})
	store := NewMySQLStore(new(mocks.DB), db)
	authStore := NewAuthStore(store)
	collab, err := authStore.IsCollaborator(1, 0, "recipe")
	assert.NoError(t, err)
	assert.True(t, collab)
}

func TestAuthStore_IsCollaborator_False(t *testing.T) {
	db := new(mocks.DBQueryExecer)
	db.On("Select", mock.AnythingOfType("*[]uint64"), mock.Anything, mock.Anything).
		Return(nil).
		Run(func(args mock.Arguments) {
			p := args.Get(0).(*[]uint64)
			*p = []uint64{1, 2, 3}
		})
	store := NewMySQLStore(new(mocks.DB), db)
	authStore := NewAuthStore(store)
	collab, err := authStore.IsCollaborator(4, 0, "recipe")
	assert.NoError(t, err)
	assert.False(t, collab)
}

func TestAuthStore_IsCollaborator_Error(t *testing.T) {
	db := new(mocks.DBQueryExecer)
	db.On("Select", mock.AnythingOfType("*[]uint64"), mock.Anything, mock.Anything).
		Return(errors.New("error"))
	store := NewMySQLStore(new(mocks.DB), db)
	authStore := NewAuthStore(store)
	collab, err := authStore.IsCollaborator(4, 0, "recipe")
	assert.Error(t, err)
	assert.False(t, collab)
}
