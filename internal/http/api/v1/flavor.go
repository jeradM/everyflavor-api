package v1

import (
	"everyflavor/internal/core"
	"everyflavor/internal/http/api/v1/view"
	"strconv"

	"github.com/gin-gonic/gin"
)

// SetupFlavorHandlers setup handlers for flavor routes
func SetupFlavorHandlers(r gin.IRouter, s core.FlavorService) {
	r.GET("/flavors", getFlavors(s))
	r.POST("/flavors", addFlavor(s))
	r.GET("/flavor/:id", getFlavor(s))
	r.GET("/flavors/listCount", listCount(s))
	r.GET("/flavors/stash", getStash(s))
	r.POST("/flavors/stash", addStash(s))
	//r.PUT("/:id", updateVendor(s))
}

func listCount(s core.FlavorService) gin.HandlerFunc {
	return func(c *gin.Context) {
		r, err := s.GetFlavorsList()
		if err != nil {
			badRequest(c, respError(err, ""))
			return
		}
		ok(c, r)
	}
}

func getFlavors(s core.FlavorService) gin.HandlerFunc {
	return func(c *gin.Context) {
		flavors, err := s.GetFlavorsList()
		if err != nil {
			serverError(c, respError(err, ""))
			return
		}
		ok(c, flavors)
	}
}

func getFlavor(s core.FlavorService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			badRequest(c, respError(err, "invalid id"))
			return
		}
		flavor, err := s.GetFlavorByID(id)
		if err != nil {
			notFound(c)
			return
		}
		ok(c, flavor)
	}
}

func addFlavor(s core.FlavorService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var flavor view.Flavor
		err := c.BindJSON(&flavor)
		if err != nil {
			badRequest(c, respError(err, "unable to parse request body"))
			return
		}
		err = s.SaveFlavor(flavor)
		if err != nil {
			serverError(c, respError(err, "an error occurred"))
			return
		}
		ok(c, &flavor)
	}
}

func getStash(s core.FlavorService) gin.HandlerFunc {
	return func(c *gin.Context) {
		u, err := GetUserFromSession(c)
		if err != nil {
			unauthorized(c, respError(err, err.Error()))
			return
		}
		f, err := s.GetStashForUser(u.ID)
		if err != nil {
			serverError(c, respError(err, err.Error()))
			return
		}
		ok(c, f)
	}
}

func addStash(s core.FlavorService) gin.HandlerFunc {
	return func(c *gin.Context) {
		u, err := GetUserFromSession(c)
		if err != nil {
			unauthorized(c, respError(err, err.Error()))
			return
		}
		var fs view.FlavorStash
		if c.ShouldBind(&fs) != nil {
			badRequest(c, respError(err, err.Error()))
			return
		}
		fs.OwnerID = u.ID
		if s.SaveStash(fs) != nil {
			serverError(c, respError(err, err.Error()))
			return
		}
		created(c, &fs)
	}
}
