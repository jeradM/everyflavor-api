package v1

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type ApiV1Response struct {
	Data    interface{} `json:"data"`
	Status  int         `json:"statusCode"`
	Message string      `json:"message"`
}

type ListResult struct {
	Results interface{} `json:"result"`
	Count   int64       `json:"count"`
}

type httpError struct {
	err error
	msg string
}

const (
	statusCodeKey    = "statusCode"
	requestURIKey    = "path"
	requestMethodKey = "method"
	requestHostKey   = "host"
	responseTimeKey  = "responseTime"
	userIDKey        = "userId"
	usernameKey      = "username"
)

var RequestLoggingEnabled = true

func respError(err error, msg string) httpError {
	return httpError{err, msg}
}

func (e httpError) Error() string {
	return errors.Wrap(e.err, e.msg).Error()
}

func respond(c *gin.Context, status int, data interface{}, message string) {
	response := ApiV1Response{Data: data, Status: status, Message: message}
	if status < http.StatusBadRequest {
		c.JSON(status, response)
	} else {
		c.AbortWithStatusJSON(status, response)
	}
}

func ok(c *gin.Context, data interface{}) {
	logEvent(c, http.StatusOK, log.Debug())
	respond(c, http.StatusOK, data, "success")
}

func created(c *gin.Context, data interface{}) {
	logEvent(c, http.StatusCreated, log.Debug())
	respond(c, http.StatusCreated, data, "created")
}

func badRequest(c *gin.Context, e httpError) {
	logEvent(c, http.StatusBadRequest, log.Warn().Err(e))
	respond(c, http.StatusBadRequest, nil, e.msg)
}

func unauthorized(c *gin.Context, e httpError) {
	logEvent(c, http.StatusUnauthorized, log.Warn().Err(e))
	respond(c, http.StatusUnauthorized, nil, "unauthorized")
}

func forbidden(c *gin.Context, e httpError) {
	logEvent(c, http.StatusForbidden, log.Warn().Err(e))
	respond(c, http.StatusForbidden, nil, e.msg)
}

func notFound(c *gin.Context) {
	logEvent(c, http.StatusNotFound, log.Info())
	respond(c, http.StatusNotFound, nil, "not found")
}

func serverError(c *gin.Context, e error) {
	logEvent(c, http.StatusInternalServerError, log.Error().Err(e))
	respond(c, http.StatusInternalServerError, nil, "An unknown error occurred")
}

func logEvent(c *gin.Context, s int, e *zerolog.Event) {
	if !RequestLoggingEnabled {
		return
	}
	u, _ := GetUserFromSession(c)
	if u != nil {
		e.Uint64(userIDKey, u.ID).Str(usernameKey, u.Username)
	}
	e = e.Int(statusCodeKey, s).
		Str(requestURIKey, c.Request.RequestURI).
		Str(requestMethodKey, c.Request.Method).
		Str(requestHostKey, c.Request.Host)

	var startTime time.Time
	st, ok := c.Get("startTime")
	if ok {
		startTime, ok = st.(time.Time)
	}
	if ok {
		e = e.Dur(responseTimeKey, time.Since(startTime))
	}
}
