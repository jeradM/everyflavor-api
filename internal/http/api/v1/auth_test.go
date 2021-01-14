package v1

import (
	"encoding/gob"
	"encoding/json"
	"everyflavor/internal/core/mocks"
	"everyflavor/internal/http/api/v1/view"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func setupAuthRouter(u *view.User) (*httptest.ResponseRecorder, *gin.Context, *gin.Engine) {
	w := httptest.NewRecorder()
	RequestLoggingEnabled = false
	gin.SetMode(gin.ReleaseMode)

	c, g := gin.CreateTestContext(w)
	g.Use(sessions.Sessions("efsession", memstore.NewStore([]byte("secret"))))

	if u != nil {
		g.Use(func(c *gin.Context) {
			s := sessions.Default(c)
			s.Set("Principal", u)
		})
	}
	return w, c, g
}

func TestRegisterSuccess(t *testing.T) {
	w, _, g := setupAuthRouter(nil)

	cmd := RegisterCmdObj{
		Username:  "testuser",
		Password:  "testpw",
		PasswordC: "testpw",
		Email:     "testuser@test.com",
	}

	s := new(mocks.UserService)
	s.On("UsernameExists", mock.Anything).Return(false)
	s.On("EmailExists", mock.Anything).Return(false)
	s.On("SaveUser", mock.AnythingOfType("view.User")).Return(nil)
	g.POST("/register", register(s))

	body, _ := json.Marshal(cmd)
	r, _ := http.NewRequest(http.MethodPost, "/register", strings.NewReader(string(body)))

	g.ServeHTTP(w, r)
	assert.Equal(t, http.StatusCreated, w.Code)

	response, err := parseRespMapString(w)
	assert.NoError(t, err)

	message, _ := response["message"]
	assert.Equal(t, message, "user created")

	username, _ := response["username"]
	assert.Equal(t, username, "testuser")
}

func TestRegisterFails_PasswordsDontMatch(t *testing.T) {
	w, _, g := setupAuthRouter(nil)

	cmd := RegisterCmdObj{
		Username:  "testuser",
		Password:  "testpw",
		PasswordC: "doesntmatch",
		Email:     "testuser@test.com",
	}

	s := new(mocks.UserService)
	g.POST("/register", register(s))

	body, _ := json.Marshal(cmd)
	r, _ := http.NewRequest(http.MethodPost, "/register", strings.NewReader(string(body)))

	g.ServeHTTP(w, r)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	response, err := parseRespMapString(w)
	assert.NoError(t, err)

	message, _ := response["message"]
	assert.Equal(t, "Passwords don't match", message)

	username, _ := response["username"]
	assert.Equal(t, "", username)
}

func TestRegisterFails_EmailExists(t *testing.T) {
	w, _, g := setupAuthRouter(nil)

	cmd := RegisterCmdObj{
		Username:  "testuser",
		Password:  "testpw",
		PasswordC: "testpw",
		Email:     "testuser@test.com",
	}

	s := new(mocks.UserService)
	s.On("UsernameExists", mock.Anything).Return(false)
	s.On("EmailExists", mock.Anything).Return(true)
	g.POST("/register", register(s))

	body, _ := json.Marshal(cmd)
	r, _ := http.NewRequest(http.MethodPost, "/register", strings.NewReader(string(body)))

	g.ServeHTTP(w, r)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	response, err := parseRespMapString(w)
	assert.NoError(t, err)

	message, _ := response["message"]
	assert.Equal(t, "Email address already taken", message)

	username, _ := response["username"]
	assert.Equal(t, "", username)
}

func TestRegisterFails_UsernameExists(t *testing.T) {
	w, _, g := setupAuthRouter(nil)

	cmd := RegisterCmdObj{
		Username:  "testuser",
		Password:  "testpw",
		PasswordC: "testpw",
		Email:     "testuser@test.com",
	}

	s := new(mocks.UserService)
	s.On("UsernameExists", mock.Anything).Return(true)
	g.POST("/register", register(s))

	body, _ := json.Marshal(cmd)
	r, _ := http.NewRequest(http.MethodPost, "/register", strings.NewReader(string(body)))

	g.ServeHTTP(w, r)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	response, err := parseRespMapString(w)
	assert.NoError(t, err)

	message, _ := response["message"]
	assert.Equal(t, "Username already taken", message)

	username, _ := response["username"]
	assert.Equal(t, "", username)
}

func TestAuthenticateSucceeds(t *testing.T) {
	gob.Register(&view.User{})
	w, _, g := setupAuthRouter(nil)
	hashedPw := "$2a$10$LC/tKbaaD/Et2/lQesHtGubbsS2giWJbSuy77FyxF8Iprgy/0Caxi" // testpw
	user := view.User{Password: hashedPw}

	s := new(mocks.UserService)
	s.On("GetUserByUsername", mock.Anything).Return(user, nil)
	s.On("GetUserByID", mock.Anything).Return(user, nil)

	g.POST("authenticate", authenticate(s))

	loginCmd := LoginCmdObj{
		Username: "testuser",
		Password: "testpw",
	}
	cmdData, _ := json.Marshal(loginCmd)
	body := strings.NewReader(string(cmdData))
	r, _ := http.NewRequest(http.MethodPost, "/authenticate", body)
	r.Header.Add("Content-Type", "application/json")

	g.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
}
func TestAuthenticateFails_WrongPassword(t *testing.T) {
	gob.Register(&view.User{})
	w, _, g := setupAuthRouter(nil)
	hashedPw := "$2a$10$LC/tKbaaD/Et2/lQesHtGubbsS2giWJbSuy77FyxF8Iprgy/0Caxi" // testpw
	user := view.User{Password: hashedPw}

	s := new(mocks.UserService)
	s.On("GetUserByUsername", mock.Anything).Return(user, nil)
	s.On("GetUserByID", mock.Anything).Return(user, nil)

	g.POST("authenticate", authenticate(s))

	loginCmd := LoginCmdObj{
		Username: "testuser",
		Password: "wrongpassword",
	}
	cmdData, _ := json.Marshal(loginCmd)
	body := strings.NewReader(string(cmdData))
	r, _ := http.NewRequest(http.MethodPost, "/authenticate", body)
	r.Header.Add("Content-Type", "application/json")

	g.ServeHTTP(w, r)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthenticateFails_NoUser(t *testing.T) {
	gob.Register(&view.User{})
	w, _, g := setupAuthRouter(nil)

	s := new(mocks.UserService)
	var u view.User
	s.On("GetUserByUsername", mock.AnythingOfType("string")).Return(u, errors.New("error"))

	g.POST("authenticate", authenticate(s))

	loginCmd := LoginCmdObj{
		Username: "non-existent-user",
	}
	cmdData, _ := json.Marshal(loginCmd)
	body := strings.NewReader(string(cmdData))
	r, _ := http.NewRequest(http.MethodPost, "/authenticate", body)
	r.Header.Add("Content-Type", "application/json")

	g.ServeHTTP(w, r)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
