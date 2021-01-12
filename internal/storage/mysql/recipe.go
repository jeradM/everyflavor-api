package mysql

import (
	"everyflavor/internal/storage"
	"everyflavor/internal/storage/model"
	"fmt"
	"strings"

	"github.com/hashicorp/go-uuid"

	"github.com/pkg/errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type recipeStore struct {
	store *Store
}

func NewRecipeStore(store *Store) storage.RecipeStore {
	return &recipeStore{store: store}
}

func (r *recipeStore) getNextVersionForUUID(uuid string, db storage.DB) (uint64, error) {
	query, args, err := sq.Select("max(version) version").
		From("recipes").
		Where(sq.Eq{"uuid": uuid}).
		ToSql()
	if err != nil {
		return 0, err
	}
	var v uint64
	err = r.store.DB().GetWithTX(db, &v, query, args)
	if err != nil {
		return 0, err
	}
	return v + 1, nil
}

func (r *recipeStore) beforeInsert(recipe *model.Recipe, db storage.DB) error {
	if recipe.UUID == "" {
		u, err := uuid.GenerateUUID()
		if err != nil {
			return errors.Wrap(err, "error generating UUID")
		}
		recipe.UUID = u
		recipe.Version = 1
	} else {
		v, err := r.getNextVersionForUUID(recipe.UUID, db)
		if err != nil {
			return errors.Wrap(err, "unable to determine max version")
		}
		recipe.Version = v
	}
	return nil
}

func (r *recipeStore) afterInsert(recipe *model.Recipe, db sqlx.Ext) error {
	q, args, err := sq.Update(recipe.TableName()).
		Set("current", 0).
		Where(sq.Eq{"uuid": recipe.UUID}).
		Where(sq.Lt{"version": recipe.Version}).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "recipeStore.afterInsert failed")
	}
	_, err = r.store.DB().ExecWithTX(db, q, args)
	return errors.Wrap(err, "recipeStore.afterInsert failed")
}

func (r *recipeStore) Get(id uint64) (*model.Recipe, error) {
	var m model.Recipe
	err := r.store.getEntityByID(nil, &m, id, func(stmt sq.SelectBuilder) sq.SelectBuilder {
		return stmt.Where(sq.Eq{"recipes.current": 1}).
			Where(sq.Expr("recipes.deleted_at IS NULL")).
			Having(sq.Expr("recipes.id > 0"))
	})
	return &m, err
}

func (r *recipeStore) ByUUID(uuid string) (*model.Recipe, error) {
	var m model.Recipe
	stmt := sq.Select(m.SelectFields()...).
		From(m.TableName())
	for _, jc := range m.SelectJoins() {
		stmt = stmt.JoinClause(jc)
	}
	q, args, err := stmt.Where(sq.Eq{"recipes.uuid": uuid}).
		Where(sq.Eq{"recipes.current": 1}).
		Where(sq.Expr("recipes.deleted_at IS NULL")).
		ToSql()
	if err != nil {
		return nil, err
	}
	err = r.store.DB().Get(&m, q, args)
	return &m, err
}

func (r *recipeStore) List(params *model.RecipeParams) ([]model.Recipe, uint64, error) {
	var recipes []model.Recipe
	params.Group = "recipes.id"
	cnt, err := r.store.listEntities(
		nil, &recipes, params, func(stmt sq.SelectBuilder) sq.SelectBuilder {
			return r.setFilters(stmt, params)
		})
	return recipes, cnt, err
}

func (r *recipeStore) Insert(recipe *model.Recipe, tx sqlx.Ext) error {
	err := r.beforeInsert(recipe, r.store.Connection())
	if err != nil {
		return errors.Wrap(err, "failed to execute beforeInsert")
	}
	err = r.store.insertEntity(tx, recipe, func(id int64) {
		recipe.ID = uint64(id)
	})
	if err != nil {
		return err
	}
	err = r.afterInsert(recipe, tx)
	return errors.Wrap(err, "failed to execute afterInsert")
}

func (r *recipeStore) Update(recipe *model.Recipe, tx sqlx.Ext) error {
	ins, args, err := sq.Update(recipe.TableName()).
		SetMap(recipe.UpdateMap()).
		Where(sq.Eq{"id": recipe.ID}).
		ToSql()
	if err != nil {
		return err
	}
	result, err := r.store.DB().ExecWithTX(tx, ins, args)
	if err != nil {
		return err
	}
	_, err = result.LastInsertId()
	return err
}

func (r *recipeStore) AddComment(c *model.RecipeComment) error {
	q, args, err := sq.Insert(c.TableName()).
		SetMap(c.InsertMap()).
		ToSql()
	if err != nil {
		return err
	}
	res, err := r.store.DB().Exec(q, args)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	c.ID = uint64(id)
	return nil
}

