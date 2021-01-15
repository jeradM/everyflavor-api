package mysql

import (
	"everyflavor/internal/storage"
	"everyflavor/internal/storage/model"

	"github.com/jmoiron/sqlx"

	sq "github.com/Masterminds/squirrel"
)

type vendorStore struct {
	store *Store
}

func NewVendorStore(store *Store) storage.VendorStore {
	return &vendorStore{store: store}
}

func (r *vendorStore) Get(id uint64) (*model.Vendor, error) {
	query, args, err := sq.
		Select("id", "created_at", "updated_at", "abbreviation", "name", "aliases").
		From("vendors").
		Where(sq.Expr("deleted_at IS NULL")).
		Where(sq.Eq{"id": id}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, err
	}
	var v model.Vendor
	err = r.store.DB().Get(&v, query, args)
	return &v, err
}

func (r *vendorStore) FindByAbbreviation(abbrev string) (*model.Vendor, error) {
	query, args, err := sq.
		Select("id", "created_at", "updated_at", "abbreviation", "name", "aliases").
		From("vendors").
		Where(sq.Expr("deleted_at IS NULL")).
		Where(sq.Eq{"abbreviation": abbrev}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, err
	}
	var v model.Vendor
	err = r.store.DB().Get(&v, query, args)
	return &v, err
}

func (r *vendorStore) FindByName(name string) (*model.Vendor, error) {
	query, args, err := sq.
		Select("id", "created_at", "updated_at", "abbreviation", "name", "aliases").
		From("vendors").
		Where(sq.Expr("deleted_at IS NULL")).
		Where(sq.Eq{"name": name}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, err
	}
	var v model.Vendor
	err = r.store.DB().Get(&v, query, args)
	return &v, err
}

func (r *vendorStore) List() ([]model.Vendor, uint64, error) {
	var count uint64
	var vendors []model.Vendor
	stmt := sq.Select().
		From("vendors").
		Where(sq.Expr("deleted_at IS NULL"))

	query, args, err := stmt.Columns("count(vendors.id)").ToSql()
	if err != nil {
		return vendors, count, err
	}
	err = r.store.DB().Get(&count, query, args)
	if err != nil {
		return vendors, count, err
	}

	query, args, err = stmt.
		Columns("id", "created_at", "updated_at", "abbreviation", "name", "aliases").
		ToSql()
	if err != nil {
		return vendors, count, err
	}
	err = r.store.DB().Select(&vendors, query, args)
	return vendors, count, err
}

func (r *vendorStore) Insert(vendor *model.Vendor, tx sqlx.Execer) error {
	query, args, err := sq.Insert("vendors").
		SetMap(map[string]interface{}{
			"abbreviation": vendor.Abbreviation,
			"name":         vendor.Name,
			"aliases":      vendor.Aliases,
		}).
		ToSql()
	if err != nil {
		return err
	}
	_, err = r.store.DB().ExecWithTX(tx, query, args)
	return err
}

func (r *vendorStore) Update(vendor *model.Vendor, tx sqlx.Execer) error {
	query, args, err := sq.Update("vendors").
		SetMap(map[string]interface{}{
			"abbreviation": vendor.Abbreviation,
			"name":         vendor.Name,
			"aliases":      vendor.Aliases,
		}).
		Where(sq.Eq{"id": vendor.ID}).
		ToSql()
	if err != nil {
		return err
	}
	_, err = r.store.DB().ExecWithTX(tx, query, args)
	return err
}
