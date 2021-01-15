package mysql

import (
	"everyflavor/internal/storage"
	"github.com/pkg/errors"

	sq "github.com/Masterminds/squirrel"
)

type authStore struct {
	store *Store
}

func NewAuthStore(store *Store) storage.AuthStore {
	return &authStore{store: store}
}

func (r *authStore) IsPublic(id uint64, table string) (bool, error) {
	q, args, err := sq.Select("public").
		From(table).
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return false, errors.Wrap(err, "failed to build select query")
	}
	var p uint64
	err = r.store.DB().Get(&p, q, args)
	if err != nil {
		return false, err
	}
	return p == 1, nil
}

func (r *authStore) IsOwner(userID, resourceID uint64, table string) (bool, error) {
	q, args, err := sq.Select("owner_id").
		From(table).
		Where(sq.Eq{"id": resourceID}).
		ToSql()
	if err != nil {
		return false, err
	}
	var resId uint64
	err = r.store.DB().Get(&resId, q, args)
	if err != nil {
		return false, err
	}
	return resId == userID, nil
}

func (r *authStore) IsCollaborator(userID, resourceID uint64, model string) (bool, error) {
	query, args, err := sq.Select("user_id").
		From(model + "_collaborators").
		Where(sq.Eq{model + "_id": resourceID}).
		ToSql()
	if err != nil {
		return false, errors.Wrap(err, "failed to build select query")
	}
	var ids []uint64
	err = r.store.DB().Select(&ids, query, args)
	if err != nil {
		return false, err
	}
	for _, id := range ids {
		if userID == id {
			return true, nil
		}
	}
	return false, nil
}
