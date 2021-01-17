package v1

import (
	"everyflavor/internal/http/api/v1/view"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupMiddlewareRouter(u *view.User) (*httptest.ResponseRecorder, *gin.Context, *gin.Engine) {
	w := httptest.NewRecorder()
	c, g := setupTestRouter(w, u)
	return w, c, g
}

func TestEnsureLoggedInSucceeds(t *testing.T) {
	w, _, g := setupMiddlewareRouter(&view.User{ID: 1})
	g.GET("EnsureLoggedIn", EnsureLoggedIn)

	r, _ := http.NewRequest("GET", "/EnsureLoggedIn", nil)
	g.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestEnsureLoggedInUnauthorized(t *testing.T) {
	w, _, g := setupMiddlewareRouter(nil)
	g.GET("EnsureLoggedIn", EnsureLoggedIn)

	r, _ := http.NewRequest("GET", "/EnsureLoggedIn", nil)
	g.ServeHTTP(w, r)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestEnsureRoleSucceeds(t *testing.T) {
	w, _, g := setupMiddlewareRouter(&view.User{Roles: []view.UserRole{{Authority: "user"}}})
	g.GET("EnsureUser", EnsureRole("user"))

	r, _ := http.NewRequest("GET", "/EnsureUser", nil)
	g.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestEnsureRoleSucceedsWithMultiple(t *testing.T) {
	w, _, g := setupMiddlewareRouter(&view.User{Roles: []view.UserRole{{Authority: "user"}, {Authority: "admin"}}})
	g.GET("EnsureAdmin", EnsureRole("admin"))

	r, _ := http.NewRequest("GET", "/EnsureAdmin", nil)
	g.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestEnsureRoleForbiddenWhenEmpty(t *testing.T) {
	w, _, g := setupMiddlewareRouter(&view.User{Roles: []view.UserRole{}})
	g.GET("EnsureRole", EnsureRole("role"))

	r, _ := http.NewRequest("GET", "/EnsureRole", nil)
	g.ServeHTTP(w, r)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestEnsureRoleForbiddenWhenNil(t *testing.T) {
	w, _, g := setupMiddlewareRouter(&view.User{Roles: nil})
	g.GET("EnsureRole", EnsureRole("role"))

	r, _ := http.NewRequest("GET", "/EnsureRole", nil)
	g.ServeHTTP(w, r)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestEnsureRoleForbidden(t *testing.T) {
	w, _, g := setupMiddlewareRouter(&view.User{Roles: []view.UserRole{{Authority: "user"}}})
	g.GET("EnsureAdmin", EnsureRole("admin"))

	r, _ := http.NewRequest("GET", "/EnsureAdmin", nil)
	g.ServeHTTP(w, r)
	assert.Equal(t, http.StatusForbidden, w.Code)
}
