package core

import (
	"everyflavor/internal/http/api/v1/view"
	"everyflavor/internal/mapper"
	"fmt"

	"github.com/pkg/errors"

	"github.com/rs/zerolog/log"
)

func (a *App) GetFlavorByID(id uint64) (view.Flavor, error) {
	f, err := a.Store.Flavor().Get(id)
	if err != nil {
		log.Error().Err(err).Msg("An error occurred: GetFlavorByID")
	}
	v, err := a.Store.Vendor().Get(f.VendorID)
	if err != nil {
		return view.Flavor{}, err
	}
	return mapper.FlavorFromModel(&f, v), err
}

func (a *App) GetFlavorsList() (view.ListResult, error) {
	var res view.ListResult
	f, cnt, err := a.Store.Flavor().List()
	if err != nil {
		return res, err
	}
	v, _, err := a.Store.Vendor().List()
	fmt.Println(cnt)
	if err != nil {
		return res, err
	}
	flavors := make([]view.Flavor, len(f))
	for idx, fl := range f {
		for _, vendor := range v {
			if vendor.ID == fl.VendorID {
				flavors[idx] = mapper.FlavorFromModel(&fl, &vendor)
				break
			}
		}
	}
	return view.ListResult{Results: flavors, Count: cnt}, err
}

func (a *App) SaveFlavor(v view.Flavor) error {
	tx, err := a.Store.Connection().Beginx()
	if err != nil {
		return errors.Wrap(err, "failed to start transaction")
	}
	flavor, _ := mapper.FlavorToModel(v)
	if v.ID > 0 {
		err = a.Store.Flavor().Update(flavor, tx)
	} else {
		err = a.Store.Flavor().Insert(flavor, tx)
	}
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (a *App) GetStashForUser(userID uint64) ([]view.FlavorStash, error) {
	fs, err := a.Store.Stash().List(userID)
	var v []view.FlavorStash
	for _, f := range fs {
		v = append(v, mapper.FlavorStashFromModel(f))
	}
	return v, err
}

func (a *App) SaveStash(v view.FlavorStash) error {
	fs := mapper.FlavorStashToModel(v)
	err := a.Store.Stash().Insert(&fs, nil)
	if err != nil {
		log.Error().Err(err).Msg("Failed to save flavor stash")
	}
	return err
}
