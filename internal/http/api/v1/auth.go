package v1

import (
	"everyflavor/internal/core"
	"everyflavor/internal/http/api/v1/view"
	"everyflavor/internal/storage/model"
	"github.com/pkg/errors"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

type RegisterCmdObj struct {
	Username  string
	Password  string
	PasswordC string
	Email     string
}

func (r RegisterCmdObj) validate(s core.UserService) error {
	if r.Username == "" {
		return errors.New("Username is required")
	}
	if r.Email == "" {
		return errors.New("Email is required")
	}
	if r.Password == "" {
		return errors.New("Password is required")
	}
	if r.Password != r.PasswordC {
		return errors.New("Passwords don't match")
	}
	if s.UsernameExists(r.Username) {
		return errors.New("Username already taken")
	}
	if s.EmailExists(r.Email) {
		return errors.New("Email address already taken")
	}
	return nil
}

type LoginCmdObj struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ChangePasswordCmdObj struct {
	CurrentPassword string
	NewPassword     string
	NewPasswordC    string
}

func (c ChangePasswordCmdObj) validate() error {
	if c.CurrentPassword == "" {
		return errors.New("Current password is required")
	}
	if c.NewPassword != c.NewPasswordC {
		return errors.New("Passwords don't match")
	}
	return nil
}

type RegisterResponse struct {
	Username *string `json:"username"`
	Message  string  `json:"message"`
}

func SetupAuthHandlers(r gin.IRouter, s *core.Server) {
	auth := r.Group("/auth")
	auth.POST("/register", register(s.App))
	auth.POST("/authenticate", authenticate(s.App))
	auth.GET("/authenticate", authenticate(s.App))
	auth.POST("/logout", logout)
	auth.POST("/changePassword", changePassword(s.App))
}

func register(s core.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var userObj RegisterCmdObj
		if err := c.BindJSON(&userObj); err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		if err := userObj.validate(s); err != nil {
			c.JSON(http.StatusBadRequest, RegisterResponse{Message: err.Error()})
			return
		}
		u := view.User{
			Username: userObj.Username,
			Email:    userObj.Email,
			Password: userObj.Password,
		}
		err := s.SaveUser(u)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, RegisterResponse{Message: err.Error()})
			return
		}
		c.JSON(http.StatusCreated, &RegisterResponse{
			Username: &u.Username,
			Message:  "user created",
		})
	}
}

func authenticate(s core.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var login LoginCmdObj
		err := c.ShouldBind(&login)
		if err != nil {
			c.JSON(http.StatusUnauthorized, "Login failed")
			return
		}
		u, err := s.GetUserByUsername(login.Username)
		if err != nil {
			c.JSON(http.StatusUnauthorized, "Login failed")
			return
		}
		err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(login.Password))
		if err != nil {
			c.JSON(http.StatusUnauthorized, "Login failed")
			return
		}
		session := sessions.Default(c)
		u.Password = ""
		session.Set("Principal", u)
		err = session.Save()
		if err != nil {
			log.Debug().Msg(err.Error())
			c.JSON(400, err.Error())
			return
		}
		user, err := s.GetUserByID(u.ID)
		if err != nil {
			log.Error().Err(err).Msg("")
			c.JSON(500, err.Error())
			return
		}
		c.JSON(http.StatusOK, user)
	}
}

func changePassword(s core.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		u, ok := session.Get("Principal").(*model.User)
		if !ok || u == nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		var o ChangePasswordCmdObj
		if err := c.ShouldBind(&o); err != nil {
			badRequest(c, respError(err, "Bad request"))
			return
		}
		if err := o.validate(); err != nil {
			badRequest(c, respError(err, err.Error()))
			return
		}
		err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(o.NewPassword))
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		err = s.UpdateUserPassword(u.ID, o.NewPassword)
		if err != nil {
			serverError(c, httpError{err: err})
			return
		}
		c.JSON(http.StatusOK, "Password changed")
	}
}

func logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Options(sessions.Options{MaxAge: -1})
	if err := session.Save(); err != nil {
		serverError(c, httpError{err: err})
		return
	}
	c.JSON(http.StatusOK, "logged out")
}
