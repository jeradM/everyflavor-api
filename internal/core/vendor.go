package core

import (
	"everyflavor/internal/storage/model"
	"github.com/pkg/errors"

	"github.com/rs/zerolog/log"
)

func (a *App) GetVendorByID(id uint64) (*model.Vendor, error) {
	vendor, err := a.Store.Vendor().Get(id)
	if err != nil {
		log.Error().Msg(err.Error())
	}
	return vendor, err
}

func (a *App) GetVendorsList() ([]model.Vendor, error) {
	vendors, _, err := a.Store.Vendor().List()
	return vendors, err
}

func (a *App) SaveVendor(v *model.Vendor) error {
	var err error
	tx, err := a.Store.Connection().Beginx()
	if err != nil {
		return errors.Wrap(err, "failed to start transaction")
	}
	if v.ID > 0 {
		err = a.Store.Vendor().Update(v, tx)
	} else {
		err = a.Store.Vendor().Insert(v, tx)
	}
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (a *App) GetVendorByAbbreviation(abbrev string) (*model.Vendor, error) {
	vendor, err := a.Store.Vendor().FindByAbbreviation(abbrev)
	if err != nil {
		log.Error().Msg(err.Error())
	}
	return vendor, err
}

func (a *App) GetVendorByName(name string) (*model.Vendor, error) {
	vendor, err := a.Store.Vendor().FindByName(name)
	if err != nil {
		log.Error().Msg(err.Error())
	}
	return vendor, err
}
