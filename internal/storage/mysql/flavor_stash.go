package mysql

import (
	"everyflavor/internal/storage"
	"everyflavor/internal/storage/model"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type flavorStashStore struct {
	store *Store
}

func NewFlavorStashStore(store *Store) storage.FlavorStashStore {
	return &flavorStashStore{store: store}
}

func (r *flavorStashStore) List(userID uint64) ([]model.FlavorStash, error) {
	stmt := sq.Select(flavorStashSelectFields...).
		From("flavor_stashes")
	for _, jc := range flavorStashSelectJoins {
		stmt = stmt.JoinClause(jc)
	}
	query, args, err := stmt.Where(sq.Eq{"flavor_stashes.owner_id": userID}).
		Where(sq.Expr("flavor_stashes.deleted_at IS NULL")).
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "failed to build select query")
	}
	var stash []model.FlavorStash
	err = r.store.DB().Select(&stash, query, args)
	return stash, errors.Wrap(err, "failed to fetch flavor stash")
}

func (r *flavorStashStore) Insert(s *model.FlavorStash, tx sqlx.Execer) error {
	query, args, err := sq.Insert("flavor_stashes").
		SetMap(flavorStashInsertMap(s)).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "failed to build insert query")
	}
	result, err := r.store.DB().ExecWithTX(tx, query, args)
	if err != nil {
		return errors.Wrap(err, "failed to insert flavor stash")
	}
	id, err := result.LastInsertId()
	if err != nil {
		return errors.Wrap(err, "failed to get LastInsertID")
	}
	s.ID = uint64(id)
	return nil
}

func (r *flavorStashStore) Delete(s *model.FlavorStash, tx sqlx.Execer) (int64, error) {
	query, args, err := sq.Delete("flavor_stashes").
		Where(sq.Eq{"id": s.ID}).
		ToSql()
	if err != nil {
		return 0, errors.Wrap(err, "failed to build delete query")
	}
	result, err := r.store.DB().ExecWithTX(tx, query, args)
	if err != nil {
		return 0, errors.Wrap(err, "failed to delete flavor stash")
	}
	n, err := result.RowsAffected()
	return n, errors.Wrap(err, "failed to get RowsAffected")
}

var (
	flavorStashSelectFields = []string{
		"flavor_stashes.id",
		"flavor_stashes.created_at",
		"flavor_stashes.updated_at",
		"flavor_stashes.on_hand_m",
		"flavor_stashes.density_m",
		"flavor_stashes.vg",
		"flavor_stashes.flavor_id",
		"flavor_stashes.owner_id",
	}

	flavorStashSelectJoins = []string{
		"LEFT JOIN recipe_flavors ON recipe_flavors.flavor_id = flavor_stashes.flavor_id",
		"LEFT JOIN recipes ON recipes.id = recipe_flavors.recipe_id AND recipes.owner_id = flavor_stashes.owner_id",
	}
)

func flavorStashInsertMap(f *model.FlavorStash) map[string]interface{} {
	return map[string]interface{}{
		"created_at": sq.Expr("NOW()"),
		"updated_at": sq.Expr("NOW()"),
		"on_hand_m":  f.OnHandM,
		"density_m":  f.DensityM,
		"vg":         f.Vg,
		"flavor_id":  f.FlavorID,
		"owner_id":   f.OwnerID,
	}
}
