package mysql

import (
	"everyflavor/internal/storage"
	"everyflavor/internal/storage/model"
	"golang.org/x/crypto/bcrypt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"

	"strings"
)

type userStore struct {
	store *Store
}

func NewUserStore(store *Store) storage.UserStore {
	return &userStore{store: store}
}

func (r *userStore) beforeInsertUser(u *model.User) error {
	pw := []byte(u.Password)
	_, err := bcrypt.Cost(pw)
	if err == nil {
		return nil
	}
	hash, err := bcrypt.GenerateFromPassword(pw, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return nil
}

func (r *userStore) Get(id uint64) (*model.User, error) {
	var u model.User
	err := r.store.getEntityByID(nil, &u, id, func(stmt sq.SelectBuilder) sq.SelectBuilder {
		return stmt.Where(sq.Expr("users.deleted_at IS NULL"))
	})
	return &u, err
}

func (r *userStore) List() ([]model.User, error) {
	p := model.BaseListParams{}
	var user []model.User
	_, err := r.store.listEntities(nil, &user, p, func(stmt sq.SelectBuilder) sq.SelectBuilder {
		return stmt.Where(sq.Expr("deleted_at IS NULL"))
	})
	return user, err
}

func (r *userStore) Insert(u *model.User, tx sqlx.Execer) error {
	err := r.beforeInsertUser(u)
	if err != nil {
		return err
	}
	return r.store.insertEntity(tx, u, func(id int64) {
		u.ID = uint64(id)
	})
}

func (r *userStore) Update(u *model.User, tx sqlx.Execer) error {
	_, err := r.store.updateEntity(tx, u)
	return err
}

func (r *userStore) Delete(u *model.User, tx sqlx.Execer) error {
	_, err := r.store.deleteEntity(tx, u)
	return err
}

func (r *userStore) UpdatePassword(id uint64, pw string) error {
	q, args, err := sq.Update("users").
		Set("password", pw).
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return err
	}
	_, err = r.store.DB().Exec(q, args)
	return err
}

func (r *userStore) FindByUsername(username string) (*model.User, error) {
	u := model.User{}
	q, args, err := sq.Select(u.SelectFields()...).
		From(u.TableName()).
		Where(sq.Eq{"username": username}).
		Where(sq.Expr("deleted_at IS NULL")).
		ToSql()
	if err != nil {
		return nil, err
	}
	err = r.store.DB().Get(&u, q, args)
	return &u, err
}

func (r *userStore) FindAllByUsernameLike(username string) ([]model.User, error) {
	var userModel model.User
	un := "%" + strings.ToLower(username) + "%"
	q, args, err := sq.Select("users.id", "users.username").
		From(userModel.TableName()).
		Where(sq.Expr("deleted_at IS NULL")).
		Where(sq.Expr("lower(username) like ?", un)).
		ToSql()
	if err != nil {
		return nil, err
	}
	var u []model.User
	err = r.store.DB().Select(&u, q, args)
	return u, err
}

func (r *userStore) GetStats(id uint64) (*model.UserStats, error) {
	u := model.UserStats{}
	pubSQ := sq.Select("count(r2.id)").
		From("recipes as r2").
		Where("r2.owner_id = u.id").
		Where("r2.public = 1")
	privSQ := sq.Select("count(id)").
		From("recipes").
		Where("owner_id = u.id").
		Where("public = 0")
	q, args, err := sq.Select().
		Column("u.id as user_id").
		Column("IFNULL(CAST(avg(rr.rating) * 1000 as UNSIGNED), 0) as avg_rating").
		Column(sq.Alias(pubSQ, "num_public")).
		Column(sq.Alias(privSQ, "num_private")).
		From("users as u").
		LeftJoin("recipes as r ON r.owner_id = u.id and r.public = 1").
		LeftJoin("recipe_ratings as rr ON rr.recipe_id = r.id").
		Where(sq.Expr("u.id = ?", id)).
		GroupBy("u.id").
		ToSql()
	if err != nil {
		return nil, err
	}
	err = r.store.DB().Get(&u, q, args)
	return &u, err
}

func (r *userStore) ListRoles(userIDs []uint64) ([]model.UserRole, error) {
	var m model.UserRole
	stmt := sq.Select(m.SelectFields()...).
		From(m.TableName())
	for _, jc := range m.SelectJoins() {
		stmt = stmt.JoinClause(jc)
	}
	q, args, err := stmt.Where(sq.Eq{"user_roles.user_id": userIDs}).
		ToSql()
	if err != nil {
		return nil, err
	}
	var roles []model.UserRole
	err = r.store.DB().Select(&roles, q, args)
	return roles, err
}

func (r *userStore) UsernameExists(username string) bool {
	query := "SELECT id FROM users WHERE username = ?"
	var id uint64
	_ = r.store.DB().Get(&id, query, []interface{}{username})
	return id != 0
}

func (r *userStore) EmailExists(email string) bool {
	query := "SELECT id FROM users WHERE email = ?"
	var id uint64
	_ = r.store.DB().Get(&id, query, []interface{}{email})
	return id != 0
}
