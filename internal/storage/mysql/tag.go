package mysql

import (
	"everyflavor/internal/storage"
	"everyflavor/internal/storage/model"

	sq "github.com/Masterminds/squirrel"
)

type tagStore struct {
	store *Store
}

func NewTagStore(store *Store) storage.TagStore {
	return &tagStore{store: store}
}

func (r *tagStore) Get(id uint64) (*model.Tag, error) {
	var m model.Tag
	err := r.store.getEntityByID(nil, &m, id, nil)
	return &m, err
}

func (r *tagStore) List() ([]model.Tag, error) {
	var m model.Tag
	l := sq.Select(m.SelectFields()...).
		From(m.TableName()).
		Where(sq.Expr("deleted_at IS NULL"))
	q, args, err := l.ToSql()
	if err != nil {
		return nil, err
	}
	var tags []model.Tag
	err = r.store.DB().Select(&tags, q, args)
	return tags, err
}
