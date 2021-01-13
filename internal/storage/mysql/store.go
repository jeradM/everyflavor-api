package mysql

import (
	"database/sql"
	"database/sql/driver"
	"everyflavor/internal/storage"
	"everyflavor/internal/storage/model"
	"reflect"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	queryLogKey    = "query"
	argsLogKey     = "args"
	durationLogKey = "duration"
)

type SQLLoggingQueryExecer struct {
	enabled    bool
	connection storage.DB
}

func (l *SQLLoggingQueryExecer) Enable(e bool) {
	l.enabled = e
}

func (l *SQLLoggingQueryExecer) log(e error, q string, a interface{}, t time.Time, fns ...func(*zerolog.Event)) {
	if l.enabled {
		stmt := log.Debug().
			Str(queryLogKey, q).
			Interface(argsLogKey, a).
			Dur(durationLogKey, time.Since(t)).
			Err(e)
		for _, fn := range fns {
			fn(stmt)
		}
		stmt.Msg("")
	}
}

func (l *SQLLoggingQueryExecer) Get(dest interface{}, query string, args []interface{}) error {
	return l.GetWithTX(l.connection, dest, query, args)
}

func (l *SQLLoggingQueryExecer) GetWithTX(tx sqlx.Queryer, dest interface{}, query string, args []interface{}) error {
	if tx == nil {
		tx = l.connection
	}
	t := reflect.TypeOf(dest)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	start := time.Now()
	err := sqlx.Get(tx, dest, query, args...)
	l.log(err, query, args, start, func(e *zerolog.Event) {
		e.Str("modelType", t.String())
	})
	return err
}

func (l *SQLLoggingQueryExecer) Select(dest interface{}, query string, args []interface{}) error {
	return l.SelectWithTX(l.connection, dest, query, args)
}

func (l *SQLLoggingQueryExecer) SelectWithTX(tx sqlx.Queryer, dest interface{}, query string, args []interface{}) error {
	if tx == nil {
		tx = l.connection
	}
	t := reflect.TypeOf(dest)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	start := time.Now()
	err := sqlx.Select(tx, dest, query, args...)
	l.log(err, query, args, start, func(e *zerolog.Event) {
		e.Str("modelType", t.String())
	})
	return err
}

func (l *SQLLoggingQueryExecer) Exec(query string, args []interface{}) (driver.Result, error) {
	return l.ExecWithTX(l.connection, query, args)
}

func (l *SQLLoggingQueryExecer) ExecWithTX(tx sqlx.Execer, query string, args []interface{}) (driver.Result, error) {
	if tx == nil {
		tx = l.connection
	}
	start := time.Now()
	res, err := tx.Exec(query, args...)
	l.log(err, query, args, start)
	return res, err
}

func (l *SQLLoggingQueryExecer) Query(tx sqlx.Queryer, query string, args []interface{}) (*sqlx.Rows, error) {
	start := time.Now()
	rows, err := tx.Queryx(query, args...)
	l.log(err, query, args, start)
	return rows, err
}

type Store struct {
	pool   storage.DB
	db     storage.DBQueryExecer
	auth   storage.AuthStore
	batch  storage.BatchStore
	flavor storage.FlavorStore
	recipe storage.RecipeStore
	stash  storage.FlavorStashStore
	tag    storage.TagStore
	user   storage.UserStore
	vendor storage.VendorStore
}

func NewMySQLStore(conn storage.DB, logQueries bool) *Store {
	store := &Store{pool: conn}
	logger := SQLLoggingQueryExecer{enabled: logQueries, connection: conn}
	store.db = &logger
	store.auth = NewAuthStore(store)
	store.batch = NewBatchStore(store)
	store.flavor = NewFlavorStore(store)
	store.recipe = NewRecipeStore(store)
	store.stash = NewFlavorStashStore(store)
	store.tag = NewTagStore(store)
	store.user = NewUserStore(store)
	store.vendor = NewVendorStore(store)
	return store
}

func (s *Store) Connection() storage.DB {
	return s.pool
}

func (s *Store) DB() storage.DBQueryExecer {
	return s.db
}

func (s *Store) Auth() storage.AuthStore {
	return s.auth
}

func (s *Store) Batch() storage.BatchStore {
	return s.batch
}

func (s *Store) Flavor() storage.FlavorStore {
	return s.flavor
}

func (s *Store) Recipe() storage.RecipeStore {
	return s.recipe
}

func (s *Store) Stash() storage.FlavorStashStore {
	return s.stash
}

func (s *Store) Tag() storage.TagStore {
	return s.tag
}

func (s *Store) User() storage.UserStore {
	return s.user
}

func (s *Store) Vendor() storage.VendorStore {
	return s.vendor
}

type SelectStmtFn func(builder sq.SelectBuilder) sq.SelectBuilder
type InsertStmtFn func(builder sq.InsertBuilder) sq.InsertBuilder
type UpdateStmtFn func(builder sq.UpdateBuilder) sq.UpdateBuilder
type ResultFn func(int64)
type ListResultFn func(uint64)

func (s *Store) getEntityByID(queryer sqlx.Queryer, entity model.Entity, id uint64, fns ...SelectStmtFn) error {
	stmt := sq.Select(entity.SelectFields()...).
		From(entity.TableName())
	for _, jc := range entity.SelectJoins() {
		stmt = stmt.JoinClause(jc)
	}
	for _, fn := range fns {
		stmt = fn(stmt)
	}
	query, args, err := stmt.Where(sq.Eq{entity.TableName() + ".id": id}).ToSql()
	if err != nil {
		return errors.Wrap(err, "failed to build SQL query string")
	}
	if queryer == nil {
		err = s.DB().Get(entity, query, args)
	} else {
		err = s.db.GetWithTX(queryer, entity, query, args)
	}
	return errors.Wrap(err, "failed to get entity")
}