func (r *recipeStore) InsertRating(m *model.RecipeRating, tx sqlx.Execer) error {
	q, args, err := sq.Insert(m.TableName()).
		SetMap(m.InsertMap()).
		ToSql()
	if err != nil {
		return err
	}
	res, err := r.store.DB().Exec(q, args)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	m.ID = uint64(id)
	return nil
}

func (r *recipeStore) ListFlavors(p model.RecipeFlavorParams) ([]model.RecipeFlavor, error) {
	var m model.RecipeFlavor
	stmt := sq.Select(m.SelectFields()...).
		From(m.TableName())
	if p.RecipeID != 0 {
		stmt = stmt.Where(sq.Eq{"recipe_flavors.recipe_id": p.RecipeID})
	} else if len(p.RecipeIDs) > 0 {
		stmt = stmt.Where(sq.Eq{"recipe_flavors.recipe_id": p.RecipeIDs})
	}
	query, args, err := stmt.ToSql()
	if err != nil {
		return nil, err
	}
	var f []model.RecipeFlavor
	err = r.store.DB().Select(&f, query, args)
	return f, err
}

func (r *recipeStore) ReplaceFlavors(recipeID uint64, flavors []model.RecipeFlavor, tx sqlx.Execer) error {
	if flavors == nil {
		return nil
	}
	var existingIDs, flavorIDs []uint64
	for _, f := range flavors {
		flavorIDs = append(flavorIDs, f.FlavorID)
		if f.ID > 0 {
			existingIDs = append(existingIDs, f.ID)
		}
	}

	delStmt := sq.Delete("recipe_flavors").
		Where(sq.Eq{"recipe_id": recipeID})

	if len(flavors) > 0 {
		delStmt = delStmt.Where(sq.NotEq{"flavor_id": flavorIDs}).
			Where(sq.NotEq{"id": existingIDs})
	}
	del, args, err := delStmt.ToSql()
	if err != nil {
		return errors.Wrap(err, "failed to build SQL delete statement")
	}
	_, err = r.store.DB().ExecWithTX(tx, del, args)
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
	stmt = stmt.SuffixExpr(sq.Expr("ON DUPLICATE KEY UPDATE " +
		"flavor_id = VALUES(flavor_id), percent_m = VALUES(percent_m)"))
	ins, args, err := stmt.ToSql()
	if err != nil {
		return errors.Wrap(err, "failed to build SQL insert statement")
	}
	_, err = r.store.DB().ExecWithTX(tx, ins, args)

	return errors.Wrap(err, "failed to insert into recipe_flavors")
}

func (r *recipeStore) ListCollaborators(p model.RecipeCollaboratorParams) ([]model.RecipeCollaborator, error) {
	var m model.RecipeCollaborator
	stmt := sq.Select(m.SelectFields()...).
		From(m.TableName())
	for _, jc := range m.SelectJoins() {
		stmt = stmt.JoinClause(jc)
	}
	if p.RecipeID != 0 {
		stmt = stmt.Where(sq.Eq{"recipe_collaborators.recipe_id": p.RecipeID})
	} else if len(p.RecipeIDs) > 0 {
		stmt = stmt.Where(sq.Eq{"recipe_collaborators.recipe_id": p.RecipeIDs})
	}
	q, args, err := stmt.ToSql()
	if err != nil {
		return nil, err
	}
	var u []model.RecipeCollaborator
	err = r.store.DB().Select(&u, q, args)
	return u, err
}

func (r *recipeStore) ReplaceCollaborators(recipeID uint64, collabs []model.RecipeCollaborator, tx sqlx.Execer) error {
	if collabs == nil {
		return nil
	}
	userIDs := make([]uint64, len(collabs))
	for idx, t := range collabs {
		userIDs[idx] = t.UserID
	}

	delStmt := sq.Delete("recipe_collaborators").
		Where(sq.Eq{"recipe_id": recipeID})
	if len(userIDs) > 0 {
		delStmt = delStmt.Where(sq.NotEq{"user_id": userIDs})
	}
	del, args, err := delStmt.ToSql()
	if err != nil {
		return errors.Wrap(err, "failed to build SQL delete statement")
	}
	_, err = r.store.DB().ExecWithTX(tx, del, args)
	if err != nil {
		return errors.Wrap(err, "failed to delete RecipeCollaborators")
	}

	if len(collabs) == 0 {
		return nil
	}

	stmt := sq.Insert("recipe_collaborators").
		Options("IGNORE").
		Columns("recipe_id", "user_id")
	for _, t := range collabs {
		stmt = stmt.Values(recipeID, t.UserID)
	}
	ins, args, err := stmt.ToSql()
	if err != nil {
		return errors.Wrap(err, "failed to build SQL insert statement")
	}
	_, err = r.store.DB().ExecWithTX(tx, ins, args)

	return errors.Wrap(err, "failed to insert into recipe_collaborators")
}

