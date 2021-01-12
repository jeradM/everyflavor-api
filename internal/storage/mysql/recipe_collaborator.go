package mysql

import (
	"everyflavor/internal/storage"
	"everyflavor/internal/storage/model"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type recipeCollaboratorStore struct {
	db *sqlx.DB
}

func NewRecipeCollaboratorStore(db *sqlx.DB) storage.RecipeCollaboratorStore {
	return &recipeCollaboratorStore{db: db}
}

func (r *recipeCollaboratorStore) Replace(recipeID uint64, collabs []model.RecipeCollaborator, tx sqlx.Execer) error {
	if tx == nil {
		tx = r.db
	}
	userIDs := make([]uint64, len(collabs))
	for idx, t := range collabs {
		userIDs[idx] = t.UserID
	}
	del, args, err := sq.Delete("recipe_collaborators").
		Where(sq.Eq{"recipe_id": recipeID}).
		Where(sq.NotEq{"user_id": userIDs}).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "failed to build SQL delete statement")
	}
	log.Debug().Interface("args", args).Msg(del)
	_, err = tx.Exec(del, args...)
	if err != nil {
		return errors.Wrap(err, "failed to delete RecipeCollaborators")
	}
	if len(collabs) == 1 {
		return nil
	}
	stmt := sq.Insert("recipe_collaborators").
		Columns("recipe_id", "user_id")
	for _, t := range collabs {
		stmt = stmt.Values(recipeID, t.UserID)
	}
	stmt = stmt.Options("IGNORE")
	stmt.Options()
	ins, args, err := stmt.ToSql()
	if err != nil {
		return errors.Wrap(err, "failed to build SQL insert statement")
	}
	log.Debug().Interface("args", args).Msg(ins)
	_, err = tx.Exec(ins, args...)
	return errors.Wrap(err, "failed to insert into recipe_collaborators")
}
