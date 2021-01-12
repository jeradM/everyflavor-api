package v1

import (
	"everyflavor/internal/core"
	"everyflavor/internal/http/api/v1/view"
	"strconv"

	"github.com/pkg/errors"

	"github.com/gin-gonic/gin"
)

// SetupBatchHandlers initialized handlers for batch routes
func SetupBatchHandlers(router gin.IRouter, s core.Application) {
	router.GET("/batch/:id", CanViewBatch(s), getBatch(s))
	router.GET("/batches", EnsureLoggedIn, getBatches(s))
	router.POST("/batches", EnsureLoggedIn, saveBatch(s))
	router.PUT("/batches/:id", CanEditBatch(s), updateBatch(s))
}

func getBatch(s core.BatchService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			badRequest(c, respError(err, "invalid id"))
			return
		}
		recipe, err := s.GetBatch(id)
		if err != nil {
			notFound(c)
			return
		}
		ok(c, &recipe)
	}
}

func getBatches(s core.BatchService) gin.HandlerFunc {
	return func(c *gin.Context) {
		u, _ := GetUserFromSession(c)
		b, err := s.GetBatchesForUser(u.ID)
		if err != nil {
			serverError(c, respError(err, ""))
			return
		}
		ok(c, ListResult{Results: b, Count: int64(len(b))})
	}
}

func saveBatch(s core.BatchService) gin.HandlerFunc {
	return func(c *gin.Context) {
		u, _ := GetUserFromSession(c)
		var batch view.Batch
		if err := c.ShouldBind(&batch); err != nil {
			badRequest(c, respError(err, ""))
			return
		}
		batch.OwnerID = u.ID
		batch, err := s.SaveBatch(batch)
		if err != nil {
			badRequest(c, respError(err, ""))
			return
		}
		ok(c, batch)
	}
}

func updateBatch(s core.BatchService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			badRequest(c, respError(err, "invalid id"))
			return
		}
		var batch view.Batch
		if err := c.ShouldBind(&batch); err != nil {
			badRequest(c, respError(err, err.Error()))
			return
		}
		if batch.ID != id {
			badRequest(c, respError(errors.New("batch id does not match request"), ""))
			return
		}
		if _, err := s.SaveBatch(batch); err != nil {
			serverError(c, respError(err, ""))
			return
		}
		ok(c, batch)
	}
}
