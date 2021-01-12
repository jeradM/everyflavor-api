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
	var m model.FlavorStash
	stmt := sq.Select(m.SelectFields()...).
		From(m.TableName())
	for _, jc := range m.SelectJoins() {
		stmt = stmt.JoinClause(jc)
	}
	query, args, err := stmt.Where(sq.Eq{"flavor_stashes.owner_id": userID}).
		Where(sq.Expr("flavor_stashes.deleted_at IS NULL")).
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "failed to build SQL query string")
	}
	stash := []model.FlavorStash{}
	err = r.store.DB().Select(&stash, query, args)
	return stash, errors.Wrap(err, "failed to fetch flavor stash")
}

func (r *flavorStashStore) Insert(s *model.FlavorStash, tx sqlx.Execer) error {
	return r.store.insertEntity(tx, s, func(id int64) {
		s.ID = uint64(id)
	})
}

func (r *flavorStashStore) Delete(s *model.FlavorStash, tx sqlx.Execer) (int64, error) {
	panic("not implemented")
}
