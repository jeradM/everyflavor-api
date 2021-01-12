package mapper

import (
	"database/sql"
	"everyflavor/internal/http/api/v1/view"
	"everyflavor/internal/storage/model"
)

func UserFromModel(u model.User, r []model.UserRole, s *model.UserStats) view.User {
	user := view.User{
		ID:        u.ID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		Username:  u.Username,
		Email:     u.Email,
		Password:  u.Password,
		Roles:     []view.UserRole{},
	}
	for _, ur := range r {
		user.Roles = append(user.Roles, view.UserRole{
			ID:        ur.RoleID,
			Authority: ur.Authority,
		})
	}
	if s != nil {
		user.Stats = view.UserStats{
			NumPublic:  s.NumPublic,
			NumPrivate: s.NumPrivate,
			AvgRating:  s.AvgRating,
		}
	}
	return user
}

func UserToModel(v view.User) (model.User, []model.UserRole) {
	u := model.User{
		ID:        v.ID,
		CreatedAt: v.CreatedAt,
		UpdatedAt: v.UpdatedAt,
		DeletedAt: sql.NullTime{},
		Username:  v.Username,
		Email:     v.Email,
		Password:  v.Password,
	}

	var roles []model.UserRole
	for _, r := range v.Roles {
		roles = append(roles, model.UserRole{
			UserID:    v.ID,
			RoleID:    r.ID,
			Authority: r.Authority,
		})
	}
	return u, roles
}
