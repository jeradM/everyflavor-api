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
	var m model.Vendor
	err := r.store.getEntityByID(nil, &m, id, func(builder sq.SelectBuilder) sq.SelectBuilder {
		return builder.Where(sq.Expr("vendors.deleted_at IS NULL"))
	})
	return &m, err
}

func (r *vendorStore) FindByAbbreviation(abbrev string) (*model.Vendor, error) {
	q, args, err := sq.
		Select("id", "created_at", "updated_at", "abbreviation", "name", "aliases").
		From("vendors").
		Where(sq.Expr("deleted_at IS NULL")).
		Where(sq.Eq{"abbreviation": abbrev}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, err
	}
	return r.fetchOne(q, args)
}

func (r *vendorStore) FindByName(name string) (*model.Vendor, error) {
	q, args, err := sq.
		Select("id", "created_at", "updated_at", "abbreviation", "name", "aliases").
		From("vendors").
		Where(sq.Expr("deleted_at IS NULL")).
		Where(sq.Eq{"name": name}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, err
	}
	return r.fetchOne(q, args)
}

func (r *vendorStore) List() ([]model.Vendor, uint64, error) {
	var vendors []model.Vendor
	cnt, err := r.store.listEntities(
		nil, &vendors, model.BaseListParams{}, func(stmt sq.SelectBuilder) sq.SelectBuilder {
			return stmt.Where("deleted_at IS NULL")
		})
	if err != nil {
		return vendors, cnt, err
	}
	return vendors, cnt, err
}

func (r *vendorStore) Insert(vendor *model.Vendor, tx sqlx.Execer) error {
	return r.store.insertEntity(tx, vendor, func(id int64) {
		vendor.ID = uint64(id)
	})
}

func (r *vendorStore) Update(vendor *model.Vendor, tx sqlx.Execer) error {
	_, err := r.store.updateEntity(tx, vendor)
	return err
}

func (r *vendorStore) fetchOne(q string, args ...interface{}) (*model.Vendor, error) {
	v := model.Vendor{}
	err := r.store.DB().Get(&v, q, args)
	return &v, err
}
