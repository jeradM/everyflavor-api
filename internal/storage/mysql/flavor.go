package mysql

import (
	"everyflavor/internal/storage"
	"everyflavor/internal/storage/model"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type flavorStore struct {
	store *Store
}

// NewFlavorStore creates and returns a new mysql-based flavor store
func NewFlavorStore(store *Store) storage.FlavorStore {
	return &flavorStore{store: store}
}

func (r *flavorStore) Get(id uint64) (model.Flavor, error) {
	var m model.Flavor
	err := r.store.getEntityByID(nil, &m, id, func(stmt sq.SelectBuilder) sq.SelectBuilder {
		return stmt.Where(sq.Expr("flavors.deleted_at IS NULL"))
	})
	return m, err
}

func (r *flavorStore) List() ([]model.Flavor, uint64, error) {
	var flavors []model.Flavor
	cnt, err := r.store.listEntities(
		nil,
		&flavors,
		model.BaseListParams{Group: "flavors.id"},
		func(stmt sq.SelectBuilder) sq.SelectBuilder {
			return stmt.Where(sq.Expr("flavors.deleted_at IS NULL"))
		})
	fmt.Println(cnt)
	if err != nil {
		return nil, 0, err
	}
	return flavors, cnt, err
}

// Insert adds a new flavor to the database
func (r *flavorStore) Insert(flavor model.Flavor, tx sqlx.Execer) error {
	return r.store.insertEntity(tx, &flavor, func(id int64) {
		flavor.ID = uint64(id)
	})
}

func (r *flavorStore) Update(flavor model.Flavor, tx sqlx.Execer) error {
	_, err := r.store.updateEntity(tx, &flavor)
	return err
}

func (r *flavorStore) SaveStash(stashes []model.FlavorStash) error {
	ins := sq.Insert("flavor_stashes").
		Columns("flavor_id", "owner_id", "on_hand_m", "density_m", "vg", "rating")
	for _, s := range stashes {
		ins = ins.Values(s.FlavorID, s.OwnerID, s.OnHandM, s.DensityM, s.Vg, s.Rating)
	}
	q, args, err := ins.ToSql()
	if err != nil {
		return err
	}
	_, err = r.store.DB().Exec(q, args)
	return err
}
