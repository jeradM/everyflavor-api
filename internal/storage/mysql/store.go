package mysql

import (
	"database/sql/driver"
	"everyflavor/internal/storage"
	"reflect"
	"time"

	"github.com/jmoiron/sqlx"
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

func NewSQLLoggingerQueryExecer(c storage.DB, e bool) SQLLoggingQueryExecer {
	return SQLLoggingQueryExecer{connection: c, enabled: e}
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

func NewMySQLStore(conn storage.DB, execer storage.DBQueryExecer) *Store {
	store := &Store{pool: conn}
	logger := execer
	store.db = logger
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
