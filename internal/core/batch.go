package core

import (
	"everyflavor/internal/http/api/v1/view"
	"everyflavor/internal/mapper"
	"everyflavor/internal/storage/model"

	"github.com/rs/zerolog/log"
)

func (a *App) GetBatch(id uint64) (*view.Batch, error) {
	b, err := a.Store.Batch().Get(id)
	if err != nil {
		return nil, err
	}
	flavors, err := a.Store.Batch().ListFlavorsForBatches([]uint64{b.ID})
	batch := mapper.BatchFromModel(b, flavors)
	return &batch, err
}

func (a *App) GetBatchesForUser(ownerId uint64) ([]model.Batch, error) {
	b, err := a.Store.Batch().List(ownerId)
	if err != nil {
		log.Error().Err(err).Msg("An error occurred: BatchesService.ForUser()")
	}
	return b, err
}

func (a *App) SaveBatch(b view.Batch) (view.Batch, error) {
	var err error
	batch, flavors := mapper.BatchToModel(b)
	tx, err := a.Store.Connection().Beginx()
	if err != nil {
		return b, err
	}
	if b.ID > 0 {
		err = a.Store.Batch().Update(&batch, tx)
	} else {
		err = a.Store.Batch().Insert(&batch, tx)
	}
	if err != nil {
		_ = tx.Rollback()
		log.Error().Err(err).Msg("An error occurred saving Recipe")
		return b, err
	}
	for _, f := range flavors {
		f.BatchID = batch.ID
		err = a.Store.Batch().InsertFlavor(&f, tx)
		if err != nil {
			_ = tx.Rollback()
			return b, err
		}
	}
	err = tx.Commit()
	return mapper.BatchFromModel(batch, flavors), err
}
