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
	query, args, err := sq.Select(userSelectFields...).
		From("users").
		Where(sq.Eq{"users.id": id}).
		ToSql()
	if err != nil {
		return &u, err
	}
	err = r.store.DB().Get(&u, query, args)
	return &u, err
}

func (r *userStore) List() ([]model.User, error) {
	var users []model.User
	query, args, err := sq.Select(userSelectFields...).
		From("users").
		Where(sq.Expr("users.deleted_at IS NULL")).
		ToSql()
	if err != nil {
		return users, err
	}
	err = r.store.DB().Select(&users, query, args)
	return users, err
}

func (r *userStore) Insert(u *model.User, tx sqlx.Execer) error {
	err := r.beforeInsertUser(u)
	if err != nil {
		return err
	}
	query, args, err := sq.Insert("users").
		SetMap(userInsertMap(u)).
		ToSql()
	if err != nil {
		return err
	}
	_, err = r.store.DB().ExecWithTX(tx, query, args)
	return err
}

func (r *userStore) Update(u *model.User, tx sqlx.Execer) error {
	query, args, err := sq.Update("users").
		SetMap(userUpdateMap(u)).
		ToSql()
	if err != nil {
		return err
	}
	_, err = r.store.DB().ExecWithTX(tx, query, args)
	return err
}

func (r *userStore) Delete(u *model.User, tx sqlx.Execer) error {
	query, args, err := sq.Delete("users").
		Where(sq.Eq{"users.id": u.ID}).
		ToSql()
	if err != nil {
		return err
	}
	_, err = r.store.DB().ExecWithTX(tx, query, args)
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
	var u model.User
	q, args, err := sq.Select(userSelectFields...).
		From("users").
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
	un := "%" + strings.ToLower(username) + "%"
	q, args, err := sq.Select("users.id", "users.username").
		From("users").
		Where(sq.Expr("users.deleted_at IS NULL")).
		Where(sq.Expr("lower(users.username) like ?", un)).
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
	stmt := sq.Select(userRoleSelectFields...).
		From("user_roles")
	for _, jc := range userRoleSelectJoins {
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

var (
	userSelectFields = []string{
		"users.id",
		"users.created_at",
		"users.updated_at",
		"users.username",
		"users.email",
		"users.password",
	}
	roleSelectFields = []string{
		"roles.id",
		"roles.authority",
	}
	userRoleSelectFields = []string{
		"user_roles.user_id",
		"user_roles.role_id",
		"roles.authority",
	}
	userRoleSelectJoins = []string{
		"JOIN roles ON roles.id = user_roles.role_id",
	}
)

func userInsertMap(u *model.User) map[string]interface{} {
	return map[string]interface{}{
		"created_at": sq.Expr("NOW()"),
		"updated_at": sq.Expr("NOW()"),
		"username":   u.Username,
		"email":      u.Email,
		"password":   u.Password,
	}
}

func userUpdateMap(u *model.User) map[string]interface{} {
	return map[string]interface{}{
		"updated_at": sq.Expr("NOW()"),
		"username":   u.Username,
		"email":      u.Email,
		"password":   u.Password,
	}
}

func roleInsertMap(r *model.Role) map[string]interface{} {
	return map[string]interface{}{"authority": r.Authority}
}

func roleUpdateMap(r *model.Role) map[string]interface{} {
	return roleInsertMap(r)
}

func userRoleInsertMap(u *model.UserRole) map[string]interface{} {
	return map[string]interface{}{
		"user_id": u.UserID,
		"role_id": u.RoleID,
	}
}

func userRoleUpdateMap(u *model.UserRole) map[string]interface{} {
	return userRoleInsertMap(u)
}
