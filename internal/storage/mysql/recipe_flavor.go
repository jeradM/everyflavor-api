package mysql

import (
	"everyflavor/internal/storage"
	"everyflavor/internal/storage/model"

	"github.com/rs/zerolog/log"

	"github.com/pkg/errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type recipeFlavorStore struct {
	db *sqlx.DB
}

func NewRecipeFlavorStore(db *sqlx.DB) storage.RecipeFlavorStore {
	return &recipeFlavorStore{db: db}
}

func (r *recipeFlavorStore) Get(id uint64) (*model.RecipeFlavor, error) {
	f := model.RecipeFlavor{}
	q, args, err := sq.Select(f.SelectFields()...).
		From(f.TableName()).
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, err
	}
	err = r.db.Get(&f, q, args...)
	return &f, err
}

func (r *recipeFlavorStore) List(p model.RecipeFlavorParams) ([]model.RecipeFlavor, error) {
	var m model.RecipeFlavor
	lst := sq.Select(m.SelectFields()...).
		From(m.TableName())
	if p.RecipeIDs != nil {
		lst = lst.Where(sq.Eq{"recipe_id": p.RecipeIDs})
	}
	q, args, err := lst.ToSql()
	if err != nil {
		return nil, err
	}
	f := []model.RecipeFlavor{}
	err = r.db.Select(&f, q, args...)
	return f, err
}

func (r *recipeFlavorStore) Insert(f *model.RecipeFlavor, tx sqlx.Execer) error {
	if tx == nil {
		tx = r.db
	}
	q, args, err := sq.Insert(f.TableName()).
		SetMap(f.InsertMap()).
		SuffixExpr(sq.Expr("ON DUPLICATE KEY UPDATE percent_m = ?", f.PercentM)).
		ToSql()
	if err != nil {
		return err
	}
	res, err := tx.Exec(q, args...)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	f.ID = uint64(id)
	return nil
}

func (r *recipeFlavorStore) Update(f *model.RecipeFlavor, tx sqlx.Execer) error {
	if tx == nil {
		tx = r.db
	}
	q, args, err := sq.Insert(f.TableName()).
		SetMap(f.UpdateMap()).
		ToSql()
	if err != nil {
		return err
	}
	_, err = tx.Exec(q, args...)
	if err != nil {
		return err
	}
	return nil
}

func (r *recipeFlavorStore) Replace(recipeID uint64, flavors []model.RecipeFlavor, tx sqlx.Execer) error {
	if tx == nil {
		tx = r.db
	}
	var existingIDs, flavorIDs []uint64
	for _, f := range flavors {
		flavorIDs = append(flavorIDs, f.FlavorID)
		if f.ID > 0 {
			existingIDs = append(existingIDs, f.ID)
		}
	}
	del, args, err := sq.Delete("recipe_flavors").
		Where(sq.Eq{"recipe_id": recipeID}).
		Where(sq.NotEq{"flavor_id": flavorIDs}).
		Where(sq.NotEq{"id": existingIDs}).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "failed to build SQL delete statement")
	}
	log.Debug().Interface("args", args).Msg(del)
	_, err = tx.Exec(del, args...)
	if err != nil {
		return errors.Wrap(err, "failed to delete RecipeFlavors")
	}
	if len(flavors) == 0 {
		return nil
	}
	stmt := sq.Insert("recipe_flavors").
		Columns("id", "recipe_id", "flavor_id", "percent_m")
	for _, f := range flavors {
		stmt = stmt.Values(f.ID, recipeID, f.FlavorID, f.PercentM)
	}
	stmt = stmt.SuffixExpr(sq.Expr("ON DUPLICATE KEY UPDATE flavor_id = VALUES(flavor_id), percent_m = VALUES(percent_m)"))
	ins, args, err := stmt.ToSql()
	if err != nil {
		return errors.Wrap(err, "failed to build SQL insert statement")
	}
	log.Debug().Interface("args", args).Msg(ins)
	_, err = tx.Exec(ins, args...)
	return errors.Wrap(err, "failed to insert into recipe_flavors")
}

func (r *recipeFlavorStore) DeleteAllByRecipeID(recipeID uint64, tx sqlx.Execer) (int64, error) {
	if tx == nil {
		tx = r.db
	}
	q, args, err := sq.Delete("recipe_flavors").
		Where(sq.Eq{"recipe_id": recipeID}).
		ToSql()
	if err != nil {
		return 0, errors.Wrap(err, "failed to build SQL delete statement")
	}
	result, err := tx.Exec(q, args...)
	if err != nil {
		return 0, errors.Wrap(err, "failed to delete from recipe_flavors")
	}
	cnt, err := result.RowsAffected()
	return cnt, errors.Wrap(err, "failed to inspect RowsAffected")
}
