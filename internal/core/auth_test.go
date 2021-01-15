package core

import (
	mocks "everyflavor/internal/storage/mockstore"
	"testing"

	"github.com/pkg/errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestApp_CanViewRecipe_False(t *testing.T) {
	authStore := new(mocks.AuthStore)
	authStore.On("IsPublic", mock.Anything, mock.Anything).Return(false, nil)
	authStore.On("IsOwner", mock.Anything, mock.Anything, mock.Anything).Return(false, nil)
	authStore.On("IsCollaborator", mock.Anything, mock.Anything, mock.Anything).Return(false, nil)
	store := new(mocks.Store)
	store.On("Auth").Return(authStore)

	a := NewApp(AppConfig{}, store)

	canView, err := a.CanViewRecipe(0, 0)
	assert.NoError(t, err)
	assert.False(t, canView)
}

func TestApp_CanViewRecipe_Error(t *testing.T) {
	authStore := new(mocks.AuthStore)
	authStore.On("IsPublic", mock.Anything, mock.Anything).Return(true, errors.New(""))
	store := new(mocks.Store)
	store.On("Auth").Return(authStore)

	a := NewApp(AppConfig{}, store)

	canView, err := a.CanViewRecipe(0, 0)
	assert.Error(t, err)
	assert.False(t, canView)
}

func TestApp_CanViewRecipe_IsCollaborator(t *testing.T) {
	authStore := new(mocks.AuthStore)
	authStore.On("IsPublic", mock.Anything, mock.Anything).Return(false, nil)
	authStore.On("IsOwner", mock.Anything, mock.Anything, mock.Anything).Return(false, nil)
	authStore.On("IsCollaborator", mock.Anything, mock.Anything, mock.Anything).Return(true, nil)
	store := new(mocks.Store)
	store.On("Auth").Return(authStore)

	a := NewApp(AppConfig{}, store)

	canView, err := a.CanViewRecipe(0, 0)
	assert.NoError(t, err)
	assert.True(t, canView)
}

func TestApp_CanViewRecipe_IsOwner(t *testing.T) {
	authStore := new(mocks.AuthStore)
	authStore.On("IsPublic", mock.Anything, mock.Anything).Return(false, nil)
	authStore.On("IsOwner", mock.Anything, mock.Anything, mock.Anything).Return(true, nil)
	store := new(mocks.Store)
	store.On("Auth").Return(authStore)

	a := NewApp(AppConfig{}, store)

	canView, err := a.CanViewRecipe(0, 0)
	assert.NoError(t, err)
	assert.True(t, canView)
}

func TestApp_CanViewRecipe_IsPublic(t *testing.T) {
	authStore := new(mocks.AuthStore)
	authStore.On("IsPublic", mock.Anything, mock.Anything).Return(true, nil)
	store := new(mocks.Store)
	store.On("Auth").Return(authStore)

	a := NewApp(AppConfig{}, store)

	canView, err := a.CanViewRecipe(0, 0)
	assert.NoError(t, err)
	assert.True(t, canView)
}

func TestApp_CanEditRecipe_False(t *testing.T) {
	authStore := new(mocks.AuthStore)
	authStore.On("IsOwner", mock.Anything, mock.Anything, mock.Anything).Return(false, nil)
	authStore.On("IsCollaborator", mock.Anything, mock.Anything, mock.Anything).Return(false, nil)
	store := new(mocks.Store)
	store.On("Auth").Return(authStore)

	a := NewApp(AppConfig{}, store)

	canEdit, err := a.CanEditRecipe(0, 0)
	assert.NoError(t, err)
	assert.False(t, canEdit)
}

func TestApp_CanEditRecipe_Error(t *testing.T) {
	authStore := new(mocks.AuthStore)
	authStore.On("IsOwner", mock.Anything, mock.Anything, mock.Anything).Return(true, errors.New(""))
	store := new(mocks.Store)
	store.On("Auth").Return(authStore)

	a := NewApp(AppConfig{}, store)

	canEdit, err := a.CanEditRecipe(0, 0)
	assert.Error(t, err)
	assert.False(t, canEdit)
}

