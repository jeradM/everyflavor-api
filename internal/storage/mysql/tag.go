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
	query, args, err := sq.Select(
		"tags.id",
		"tags.created_at",
		"tags.updated_at",
		"tags.tag").
		From("tags").
		Where(sq.Eq{"tags.id": id}).
		ToSql()
	if err != nil {
		return &m, err
	}
	err = r.store.DB().Get(&m, query, args)
	return &m, err
}

func (r *tagStore) List() ([]model.Tag, error) {
	var tags []model.Tag
	query, args, err := sq.Select(
		"tags.id",
		"tags.created_at",
		"tags.updated_at",
		"tags.tag").
		From("tags").
		Where(sq.Expr("deleted_at IS NULL")).
		ToSql()
	if err != nil {
		return tags, err
	}
	err = r.store.DB().Select(&tags, query, args)
	return tags, err
}
