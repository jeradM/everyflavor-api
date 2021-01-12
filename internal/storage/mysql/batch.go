package mysql

import (
	"everyflavor/internal/storage"
	"everyflavor/internal/storage/model"

	"github.com/pkg/errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type batchStore struct {
	db    *sqlx.DB
	store *Store
}

// NewBatchStore creates a new batchStore
func NewBatchStore(store *Store) storage.BatchStore {
	return &batchStore{store: store}
}

func (r *batchStore) Get(id uint64) (model.Batch, error) {
	var batch model.Batch
	err := r.store.getEntityByID(nil, &batch, id)
	return batch, err
}

func (s *batchStore) List(userID uint64) ([]model.Batch, error) {
	var m model.Batch
	stmt := sq.Select(m.SelectFields()...).From(m.TableName())
	for _, jc := range m.SelectJoins() {
		stmt = stmt.JoinClause(jc)
	}
	q, args, err := stmt.Where(sq.Eq{"owner_id": userID}).ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "failed to build query")
	}
	var batches []model.Batch
	err = s.store.DB().Select(&batches, q, args)
	return batches, errors.Wrap(err, "failed to fetch record")
}

func (r *batchStore) Insert(b *model.Batch, tx sqlx.Execer) error {
	return r.store.insertEntity(tx, b, func(id int64) {
		b.ID = uint64(id)
	})
}

func (r *batchStore) Update(b *model.Batch, tx sqlx.Execer) error {
	_, err := r.store.updateEntity(tx, b)
	return err
}

func (r *batchStore) ListFlavorsForBatches(batchIDs []uint64) ([]model.BatchFlavor, error) {
	var m model.BatchFlavor
	query, args, err := sq.Select(m.SelectFields()...).
		From(m.TableName()).
		Where(sq.Eq{"batch_id": batchIDs}).
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate SQL query string")
	}
	var flavors []model.BatchFlavor
	err = r.store.DB().Select(&flavors, query, args)
	return flavors, errors.Wrap(err, "failed to fetch BatchFlavors")
}

func (r *batchStore) InsertFlavor(m *model.BatchFlavor, tx sqlx.Execer) error {
	var fn InsertStmtFn = func(stmt sq.InsertBuilder) sq.InsertBuilder {
		return stmt.SuffixExpr(sq.Expr("ON DUPLICATE KEY UPDATE percent_m = ?, vg = ?", m.PercentM, m.Vg))
	}
	err := r.store.insertEntity(tx, m, func(id int64) { m.ID = uint64(id) }, fn)
	return errors.Wrap(err, "InsertFlavor() failed")
}