func TestApp_CanEditRecipe_Owner(t *testing.T) {
	authStore := new(mocks.AuthStore)
	authStore.On("IsOwner", mock.Anything, mock.Anything, mock.Anything).Return(true, nil)
	store := new(mocks.Store)
	store.On("Auth").Return(authStore)

	a := NewApp(AppConfig{}, store)

	canEdit, err := a.CanEditRecipe(0, 0)
	assert.NoError(t, err)
	assert.True(t, canEdit)
}

func TestApp_CanEditRecipe_Collaborator(t *testing.T) {
	authStore := new(mocks.AuthStore)
	authStore.On("IsOwner", mock.Anything, mock.Anything, mock.Anything).Return(false, nil)
	authStore.On("IsCollaborator", mock.Anything, mock.Anything, mock.Anything).Return(true, nil)
	store := new(mocks.Store)
	store.On("Auth").Return(authStore)

	a := NewApp(AppConfig{}, store)

	canEdit, err := a.CanEditRecipe(0, 0)
	assert.NoError(t, err)
	assert.True(t, canEdit)
}

func TestApp_CanViewBatch_False(t *testing.T) {
	authStore := new(mocks.AuthStore)
	authStore.On("IsOwner", mock.Anything, mock.Anything, mock.Anything).Return(false, nil)
	store := new(mocks.Store)
	store.On("Auth").Return(authStore)

	a := NewApp(AppConfig{}, store)

	canView, err := a.CanViewBatch(0, 0)
	assert.NoError(t, err)
	assert.False(t, canView)
}

func TestApp_CanViewBatch_Error(t *testing.T) {
	authStore := new(mocks.AuthStore)
	authStore.On("IsOwner", mock.Anything, mock.Anything, mock.Anything).Return(true, errors.New(""))
	store := new(mocks.Store)
	store.On("Auth").Return(authStore)

	a := NewApp(AppConfig{}, store)

	canView, err := a.CanViewBatch(0, 0)
	assert.Error(t, err)
	assert.False(t, canView)
}

func TestApp_CanViewBatch_True(t *testing.T) {
	authStore := new(mocks.AuthStore)
	authStore.On("IsOwner", mock.Anything, mock.Anything, mock.Anything).Return(true, nil)
	store := new(mocks.Store)
	store.On("Auth").Return(authStore)

	a := NewApp(AppConfig{}, store)

	canView, err := a.CanViewBatch(0, 0)
	assert.NoError(t, err)
	assert.True(t, canView)
}

func TestApp_CanEditBatch_False(t *testing.T) {
	authStore := new(mocks.AuthStore)
	authStore.On("IsOwner", mock.Anything, mock.Anything, mock.Anything).Return(false, nil)
	store := new(mocks.Store)
	store.On("Auth").Return(authStore)

	a := NewApp(AppConfig{}, store)

	canEdit, err := a.CanEditBatch(0, 0)
	assert.NoError(t, err)
	assert.False(t, canEdit)
}

func TestApp_CanEditBatch_Error(t *testing.T) {
	authStore := new(mocks.AuthStore)
	authStore.On("IsOwner", mock.Anything, mock.Anything, mock.Anything).Return(true, errors.New(""))
	store := new(mocks.Store)
	store.On("Auth").Return(authStore)

	a := NewApp(AppConfig{}, store)

	canEdit, err := a.CanEditBatch(0, 0)
	assert.Error(t, err)
	assert.False(t, canEdit)
}

func TestApp_CanEditBatch_Owner(t *testing.T) {
	authStore := new(mocks.AuthStore)
	authStore.On("IsOwner", mock.Anything, mock.Anything, mock.Anything).Return(true, nil)
	store := new(mocks.Store)
	store.On("Auth").Return(authStore)

	a := NewApp(AppConfig{}, store)

	canEdit, err := a.CanEditBatch(0, 0)
	assert.NoError(t, err)
	assert.True(t, canEdit)
}
