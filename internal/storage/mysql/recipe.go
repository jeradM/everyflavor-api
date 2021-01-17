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
	q, args, err := sq.Update(recipeTableName).
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
	stmt := sq.Select(recipeSelectFields...).
		From(recipeTableName).
		Where(sq.Eq{"recipes.id": id}).
		Where(sq.Expr("recipes.deleted_at IS NULL")).
		Having(sq.Expr("recipes.id > 0"))
	for _, jc := range recipeSelectJoins {
		stmt = stmt.JoinClause(jc)
	}
	query, args, err := stmt.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "failed to build select query")
	}
	err = r.store.DB().Get(&m, query, args)
	return &m, errors.Wrap(err, "failed to fetch recipe")
}

func (r *recipeStore) ByUUID(uuid string) (*model.Recipe, error) {
	var m model.Recipe
	stmt := sq.Select(recipeSelectFields...).
		From(recipeTableName).
		Where(sq.Eq{"recipes.uuid": uuid}).
		Where(sq.Eq{"recipes.current": 1}).
		Where(sq.Expr("recipes.deleted_at IS NULL")).
		Having(sq.Expr("recipes.id > 0"))
	for _, jc := range recipeSelectJoins {
		stmt = stmt.JoinClause(jc)
	}
	query, args, err := stmt.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "failed to build select query")
	}
	err = r.store.DB().Get(&m, query, args)
	return &m, errors.Wrap(err, "failed to fetch recipe")
}

func (r *recipeStore) List(params *model.RecipeParams) ([]model.Recipe, uint64, error) {
	var count uint64
	var recipes []model.Recipe
	stmt := sq.Select().From(recipeTableName)
	for _, jc := range recipeSelectJoins {
		stmt = stmt.JoinClause(jc)
	}
	params.Group = "recipes.id"
	stmt = setRecipeFilters(stmt, params)

	query, args, err := stmt.Columns("count(recipes.id)").ToSql()
	if err != nil {
		return recipes, count, errors.Wrap(err, "failed to build count query")
	}
	err = r.store.DB().Get(&count, query, args)
	if err != nil {
		return recipes, count, errors.Wrap(err, "failed to execute count query")
	}

	stmt = setRecipePagingParams(stmt, params)
	query, args, err = stmt.Columns(recipeSelectFields...).ToSql()
	if err != nil {
		return recipes, count, errors.Wrap(err, "failed to build select query")
	}
	err = r.store.DB().Select(&recipes, query, args)
	return recipes, count, errors.Wrap(err, "failed to execute select query")
}

func (r *recipeStore) Insert(recipe *model.Recipe, tx sqlx.Ext) error {
	err := r.beforeInsert(recipe, r.store.Connection())
	if err != nil {
		return errors.Wrap(err, "failed to execute beforeInsert")
	}

	query, args, err := sq.Insert(recipeTableName).
		SetMap(recipeInsertMap(recipe)).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "failed to build insert query")
	}

	result, err := r.store.DB().ExecWithTX(tx, query, args)
	if err != nil {
		return errors.Wrap(err, "failed to execute insert query")
	}
	id, err := result.LastInsertId()
	if err != nil {
		return errors.Wrap(err, "failed to get LastInsertID")
	}
	recipe.ID = uint64(id)

	err = r.afterInsert(recipe, tx)
	return errors.Wrap(err, "failed to execute afterInsert")
}

func (r *recipeStore) Update(recipe *model.Recipe, tx sqlx.Ext) error {
	query, args, err := sq.Update(recipeTableName).
		SetMap(recipeUpdateMap(recipe)).
		Where(sq.Eq{"recipes.id": recipe.ID}).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "failed to build update query")
	}

	_, err = r.store.DB().ExecWithTX(tx, query, args)
	return errors.Wrap(err, "failed to execute update query")
}

func (r *recipeStore) Publish(recipeID uint64, tx sqlx.Execer) error {
	query, args, err := sq.Update(recipeTableName).
		Set("public", 1).
		Set("published_at", sq.Expr("NOW()")).
		Where(sq.Eq{"id": recipeID}).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "failed to build update query")
	}
	_, err = r.store.DB().ExecWithTX(tx, query, args)
	return errors.Wrap(err, "failed to update recipe")
}

