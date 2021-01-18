package mysql

import (
	"everyflavor/internal/storage"
	"everyflavor/internal/storage/model"

	"github.com/pkg/errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type flavorStore struct {
	store *Store
}

// NewFlavorStore creates and returns a new mysql-based flavor store
func NewFlavorStore(store *Store) storage.FlavorStore {
	return &flavorStore{store: store}
}

func (r *flavorStore) Get(id uint64) (model.Flavor, error) {
	var m model.Flavor
	stmt := sq.Select(flavorSelectFields...).
		From("flavors").
		Where(sq.Eq{"id": id})
	for _, jc := range flavorSelectJoins {
		stmt = stmt.JoinClause(jc)
	}
	query, args, err := stmt.ToSql()
	if err != nil {
		return m, errors.Wrap(err, "failed to build select query")
	}
	err = r.store.DB().Get(&m, query, args)
	return m, errors.Wrap(err, "failed to fetch Flavor")
}

func (r *flavorStore) List() ([]model.Flavor, uint64, error) {
	var flavors []model.Flavor
	var cnt uint64
	stmt := sq.Select().
		From("flavors").
		Where(sq.Expr("flavors.deleted_at IS NULL")).
		GroupBy("flavors.id")
	for _, jc := range flavorSelectJoins {
		stmt = stmt.JoinClause(jc)
	}

	query, args, err := stmt.Columns("count(flavors.id)").ToSql()
	if err != nil {
		return flavors, 0, errors.Wrap(err, "failed to build select count statement")
	}
	err = r.store.DB().Get(&cnt, query, args)
	if err != nil {
		return flavors, cnt, errors.Wrap(err, "failed to fetch flavor count")
	}

	query, args, err = stmt.Columns(flavorSelectFields...).ToSql()
	if err != nil {
		return flavors, 0, errors.Wrap(err, "failed to build select statement")
	}
	err = r.store.DB().Select(&flavors, query, args)

	return flavors, cnt, errors.Wrap(err, "failed to fetch flavors")
}

// Insert adds a new flavor to the database
func (r *flavorStore) Insert(flavor *model.Flavor, tx sqlx.Execer) error {
	query, args, err := sq.Insert("flavors").
		SetMap(flavorInsertMap(flavor)).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "failed to build insert query")
	}
	result, err := r.store.DB().ExecWithTX(tx, query, args)
	if err != nil {
		return errors.Wrap(err, "failed to insert flavor")
	}
	id, err := result.LastInsertId()
	if err != nil {
		return errors.Wrap(err, "failed to get LastInsertID")
	}
	flavor.ID = uint64(id)
	return nil
}

func (r *flavorStore) Update(flavor *model.Flavor, tx sqlx.Execer) error {
	query, args, err := sq.Update("flavors").
		SetMap(flavorUpdateMap(flavor)).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "failed to build update query")
	}
	_, err = r.store.DB().ExecWithTX(tx, query, args)
	return errors.Wrap(err, "failed to update flavor")
}

var (
	flavorSelectFields = []string{
		"flavors.id",
		"flavors.created_at",
		"flavors.updated_at",
		"flavors.vendor_id",
		"flavors.name",
		"flavors.aliases",
		"count(recipes.id) recipe_count",
		"IFNULL(CAST(avg(recipe_flavors.percent_m) as UNSIGNED), 0) avg_percent",
	}

	flavorSelectJoins = []string{
		"LEFT JOIN recipe_flavors ON recipe_flavors.flavor_id = flavors.id",
		"LEFT JOIN recipes ON recipes.id = recipe_flavors.recipe_id",
	}
)

func flavorInsertMap(f *model.Flavor) map[string]interface{} {
	return map[string]interface{}{
		"flavors.created_at": sq.Expr("NOW()"),
		"flavors.updated_at": sq.Expr("NOW()"),
		"flavors.vendor_id":  f.VendorID,
		"flavors.name":       f.Name,
		"flavors.aliases":    f.Aliases,
	}
}

func flavorUpdateMap(f *model.Flavor) map[string]interface{} {
	return map[string]interface{}{
		"flavors.updated_at": sq.Expr("NOW()"),
		"flavors.vendor_id":  f.VendorID,
		"flavors.name":       f.Name,
		"flavors.aliases":    f.Aliases,
	}
}
