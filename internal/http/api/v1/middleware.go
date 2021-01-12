package v1

import (
	"database/sql"
	"errors"
	"everyflavor/internal/core"
	"everyflavor/internal/http/api/v1/view"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// GetUserFromSession extracts the current user from default session
// if one exists, or returns an error if current session is unauthenticated
func GetUserFromSession(c *gin.Context) (*view.User, error) {
	s := sessions.Default(c)
	u, ok := s.Get("Principal").(*view.User)
	if !ok || u == nil {
		return nil, errors.New("user not found")
	}
	return u, nil
}

// EnsureLoggedIn checks that the current session is authenticated
func EnsureLoggedIn(c *gin.Context) {
	_, err := GetUserFromSession(c)
	if err != nil {
		unauthorized(c, respError(errors.New("not logged in"), ""))
	}
}

// EnsureRole checks that the current user has been granted a specific role
func EnsureRole(r string) gin.HandlerFunc {
	return func(c *gin.Context) {
		u, err := GetUserFromSession(c)
		if err != nil {
			unauthorized(c, respError(err, err.Error()))
			return
		}
		hasRole := false
		for _, role := range u.Roles {
			if role.Authority == r {
				hasRole = true
				break
			}
		}
		if !hasRole {
			forbidden(c, respError(errors.New("forbidden"), ""))
		}
	}
}

func CanViewRecipe(s core.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			badRequest(c, respError(err, err.Error()))
			return
		}
		u, _ := GetUserFromSession(c)
		var userID uint64
		if u == nil {
			userID = 0
		} else {
			userID = u.ID
		}
		canView, err := s.CanViewRecipe(userID, id)
		if err != nil {
			if err == sql.ErrNoRows {
				notFound(c)
				return
			}
			serverError(c, respError(err, "an unknown error occurred"))
			return
		}
		if !canView {
			forbidden(c, respError(errors.New("forbidden"), ""))
		}
	}
}

func CanEditRecipe(s core.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			badRequest(c, respError(err, err.Error()))
			return
		}
		u, _ := GetUserFromSession(c)
		var userID uint64
		if u == nil {
			userID = 0
		} else {
			userID = u.ID
		}
		canEdit, err := s.CanEditRecipe(userID, id)
		if err != nil {
			if err == sql.ErrNoRows {
				notFound(c)
				return
			}
			serverError(c, respError(err, "an unknown error occurred"))
			return
		}
		if !canEdit {
			forbidden(c, respError(errors.New("forbidden"), ""))
		}
	}
}

func CanViewBatch(s core.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			badRequest(c, respError(err, err.Error()))
			return
		}
		u, _ := GetUserFromSession(c)
		var userID uint64
		if u == nil {
			userID = 0
		} else {
			userID = u.ID
		}
		canView, err := s.CanViewBatch(userID, id)
		if err != nil {
			serverError(c, respError(err, "an unknown error occurred"))
			return
		}
		if !canView {
			forbidden(c, respError(errors.New("forbidden"), ""))
		}
	}
}

func CanEditBatch(s core.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			badRequest(c, respError(err, err.Error()))
			return
		}
		u, _ := GetUserFromSession(c)
		var userID uint64
		if u == nil {
			userID = 0
		} else {
			userID = u.ID
		}
		canEdit, err := s.CanEditBatch(userID, id)
		if err != nil {
			if err == sql.ErrNoRows {
				notFound(c)
				return
			}
			serverError(c, respError(err, "an unknown error occurred"))
			return
		}
		if !canEdit {
			forbidden(c, respError(errors.New("forbidden"), ""))
		}
	}
}