func (r *recipeStore) AddComment(c *model.RecipeComment) error {
	q, args, err := sq.Insert("recipe_comments").
		SetMap(map[string]interface{}{
			"content":     c.Content,
			"recipe_id":   c.RecipeID,
			"owner_id":    c.OwnerID,
			"reply_to_id": c.ReplyToID,
		}).
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
	q, args, err := sq.Insert("recipe_ratings").
		SetMap(map[string]interface{}{
			"rating":    m.Rating,
			"recipe_id": m.RecipeID,
			"owner_id":  m.OwnerID,
		}).
		ToSql()
	if err != nil {
		return err
	}
	res, err := r.store.DB().ExecWithTX(tx, q, args)
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
	stmt := sq.Select(recipeFlavorsSelectFields...).
		From(recipeFlavorsTableName)
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
	stmt := sq.Select(recipeCollaboratorsSelectFields...).
		From(recipeCollaboratorsTableName)
	for _, jc := range recipeCollaboratorsSelectJoins {
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
	stmt := sq.Select(recipeTagsSelectFields...).
		From(recipeTagsTableName)
	for _, jc := range recipeTagsSelectJoins {
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

func setRecipeFilters(q sq.SelectBuilder, p *model.RecipeParams) sq.SelectBuilder {
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

func setRecipePagingParams(q sq.SelectBuilder, p *model.RecipeParams) sq.SelectBuilder {
	if p.GetGroup() != "" {
		q = q.GroupBy(p.GetGroup())
	}
	if p.GetLimit() > 0 {
		q = q.Limit(p.GetLimit())
	}
	if p.GetOffset() > 0 {
		q = q.Offset(p.GetOffset())
	}
	if p.GetSort() != "" {
		q = q.OrderBy(p.GetSort())
	}
	return q
}

var (
	recipeTableName    = "recipes"
	recipeSelectFields = []string{
		"recipes.id",
		"recipes.created_at",
		"recipes.updated_at",
		"recipes.owner_id",
		"recipes.current",
		"recipes.description",
		"recipes.public",
		"recipes.published_at",
		"recipes.remix_of_id",
		"recipes.snv",
		"recipes.steep_days",
		"recipes.temp_f",
		"recipes.title",
		"recipes.uuid",
		"recipes.version",
		"recipes.vg_percent_m",
		"recipes.wip",
		"CAST(IFNULL(avg(recipe_ratings.rating) * 1000, 0) as UNSIGNED) as avg_rating",
		"users.username owner_username",
		"remix.title remix_of_title",
	}
	recipeSelectJoins = []string{
		"LEFT JOIN users ON users.id = recipes.owner_id",
		"LEFT JOIN recipe_ratings ON recipe_ratings.recipe_id = recipes.id",
		"LEFT JOIN recipes remix ON remix.id = recipes.remix_of_id",
	}

	recipeFlavorsTableName    = "recipe_flavors"
	recipeFlavorsSelectFields = []string{
		"recipe_flavors.id",
		"recipe_flavors.percent_m",
		"recipe_flavors.flavor_id",
		"recipe_flavors.recipe_id",
	}

	recipeCollaboratorsTableName    = "recipe_collaborators"
	recipeCollaboratorsSelectFields = []string{
		"recipe_collaborators.user_id",
		"users.username",
		"recipe_collaborators.recipe_id",
	}
	recipeCollaboratorsSelectJoins = []string{
		"JOIN users ON users.id = recipe_collaborators.user_id",
	}

	recipeTagsTableName    = "recipe_tags"
	recipeTagsSelectFields = []string{
		"recipe_tags.recipe_id",
		"recipe_tags.tag_id",
		"tags.tag",
	}
	recipeTagsSelectJoins = []string{
		"JOIN tags ON tags.id = recipe_tags.tag_id",
	}
)

func recipeInsertMap(r *model.Recipe) map[string]interface{} {
	return map[string]interface{}{
		"owner_id":     r.OwnerID,
		"current":      1,
		"description":  r.Description,
		"public":       r.Public,
		"remix_of_id":  r.RemixOfID,
		"snv":          r.Snv,
		"steep_days":   r.SteepDays,
		"temp_f":       r.TempF,
		"title":        r.Title,
		"uuid":         r.UUID,
		"version":      r.Version,
		"vg_percent_m": r.VgPercentM,
		"wip":          r.Wip,
	}
}

func recipeUpdateMap(r *model.Recipe) map[string]interface{} {
	return map[string]interface{}{
		"owner_id":     r.OwnerID,
		"description":  r.Description,
		"public":       r.Public,
		"remix_of_id":  r.RemixOfID,
		"snv":          r.Snv,
		"steep_days":   r.SteepDays,
		"temp_f":       r.TempF,
		"title":        r.Title,
		"vg_percent_m": r.VgPercentM,
		"wip":          r.Wip,
	}
}

func recipeFlavorsInsertMap(r *model.RecipeFlavor) map[string]interface{} {
	return map[string]interface{}{
		"percent_m": r.PercentM,
		"flavor_id": r.FlavorID,
		"recipe_id": r.RecipeID,
	}
}

func recipeFlavorUpdateMap(r *model.RecipeFlavor) map[string]interface{} {
	return recipeFlavorsInsertMap(r)
}

func recipeCollaboratorsInsertMap(r *model.RecipeCollaborator) map[string]interface{} {
	return map[string]interface{}{
		"recipe_id": r.RecipeID,
		"user_id":   r.UserID,
	}
}

func recipeCollaboratorsUpdateMap(r *model.RecipeCollaborator) map[string]interface{} {
	return recipeCollaboratorsInsertMap(r)
}

func recipeTagInsertMap(r *model.RecipeTag) map[string]interface{} {
	return map[string]interface{}{
		"recipe_id": r.RecipeID,
		"tag_id":    r.TagID,
	}
}

func recipeTagUpdateMap(r *model.RecipeTag) map[string]interface{} {
	return recipeTagInsertMap(r)
}
