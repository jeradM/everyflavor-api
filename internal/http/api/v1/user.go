package v1

import (
	"everyflavor/internal/core"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func SetupUserHandlers(r gin.IRouter, s core.UserService) {
	r.GET("/user", me)
	r.GET("/user/:id", getUser(s))
	r.GET("/user/:id/stats", stats(s))
	r.GET("/users/search", searchUsers(s))
}

func getUser(s core.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			badRequest(c, respError(err, "invalid id"))
			return
		}
		user, err := s.GetUserByID(id)
		if err != nil {
			notFound(c)
			return
		}
		u, err := GetUserFromSession(c)
		if err != nil || u.ID != id {
			user.Email = ""
			user.Roles = nil
		}
		ok(c, &user)
	}
}

func searchUsers(s core.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		un := c.Query("username")
		if un == "" {
			badRequest(c, respError(errors.New("invalid username"), ""))
			return
		}
		u, err := s.SearchUsersByUsername(un)
		if err != nil {
			serverError(c, respError(err, err.Error()))
			return
		}

		type collab struct {
			ID       uint64 `json:"id"`
			Username string `json:"username"`
		}
		s := make([]collab, len(u))
		for idx, user := range u {
			s[idx] = collab{
				ID:       user.ID,
				Username: user.Username,
			}
		}
		ok(c, s)
	}
}

func me(c *gin.Context) {
	u, _ := GetUserFromSession(c)
	ok(c, u)
}

func stats(s core.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			badRequest(c, respError(err, "invalid id"))
			return
		}
		stats, err := s.GetUserStatsByID(id)
		if err != nil {
			serverError(c, respError(err, err.Error()))
			return
		}
		ok(c, stats)
	}
}
