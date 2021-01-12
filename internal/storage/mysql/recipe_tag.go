package mysql

import (
	"everyflavor/internal/storage"
	"everyflavor/internal/storage/model"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type recipeTagStore struct {
	db *sqlx.DB
}

func NewRecipeTagStore(db *sqlx.DB) storage.RecipeTagStore {
	return &recipeTagStore{db: db}
}

func (r *recipeTagStore) List(p model.RecipeTagParams) ([]model.RecipeTag, error) {
	var m model.RecipeTag
	stmt := sq.Select(m.SelectFields()...).
		From(m.TableName())
	for _, jc := range m.SelectJoins() {
		stmt = stmt.JoinClause(jc)
	}
	if p.RecipeID != 0 {
		stmt = stmt.Where(sq.Eq{"id": p.RecipeID})
	} else if len(p.RecipeIDs) > 0 {
		stmt = stmt.Where(sq.Eq{"id": p.RecipeIDs})
	}
	q, args, err := stmt.ToSql()
	if err != nil {
		return nil, err
	}
	tags := []model.RecipeTag{}
	err = r.db.Select(&tags, q, args...)
	return tags, err
}

func (r *recipeTagStore) Insert(t *model.RecipeTag, tx sqlx.Execer) error {
	var e sqlx.Execer
	if tx == nil {
		e = r.db
	} else {
		e = tx
	}
	q, args, err := sq.Insert(t.TableName()).
		SetMap(t.InsertMap()).
		ToSql()
	if err != nil {
		return err
	}
	_, err = e.Exec(q, args...)
	return err
}

func (r *recipeTagStore) DeleteAllByRecipeID(recipeID uint64, tx sqlx.Execer) (int64, error) {
	if tx == nil {
		tx = r.db
	}
	var m model.RecipeTag
	q, args, err := sq.Delete(m.TableName()).
		Where(sq.Eq{"recipe_id": recipeID}).
		ToSql()
	if err != nil {
		return 0, err
	}
	result, err := tx.Exec(q, args...)
	if err != nil {
		return 0, errors.Wrap(err, "failed to delete from recipe_tags")
	}
	cnt, err := result.RowsAffected()
	return cnt, errors.Wrap(err, "failed to inspect RowsAffected")
}

func (r *recipeTagStore) Replace(recipeID uint64, tags []model.RecipeTag, tx sqlx.Execer) error {
	if tx == nil {
		tx = r.db
	}
	tagIDs := make([]uint64, len(tags))
	for idx, t := range tags {
		tagIDs[idx] = t.TagID
	}
	del, args, err := sq.Delete("recipe_tags").
		Where(sq.Eq{"recipe_id": recipeID}).
		Where(sq.NotEq{"tag_id": tagIDs}).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "failed to build SQL delete statement")
	}
	log.Debug().Interface("args", args).Msg(del)
	_, err = tx.Exec(del, args...)
	if err != nil {
		return errors.Wrap(err, "failed to delete RecipeTags")
	}
	if len(tags) == 1 {
		return nil
	}
	stmt := sq.Insert("recipe_tags").
		Columns("recipe_id", "tag_id")
	for _, t := range tags {
		stmt = stmt.Values(recipeID, t.TagID)
	}
	stmt = stmt.Options("IGNORE")
	ins, args, err := stmt.ToSql()
	if err != nil {
		return errors.Wrap(err, "failed to build SQL insert statement")
	}
	log.Debug().Interface("args", args).Msg(ins)
	_, err = tx.Exec(ins, args...)
	return errors.Wrap(err, "failed to insert into recipe_tags")
}
