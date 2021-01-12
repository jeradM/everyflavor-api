package v1

import (
	"everyflavor/internal/core"
	"everyflavor/internal/http/api/v1/view"
	"everyflavor/internal/storage/model"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type RecipeSearchParams struct {
	Title string `form:"title"`
}

func setupRecipeHandlers(router gin.IRouter, s core.Application) {
	router.GET("/recipe/:id", CanViewRecipe(s), getRecipe(s))
	router.POST("/recipe/:id/comment", EnsureLoggedIn, postComment(s))
	router.PUT("/recipe/:id", CanEditRecipe(s), updateRecipe(s))
	router.GET("/recipes", getRecipes(s, false))
	router.POST("/recipes", EnsureLoggedIn, saveRecipe(s))
	router.GET("/recipes/mine", EnsureLoggedIn, getRecipes(s, true))
	router.GET("/tags", getTags(s))
}

func getRecipes(app core.RecipeService, mine bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		p := &model.RecipeParams{}
		if c.ShouldBind(p) != nil {
			badRequest(c, respError(errors.New("invalid parameters"), ""))
			return
		}
		if mine {
			user, _ := GetUserFromSession(c)
			p.UserID = &user.ID
			p.Public = false
		} else {
			p.Public = true
		}
		p.Current = true
		result, err := app.GetRecipesList(p)
		if err != nil {
			badRequest(c, respError(err, err.Error()))
			return
		}
		ok(c, result)
	}
}

func getRecipe(s core.RecipeService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			badRequest(c, respError(err, "invalid id"))
			return
		}
		recipe, err := s.GetRecipeByID(id)
		if err != nil {
			notFound(c)
			return
		}
		ok(c, &recipe)
	}
}

func saveRecipe(s core.RecipeService) gin.HandlerFunc {
	return func(c *gin.Context) {
		u, err := GetUserFromSession(c)
		if err != nil {
			unauthorized(c, respError(err, err.Error()))
			return
		}
		var recipe view.Recipe
		if err := c.ShouldBind(&recipe); err != nil {
			badRequest(c, respError(err, err.Error()))
			return
		}
		recipe.OwnerID = u.ID
		recipe, err = s.SaveRecipe(recipe)
		if err != nil {
			badRequest(c, respError(err, err.Error()))
			return
		}
		ok(c, recipe)
	}
}

func updateRecipe(s core.RecipeService) gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err := GetUserFromSession(c)
		if err != nil {
			unauthorized(c, respError(err, err.Error()))
			return
		}
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			badRequest(c, respError(err, "invalid id"))
			return
		}
		var recipe view.Recipe
		if err := c.ShouldBind(&recipe); err != nil {
			badRequest(c, respError(err, err.Error()))
			return
		}
		if recipe.ID != id {
			badRequest(c, respError(errors.New("recipe id does not match request"), ""))
			return
		}
		recipe, err = s.UpdateRecipe(recipe)
		if err != nil {
			serverError(c, respError(err, "failed to save recipe"))
			return
		}
		ok(c, recipe)
	}
}

func postComment(s core.RecipeService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			badRequest(c, respError(err, "invalid id"))
			return
		}
		u, _ := GetUserFromSession(c)
		rc := model.RecipeComment{OwnerID: u.ID, RecipeID: id}
		if err = c.ShouldBind(&rc); err != nil {
			badRequest(c, respError(err, err.Error()))
			return
		}
		//err = s.(&rc)
		//if err != nil {
		//	serverError(c, err)
		//	return
		//}
		created(c, nil)
	}
}

func getTags(s core.TagService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tags, err := s.ListTags()
		if err != nil {
			serverError(c, respError(err, "error fetching tags"))
			return
		}
		ok(c, tags)
	}
}
