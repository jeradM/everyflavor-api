package core

import (
	"everyflavor/internal/http/api/v1/view"
	"everyflavor/internal/mapper"
	"everyflavor/internal/storage/model"

	"github.com/pkg/errors"
)

func (a *App) GetUserByID(id uint64) (view.User, error) {
	u, err := a.Store.User().Get(id)
	if err != nil {
		return view.User{}, err
	}
	stats, err := a.GetUserStatsByID(id)
	if err != nil {
		return view.User{}, err
	}
	roles, err := a.GetRolesForUsers([]uint64{id})
	if err != nil {
		return view.User{}, err
	}
	return mapper.UserFromModel(*u, roles, stats), nil
}

func (a *App) GetUserList() ([]model.User, error) {
	return a.Store.User().List()
}

func (a *App) SaveUser(v view.User) error {
	tx, err := a.Store.Connection().Beginx()
	if err != nil {
		return errors.Wrap(err, "failed to start transaction")
	}
	u, _ := mapper.UserToModel(v)
	if u.ID > 0 {
		err = a.Store.User().Update(&u, tx)
	} else {
		err = a.Store.User().Insert(&u, tx)
	}
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (a *App) UpdateUserPassword(id uint64, pw string) error {
	return a.Store.User().UpdatePassword(id, pw)
}

func (a *App) GetUserByUsername(username string) (view.User, error) {
	u, err := a.Store.User().FindByUsername(username)
	if err != nil {
		return view.User{}, err
	}
	v := mapper.UserFromModel(*u, nil, nil)
	return v, err
}

func (a *App) SearchUsersByUsername(username string) ([]view.User, error) {
	u, err := a.Store.User().FindAllByUsernameLike(username)
	if err != nil {
		return nil, err
	}
	users := make([]view.User, len(u))
	for idx, user := range u {
		users[idx] = mapper.UserFromModel(user, nil, nil)
		users[idx].Password = ""
	}
	return users, err
}

func (a *App) GetUserStatsByID(id uint64) (*model.UserStats, error) {
	return a.Store.User().GetStats(id)
}

func (a *App) GetRolesForUsers(ids []uint64) ([]model.UserRole, error) {
	return a.Store.User().ListRoles(ids)
}
