package v1

import (
	"everyflavor/internal/core"
	"everyflavor/internal/storage/model"
	"strconv"

	"github.com/gin-gonic/gin"
)

func SetupVendorHandlers(r gin.IRouter, s core.VendorService) {
	r.GET("/vendors", getVendors(s))
	r.POST("/vendor", EnsureRole("admin"), addVendor(s))
	r.GET("/vendor/:id", getVendor(s))
	r.PUT("/vendor/:id", EnsureRole("admin"), updateVendor(s))
}

func addVendor(s core.VendorService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var vendor model.Vendor
		err := c.BindJSON(&vendor)
		if err != nil {
			badRequest(c, respError(err, "unable to parse request body"))
			return
		}
		err = s.SaveVendor(&vendor)
		if err != nil {
			serverError(c, respError(err, "insert failed"))
			return
		}

		created(c, &vendor)
	}
}

func getVendor(s core.VendorService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			badRequest(c, respError(err, "invalid id"))
			return
		}
		vendor, err := s.GetVendorByID(id)
		if err != nil {
			notFound(c)
			return
		}
		ok(c, &vendor)
	}
}

func getVendors(s core.VendorService) gin.HandlerFunc {
	return func(c *gin.Context) {
		vendors, err := s.GetVendorsList()
		if err != nil {
			badRequest(c, respError(err, err.Error()))
			return
		}
		ok(c, vendors)
	}
}

func updateVendor(s core.VendorService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			badRequest(c, respError(err, "invalid id"))
			return
		}
		var vendor model.Vendor
		if err = c.BindJSON(&vendor); err != nil {
			badRequest(c, respError(err, "unable to parse request body"))
			return
		}
		vendor.ID = id
		if err = s.SaveVendor(&vendor); err != nil {
			serverError(c, respError(err, "update failed"))
			return
		}
		ok(c, &vendor)
	}
}
