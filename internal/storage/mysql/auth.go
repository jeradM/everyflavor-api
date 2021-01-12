package mysql

import (
	"everyflavor/internal/storage"
	"everyflavor/internal/storage/model"

	sq "github.com/Masterminds/squirrel"
)

type authStore struct {
	store *Store
}

func NewAuthStore(store *Store) storage.AuthStore {
	return &authStore{store: store}
}

func (r *authStore) IsPublic(resourceID uint64, t model.PublishableEntity) (bool, error) {
	q, args, err := sq.Select(t.PublicField()).
		From(t.TableName()).
		Where(sq.Eq{"id": resourceID}).
		ToSql()
	if err != nil {
		return false, err
	}
	var p uint64
	err = r.store.DB().Get(&p, q, args)
	if err != nil {
		return false, err
	}
	return p == 1, nil
}

func (r *authStore) IsOwner(userID, resourceID uint64, t model.OwnedEntity) (bool, error) {
	q, args, err := sq.Select(t.OwnerField()).
		From(t.TableName()).
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

func (r *authStore) IsCollaborator(userID, resourceID uint64, entity model.SharedEntity) (bool, error) {
	query, args := entity.CollaboratorIDsQuery(resourceID)
	ids := []uint64{}
	err := r.store.DB().Select(&ids, query, args)
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
