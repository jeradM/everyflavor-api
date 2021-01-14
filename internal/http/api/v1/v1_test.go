package v1

import (
	"encoding/json"
	"everyflavor/internal/http/api/v1/view"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
)

func setupTestRouter(w http.ResponseWriter, u *view.User) (*gin.Context, *gin.Engine) {
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
	return c, g
}

func parseRespMapString(w *httptest.ResponseRecorder) (map[string]string, error) {
	var resp map[string]string
	err := json.Unmarshal([]byte(w.Body.String()), &resp)
	return resp, err
}
