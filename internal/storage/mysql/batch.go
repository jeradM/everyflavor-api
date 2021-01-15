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
	query, args, err := sq.Select(batchSelectFields...).
		From("batches").
		Where(sq.Eq{"batches.id": id}).
		ToSql()
	if err != nil {
		return batch, errors.Wrap(err, "failed to build select query")
	}
	err = r.store.DB().Get(&batch, query, args)
	return batch, err
}

func (s *batchStore) List(userID uint64) ([]model.Batch, error) {
	var batches []model.Batch
	query, args, err := sq.Select(batchSelectFields...).
		From("batches").
		Where(sq.Expr("batches.deleted_at IS NULL")).
		Where(sq.Eq{"batches.owner_id": userID}).
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "failed to build select query")
	}
	err = s.store.DB().Select(&batches, query, args)
	return batches, errors.Wrap(err, "failed to fetch records")
}

func (r *batchStore) Insert(b *model.Batch, tx sqlx.Execer) error {
	query, args, err := sq.Insert("batches").
		SetMap(batchInsertMap(b)).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "failed to build insert query")
	}
	result, err := r.store.DB().ExecWithTX(tx, query, args)
	if err != nil {
		return errors.Wrap(err, "insert failed: Batch")
	}
	i, err := result.LastInsertId()
	if err != nil {
		return errors.Wrap(err, "failed to get LastInsertID")
	}
	b.ID = uint64(i)
	return nil
}

func (r *batchStore) Update(b *model.Batch, tx sqlx.Execer) error {
	query, args, err := sq.Update("batches").
		SetMap(batchUpdateMap(b)).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "failed to build update query")
	}
	result, err := r.store.DB().ExecWithTX(tx, query, args)
	if err != nil {
		return errors.Wrap(err, "update failed: Batch")
	}
	_, err = result.RowsAffected()
	return errors.Wrap(err, "failed to get RowsAffected")
}

func (r *batchStore) ListFlavorsForBatches(batchIDs []uint64) ([]model.BatchFlavor, error) {
	var flavors []model.BatchFlavor
	query, args, err := sq.Select(batchFlavorSelectFields...).
		From("batch_flavors").
		Where(sq.Eq{"batch_flavors.batch_id": batchIDs}).
		ToSql()
	if err != nil {
		return flavors, errors.Wrap(err, "failed to generate SQL query string")
	}
	err = r.store.DB().Select(&flavors, query, args)
	return flavors, errors.Wrap(err, "failed to fetch BatchFlavors")
}

func (r *batchStore) InsertFlavor(flavor *model.BatchFlavor, tx sqlx.Execer) error {
	query, args, err := sq.Insert("batch_flavors").
		SetMap(batchFlavorInsertMap(flavor)).
		SuffixExpr(sq.Expr("ON DUPLICATE KEY UPDATE percent_m = ?, vg = ?", flavor.PercentM, flavor.Vg)).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "failed to build insert query")
	}
	result, err := r.store.DB().ExecWithTX(tx, query, args)
	if err != nil {
		return errors.Wrap(err, "failed to insert BatchFlavor")
	}
	i, err := result.LastInsertId()
	if err != nil {
		return errors.Wrap(err, "failed to get LastInsertID")
	}
	flavor.ID = uint64(i)
	return nil
}

var (
	batchSelectFields = []string{
		"batches.id",
		"batches.created_at",
		"batches.updated_at",
		"batches.batch_size_m",
		"batches.batch_strength",
		"batches.batch_vg_m",
		"batches.max_vg",
		"batches.nic_strength",
		"batches.nic_vg_m",
		"batches.recipe_id",
		"batches.owner_id",
		"batches.use_nic",
	}

	batchFlavorSelectFields = []string{
		"batch_flavors.id",
		"batch_flavors.percent_m",
		"batch_flavors.vg",
		"batch_flavors.flavor_id",
		"batch_flavors.batch_id",
	}
)

func batchInsertMap(b *model.Batch) map[string]interface{} {
	return map[string]interface{}{
		"batch_size_m":   b.BatchSizeM,
		"batch_strength": b.BatchStrength,
		"batch_vg_m":     b.BatchVgM,
		"max_vg":         b.MaxVg,
		"nic_strength":   b.NicStrength,
		"nic_vg_m":       b.NicVgM,
		"recipe_id":      b.RecipeID,
		"owner_id":       b.OwnerID,
		"use_nic":        b.UseNic,
	}
}

func batchUpdateMap(b *model.Batch) map[string]interface{} {
	return map[string]interface{}{
		"batch_size_m":   b.BatchSizeM,
		"batch_strength": b.BatchStrength,
		"batch_vg_m":     b.BatchVgM,
		"max_vg":         b.MaxVg,
		"nic_strength":   b.NicStrength,
		"nic_vg_m":       b.NicVgM,
		"use_nic":        b.UseNic,
	}
}

func batchFlavorInsertMap(b *model.BatchFlavor) map[string]interface{} {
	return map[string]interface{}{
		"percent_m": b.PercentM,
		"vg":        b.Vg,
		"flavor_id": b.FlavorID,
		"batch_id":  b.BatchID,
	}
}

func batchFlavorUpdateMap(b *model.BatchFlavor) map[string]interface{} {
	return map[string]interface{}{
		"percent_m": b.PercentM,
		"vg":        b.Vg,
	}
}