func (r *recipeStore) ListTags(p model.RecipeTagParams) ([]model.RecipeTag, error) {
	var m model.RecipeTag
	stmt := sq.Select(m.SelectFields()...).
		From(m.TableName())
	for _, jc := range m.SelectJoins() {
		stmt = stmt.JoinClause(jc)
	}
	if p.RecipeID != 0 {
		stmt = stmt.Where(sq.Eq{"recipe_tags.recipe_id": p.RecipeID})
	} else if len(p.RecipeIDs) > 0 {
		stmt = stmt.Where(sq.Eq{"recipe_tags.recipe_id": p.RecipeIDs})
	}
	q, args, err := stmt.ToSql()
	if err != nil {
		return nil, err
	}
	var tags []model.RecipeTag
	err = r.store.DB().Select(&tags, q, args)
	return tags, err
}

func (r *recipeStore) ReplaceTags(recipeID uint64, tags []model.RecipeTag, tx sqlx.Execer) error {
	if tags == nil {
		return nil
	}
	tagIDs := make([]uint64, len(tags))
	for idx, t := range tags {
		tagIDs[idx] = t.TagID
	}

	delStmt := sq.Delete("recipe_tags").
		Where(sq.Eq{"recipe_id": recipeID})
	if len(tagIDs) > 0 {
		delStmt = delStmt.Where(sq.NotEq{"tag_id": tagIDs})
	}
	del, args, err := delStmt.ToSql()
	if err != nil {
		return errors.Wrap(err, "failed to build SQL delete statement")
	}
	_, err = r.store.DB().ExecWithTX(tx, del, args)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to delete tags for recipe: %d", recipeID))
	}

	if len(tags) == 0 {
		return nil
	}

	stmt := sq.Insert("recipe_tags").
		Options("IGNORE").
		Columns("recipe_id", "tag_id")
	for _, t := range tags {
		stmt = stmt.Values(recipeID, t.TagID)
	}
	ins, args, err := stmt.ToSql()
	if err != nil {
		return errors.Wrap(err, "failed to build SQL insert statement")
	}
	_, err = r.store.DB().ExecWithTX(tx, ins, args)

	return errors.Wrap(err, "failed to insert into recipe_tags")
}

func (r *recipeStore) setFilters(q sq.SelectBuilder, p *model.RecipeParams) sq.SelectBuilder {
	if p.Title != "" {
		q = q.Where("lower(recipes.title) like ?", "%"+strings.ToLower(p.Title)+"%")
	}
	if p.UserID != nil {
		q = q.Where(sq.Eq{"recipes.owner_id": p.UserID})
	}
	if p.ExcludeFlavors != nil {
		exc := sq.Select("recipe_flavors as rf1").
			Where(sq.Expr("rf1.recipe_id = r.id")).
			Where(sq.Eq{"rf1.flavor_id": p.ExcludeFlavors})
		q = q.Where(sq.Expr("not exists(?)", exc))
	}
	if p.IncludeFlavors != nil {
		inc := sq.Select("count(id)").
			From("recipe_flavors as rf2").
			Where(sq.Expr("rf2.recipe_id = r.id")).
			Where(sq.Eq{"rf2.flavor_id": p.IncludeFlavors})
		q = q.Where("(?) = ?", inc, len(p.IncludeFlavors))
	}
	if p.CreatedFrom != nil {
		q = q.Where(sq.Expr("recipes.created_at >= ?", p.CreatedFrom))
	}
	if p.CreatedTo != nil {
		q = q.Where(sq.Expr("recipes.created_at <= ?", p.CreatedTo))
	}
	if p.Snv != nil {
		q = q.Where(sq.Expr("recipes.snv = ?", p.Snv))
	}
	if p.Public {
		q = q.Where(sq.Expr("recipes.public = ?", 1))
	}
	if p.Current {
		q = q.Where(sq.Expr("recipes.current = 1"))
	}
	return q
}

func (r *recipeStore) setPagingParams(q sq.SelectBuilder, p *model.RecipeParams) sq.SelectBuilder {
	var offset uint64
	var limit uint64
	var order string
	if p.Limit <= 0 {
		limit = 25
	} else if p.Limit > 100 {
		limit = 100
	} else {
		limit = p.Limit
	}
	if p.Offset < 0 {
		offset = 0
	} else {
		offset = p.Offset
	}
	q = q.Limit(limit).Offset(offset)
	if strings.ToLower(p.Order) != "desc" {
		order = "asc"
	} else {
		order = "desc"
	}
	if p.Sort != "" {
		sort := p.Sort
		if sort == "rating" {
			sort = "avg_rating"
		} else if !strings.Contains(sort, ".") {
			sort = "recipes." + sort
		}
		q = q.OrderBy(sort + " " + order)
	}

	return q
}
