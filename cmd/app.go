package cmd

import (
	"encoding/gob"
	"everyflavor/internal/core"
	"everyflavor/internal/http/api"
	"everyflavor/internal/http/api/v1/view"
	"everyflavor/internal/storage/mysql"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	sredis "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	cors "github.com/rs/cors/wrapper/gin"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	lredis "github.com/ulule/limiter/v3/drivers/store/redis"
)

var appCmd = &cobra.Command{
	Use:   "app",
	Short: "Run everyflavor web application",
	Run: func(cmd *cobra.Command, args []string) {
		rand.Seed(time.Now().UnixNano())

		config := mustLoadConfig()
		logger := core.MustInitLogging(config)
		db := core.MustSetupMySQLDatabase(config)
		app := core.NewApp(config, mysql.NewMySQLStore(db, config.ShowSQL))
		server := &core.Server{Config: &config, App: app, Logger: logger}

		mustInitRedis(server)
		initHttp(server)

		gob.Register(&view.User{})
		log.Panic().Err(server.Run(config.ServerAddr)).Msg("")
	},
}

func mustLoadConfig() core.AppConfig {
	var config core.AppConfig
	if err := viper.Unmarshal(&config); err != nil {
		log.Panic().Err(err).Msg("Failed to parse config")
	}
	return config
}

func mustInitRedis(s *core.Server) {
	if s.Config.RedisURL == "" {
		log.Panic().Msg("Redis URL not defined")
	}
	opts := redis.Options{Network: "tcp", Addr: s.Config.RedisURL, DB: s.Config.RedisDB}
	client := redis.NewClient(&opts)
	s.Redis = client
}

func initHttp(s *core.Server) {
	keys := make([][]byte, len(s.Config.RedisKeys))
	for idx, key := range s.Config.RedisKeys {
		keys[idx] = []byte(key)
	}
	store, _ := sredis.NewStore(10, "tcp", s.Config.RedisURL, "", keys...)
	store.Options(sessions.Options{MaxAge: 3600, Path: "/"})
	s.Router = gin.New()
	s.Router.Use(
		mustInitRateLimitMiddleware(s.Redis),
		sessions.Sessions("efsession", store),
		gin.Recovery(),
		cors.New(cors.Options{
			AllowedOrigins:   s.Config.CorsAllowedOrigins,
			AllowCredentials: true,
			AllowedMethods: []string{
				http.MethodHead,
				http.MethodGet,
				http.MethodPost,
				http.MethodPut,
				http.MethodDelete,
			},
		}),
		func(c *gin.Context) {
			c.Set("startTime", time.Now())
		},
	)
	api.SetupHandlers(s)
}

func mustInitRateLimitMiddleware(r *redis.Client) gin.HandlerFunc {
	store, err := lredis.NewStoreWithOptions(r, limiter.StoreOptions{
		Prefix:   "rate_limit",
		MaxRetry: 3,
	})
	if err != nil {
		log.Panic().Msg("unable to create Redis rate-limit store")
	}
	rate := limiter.Rate{
		Period: 1 * time.Second,
		Limit:  10,
	}
	// Create a new middleware with the limiter instance.
	return mgin.NewMiddleware(limiter.New(store, rate))
}
