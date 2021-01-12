package v1

import (
	"encoding/json"
	"everyflavor/internal/core"
	"everyflavor/internal/core/mocks"
	"everyflavor/internal/http/api/v1/view"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

var testUser = view.User{
	ID:       1,
	Username: "testuser",
	Email:    "testuser@test.com",
	Password: "skdjflsdfj",
	Roles: []view.UserRole{
		{ID: 1, Authority: "user"},
	},
}

func setupRecipeRouter(app core.Application, u *view.User) (*httptest.ResponseRecorder, *gin.Context, *gin.Engine) {
	w := httptest.NewRecorder()
	c, g := setupTestRouter(w, u)
	recipeGroup := g.Group("/api/v1")
	setupRecipeHandlers(recipeGroup, app)
	return w, c, g
}

func parseJSONResponse(w *httptest.ResponseRecorder) (map[string]interface{}, error) {
	var response map[string]interface{}
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	return response, err
}

func TestGetRecipe(t *testing.T) {
	app := new(mocks.Application)
	app.On("GetRecipeByID", uint64(1)).Return(&view.Recipe{OwnerID: 1, Public: true, Title: "test title"}, nil)
	app.On("CanViewRecipe", uint64(1), uint64(1)).Return(true, nil)

	w, _, router := setupRecipeRouter(app, &testUser)

	r, _ := http.NewRequest("GET", "/api/v1/recipe/1", nil)
	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)

	response, err := parseJSONResponse(w)
	assert.NoError(t, err)

	data, exists := response["data"]
	require.True(t, exists, "missing response.data")
	dataMap, ok := data.(map[string]interface{})
	require.True(t, ok)

	title, exists := dataMap["title"]
	require.True(t, exists, "missing response.title")
	titleStr, ok := title.(string)
	require.True(t, ok, "response.title not string")
	assert.Equal(t, "test title", titleStr)
}

func TestRecipeNotFound(t *testing.T) {
	app := new(mocks.Application)
	app.On("CanViewRecipe", mock.Anything, mock.Anything).Return(true, nil)
	app.On("GetRecipeByID", uint64(1)).Return(&view.Recipe{}, errors.New("error"))

	w, _, router := setupRecipeRouter(app, &testUser)

	r, _ := http.NewRequest("GET", "/api/v1/recipe/1", nil)
	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetRecipeForbidden(t *testing.T) {
	app := new(mocks.Application)
	app.On("CanViewRecipe", mock.Anything, mock.Anything).Return(false, nil)

	w, _, router := setupRecipeRouter(app, &testUser)

	r, _ := http.NewRequest("GET", "/api/v1/recipe/1", nil)
	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusForbidden, w.Code)
}