func getEntityFromDest(dest interface{}) (model.Entity, error) {
	val := reflect.ValueOf(dest)
	if val.Kind() != reflect.Ptr {
		return nil, errors.New("destination must be a *slice")
	}
	val = val.Elem()
	if val.Kind() != reflect.Slice {
		return nil, errors.New("destination must be a *slice")
	}
	t := val.Type().Elem()
	z := reflect.Zero(t).Interface()
	ent, ok := z.(model.Entity)
	if !ok {
		return nil, errors.New("destination must implement model.Entity")
	}
	return ent, nil
}

func (s *Store) listEntities(queryer sqlx.Queryer, dest interface{}, params model.ListParams, fns ...SelectStmtFn) (uint64, error) {
	entity, err := getEntityFromDest(dest)
	if err != nil {
		return 0, err
	}

	// Build base select for count and list queries
	stmt := sq.Select().From(entity.TableName())
	for _, jc := range entity.SelectJoins() {
		stmt = stmt.JoinClause(jc)
	}
	for _, fn := range fns {
		stmt = fn(stmt)
	}
	// Get count
	cnt := 0
	cntQuery, args, err := stmt.Columns("count(distinct " + entity.TableName() + ".id)").ToSql()
	if err != nil {
		return uint64(cnt), errors.Wrap(err, "failed to build SQL count query string")
	}
	if queryer == nil {
		err = s.DB().Get(&cnt, cntQuery, args)
	} else {
		err = s.db.GetWithTX(queryer, &cnt, cntQuery, args)
	}
	if err != nil {
		return uint64(cnt), errors.Wrap(err, "failed to fetch count")
	}

	// Get list of entities
	stmt = setListParams(stmt, params)
	lstQuery, args, err := stmt.Columns(entity.SelectFields()...).ToSql()
	if err != nil {
		return uint64(cnt), errors.Wrap(err, "failed to build SQL select query string")
	}
	if queryer == nil {
		err = s.DB().Select(dest, lstQuery, args)
	} else {
		err = s.db.GetWithTX(queryer, dest, lstQuery, args)
	}

	return uint64(cnt), errors.Wrap(err, "failed to fetch entities")
}

func (s *Store) insertEntity(e sqlx.Execer, m model.Entity, handleID ResultFn, fns ...InsertStmtFn) error {
	stmt := sq.Insert(m.TableName()).
		SetMap(m.InsertMap())
	for _, fn := range fns {
		stmt = fn(stmt)
	}
	q, args, err := stmt.ToSql()
	if err != nil {
		return err
	}
	var result sql.Result
	if e == nil {
		result, err = s.db.Exec(q, args)
	} else {
		result, err = s.db.ExecWithTX(e, q, args)
	}
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return errors.Wrap(err, "failed to fetch last insert id")
	}
	if handleID != nil {
		handleID(id)
	}
	return nil
}

func (s *Store) updateEntity(tx sqlx.Execer, m model.Entity, stmtFns ...UpdateStmtFn) (int64, error) {
	stmt := sq.Update(m.TableName()).
		SetMap(m.UpdateMap()).
		Where(sq.Eq{"id": m.GetID()})
	for _, fn := range stmtFns {
		stmt = fn(stmt)
	}
	query, args, err := stmt.ToSql()
	if err != nil {
		return 0, errors.Wrap(err, "failed to build SQL update statement")
	}
	var result sql.Result
	if tx == nil {
		result, err = s.db.Exec(query, args)
	} else {
		result, err = s.db.ExecWithTX(tx, query, args)
	}
	if err != nil {
		return 0, errors.Wrap(err, "failed to execute update statement")
	}
	n, err := result.RowsAffected()
	return n, errors.Wrap(err, "failed to inspect number of affected rows")
}

func (s *Store) deleteEntity(tx sqlx.Execer, m model.Entity) (int64, error) {
	query, args, err := sq.Delete(m.TableName()).
		Where(sq.Eq{"id": m.GetID()}).
		ToSql()
	if err != nil {
		return 0, errors.Wrap(err, "failed to build SQL delete statement")
	}
	result, err := s.DB().ExecWithTX(tx, query, args)
	if err != nil {
		return 0, errors.Wrap(err, "failed to execute SQL delete statement")
	}
	n, err := result.RowsAffected()
	return n, errors.Wrap(err, "failed to inspect number of affected rows")
}

func (s *Store) softDeleteEntity(tx sqlx.Execer, m model.Entity) (int64, error) {
	query, args, err := sq.Update(m.TableName()).
		Set("deleted_at", sq.Expr("NOW")).
		Where(sq.Eq{"id": m.GetID()}).
		ToSql()
	if err != nil {
		return 0, errors.Wrap(err, "failed to build SQL update statement")
	}
	result, err := s.DB().ExecWithTX(tx, query, args)
	if err != nil {
		return 0, errors.Wrap(err, "failed to execute SQL update statement")
	}
	n, err := result.RowsAffected()
	return n, errors.Wrap(err, "failed to inspect number of affected rows")
}

func setListParams(q sq.SelectBuilder, p model.ListParams) sq.SelectBuilder {
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
