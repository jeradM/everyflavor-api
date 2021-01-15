package core

import (
	"database/sql"
	"everyflavor/internal/http/api/v1/view"
	mocks "everyflavor/internal/storage/mockstore"
	"everyflavor/internal/storage/model"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func dummyBatchModel() model.Batch {
	return model.Batch{
		ID:            0,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		DeletedAt:     sql.NullTime{Time: time.Now(), Valid: false},
		BatchSizeM:    10000,
		BatchStrength: 6,
		BatchVgM:      70000,
		MaxVg:         false,
		NicStrength:   100,
		NicVgM:        0,
		RecipeID:      1,
		OwnerID:       1,
		UseNic:        true,
	}
}

func dummyBatchView() view.Batch {
	return view.Batch{
		ID:            0,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		BatchSizeM:    10000,
		BatchStrength: 6,
		BatchVgM:      70000,
		MaxVg:         false,
		NicStrength:   100,
		NicVgM:        0,
		RecipeID:      1,
		OwnerID:       1,
		UseNic:        true,
		Flavors: []view.BatchFlavor{
			{FlavorID: 1, PercentM: 5000, Vg: false},
			{FlavorID: 2, PercentM: 2000, Vg: false},
		},
	}
}

func TestApp_GetBatch_Error(t *testing.T) {
	batchStore := new(mocks.BatchStore)
	batchStore.On("Get", mock.AnythingOfType("uint64")).Return(dummyBatchModel(), errors.New(""))
	store := new(mocks.Store)
	store.On("Batch").Return(batchStore)

	a := NewApp(AppConfig{}, store)
	batch, err := a.GetBatch(1)
	assert.Error(t, err)
	assert.Nil(t, batch)
}

func TestApp_GetBatch(t *testing.T) {
	batchStore := new(mocks.BatchStore)
	batchStore.On("Get", mock.AnythingOfType("uint64")).
		Return(dummyBatchModel(), nil)
	batchStore.On("ListFlavorsForBatches", mock.AnythingOfType("[]uint64")).
		Return([]model.BatchFlavor{}, nil)
	store := new(mocks.Store)
	store.On("Batch").Return(batchStore)

	a := NewApp(AppConfig{}, store)
	batch, err := a.GetBatch(1)
	assert.NoError(t, err)
	assert.NotNil(t, batch)
	assert.Equal(t, uint64(0), batch.ID)
	assert.Equal(t, uint64(1), batch.OwnerID)
	assert.IsType(t, &view.Batch{}, batch)
}

func TestSaveBatch_Insert(t *testing.T) {
	db, dbMock, _ := sqlmock.New()
	dbMock.ExpectBegin()
	dbMock.ExpectCommit()
	database := sqlx.NewDb(db, "sqlmock")
	batchStore := new(mocks.BatchStore)
	batchStore.On("Insert", mock.AnythingOfType("*model.Batch"), mock.AnythingOfType("*sqlx.Tx")).
		Return(nil).
		Run(func(args mock.Arguments) {
			b := args.Get(0).(*model.Batch)
			b.ID = 1
		})
	batchStore.On("InsertFlavor", mock.Anything, mock.Anything).Return(nil)
	store := new(mocks.Store)
	store.On("Batch").Return(batchStore)
	store.On("Connection").Return(database)

	a := NewApp(AppConfig{}, store)
	v, err := a.SaveBatch(dummyBatchView())
	assert.NoError(t, err)
	assert.NoError(t, dbMock.ExpectationsWereMet())
	assert.IsType(t, view.Batch{}, v)
	batchStore.AssertNumberOfCalls(t, "Insert", 1)
	batchStore.AssertNumberOfCalls(t, "InsertFlavor", 2)
}

func TestSaveBatch_InsertFail(t *testing.T) {
	db, dbMock, _ := sqlmock.New()
	dbMock.ExpectBegin()
	dbMock.ExpectRollback()
	database := sqlx.NewDb(db, "sqlmock")
	batchStore := new(mocks.BatchStore)
	batchStore.On("Insert", mock.AnythingOfType("*model.Batch"), mock.AnythingOfType("*sqlx.Tx")).
		Return(errors.New(""))
	store := new(mocks.Store)
	store.On("Batch").Return(batchStore)
	store.On("Connection").Return(database)

	a := NewApp(AppConfig{}, store)
	v, err := a.SaveBatch(dummyBatchView())
	assert.Error(t, err)
	assert.NoError(t, dbMock.ExpectationsWereMet())
	assert.IsType(t, view.Batch{}, v)
	batchStore.AssertNumberOfCalls(t, "Insert", 1)
	batchStore.AssertNumberOfCalls(t, "InsertFlavor", 0)
}

func TestSaveBatch_Update(t *testing.T) {
	db, dbMock, _ := sqlmock.New()
	dbMock.ExpectBegin()
	dbMock.ExpectCommit()
	database := sqlx.NewDb(db, "sqlmock")
	batchStore := new(mocks.BatchStore)
	batchStore.On("Update", mock.AnythingOfType("*model.Batch"), mock.AnythingOfType("*sqlx.Tx")).
		Return(nil)
	batchStore.On("InsertFlavor", mock.Anything, mock.Anything).Return(nil)
	store := new(mocks.Store)
	store.On("Batch").Return(batchStore)
	store.On("Connection").Return(database)

	a := NewApp(AppConfig{}, store)
	b := dummyBatchView()
	b.ID = 1
	v, err := a.SaveBatch(b)
	assert.NoError(t, err)
	assert.NoError(t, dbMock.ExpectationsWereMet())
	assert.IsType(t, view.Batch{}, v)
	batchStore.AssertNumberOfCalls(t, "Update", 1)
	batchStore.AssertNumberOfCalls(t, "InsertFlavor", 2)
}

func TestSaveBatch_UpdateFail(t *testing.T) {
	db, dbMock, _ := sqlmock.New()
	dbMock.ExpectBegin()
	dbMock.ExpectRollback()
	database := sqlx.NewDb(db, "sqlmock")
	batchStore := new(mocks.BatchStore)
	batchStore.On("Update", mock.AnythingOfType("*model.Batch"), mock.AnythingOfType("*sqlx.Tx")).
		Return(errors.New(""))
	store := new(mocks.Store)
	store.On("Batch").Return(batchStore)
	store.On("Connection").Return(database)

	a := NewApp(AppConfig{}, store)
	b := dummyBatchView()
	b.ID = 1
	v, err := a.SaveBatch(b)
	assert.Error(t, err)
	assert.NoError(t, dbMock.ExpectationsWereMet())
	assert.IsType(t, view.Batch{}, v)
	batchStore.AssertNumberOfCalls(t, "Update", 1)
	batchStore.AssertNumberOfCalls(t, "InsertFlavor", 0)
}

func TestSaveBatch_Update_InsertFlavorFail(t *testing.T) {
	db, dbMock, _ := sqlmock.New()
	dbMock.ExpectBegin()
	dbMock.ExpectRollback()
	database := sqlx.NewDb(db, "sqlmock")
	batchStore := new(mocks.BatchStore)
	batchStore.On("Insert", mock.AnythingOfType("*model.Batch"), mock.AnythingOfType("*sqlx.Tx")).
		Return(nil)
	batchStore.On("InsertFlavor", mock.Anything, mock.Anything).
		Return(errors.New(""))
	store := new(mocks.Store)
	store.On("Batch").Return(batchStore)
	store.On("Connection").Return(database)

	a := NewApp(AppConfig{}, store)
	v, err := a.SaveBatch(dummyBatchView())
	assert.Error(t, err)
	assert.NoError(t, dbMock.ExpectationsWereMet())
	assert.IsType(t, view.Batch{}, v)
	batchStore.AssertNumberOfCalls(t, "Insert", 1)
	batchStore.AssertNumberOfCalls(t, "InsertFlavor", 1)
}
